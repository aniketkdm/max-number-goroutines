package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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
	writeBatch(ch, wg)

	if err := wg.Wait(); err != nil {
		logrus.Fatal(err)
	}
}

func readBatch(scanner *bufio.Scanner, wg *errgroup.Group) <-chan string {
	ch := make(chan string, 100)

	wg.Go(func() error {
		for scanner.Scan() {
			ch <- scanner.Text()
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

func writeBatch(ch <-chan string, wg *errgroup.Group) error {
	wg.Go(func() error {
		for str := range ch {
			fmt.Printf("processing: %v\n", str)
		}
		return nil
	})
	return nil
}
