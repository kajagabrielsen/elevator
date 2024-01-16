package main

import (
	"fmt"
	"net"
)

func main() {
	localPort := 20010 // Replace with the port you want to listen on
	clientIP := getLocalIP()

	// Start the server in a separate goroutine
	go startServer(localPort)

	// Connect to the TCP server
	serverPort := 34933 // Change this to the port your server is listening on
	serverIP := "10.100.23.129"
	connectBackMessage := fmt.Sprintf("Connect to: %s:%d\x00", clientIP, localPort)
	connectToServer(serverIP, serverPort, connectBackMessage)
}

func startServer(localPort int) {
	// Bind to a local address and start listening
	addr := fmt.Sprintf(":%d", localPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error binding to address:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server listening on port %d\n", localPort)

	// Accept incoming connections and handle them
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Send a welcome message to the client
	welcomeMessage := "Welcome to the server!"
	conn.Write([]byte(welcomeMessage))

	// Receive and echo messages
	for {
		buffer := make([]byte, 1024)
		numBytesReceived, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}

		receivedMessage := string(buffer[:numBytesReceived])
		fmt.Printf("Received message: %s\n", receivedMessage)

		// Echo the message back to the client
		_, err = conn.Write([]byte(receivedMessage))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}

func connectToServer(serverIP string, serverPort int, message string) {
	// Connect to the TCP server
	addr := fmt.Sprintf("%s:%d", serverIP, serverPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server!")

	// Send and receive messages
	sendMessage(conn, message)
	receiveMessage(conn)
}

func sendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	fmt.Printf("Sent message: %s\n", message)
}

func receiveMessage(conn net.Conn) {
	buffer := make([]byte, 1024)
	numBytesReceived, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving message:", err)
		return
	}

	receivedMessage := string(buffer[:numBytesReceived])
	fmt.Printf("Received message: %s\n", receivedMessage)
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80") // Use a valid external IP as a dummy address
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
