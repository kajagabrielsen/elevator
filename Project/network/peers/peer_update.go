package peers

import (
	"Elevator/driver_go_master/elevio"
	"Elevator/elevator/initialize"
	"Elevator/hallassign/call_handling"
	"Elevator/hallassign/assign_hall_request"
	"Elevator/network/list"
	"fmt"
)

var DeadElevatorsID []string

var AliveElevatorsID []string

type HelloMsg struct {
	Elevator initialize.Elevator
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

			//fjerner lost peers fra ListOfElevators
			var result []initialize.Elevator

			for _, elevator := range list.ListOfElevators {
				found := false

				for _, deadID := range DeadElevatorsID {
					if elevator.ID == deadID {
						found = true
					}
				}

				if !found {
					result = append(result, elevator)
				}
			}
			list.ListOfElevators = result
			assign.AssignHallRequest()

		case elev := <-helloRx:

			if elev.Elevator.Obstructed {
				list.RemoveFromListOfElevators(list.ListOfElevators, elev.Elevator)
			}


			flag := 0
			for i, element := range list.ListOfElevators {
				if element.ID == elev.Elevator.ID {
					list.ListOfElevators[i] = elev.Elevator
					flag = 1
				}
			}
			if flag == 0 && !elev.Elevator.Obstructed{
				list.ListOfElevators = append(list.ListOfElevators, elev.Elevator)				
			}

			fmt.Printf("Received: %#v\n", elev)
		case btn := <-drv_buttons:
			initialize.ElevatorGlob.Requests[btn.Floor][btn.Button] = true
		}
	}

}
