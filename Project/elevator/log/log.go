package log

import (
	"Elevator/driver_go_master/elevio"
	"Elevator/elevator/initialize"
	"fmt"
)

// ebToString converts ElevatorBehaviour to a string.
func EbToString(eb initialize.ElevatorBehaviour) string {
	switch eb {
	case initialize.EB_Idle:
		return "EB_Idle"
	case initialize.EB_DoorOpen:
		return "EB_DoorOpen"
	case initialize.EB_Moving:
		return "EB_Moving"
	default:
		return "EB_UNDEFINED"
	}
}

// DirnToString converts Direction to a string.
func DirnToString(d elevio.MotorDirection) string {
	switch d {
	case elevio.MDUp:
		return "D_Up"
	case elevio.MDDown:
		return "D_Down"
	case elevio.MDStop:
		return "D_Stop"
	default:
		return "D_UNDEFINED"
	}
}

// elevatorPrint prints the state of the elevator.
func ElevatorLog(es initialize.Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n",
		es.Floor,
		DirnToString(es.Dirn), 
		EbToString(es.Behaviour),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := initialize.N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < initialize.N_BUTTONS; btn++ {
			if (f == initialize.N_FLOORS-1 && btn == int(elevio.BTHallUp)) ||
				(f == 0 && btn == int(elevio.BTHallDown)) {
				fmt.Print("|     ")
			} else {
				fmt.Print("|", es.Requests[f][btn])
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

func ElevatorInitialized() initialize.Elevator {
	return initialize.Elevator{
		Floor:     1,
		Dirn:      elevio.MDStop,
		Behaviour: initialize.EB_Idle,
		ClearRequestVariant: initialize.CV_All,
		DoorOpenDuration:   3.0,
		ID: "5",
	}
}