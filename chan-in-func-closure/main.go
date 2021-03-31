package main

import (
	"fmt"
	"sync"
	"time"
)

/*
	Style 3: Using a Channel Inside a Function Closure

	This code will look different from the previous examples, because the `worker` function won’t know anything about a return channel.
	That fact is both good and bad. On the up side, you can change your code later to do anything you want inside that function
	(e.g. good for testing!), but on the down side, you can’t pass dynamic `time.Durations` into the `worker` goroutines, as you could in the previous example.
*/

func main() {
	// start up the worker goroutines
	broadcastChan := make(chan func(time.Duration))
	// start workers
	for i := 1; i <= 10; i++ {
		go worker(i, time.Duration(i+1)*time.Second, broadcastChan)
	}

	// prepare an array of ack channels
	ackChannels := make([]chan time.Duration, 10)
	// send ack function to each worker
	for i := 0; i < 10; i++ {
		ackChan := make(chan time.Duration)
		ackChannels[i] = ackChan
		// create functions that close over each channel, then send them to each worker
		ackFunc := func(sleepTime time.Duration) {
			ackChan <- sleepTime
		}
		broadcastChan <- ackFunc
	}

	// wait for workers ack
	var waitGr sync.WaitGroup // use this to block until all goroutines have completed their stuffs
	for i := 1; i <= 10; i++ {
		go ackListener(i, &waitGr, ackChannels[i-1])
	}
	waitGr.Wait()

	fmt.Println("All Workers COMPLETED")
}

// worker will receive an ack function 'ackFunc' to invoke in order to communicate that job is completed
func worker(id int, sleepTime time.Duration, broadcastChan <-chan func(time.Duration)) {
	fmt.Printf("Worker %d START \n", id)
	defer fmt.Printf("Worker %d COMPLETED \n", id)

	ackFunc := <-broadcastChan
	fmt.Printf("Worker %d received ack funcion \n", id)
	time.Sleep(sleepTime)
	ackFunc(sleepTime)
}

// ackListener receives the ack that workers send back
func ackListener(id int, waitGr *sync.WaitGroup, ackChan chan time.Duration) {
	fmt.Printf("AckListener %d START \n", id)
	defer fmt.Printf("AckListener %d COMPLETED \n", id)

	waitGr.Add(1)
	defer waitGr.Done()

	sleepTime := <-ackChan
	fmt.Printf("AckListener %d: a worker slept for %.0f seconds \n", id, sleepTime.Seconds())
}
