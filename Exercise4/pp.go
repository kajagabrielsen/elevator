package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	count      int
	backupLock sync.Mutex
)

func primary() {
	for {
		// Print and increment the count
		fmt.Println(count)
		count++

		// Simulate some work
		time.Sleep(time.Second)
	}
}

func backup() {
	// Monitor the primary process
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	for {
		<-ch
		// Primary process has terminated, become the new primary
		backupLock.Lock()
		go primary()
		backupLock.Unlock()

		// Create a new backup
		go backup()
	}
}

func main() {
	// Start the primary process
	go primary()

	// Start the backup process
	go backup()

	// Wait for termination signal
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGTERM, syscall.SIGINT)
	<-terminate
}
