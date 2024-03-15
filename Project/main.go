package main

import (
	"Elevator/driver_go_master/elevio"
	"Elevator/elevator/initialize"
	"Elevator/elevator/fsm_func"
	"Elevator/hallassign/fsm_implementation"
	"Elevator/hallassign/motor_stop"
	"Elevator/hallassign/call_handling"
	"Elevator/network/bcast"
	"Elevator/network/peers"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	var id string = os.Args[1]
	id_int, _ := strconv.Atoi(id)
	port := 15656 + id_int

	elevio.InitDriver("localhost:"+strconv.Itoa((port)), initialize.N_FLOORS)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)


	fsm.FsmOnInitBetweenFloors()
	drv_buttons2 := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(drv_buttons2)

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15611, id, peerTxEnable)
	go peers.Receiver(15611, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	helloTx := make(chan peers.HelloMsg)
	helloRx := make(chan peers.HelloMsg)

	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16578, helloTx)
	go bcast.Receiver(16578, helloRx)

	// The example message. We just send one of these every second.
	go func() {
		initialize.ElevatorGlob.ID = id
		OneElevCabCalls, _ := call.GetCabCalls(initialize.ElevatorGlob)
		for i := range initialize.ElevatorGlob.Requests{
			initialize.ElevatorGlob.Requests[i][2] = OneElevCabCalls[i]
		}
		e := initialize.ElevatorGlob
		helloMsg := peers.HelloMsg{
			Elevator: e,
			Iter:     0,
		}
		for {
			helloMsg.Elevator = initialize.ElevatorGlob
			helloTx <- helloMsg
			helloMsg.Iter++
			time.Sleep(1 * time.Second)
		}
	}()
	fmt.Println("Started")

	go motorstop.DetectMotorStop()

	go implfsm.FSM(drv_buttons, drv_floors, drv_obstr, drv_stop)

	go peers.PeersUpdate(drv_buttons2, peerUpdateCh, helloRx)

	select {}
}
