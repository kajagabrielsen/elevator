package assign

import (
	"Elevator/elevator/initial"
	call "Elevator/hallassign/call_handling"
	"Elevator/hallassign/cost"
	"Elevator/network/list"
)

var OneElevRequests = [initial.NFloors][initial.NButtons]bool{}

func AssignHallRequest() {
	AssignedHallCalls :=cost.CalculateCostFunc(list.ListOfElevators)
	OneElevCabCalls,_ := call.GetCabCalls(initial.ElevatorGlob)
	OneElevHallCalls := AssignedHallCalls[initial.ElevatorGlob.ID]

	for floor := 0; floor < initial.NFloors; floor++ {
		OneElevRequests[floor][0] = OneElevHallCalls[floor][0]
		OneElevRequests[floor][1] = OneElevHallCalls[floor][1]
		OneElevRequests[floor][2] = OneElevCabCalls[floor]
	}
	
}
