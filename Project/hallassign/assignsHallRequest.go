package hallassign

import (
	"Elevator/networkcom"
	"Elevator/utils"
	"strconv"
)

func AssignHallRequest(e utils.Elevator) [utils.N_FLOORS][utils.N_BUTTONS] bool{
	ListOfElevators := network.GetListOfElevators()
	AssignedHallCalls := *CalculateCostFunc(ListOfElevators)
	OneElevCabCalls := GetCabCalls(ListOfElevators[e.ID])
	OneElevHallCalls := AssignedHallCalls[strconv.Itoa(e.ID)]

	OneElevRequests := [utils.N_FLOORS][utils.N_BUTTONS] bool {}

	for floor := 0; floor < utils.N_FLOORS ; floor++{
		OneElevRequests[floor][0] = OneElevHallCalls[floor][0]
		OneElevRequests[floor][1] = OneElevHallCalls[floor][1]
		OneElevRequests[floor][2] = OneElevCabCalls[floor]
	}
	return OneElevRequests
}