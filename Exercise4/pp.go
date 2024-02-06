package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var (
	count      int
	backupLock sync.Mutex
)

func counter() {

	localhost, _ := net.ResolveUDPAddr("udp", "localhost:36243")
	udpsocket, _ := net.ListenUDP("udp", localhost)
	remoteIP := "localhost"
	remotePort := 36243
	i := 0

	for {
		buffer := make([]byte, 1024)
		udpsocket.SetReadDeadline(time.Now().Add(2 * time.Second))

		numBytesReceived, _, err := udpsocket.ReadFromUDP(buffer)
		receivedData := string(buffer[:numBytesReceived])
		if err != nil {
			break
		}

		i, _ = strconv.Atoi(receivedData)

	}

	exec.Command("gnome-terminal", "--", "go", "run", "pp.go").Run()

	for {

		// Print and increment the count
		fmt.Println(i)
		i++

		addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", remoteIP, remotePort))
		sendSock, _ := net.DialUDP("udp", nil, addr)

		// Send the message
		_, _ = sendSock.Write([]byte(string(i)))

		// Simulate some work
		time.Sleep(time.Second)
	}

}

func main() {
	go counter()
}
