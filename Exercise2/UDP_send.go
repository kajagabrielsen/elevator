package main

import (
	"fmt"
	"net"
)

func main() {

	// Define the address to listen on
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:20010")
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

	// Example usage for sending directly to a single remote machine
	remoteIP := "10.100.23.129" // Replace with the actual remote IP
	remotePort := 20010         // Replace with the actual remote port
	message := "Vi er gruppe 10, håper ale er bra med deg kjære server!"

	sendDirectly(remoteIP, remotePort, message)

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

func sendDirectly(remoteIP string, remotePort int, message string) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", remoteIP, remotePort))
	if err != nil {
		fmt.Println("Error resolving address for direct send:", err)
		return
	}

	sendSock, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error creating sending socket for direct send:", err)
		return
	}

	// Send the message
	_, err = sendSock.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	fmt.Printf("Sent directly to %s:%d: %s\n", remoteIP, remotePort, message)
}
