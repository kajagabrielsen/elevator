package hallassign

import (
	"Elevator/elevator/initialize"
	"Elevator/network/list"
	"time"
)

var OneElevRequests = [initialize.N_FLOORS][initialize.N_BUTTONS]bool{}

var PrevElevatorRequests [initialize.N_FLOORS][initialize.N_BUTTONS]bool

var Stopped bool

var NoChangeInRequestsTimer int

func AssignHallRequest() {
	AssignedHallCalls := CalculateCostFunc(list.ListOfElevators)
	OneElevCabCalls,_ := GetCabCalls(initialize.ElevatorGlob)
	OneElevHallCalls := AssignedHallCalls[initialize.ElevatorGlob.ID]

	for floor := 0; floor < initialize.N_FLOORS; floor++ {
		OneElevRequests[floor][0] = OneElevHallCalls[floor][0]
		OneElevRequests[floor][1] = OneElevHallCalls[floor][1]
		OneElevRequests[floor][2] = OneElevCabCalls[floor]
	}
	
}

func DetectMotorStop(){
    for{
    for floor := 0; floor < initialize.N_FLOORS; floor++ {
        for button := 0; button < initialize.N_BUTTONS; button++ {
            if PrevElevatorRequests[floor][button] != initialize.ElevatorGlob.Requests[floor][button] { 
                NoChangeInRequestsTimer = 0
            } else if  PrevElevatorRequests[floor][button] == initialize.ElevatorGlob.Requests[floor][button]  && PrevElevatorRequests[floor][button] {
                NoChangeInRequestsTimer += 1
            }
        }
    }
    if NoChangeInRequestsTimer > 3 {
        Stopped = true
    }else {
        Stopped = false
    }
    PrevElevatorRequests = initialize.ElevatorGlob.Requests
    time.Sleep(1 * time.Second)
}
}