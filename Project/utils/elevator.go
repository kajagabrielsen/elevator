package utils

import (
	"fmt"
)

// ebToString converts ElevatorBehaviour to a string.
func EbToString(eb ElevatorBehaviour) string {
	switch eb {
	case EB_Idle:
		return "EB_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	default:
		return "EB_UNDEFINED"
	}
}

// elevatorPrint prints the state of the elevator.
func ElevatorPrint(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n",
		es.Floor,
		DirnToString(es.Dirn), // Replace with the Go equivalent of elevio_dirn_toString(es.Dirn)
		EbToString(es.Behaviour),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < N_BUTTONS; btn++ {
			if (f == N_FLOORS-1 && btn == int(B_HallUp)) ||
				(f == 0 && btn == int(B_HallDown)) {
				fmt.Print("|     ")
			} else {
				fmt.Print("|", es.Requests[f][btn])  // Replace with the Go equivalent for requests[f][btn]
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

// elevatorUninitialized initializes and returns an uninitialized elevator.
func ElevatorUninitialized() Elevator {
	return Elevator{
		Floor:     -1,
		Dirn:      D_Stop,
		Behaviour: EB_Idle,
		ClearRequestVariant: CV_All,
		DoorOpenDuration_s:   3.0,
	}
}