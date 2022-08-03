package main

import (
	"fmt"
	"time"
)

func main() {
	maxWorker := 5
	numberOfRequest := 100
	queueChan := make(chan int, numberOfRequest)
	doneChan := make(chan int)

	for i := 1; i <= maxWorker; i++ {
		go func(name string) {
			for v := range queueChan {
				// the for loop is run real fast
				// therefore the process with be hold and sometime only 1 process will be executed
				crawl(name, v)
			}

			// This print will never run if the channels not being closed
			fmt.Printf("%s is done\n", name)
			doneChan <- 1
		}(fmt.Sprintf("%d", i))

	}

	/*
		Place this loop after the async loop because i want to set up first
		and then push value to channel and run later
	*/
	for i := 1; i <= numberOfRequest; i++ {
		queueChan <- i
	}

	//if not close of the process in these channels will be hold forever
	close(queueChan)

	// this will hold the main thread to wait for all the worker is done
	for i := 1; i <= maxWorker; i++ {
		<-doneChan
	}
}

func crawl(name string, value int) {
	time.Sleep(time.Millisecond * 200)
	fmt.Printf("Worker %s is Crawling %d \n", name, value)
}
