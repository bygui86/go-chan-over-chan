package main

import (
	"fmt"
	"sync"
	"time"
)

/*
	Style 2: Using a Channel Stored Inside a Struct

	This code will look almost identical to the previous snippet, with 2 exceptions:
		- The ack channel will be stored inside a `struct`
		- The sleep time will be stored inside that same `struct`, so we can pass it over the `channel`
    		- This makes the code more flexible, because we can tell `worker` how long to sleep when we send to it, rather than when we start it
*/
func main() {
	// create the broadcast channel (channel-over-channel)
	broadcastChan := make(chan *process)
	// start workers
	for i := 1; i <= 10; i++ {
		go worker(i, broadcastChan)
	}

	// prepare an array of processes
	processes := make([]*process, 10)
	// send process to each worker
	for i := 1; i <= 10; i++ {
		processes[i-1] = &process{
			id:        i,
			sleepTime: time.Duration(i+1) * time.Second,
			ackChan:   make(chan time.Duration),
		}
		broadcastChan <- processes[i-1]
	}

	// wait for workers ack
	var waitGr sync.WaitGroup // use this to block until all goroutines have completed their stuffs
	for i := 1; i <= 10; i++ {
		go ackListener(i, &waitGr, processes[i-1].ackChan)
	}
	waitGr.Wait()

	fmt.Println("All Workers COMPLETED")
}

// process will be passed over the broadcast channel to a worker
type process struct {
	id        int
	sleepTime time.Duration
	ackChan   chan time.Duration
}

// worker will receive a 'process' struct on 'broadcastChan', sleep for 'proc.sleepTime', then send 'proc.sleepTime' on 'proc.ackChan'
func worker(id int, broadcastChan <-chan *process) {
	fmt.Printf("Worker %d START \n", id)
	defer fmt.Printf("Worker %d COMPLETED \n", id)

	proc := <-broadcastChan
	fmt.Printf("Worker %d received process with id %d \n", id, proc.id)
	time.Sleep(proc.sleepTime)
	proc.ackChan <- proc.sleepTime
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
