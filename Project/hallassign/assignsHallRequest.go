package hallassign

import (
	"Elevator/networkcom"
	"Elevator/utils"
)

var OneElevRequests = [utils.N_FLOORS][utils.N_BUTTONS]bool{}

func AssignHallRequest() {
	AssignedHallCalls := CalculateCostFunc(network.ListOfElevators)
	OneElevCabCalls,_ := GetCabCalls(utils.ElevatorGlob)
	OneElevHallCalls := AssignedHallCalls[utils.ElevatorGlob.ID]

	for floor := 0; floor < utils.N_FLOORS; floor++ {
		OneElevRequests[floor][0] = OneElevHallCalls[floor][0]
		OneElevRequests[floor][1] = OneElevHallCalls[floor][1]
		OneElevRequests[floor][2] = OneElevCabCalls[floor]
	}
}
