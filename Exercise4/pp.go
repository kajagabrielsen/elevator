package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

const backupFile = "backup.txt"

func programA() {
	fmt.Println("Program A is running...")

	// Simulate some work
	time.Sleep(5 * time.Second)

	fmt.Println("Program A terminates.")
}

func programB() {
	fmt.Println("Program B is starting Program A...")

	// Start Program A
	cmd := exec.Command("go", "run", "programA.go")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting Program A:", err)
	}

	// Communication with Program A
	data := "Hello from Program B!"
	err = ioutil.WriteFile(backupFile, []byte(data), 0644)
	if err != nil {
		fmt.Println("Error writing to backup file:", err)
	}
}

func main() {
	for {
		// Check if the primary process is running
		_, err := os.Stat(backupFile)
		isPrimary := err == nil

		if isPrimary {
			// Primary loop for doing work
			fmt.Println("Program A (Primary) is running...")

			// Simulate some work
			time.Sleep(3 * time.Second)

		} else {
			// Create a new backup
			fmt.Println("Primary not found. Creating a new backup (Program B)...")
			go programB()

			// Wait for the backup to be created
			time.Sleep(1 * time.Second)
		}
	}
}
