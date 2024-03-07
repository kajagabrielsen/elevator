package network

import (
	"Elevator/networkcom/network/bcast"
	//"Elevator/networkcom/network/localip"
	"Elevator/networkcom/network/peers"
	"Elevator/utils"
	"flag"
	"fmt"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
type HelloMsg struct {
	Elevator utils.Elevator
	Iter    int
}

var ListOfElevators [utils.N_ELEVATORS]utils.Elevator

func GetListOfElevators () [utils.N_ELEVATORS]utils.Elevator {
	return ListOfElevators
}

func InitNetwork() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	/*if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}*/

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15611, id, peerTxEnable)
	go peers.Receiver(15611, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16578, helloTx)
	go bcast.Receiver(16578, helloRx)

	// The example message. We just send one of these every second.
	go func() {
		helloMsg := HelloMsg{utils.GetElevator(), 0}
		for {
			helloMsg.Iter++
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Started")
	//ListOfElevators := [3]utils.Elevator{}

	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case elev := <-helloRx:
			ListOfElevators[elev.Elevator.ID] = elev.Elevator
			fmt.Printf("Received: %#v\n", elev)
		}
	}

}
