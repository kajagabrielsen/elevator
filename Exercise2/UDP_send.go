package main

import (
	"fmt"
	"net"
)

func main() {
	// Example usage for sending directly to a single remote machine
	remoteIP := "127.0.0.1" // Replace with the actual remote IP
	remotePort := 20000     // Replace with the actual remote port
	message := "Hello, server!"

	sendDirectly(remoteIP, remotePort, message)

	// Example usage for sending on broadcast
	broadcastIP := "255.255.255.255" // Replace with the appropriate broadcast IP
	broadcastPort := 20000           // Replace with the appropriate port
	message = "Broadcast message!"

	sendBroadcast(broadcastIP, broadcastPort, message)
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
	defer sendSock.Close()

	// Send the message
	_, err = sendSock.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	fmt.Printf("Sent directly to %s:%d: %s\n", remoteIP, remotePort, message)
}

func sendBroadcast(broadcastIP string, broadcastPort int, message string) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", broadcastIP, broadcastPort))
	if err != nil {
		fmt.Println("Error resolving address for broadcast:", err)
		return
	}

	sendSock, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error creating sending socket for broadcast:", err)
		return
	}
	defer sendSock.Close()

	// Enable broadcast option
	err = sendSock.SetBroadcast(true)
	if err != nil {
		fmt.Println("Error setting broadcast option:", err)
		return
	}

	// Send the message
	_, err = sendSock.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending broadcast message:", err)
		return
	}

	fmt.Printf("Sent broadcast to %s:%d: %s\n", broadcastIP, broadcastPort, message)
}
