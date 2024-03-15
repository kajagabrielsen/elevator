package motorstop

import(
	"Elevator/elevator/initialize"
	"time"
)
var PrevElevatorRequests [initialize.N_FLOORS][initialize.N_BUTTONS]bool

var Stopped bool

var NoChangeInRequestsTimer int

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