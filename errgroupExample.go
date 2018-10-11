package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	batchSize     = 100
	maxGoroutines = 7
)

func main() {
	pCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg, _ := errgroup.WithContext(pCtx)

	file, err := os.Open("exampleFile.out")
	if err != nil {
		logrus.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	ch := readBatch(scanner, wg)

	for i := 0; i < maxGoroutines; i++ {
		wg.Go(func() error {
			return writeBatch(ch)
		})
	}

	if err := wg.Wait(); err != nil {
		logrus.Fatal(err)
	}
}

func readBatch(scanner *bufio.Scanner, wg *errgroup.Group) <-chan []string {
	ch := make(chan []string, 5)
	var line string
	var batch []string

	wg.Go(func() error {
		for scanner.Scan() {
			line = scanner.Text()
			batch = append(batch, line)
			if len(batch) == 100 {
				ch <- batch
				batch = []string{}
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			close(ch)
			return err
		}
		close(ch)
		return nil
	})
	return ch
}

func writeBatch(ch <-chan []string) error {
	for batch := range ch {
		fmt.Println("received a batch")
		for _, str := range batch {
			fmt.Printf("processing: %v\n", str)
		}
	}
	return nil

	return nil
}

// with just batching: go run errgroupExample.go  0.25s user 0.29s system 52% cpu 1.043 total
// parallel of 3: go run errgroupExample.go  0.28s user 0.32s system 43% cpu 1.369 total
// parallel of 3: go run errgroupExample.go  0.30s user 0.35s system 50% cpu 1.284 total
