package hallassign

import (
	"Elevator/networkcom"
	"Elevator/utils"
	"time"
	"fmt"
)

var OneElevRequests = [utils.N_FLOORS][utils.N_BUTTONS]bool{}

var PrevElevatorRequests [utils.N_FLOORS][utils.N_BUTTONS]bool

var Stopped bool

var NoChangeInRequestsTimer int

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

func DetectMotorStop(){
    for{
    for floor := 0; floor < utils.N_FLOORS; floor++ {
        for button := 0; button < utils.N_BUTTONS; button++ {
            if PrevElevatorRequests[floor][button] != utils.ElevatorGlob.Requests[floor][button] { 
                NoChangeInRequestsTimer = 0
            } else if  PrevElevatorRequests[floor][button] == utils.ElevatorGlob.Requests[floor][button]  && PrevElevatorRequests[floor][button] {
                NoChangeInRequestsTimer += 1
            }
        }
    }
    if NoChangeInRequestsTimer > 3 {
        Stopped = true
    }else {
        Stopped = false
    }
    PrevElevatorRequests = utils.ElevatorGlob.Requests
    time.Sleep(1 * time.Second)
    fmt.Printf("----------------------------------------")
    fmt.Println(Stopped)
}
}