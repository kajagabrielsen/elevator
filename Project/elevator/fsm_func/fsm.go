package fsm

import (
	"Elevator/driver_go_master/elevio"
	"Elevator/elevator/initial"
	"Elevator/elevator/log"
	"Elevator/elevator/request"
	"fmt"
)

var GlobalHallCalls = [initial.NFloors][2]bool{}

func ButtonToString(b elevio.ButtonType) string {
	switch b {
	case elevio.BTHallUp:
		return "B_HallUp"
	case elevio.BTHallDown:
		return "B_HallDown"
	case elevio.BTCab:
		return "B_Cab"
	default:
		return "B_UNDEFINED"
	}
}

func SetAllLights(es initial.Elevator) {
	for floor := 0; floor < initial.NFloors; floor++ {
		for btn := 0; btn < initial.NButtons-1; btn++ {
			var B elevio.ButtonType = elevio.ButtonType(btn)
			initial.OutputDevice.RequestButtonLight(floor, B, GlobalHallCalls[floor][btn])
		}
	}

	for f := 0; f < initial.NFloors; f++ {
		var b elevio.ButtonType = elevio.ButtonType(2)
		initial.OutputDevice.RequestButtonLight(f, b, es.Requests[f][2])
	}
}


func FsmOnInitBetweenFloors() {
	initial.OutputDevice.MotorDirection(elevio.MDDown)
	initial.ElevatorGlob.Dirn = elevio.MDDown
	initial.ElevatorGlob.Behaviour = initial.EBMoving
}

func FsmOnRequestButtonPress(btnFloor int, btnType elevio.ButtonType) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsmOnRequestButtonPress", btnFloor, ButtonToString(btnType))
	log.ElevatorLog(initial.ElevatorGlob)

	switch initial.ElevatorGlob.Behaviour {
	case initial.EBDoorOpen:
		if !request.RequestsShouldClearImmediately(initial.ElevatorGlob, btnFloor, btnType) {
			initial.ElevatorGlob.Requests[btnFloor][btnType] = true
		}

	case initial.EBMoving:
		initial.ElevatorGlob.Requests[btnFloor][btnType] = true

	case initial.EBIdle:
		initial.ElevatorGlob.Requests[btnFloor][btnType] = true
		pair := request.RequestsChooseDirection(initial.ElevatorGlob)
		initial.ElevatorGlob.Dirn = pair.Dirn
		initial.ElevatorGlob.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case initial.EBDoorOpen:
			initial.OutputDevice.DoorLight(true)
			initial.ElevatorGlob = request.RequestsClearAtCurrentFloor(initial.ElevatorGlob)
		case initial.EBMoving:
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(initial.ElevatorGlob.Dirn)
			initial.OutputDevice.MotorDirection(mot_dir)

		case initial.EBIdle:
		}
	}

	SetAllLights(initial.ElevatorGlob)

	fmt.Println("\nNew state:")
	log.ElevatorLog(initial.ElevatorGlob)
}

func FsmOnFloorArrival(newFloor int, elevators []initial.Elevator) {
	fmt.Printf("\n\n%s(%d)\n", "fsmOnFloorArrival", newFloor)
	log.ElevatorLog(initial.ElevatorGlob)

	initial.ElevatorGlob.Floor = newFloor

	initial.OutputDevice.FloorIndicator(initial.ElevatorGlob.Floor)

	switch initial.ElevatorGlob.Behaviour {
	case initial.EBMoving:
		if request.RequestsShouldStop(initial.ElevatorGlob) {
			initial.OutputDevice.MotorDirection(elevio.MDStop)
			initial.OutputDevice.DoorLight(true)
			initial.ElevatorGlob = request.RequestsClearAtCurrentFloor(initial.ElevatorGlob)
			SetAllLights(initial.ElevatorGlob)
			initial.ElevatorGlob.Behaviour = initial.EBDoorOpen
		}
	default:
		break
	}

	fmt.Println("\nNew state:")
	log.ElevatorLog(initial.ElevatorGlob)

	for _,elev := range elevators {
		if elev.ID == initial.ElevatorGlob.ID {
			elev.Requests = initial.ElevatorGlob.Requests
		}
	}
}

func FsmOnDoorTimeout() {
	fmt.Printf("\n\n%s()\n", "fsmOnDoorTimeout")
	log.ElevatorLog(initial.ElevatorGlob)

	switch initial.ElevatorGlob.Behaviour {
	case initial.EBDoorOpen:
		if !initial.ElevatorGlob.Obstructed{
		pair := request.RequestsChooseDirection(initial.ElevatorGlob)
		initial.ElevatorGlob.Dirn = pair.Dirn
		initial.ElevatorGlob.Behaviour = pair.Behaviour
		}
		switch initial.ElevatorGlob.Behaviour {
		case initial.EBDoorOpen:
			initial.ElevatorGlob = request.RequestsClearAtCurrentFloor(initial.ElevatorGlob)
			SetAllLights(initial.ElevatorGlob)
		case initial.EBMoving:
			initial.OutputDevice.DoorLight(false)
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(initial.ElevatorGlob.Dirn)
			initial.OutputDevice.MotorDirection(mot_dir)
		case initial.EBIdle:
			initial.OutputDevice.DoorLight(false)
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(initial.ElevatorGlob.Dirn)
			initial.OutputDevice.MotorDirection(mot_dir)
		}

	default:
		break
	}

	fmt.Println("\nNew state:")
	log.ElevatorLog(initial.ElevatorGlob)
}
