package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func getEventIDs(constructs string) []string {
	return []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22",
		"23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42"}
}

func main() {
	// get eventsByConstructName
	eventIDs := getEventIDs("c1")

	var wg sync.WaitGroup

	maxRoutines := 10

	numRoutines := func() int {
		if len(eventIDs) < maxRoutines {
			return len(eventIDs)
		}
		return maxRoutines
	}()

	ch := make(chan string, numRoutines)

	wg.Add(numRoutines)
	for i := 1; i <= numRoutines; i++ {
		go func(iInRoutine int) {
			for {
				id, ok := <-ch
				if !ok {
					logrus.Printf("all IDs are processed. Closing the routine %v", iInRoutine)
					wg.Done()
					return
				}
				logrus.Printf("processing ID %v in go routine %v", id, iInRoutine)
				processID(id)
				logrus.Printf("done processing ID %v in go routine %v", id, iInRoutine)
			}
		}(i)
	}

	for i := 0; i < len(eventIDs); i++ {
		logrus.Printf("sending %v: %v", i, eventIDs[i])
		ch <- eventIDs[i]
		logrus.Printf("sent %v: %v", i, eventIDs[i])
	}

	close(ch)
	wg.Wait()
}

func processID(id string) {
	sleepDur := time.Duration(rand.Intn(10)) * time.Second
	time.Sleep(sleepDur)
	logrus.Printf("sleeping for :%v seconds", sleepDur)
	logrus.Printf("done processing id: %v", id)
	return
}
