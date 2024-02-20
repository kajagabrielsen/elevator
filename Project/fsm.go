package main

import (
	"fmt"
	"Elevator/driver-go-master/elevio"
)


var (
	elevator     Elevator
	outputDevice ElevOutputDevice
)

func init() {
	elevator = elevatorUninitialized()

	// Load configuration from your config.go file or define them here directly
	// elevator.Config = ElevatorConfig{
	// 	DoorOpenDurationS: doorOpenDuration,
	// 	ClearRequestVariant: clearRequestVariant,
	// 	InputPollRate: inputPollRate,
	// }
	outputDevice = GetOutputDevice()
}

func setAllLights(es Elevator) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			var BTN elevio.ButtonType = elevio.ButtonType(btn)
			outputDevice.RequestButtonLight(floor, BTN, es.Requests[floor][btn])
		}
	}
}

func fsmOnInitBetweenFloors() {
	outputDevice.MotorDirection(elevio.MD_Down)
	elevator.Dirn = D_Down
	elevator.Behaviour = EB_Moving
}

func fsmOnRequestButtonPress(btnFloor int, btnType Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsmOnRequestButtonPress", btnFloor, ButtonToString(btnType))
	elevatorPrint(elevator)

	switch elevator.Behaviour {
	case EB_DoorOpen:
		if requestsShouldClearImmediately(elevator, btnFloor, btnType) {
			timer_start(elevator.DoorOpenDuration_s)
		} else {
			elevator.Requests[btnFloor][btnType] = true
		}
		break

	case EB_Moving:
		elevator.Requests[btnFloor][btnType] = true
		break

	case EB_Idle:
		elevator.Requests[btnFloor][btnType] = true
		pair := requests_chooseDirection(elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case EB_DoorOpen:
			outputDevice.DoorLight(true)
			timer_start(elevator.DoorOpenDuration_s)
			elevator = requestsClearAtCurrentFloor(elevator)
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

	setAllLights(elevator)

	fmt.Println("\nNew state:")
	elevatorPrint(elevator)
}

func fsmOnFloorArrival(newFloor int) {
	fmt.Printf("\n\n%s(%d)\n", "fsmOnFloorArrival", newFloor)
	elevatorPrint(elevator)

	elevator.Floor = newFloor

	outputDevice.FloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case EB_Moving:
		if requestsShouldStop(elevator) {
			outputDevice.MotorDirection(D_Stop)
			outputDevice.DoorLight(true)
			elevator = requestsClearAtCurrentFloor(elevator)
			timer_start(elevator.DoorOpenDuration_s)
			setAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen
		}
		break
	default:
		break
	}

	fmt.Println("\nNew state:")
	elevatorPrint(elevator)
}

func fsmOnDoorTimeout() {
	fmt.Printf("\n\n%s()\n", "fsmOnDoorTimeout")
	elevatorPrint(elevator)

	switch elevator.Behaviour {
	case EB_DoorOpen:
		pair := requests_chooseDirection(elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour

		switch elevator.Behaviour {
		case EB_DoorOpen:
			timer_start(elevator.DoorOpenDuration_s)
			elevator = requestsClearAtCurrentFloor(elevator)
			setAllLights(elevator)
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
	elevatorPrint(elevator)
}

// You need to implement the missing functions such as elevatorUninitialized,
// elevioGetOutputDevice, timerStart, elevatorPrint, requestsShouldClearImmediately,
// requestsChooseDirection, requestsClearAtCurrentFloor, requestsShouldStop,
// and other related functions based on your specific implementation.
