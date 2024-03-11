package main

import (
	"Elevator/driver-go-master/elevio"
	"Elevator/hallassign"
	"Elevator/networkcom"
	"Elevator/networkcom/network/bcast"
	"Elevator/networkcom/network/peers"
	"Elevator/utils"
	"fmt"
	"os"
	"time"
)

func main() {

	elevio.Init("localhost:15658", utils.N_FLOORS)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	utils.FsmOnInitBetweenFloors()
    drv_buttons2 := make(chan elevio.ButtonEvent)
    go elevio.PollButtons(drv_buttons2)

	//go network.InitNetwork() //Start network
//////////////////////////////////////////////

    var id string = os.Args[1]

	// var id string
	// flag.StringVar(&id, "id", "", "id of this peer")
	// flag.Parse()

	// ... or alternatively, we can use the localListOfElevators IP address.
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
	helloTx := make(chan network.HelloMsg)
	helloRx := make(chan network.HelloMsg)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16578, helloTx)
	go bcast.Receiver(16578, helloRx)

	// The example message. We just send one of these every second.
	go func() {
		e := utils.Elevator_glob
		helloMsg := network.HelloMsg{e, 0}
		for {
			helloMsg.Iter++
            utils.Elevator_glob.ID = id
            helloMsg.Elevator = utils.Elevator_glob   
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()
    ButtonPressCh := make(chan elevio.ButtonEvent)
    go hallassign.HandleButtonPressUpdate(drv_buttons2 )
	fmt.Println("Started")

    go hallassign.FSM(ButtonPressCh, drv_buttons,  drv_floors, drv_obstr, drv_stop)

    go peers.PeersUpdate(drv_buttons, peerUpdateCh, helloRx)

    select{}
}
