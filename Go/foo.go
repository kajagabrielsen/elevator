// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
)

var i = 0

func server(inc chan bool, dec chan bool, done chan bool, quit chan bool) {
	for {
		select {
		case <-inc:
			i++
		case <-dec:
			i--
		case <-quit:
			Println("The magic number is:", i)
			done <- true
		}
	}
}

func incrementing(inc chan bool, done chan bool) {
	for b := 0; b < 1000000; b++ {
		inc <- true
	}
	done <- true
}

func decrementing(dec chan bool, done chan bool) {
	for b := 0; b < 1000000; b++ {
		dec <- true
	}
	done <- true
}

func main() {
	runtime.GOMAXPROCS(2)
	inc := make(chan bool)
	dec := make(chan bool)
	done := make(chan bool)
	quit := make(chan bool)
	go server(inc, dec, done, quit)
	go incrementing(inc, done)
	go decrementing(dec, done)

	<-done
	<-done
	quit <- true
	<-done

}
