package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	useChanelToPreventDataRacingIssue()
	fetchMultiAPIConcurrently()
	getFirstDoneRouter()

	gotDataRacingIssue()
	useMutexLockToPreventDataRacing()

}

func gotDataRacingIssue() {
	var count int = 0

	for i := 1; i <= 5; i++ {
		go func() {
			// Data racing right here because multiple routiens is writting on 1 variable
			for j := 1; j <= 10000; j++ {
				count += 1
			}
		}()
	}

	time.Sleep(time.Second * 7)
	fmt.Println("Count (Data racing issue => never be 50k): ", count)
}

func useMutexLockToPreventDataRacing() {
	var count int = 0
	/**
		RWMutex: block all the read and write Goroutines
		RLock: block all writting Goroutines but allow reading Goroutine to access (shared lock)
	**/
	lock := new(sync.RWMutex)

	for i := 1; i <= 5; i++ {
		go func() {
			for j := 1; j <= 10000; j++ {
				// this will lock the other routien try to get and update the "count" variable
				// Make the other wait util it been unlocked
				lock.Lock()
				count += 1
				// Unlock after update "count" variable so the other can get the updated value
				// and start to run
				lock.Unlock()
			}
		}()
	}

	time.Sleep(time.Second * 7)
	fmt.Println("Count:(gonna be 50k): ", count)
}

/**
	<- chan: return a chanel only use to read data
	chan <-: return a chanel only use to write data
**/
func startSender(name string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 1; i <= 5; i++ {
			c <- (name + " hello")
			time.Sleep(time.Millisecond * 200)
		}
	}()
	return c
}

func useChanelToPreventDataRacingIssue() {
	// Use chanel to handle data racing (Recommended)
	sender := startSender("Ti")
	for i := 1; i <= 5; i++ {
		fmt.Println(<-sender)
	}
}

func fetchAPI(model string) string {
	time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	return model
}

func fetchMultiAPIConcurrently() {
	responseChan := make(chan string)
	var results []string

	// Call multi apis concurrently
	go func() { responseChan <- fetchAPI("users") }()
	go func() { responseChan <- fetchAPI("categories") }()
	go func() { responseChan <- fetchAPI("products") }()

	for i := 1; i <= 3; i++ {
		results = append(results, <-responseChan)
	}
	fmt.Println(results)
}

func query(url string) string {
	time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
	return url
}

func queryFirst(servers ...string) <-chan string {
	c := make(chan string)
	for _, serv := range servers {
		go func(s string) { c <- query(s) }(serv)
	}
	return c
}

func getFirstDoneRouter() {
	result := queryFirst("server 1", "server 2", "server 3")
	fmt.Println(<-result)
}
