package peers

import (
	//"Elevator/driver-go-master/elevio"
	"Elevator/driver-go-master/elevio"
	"Elevator/hallassign"
	"Elevator/networkcom"
	"Elevator/utils"
	"fmt"
)

func PeersUpdate(drv_buttons chan elevio.ButtonEvent, peerUpdateCh chan PeerUpdate, helloRx chan network.HelloMsg){
	fmt.Printf("peers")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			network.AliveElevatorsID = p.Peers

		case elev := <-helloRx:
			flag := 0
			for i, element := range network.ListOfElevators{
				if element.ID == elev.Elevator.ID{
					network.ListOfElevators[i] = elev.Elevator
					//element = elev.Elevator
					flag = 1
				}
			}
			if flag == 0 {
				network.ListOfElevators = append(network.ListOfElevators, elev.Elevator)

			}
			fmt.Println(network.ListOfElevators)
			fmt.Printf("Received: %#v\n", elev)
			hallassign.AssignHallRequest()
			for floor_num, floor := range utils.Elevator_glob.Requests{
				for btn_num, _ := range floor {
					if utils.Elevator_glob.Requests[floor_num][btn_num]{
						utils.FsmOnRequestButtonPress(floor_num, utils.Button(btn_num))
					}
				}
			}
		case btn := <-drv_buttons:
			utils.Elevator_glob.Requests[btn.Floor][btn.Button] = true
		}
	}

}