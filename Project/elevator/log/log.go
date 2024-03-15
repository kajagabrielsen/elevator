package log

import (
	"Elevator/driver_go_master/elevio"
	"Elevator/elevator/initial"
	"fmt"
)

// ebToString converts ElevatorBehaviour to a string.
func EbToString(eb initial.ElevatorBehaviour) string {
	switch eb {
	case initial.EBIdle:
		return "EB_Idle"
	case initial.EBDoorOpen:
		return "EB_DoorOpen"
	case initial.EBMoving:
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
func ElevatorLog(es initial.Elevator) {
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
	for f := initial.N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < initial.N_BUTTONS; btn++ {
			if (f == initial.N_FLOORS-1 && btn == int(elevio.BTHallUp)) ||
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

func ElevatorInitialized() initial.Elevator {
	return initial.Elevator{
		Floor:     1,
		Dirn:      elevio.MDStop,
		Behaviour: initial.EBIdle,
		ClearRequestVariant: initial.CV_All,
		DoorOpenDuration:   3.0,
		ID: "5",
	}
}