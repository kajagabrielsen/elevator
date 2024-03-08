package peers

import (
	"Elevator/networkcom"
	"fmt"

)

func PeersUpdate(peerUpdateCh chan PeerUpdate, helloRx chan network.HelloMsg){
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
			for _, element := range network.ListOfElevators{
				if element.ID == elev.Elevator.ID{
					element = elev.Elevator
					flag = 1
				}
			}
			if flag == 0 {
				network.ListOfElevators = append(network.ListOfElevators, elev.Elevator)

			}
			fmt.Printf("Received: %#v\n", elev)
		}
	}

}