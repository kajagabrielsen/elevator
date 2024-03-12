package hallassign

import (
	network "Elevator/networkcom"
	"Elevator/utils"
)

var OneElevRequests = [utils.N_FLOORS][utils.N_BUTTONS]bool{}

func AssignHallRequest() {
	ListOfElevators := network.ListOfElevators
	AssignedHallCalls := CalculateCostFunc(ListOfElevators)
	OneElevCabCalls := GetCabCalls(utils.ElevatorGlob)
	OneElevHallCalls := AssignedHallCalls[utils.ElevatorGlob.ID]

	for floor := 0; floor < utils.N_FLOORS; floor++ {
		OneElevRequests[floor][0] = OneElevHallCalls[floor][0]
		OneElevRequests[floor][1] = OneElevHallCalls[floor][1]
		OneElevRequests[floor][2] = OneElevCabCalls[floor]
	}
}
