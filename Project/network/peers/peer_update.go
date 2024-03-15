package peers

import (
	"Elevator/driver_go_master/elevio"
	"Elevator/elevator/initial"
	"Elevator/hallassign/assign_hall_request"
	"Elevator/hallassign/call_handling"
	"Elevator/network/list"
	"fmt"
)

var DeadElevatorsID []string

var AliveElevatorsID []string

type HelloMsg struct {
	Elevator initial.Elevator
	Iter     int
}

func PeersUpdate(drv_buttons chan elevio.ButtonEvent, peerUpdateCh chan PeerUpdate, helloRx chan HelloMsg) {
	fmt.Printf("peers")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			AliveElevatorsID = p.Peers
			DeadElevatorsID = p.Lost
			call.UpdateGlobalHallCalls(list.ListOfElevators)

			for _, deadID := range DeadElevatorsID {
				list.RemoveFromListOfElevators(list.ListOfElevators,deadID)
			}
			assign.AssignHallRequest()

		case elev := <-helloRx:

			if elev.Elevator.Obstructed {
				list.RemoveFromListOfElevators(list.ListOfElevators, elev.Elevator.ID)
			}

			list.AddToListOfElevators(list.ListOfElevators,elev.Elevator)

		case btn := <-drv_buttons:
			initial.ElevatorGlob.Requests[btn.Floor][btn.Button] = true
		}
	}

}
