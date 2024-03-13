package utils

import (
	"Elevator/driver-go-master/elevio"
	"fmt"
)

var (
	ElevatorGlob Elevator
	outputDevice ElevOutputDevice
)

func init() {
	ElevatorGlob = ElevatorInitialized()
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
	ElevatorGlob.Dirn = D_Down
	ElevatorGlob.Behaviour = EB_Moving
}

func FsmOnRequestButtonPress(btnFloor int, btnType Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsmOnRequestButtonPress", btnFloor, ButtonToString(btnType))
	ElevatorPrint(ElevatorGlob)

	switch ElevatorGlob.Behaviour {
	case EB_DoorOpen:
		if RequestsShouldClearImmediately(ElevatorGlob, btnFloor, btnType) {
			TimerStart(ElevatorGlob.DoorOpenDuration_s)
		} else {
			ElevatorGlob.Requests[btnFloor][btnType] = true
		}

	case EB_Moving:
		ElevatorGlob.Requests[btnFloor][btnType] = true

	case EB_Idle:
		ElevatorGlob.Requests[btnFloor][btnType] = true
		//la til den under (gjør at lyset ikke blinker)
		//ElevatorGlob = RequestsClearAtCurrentFloor(ElevatorGlob)
		pair := RequestsChooseDirection(ElevatorGlob)
		ElevatorGlob.Dirn = pair.Dirn
		ElevatorGlob.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case EB_DoorOpen:
			outputDevice.DoorLight(true)
			TimerStart(ElevatorGlob.DoorOpenDuration_s)
			ElevatorGlob = RequestsClearAtCurrentFloor(ElevatorGlob)

		case EB_Moving:
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(ElevatorGlob.Dirn)
			outputDevice.MotorDirection(mot_dir)

		case EB_Idle:
		}
	}

	SetAllLights(ElevatorGlob)

	fmt.Println("\nNew state:")
	ElevatorPrint(ElevatorGlob)
}

func FsmOnFloorArrival(newFloor int, elevators []Elevator) {
	fmt.Printf("\n\n%s(%d)\n", "fsmOnFloorArrival", newFloor)
	ElevatorPrint(ElevatorGlob)

	ElevatorGlob.Floor = newFloor

	outputDevice.FloorIndicator(ElevatorGlob.Floor)

	switch ElevatorGlob.Behaviour {
	case EB_Moving:
		if RequestsShouldStop(ElevatorGlob) {
			outputDevice.MotorDirection(D_Stop)
			outputDevice.DoorLight(true)
			ElevatorGlob = RequestsClearAtCurrentFloor(ElevatorGlob)
			TimerStart(ElevatorGlob.DoorOpenDuration_s)
			SetAllLights(ElevatorGlob)
			ElevatorGlob.Behaviour = EB_DoorOpen
		}
	default:
		break
	}

	fmt.Println("\nNew state:")
	ElevatorPrint(ElevatorGlob)

	//la til denne under
	for _,elev := range elevators {
		if elev.ID == ElevatorGlob.ID {
			elev.Requests = ElevatorGlob.Requests
		}
	}
}

func FsmOnDoorTimeout() {
	fmt.Printf("\n\n%s()\n", "fsmOnDoorTimeout")
	ElevatorPrint(ElevatorGlob)

	switch ElevatorGlob.Behaviour {
	case EB_DoorOpen:
		if Obstructed{
			TimerStart(ElevatorGlob.DoorOpenDuration_s) 
		}else{
		pair := RequestsChooseDirection(ElevatorGlob)
		ElevatorGlob.Dirn = pair.Dirn
		ElevatorGlob.Behaviour = pair.Behaviour
		}
		switch ElevatorGlob.Behaviour {
		case EB_DoorOpen:
			TimerStart(ElevatorGlob.DoorOpenDuration_s)
			ElevatorGlob = RequestsClearAtCurrentFloor(ElevatorGlob)
			SetAllLights(ElevatorGlob)
		case EB_Moving:
			outputDevice.DoorLight(false)
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(ElevatorGlob.Dirn)
			outputDevice.MotorDirection(mot_dir)
		case EB_Idle:
			outputDevice.DoorLight(false)
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(ElevatorGlob.Dirn)
			outputDevice.MotorDirection(mot_dir)
		}

	default:
		break
	}

	fmt.Println("\nNew state:")
	ElevatorPrint(ElevatorGlob)
}
