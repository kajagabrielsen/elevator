package send

import (
	"Elevator/elevator/initial"
	"Elevator/hallassign/call_handling"
	"Elevator/network/peers"
	"time"

)

func SendElevator(id string, helloTx chan peers.HelloMsg) {
	initial.ElevatorGlob.ID = id
	OneElevCabCalls, _ := call.GetCabCalls(initial.ElevatorGlob)
	for i := range initial.ElevatorGlob.Requests{
		initial.ElevatorGlob.Requests[i][2] = OneElevCabCalls[i]
	}
	e := initial.ElevatorGlob
	helloMsg := peers.HelloMsg{
		Elevator: e,
		Iter:     0,
	}
	for {
		helloMsg.Elevator = initial.ElevatorGlob
		helloTx <- helloMsg
		helloMsg.Iter++
		time.Sleep(1 * time.Second)
	}
}