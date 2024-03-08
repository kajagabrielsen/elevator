package utils

import (
	"Elevator/driver-go-master/elevio"
	"fmt"
)

var (
	Elevator_glob Elevator
	outputDevice  ElevOutputDevice
)

func init() {
	Elevator_glob = ElevatorInitialized()
	fmt.Printf("init")

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
	Elevator_glob.Dirn = D_Down
	Elevator_glob.Behaviour = EB_Moving
}

func FsmOnRequestButtonPress(btnFloor int, btnType Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsmOnRequestButtonPress", btnFloor, ButtonToString(btnType))
	ElevatorPrint(Elevator_glob)

	switch Elevator_glob.Behaviour {
	case EB_DoorOpen:
		if RequestsShouldClearImmediately(Elevator_glob, btnFloor, btnType) {
			TimerStart(Elevator_glob.DoorOpenDuration_s)
		} else {
			Elevator_glob.Requests[btnFloor][btnType] = true
		}

	case EB_Moving:
		Elevator_glob.Requests[btnFloor][btnType] = true

	case EB_Idle:
		Elevator_glob.Requests[btnFloor][btnType] = true
		pair := RequestsChooseDirection(Elevator_glob)
		Elevator_glob.Dirn = pair.Dirn
		Elevator_glob.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case EB_DoorOpen:
			outputDevice.DoorLight(true)
			TimerStart(Elevator_glob.DoorOpenDuration_s)
			Elevator_glob = RequestsClearAtCurrentFloor(Elevator_glob)

		case EB_Moving:
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(Elevator_glob.Dirn)
			outputDevice.MotorDirection(mot_dir)

		case EB_Idle:
		}
	}

	SetAllLights(Elevator_glob)

	fmt.Println("\nNew state:")
	ElevatorPrint(Elevator_glob)
}

func FsmOnFloorArrival(newFloor int) {
	fmt.Printf("\n\n%s(%d)\n", "fsmOnFloorArrival", newFloor)
	ElevatorPrint(Elevator_glob)

	Elevator_glob.Floor = newFloor

	outputDevice.FloorIndicator(Elevator_glob.Floor)

	switch Elevator_glob.Behaviour {
	case EB_Moving:
		if RequestsShouldStop(Elevator_glob) {
			outputDevice.MotorDirection(D_Stop)
			outputDevice.DoorLight(true)
			Elevator_glob = RequestsClearAtCurrentFloor(Elevator_glob)
			TimerStart(Elevator_glob.DoorOpenDuration_s)
			SetAllLights(Elevator_glob)
			Elevator_glob.Behaviour = EB_DoorOpen
		}
	default:
		break
	}

	fmt.Println("\nNew state:")
	ElevatorPrint(Elevator_glob)
}

func FsmOnDoorTimeout() {
	fmt.Printf("\n\n%s()\n", "fsmOnDoorTimeout")
	ElevatorPrint(Elevator_glob)

	switch Elevator_glob.Behaviour {
	case EB_DoorOpen:
		pair := RequestsChooseDirection(Elevator_glob)
		Elevator_glob.Dirn = pair.Dirn
		Elevator_glob.Behaviour = pair.Behaviour

		switch Elevator_glob.Behaviour {
		case EB_DoorOpen:
			TimerStart(Elevator_glob.DoorOpenDuration_s)
			Elevator_glob = RequestsClearAtCurrentFloor(Elevator_glob)
			SetAllLights(Elevator_glob)
		case EB_Moving:
			outputDevice.DoorLight(false)
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(Elevator_glob.Dirn)
			outputDevice.MotorDirection(mot_dir)
		case EB_Idle:
			outputDevice.DoorLight(false)
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(Elevator_glob.Dirn)
			outputDevice.MotorDirection(mot_dir)
		}

	default:
		break
	}

	fmt.Println("\nNew state:")
	ElevatorPrint(Elevator_glob)
}


