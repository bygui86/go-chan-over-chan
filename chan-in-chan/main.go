package main

import (
	"fmt"
	"sync"
	"time"
)

/*
	Style 1: Using a Channel Inside a Channel

	Here's the simplest of the patterns in action. Generally this style will be easiest to read and understand, but it has some limits:
		- Each `worker` goroutine sleeps for a set amount of time. You can't change the sleep time when you send on `ch`
		- Each `worker` goroutine can only receive a `chan time.Duration` – no more data than that
*/

func main() {
	// create the broadcast channel (channel-over-channel)
	broadcastChan := make(chan chan time.Duration)
	// start workers
	for i := 1; i <= 10; i++ {
		go worker(i, time.Duration(i+1)*time.Second, broadcastChan)
	}

	// prepare ack channel
	ackChan := make(chan time.Duration)
	// send channel to each worker
	for i := 0; i < 10; i++ {
		broadcastChan <- ackChan
	}

	// wait for workers ack
	var waitGr sync.WaitGroup // use this to block until all goroutines have completed their stuffs
	for i := 1; i <= 10; i++ {
		go ackListener(i, &waitGr, ackChan)
	}
	waitGr.Wait()

	fmt.Println("All Workers COMPLETED")
}

// worker receives a channel on 'broadcastChan', sleeps for 'sleepTime', then sends 'sleepTime' on the channel it received
func worker(id int, sleepTime time.Duration, broadcastChan <-chan chan time.Duration) {
	fmt.Printf("Worker %d START \n", id)
	defer fmt.Printf("Worker %d COMPLETED \n", id)

	ackChan := <-broadcastChan
	fmt.Printf("Worker %d received channel! \n", id)
	time.Sleep(sleepTime)
	ackChan <- sleepTime
}

// ackListener receives the ack that workers send back
func ackListener(id int, waitGr *sync.WaitGroup, ackChan chan time.Duration) {
	fmt.Printf("AckListener %d START \n", id)
	defer fmt.Printf("AckListener %d COMPLETED \n", id)

	waitGr.Add(1)
	defer waitGr.Done()

	dur := <-ackChan
	fmt.Printf("AckListener %d: a worker slept for %.0f seconds \n", id, dur.Seconds())
}
