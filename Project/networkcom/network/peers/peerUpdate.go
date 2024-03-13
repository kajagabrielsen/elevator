package peers

import (
	"Elevator/driver-go-master/elevio"
	//"Elevator/hallassign"
	"Elevator/networkcom"
	"Elevator/utils"
	"fmt"
)

var DeadElevatorsID []string

func PeersUpdate(drv_buttons chan elevio.ButtonEvent, peerUpdateCh chan PeerUpdate, helloRx chan network.HelloMsg) {
	fmt.Printf("peers")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			network.AliveElevatorsID = p.Peers
			DeadElevatorsID = p.Lost
			//hallassign.UpdateGlobalHallCalls()

			//fjerner lost peers fra ListOfElevators
			var result []utils.Elevator

			for _, elevator := range network.ListOfElevators {
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
			network.ListOfElevators = result
			//hallassign.AssignHallRequest()

		case elev := <-helloRx:
			flag := 0
			for i, element := range network.ListOfElevators {
				if element.ID == elev.Elevator.ID {
					network.ListOfElevators[i] = elev.Elevator
					flag = 1
				}
			}
			if flag == 0 {
				network.ListOfElevators = append(network.ListOfElevators, elev.Elevator)

			}

			fmt.Printf("Received: %#v\n", elev)
		case btn := <-drv_buttons:
			utils.ElevatorGlob.Requests[btn.Floor][btn.Button] = true
		}
	}

}
