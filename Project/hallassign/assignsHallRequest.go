package hallassign

import (
	network "Elevator/networkcom"
	"Elevator/utils"
)

func GetIndex(key string, list []string) int {
	for i, value := range list {
		if value == key {
			return i
		}
	}

	return 0
}

func AssignHallRequest(e utils.Elevator) {
	ListOfElevators := network.ListOfElevators
	AssignedHallCalls := CalculateCostFunc(ListOfElevators)
	OneElevCabCalls := GetCabCalls(ListOfElevators[GetIndex(e.ID, network.GetAliveElevatorsID())])
	OneElevHallCalls := AssignedHallCalls[e.ID]

	OneElevRequests := [utils.N_FLOORS][utils.N_BUTTONS]bool{}

	for floor := 0; floor < utils.N_FLOORS; floor++ {
		OneElevRequests[floor][0] = OneElevHallCalls[floor][0]
		OneElevRequests[floor][1] = OneElevHallCalls[floor][1]
		OneElevRequests[floor][2] = OneElevCabCalls[floor]
	}
	utils.Elevator_glob.Requests = OneElevRequests
}