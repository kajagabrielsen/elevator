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
	elevator = ElevatorInitialized()

	outputDevice = GetOutputDevice()
}

func GetElevator() Elevator{
	return elevator
}

func SetElevator(floor int, dirn Dirn, request [N_FLOORS][N_BUTTONS]bool, behaviour ElevatorBehaviour, clear ClearRequestVariantInt, duration float64, id string) {
	elevator.Floor = floor
	elevator.Dirn = dirn
	elevator.Requests = request
	elevator.Behaviour = behaviour
	elevator.ClearRequestVariant = clear
	elevator.DoorOpenDuration_s = duration
	elevator.ID = id
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
		fmt.Printf("fsm_door")
		if RequestsShouldClearImmediately(elevator, btnFloor, btnType) {
			TimerStart(elevator.DoorOpenDuration_s)
		} else {
			elevator.Requests[btnFloor][btnType] = true
		}

	case EB_Moving:
		elevator.Requests[btnFloor][btnType] = true
		fmt.Printf("fsm_move")

	case EB_Idle:
		fmt.Printf("fsm_idle")
		elevator.Requests[btnFloor][btnType] = true
		pair := RequestsChooseDirection(elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case EB_DoorOpen:
			fmt.Printf("fsm_idle_door")
			outputDevice.DoorLight(true)
			TimerStart(elevator.DoorOpenDuration_s)
			elevator = RequestsClearAtCurrentFloor(elevator)

		case EB_Moving:
			fmt.Printf("fsm_idle_move")
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(elevator.Dirn)
			outputDevice.MotorDirection(mot_dir)

		case EB_Idle:
			fmt.Printf("fsm_idle_idle")
		}
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
		fmt.Printf("ofa_move")
		if RequestsShouldStop(elevator) {
			outputDevice.MotorDirection(D_Stop)
			outputDevice.DoorLight(true)
			elevator = RequestsClearAtCurrentFloor(elevator)
			TimerStart(elevator.DoorOpenDuration_s)
			SetAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen
		}
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
		fmt.Printf("odt_door")
		pair := RequestsChooseDirection(elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour

		switch elevator.Behaviour {
		case EB_DoorOpen:
			fmt.Printf("odt_door_door")
			TimerStart(elevator.DoorOpenDuration_s)
			elevator = RequestsClearAtCurrentFloor(elevator)
			SetAllLights(elevator)
		case EB_Moving:
			fmt.Printf("odt_door_move")
			outputDevice.DoorLight(false)
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(elevator.Dirn)
			outputDevice.MotorDirection(mot_dir)
		case EB_Idle:
			fmt.Printf("odt_door_idle")
			outputDevice.DoorLight(false)
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(elevator.Dirn)
			outputDevice.MotorDirection(mot_dir)
		}

	default:
		break
	}

	fmt.Println("\nNew state:")
	ElevatorPrint(elevator)
}

