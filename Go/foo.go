// Use `go run foo.go` to run your program

package main

import (
	"fmt"
	"runtime"
	"sync"
)

var i = 0
var wg sync.WaitGroup
var lock sync.Mutex

func incrementing(ch chan bool) {
	for b := 0; b < 1000000; b++ {
		lock.Lock()
		i += 1
		lock.Unlock()
	}
	ch <- true
}

func decrementing(ch chan bool) {
	for b := 0; b < 1000000; b++ {
		lock.Lock()
		i -= 1
		lock.Unlock()
	}
	ch <- true
}

func main() {
	runtime.GOMAXPROCS(2)

	// Increment the wait group counter
	wg.Add(2)

	// Create a channel to signal completion of each goroutine
	ch := make(chan bool, 2)

	// Spawn both functions as goroutines
	go func() {
		defer wg.Done()
		incrementing(ch)
	}()

	go func() {
		defer wg.Done()
		decrementing(ch)
	}()

	// Use select to wait for any channel to receive a value
	select {
	case <-ch:
		// One channel received a value

		// Use select again to wait for the second channel to receive a value
		select {
		case <-ch:
			// Both channels received values
		}
	}

	// Close the channel to signal that all goroutines have completed
	close(ch)

	fmt.Println("The magic number is:", i)
}
