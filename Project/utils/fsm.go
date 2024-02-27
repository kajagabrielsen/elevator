package utils

import (
	"fmt"
	"Elevator/driver-go-master/elevio"
)


var (
	elevator     Elevator
	outputDevice ElevOutputDevice
)

func init() {
	elevator = ElevatorUninitialized()

	// Load configuration from your config.go file or define them here directly
	// elevator.Config = ElevatorConfig{
	// 	DoorOpenDurationS: doorOpenDuration,
	// 	ClearRequestVariant: clearRequestVariant,
	// 	InputPollRate: inputPollRate,
	// }
	outputDevice = GetOutputDevice()
}

func SetAllLights(es Elevator) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			var BTN elevio.ButtonType = elevio.ButtonType(btn)
			outputDevice.RequestButtonLight(floor, BTN, es.Requests[floor][btn])
		}
	}
}

func FsmOnInitBetweenFloors() {
	outputDevice.MotorDirection(elevio.MD_Down)
	elevator.Dirn = D_Down
	elevator.Behaviour = EB_Moving
}

func FsmOnRequestButtonPress(btnFloor int, btnType Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsmOnRequestButtonPress", btnFloor, ButtonToString(btnType))
	ElevatorPrint(elevator)

	switch elevator.Behaviour {
	case EB_DoorOpen:
		if RequestsShouldClearImmediately(elevator, btnFloor, btnType) {
			Timer_start(elevator.DoorOpenDuration_s)
		} else {
			elevator.Requests[btnFloor][btnType] = true
		}
		break

	case EB_Moving:
		elevator.Requests[btnFloor][btnType] = true
		fmt.Printf("move")
		break

	case EB_Idle:
		fmt.Printf("idle")
		elevator.Requests[btnFloor][btnType] = true
		pair := Requests_chooseDirection(elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case EB_DoorOpen:
			outputDevice.DoorLight(true)
			Timer_start(elevator.DoorOpenDuration_s)
			elevator = RequestsClearAtCurrentFloor(elevator)
			break

		case EB_Moving:
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(elevator.Dirn)
			outputDevice.MotorDirection(mot_dir)
			break

		case EB_Idle:
			break
		}
		break
	}

	SetAllLights(elevator)

	fmt.Println("\nNew state:")
	ElevatorPrint(elevator)
}

func FsmOnFloorArrival(newFloor int) {
	fmt.Printf("\n\n%s(%d)\n", "fsmOnFloorArrival", newFloor)
	ElevatorPrint(elevator)

	elevator.Floor = newFloor

	outputDevice.FloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case EB_Moving:
		if RequestsShouldStop(elevator) {
			outputDevice.MotorDirection(D_Stop)
			outputDevice.DoorLight(true)
			elevator = RequestsClearAtCurrentFloor(elevator)
			Timer_start(elevator.DoorOpenDuration_s)
			SetAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen
		}
		break
	default:
		break
	}

	fmt.Println("\nNew state:")
	ElevatorPrint(elevator)
}

func FsmOnDoorTimeout() {
	fmt.Printf("\n\n%s()\n", "fsmOnDoorTimeout")
	ElevatorPrint(elevator)

	switch elevator.Behaviour {
	case EB_DoorOpen:
		pair := Requests_chooseDirection(elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour

		switch elevator.Behaviour {
		case EB_DoorOpen:
			Timer_start(elevator.DoorOpenDuration_s)
			elevator = RequestsClearAtCurrentFloor(elevator)
			SetAllLights(elevator)
			break
		case EB_Moving:
		case EB_Idle:
			outputDevice.DoorLight(false)
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(elevator.Dirn)
			outputDevice.MotorDirection(mot_dir)
			break
		}

		break
	default:
		break
	}

	fmt.Println("\nNew state:")
	ElevatorPrint(elevator)
}

// You need to implement the missing functions such as elevatorUninitialized,
// elevioGetOutputDevice, timerStart, elevatorPrint, requestsShouldClearImmediately,
// requestsChooseDirection, requestsClearAtCurrentFloor, requestsShouldStop,
// and other related functions based on your specific implementation.
