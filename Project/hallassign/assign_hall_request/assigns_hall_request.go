package assign

import (
	"Elevator/elevator/initialize"
	"Elevator/hallassign/call_handling"
	"Elevator/hallassign/cost"
	"Elevator/network/list"
)

var OneElevRequests = [initialize.N_FLOORS][initialize.N_BUTTONS]bool{}

func AssignHallRequest() {
	AssignedHallCalls :=cost.CalculateCostFunc(list.ListOfElevators)
	OneElevCabCalls,_ := call.GetCabCalls(initialize.ElevatorGlob)
	OneElevHallCalls := AssignedHallCalls[initialize.ElevatorGlob.ID]

	for floor := 0; floor < initialize.N_FLOORS; floor++ {
		OneElevRequests[floor][0] = OneElevHallCalls[floor][0]
		OneElevRequests[floor][1] = OneElevHallCalls[floor][1]
		OneElevRequests[floor][2] = OneElevCabCalls[floor]
	}
	
}
