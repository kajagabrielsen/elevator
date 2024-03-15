package fsm

import (
	"Elevator/driver_go_master/elevio"
	"Elevator/elevator/initialize"
	"Elevator/elevator/log"
	"Elevator/elevator/request"
	"Elevator/utils"
	"fmt"
)

var GlobalHallCalls = [initialize.N_FLOORS][2]bool{}

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

func SetAllLights(es initialize.Elevator) {
	for floor := 0; floor < initialize.N_FLOORS; floor++ {
		for btn := 0; btn < initialize.N_BUTTONS-1; btn++ {
			var B elevio.ButtonType = elevio.ButtonType(btn)
			initialize.OutputDevice.RequestButtonLight(floor, B, GlobalHallCalls[floor][btn])
		}
	}

	for f := 0; f < initialize.N_FLOORS; f++ {
		var b elevio.ButtonType = elevio.ButtonType(2)
		initialize.OutputDevice.RequestButtonLight(f, b, es.Requests[f][2])
	}
}


func FsmOnInitBetweenFloors() {
	initialize.OutputDevice.MotorDirection(elevio.MDDown)
	initialize.ElevatorGlob.Dirn = elevio.MDDown
	initialize.ElevatorGlob.Behaviour = initialize.EB_Moving
}

func FsmOnRequestButtonPress(btnFloor int, btnType elevio.ButtonType) {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsmOnRequestButtonPress", btnFloor, ButtonToString(btnType))
	log.ElevatorLog(initialize.ElevatorGlob)

	switch initialize.ElevatorGlob.Behaviour {
	case initialize.EB_DoorOpen:
		if request.RequestsShouldClearImmediately(initialize.ElevatorGlob, btnFloor, btnType) {
			utils.TimerStart(initialize.ElevatorGlob.DoorOpenDuration)
		} else {
			initialize.ElevatorGlob.Requests[btnFloor][btnType] = true
		}

	case initialize.EB_Moving:
		initialize.ElevatorGlob.Requests[btnFloor][btnType] = true

	case initialize.EB_Idle:
		initialize.ElevatorGlob.Requests[btnFloor][btnType] = true
		pair := request.RequestsChooseDirection(initialize.ElevatorGlob)
		initialize.ElevatorGlob.Dirn = pair.Dirn
		initialize.ElevatorGlob.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case initialize.EB_DoorOpen:
			initialize.OutputDevice.DoorLight(true)
			utils.TimerStart(initialize.ElevatorGlob.DoorOpenDuration)
			initialize.ElevatorGlob = request.RequestsClearAtCurrentFloor(initialize.ElevatorGlob)

		case initialize.EB_Moving:
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(initialize.ElevatorGlob.Dirn)
			initialize.OutputDevice.MotorDirection(mot_dir)

		case initialize.EB_Idle:
		}
	}

	SetAllLights(initialize.ElevatorGlob)

	fmt.Println("\nNew state:")
	log.ElevatorLog(initialize.ElevatorGlob)
}

func FsmOnFloorArrival(newFloor int, elevators []initialize.Elevator) {
	fmt.Printf("\n\n%s(%d)\n", "fsmOnFloorArrival", newFloor)
	log.ElevatorLog(initialize.ElevatorGlob)

	initialize.ElevatorGlob.Floor = newFloor

	initialize.OutputDevice.FloorIndicator(initialize.ElevatorGlob.Floor)

	switch initialize.ElevatorGlob.Behaviour {
	case initialize.EB_Moving:
		if request.RequestsShouldStop(initialize.ElevatorGlob) {
			initialize.OutputDevice.MotorDirection(elevio.MDStop)
			initialize.OutputDevice.DoorLight(true)
			initialize.ElevatorGlob = request.RequestsClearAtCurrentFloor(initialize.ElevatorGlob)
			utils.TimerStart(initialize.ElevatorGlob.DoorOpenDuration)
			SetAllLights(initialize.ElevatorGlob)
			initialize.ElevatorGlob.Behaviour = initialize.EB_DoorOpen
		}
	default:
		break
	}

	fmt.Println("\nNew state:")
	log.ElevatorLog(initialize.ElevatorGlob)

	for _,elev := range elevators {
		if elev.ID == initialize.ElevatorGlob.ID {
			elev.Requests = initialize.ElevatorGlob.Requests
		}
	}
}

func FsmOnDoorTimeout() {
	fmt.Printf("\n\n%s()\n", "fsmOnDoorTimeout")
	log.ElevatorLog(initialize.ElevatorGlob)

	switch initialize.ElevatorGlob.Behaviour {
	case initialize.EB_DoorOpen:
		if initialize.ElevatorGlob.Obstructed{
			utils.TimerStart(initialize.ElevatorGlob.DoorOpenDuration) 
		}else{
		pair := request.RequestsChooseDirection(initialize.ElevatorGlob)
		initialize.ElevatorGlob.Dirn = pair.Dirn
		initialize.ElevatorGlob.Behaviour = pair.Behaviour
		}
		switch initialize.ElevatorGlob.Behaviour {
		case initialize.EB_DoorOpen:
			utils.TimerStart(initialize.ElevatorGlob.DoorOpenDuration)
			initialize.ElevatorGlob = request.RequestsClearAtCurrentFloor(initialize.ElevatorGlob)
			SetAllLights(initialize.ElevatorGlob)
		case initialize.EB_Moving:
			initialize.OutputDevice.DoorLight(false)
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(initialize.ElevatorGlob.Dirn)
			initialize.OutputDevice.MotorDirection(mot_dir)
		case initialize.EB_Idle:
			initialize.OutputDevice.DoorLight(false)
			var mot_dir elevio.MotorDirection = elevio.MotorDirection(initialize.ElevatorGlob.Dirn)
			initialize.OutputDevice.MotorDirection(mot_dir)
		}

	default:
		break
	}

	fmt.Println("\nNew state:")
	log.ElevatorLog(initialize.ElevatorGlob)
}
