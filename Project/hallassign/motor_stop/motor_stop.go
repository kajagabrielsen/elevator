package motorstop

import (
	"Elevator/elevator/initial"
	"time"
)
var PrevElevatorRequests [initial.N_FLOORS][initial.N_BUTTONS]bool

var Stopped bool

var NoChangeInRequestsTimer int

func DetectMotorStop(){
    for{
    for floor := 0; floor < initial.N_FLOORS; floor++ {
        for button := 0; button < initial.N_BUTTONS; button++ {
            if PrevElevatorRequests[floor][button] != initial.ElevatorGlob.Requests[floor][button] { 
                NoChangeInRequestsTimer = 0
            } else if  PrevElevatorRequests[floor][button] == initial.ElevatorGlob.Requests[floor][button]  && PrevElevatorRequests[floor][button] {
                NoChangeInRequestsTimer += 1
            }
        }
    }
    if NoChangeInRequestsTimer > 3 {
        Stopped = true
    }else {
        Stopped = false
    }
    PrevElevatorRequests = initial.ElevatorGlob.Requests
    time.Sleep(1 * time.Second)
}
}