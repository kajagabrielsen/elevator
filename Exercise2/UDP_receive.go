package main

import (
	"fmt"
	"net"
)

func main() {
	// Define the address to listen on
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:30000")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// Create a UDP socket
	recvSock, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error creating socket:", err)
		return
	}
	defer recvSock.Close()

	// Buffer to store received data
	buffer := make([]byte, 1024)

	// Loop to continuously receive data
	for {
		// Receive data on the socket
		numBytesReceived, fromWho, err := recvSock.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			continue
		}

		// Convert the received bytes to a string
		receivedData := string(buffer[:numBytesReceived])

		// Optional: Filter out messages from ourselves (assuming localIP is a string)
		localIP := "127.0.0.1"
		if fromWho.IP.String() != localIP {
			// Do stuff with the received data
			fmt.Printf("Received %d bytes from %v: %s\n", numBytesReceived, fromWho, receivedData)
		}
	}
}
