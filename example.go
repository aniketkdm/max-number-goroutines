package main

import (
	"fmt"
	"sync"
)

const xthreads = 5 // Total number of threads to use, excluding the main() thread

func doSomething(a int) {
	fmt.Println("My job is", a)
	return
}

func main() {
	var ch = make(chan int, 50) // This number 50 can be anything as long as it's larger than xthreads
	var wg sync.WaitGroup

	// This starts xthreads number of goroutines that wait for something to do
	wg.Add(xthreads)
	for i := 0; i < xthreads; i++ {
		go func() {
			for {
				a, ok := <-ch
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				doSomething(a) // do the thing
			}
		}()
	}

	// Now the jobs can be added to the channel, which is used as a queue
	for i := 0; i < 50; i++ {
		ch <- i // add i to the queue
	}

	close(ch) // This tells the goroutines there's nothing else to do
	wg.Wait() // Wait for the threads to finish
}
