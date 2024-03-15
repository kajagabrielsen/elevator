package main

import (
	"Elevator/driver_go_master/elevio"
	fsm "Elevator/elevator/fsm_func"
	"Elevator/elevator/initial"
	implfsm "Elevator/hallassign/fsm_implementation"
	motorstop "Elevator/hallassign/motor_stop"
	"Elevator/network/bcast"
	"Elevator/network/peers"
	send "Elevator/network/sendig_elevator"
	"fmt"
	"os"
)

func main() {

	var id string = os.Args[1]

	elevio.InitDriver("localhost:15657", initial.NFloors)

	fsm.FsmOnInitBetweenFloors()

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	drv_buttons2 := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(drv_buttons2)

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15611, id, peerTxEnable)
	go peers.Receiver(15611, peerUpdateCh)

	helloTx := make(chan peers.HelloMsg)
	helloRx := make(chan peers.HelloMsg)
	go bcast.Transmitter(16578, helloTx)
	go bcast.Receiver(16578, helloRx)

	go send.SendElevator(id,helloTx)
	fmt.Println("Started")

	go motorstop.DetectMotorStop()

	go implfsm.FSM(drv_buttons, drv_floors, drv_obstr, drv_stop)

	go peers.PeersUpdate(drv_buttons2, peerUpdateCh, helloRx)

	select {}
}
