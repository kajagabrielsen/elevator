package hallassign

import (
	"Elevator/driver_go_master/elevio"
	"Elevator/elevator/initialize"
	"Elevator/elevator/fsm_func"
	"Elevator/network/list"
	"fmt"
	"time"
)

func FSM(drv_buttons   chan elevio.ButtonEvent, 
		drv_floors    chan int, 
		drv_obstr 	  chan bool, 
		drv_stop      chan bool) {
	for {
		select {
		case E := <-drv_buttons:
			initialize.ElevatorGlob.Requests = OneElevRequests
			initialize.ElevatorGlob.Requests[E.Floor][E.Button] = true
			UpdateCabCalls(initialize.ElevatorGlob.Requests)
			UpdateGlobalHallCalls(list.ListOfElevators)
		case F := <-drv_floors:
			fsm.FsmOnFloorArrival(F, list.ListOfElevators)
			UpdateCabCalls(initialize.ElevatorGlob.Requests)
			UpdateGlobalHallCalls(list.ListOfElevators)
		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				initialize.ElevatorGlob.Obstructed = true
				list.RemoveFromListOfElevators(list.ListOfElevators, initialize.ElevatorGlob)

			} else {
				initialize.ElevatorGlob.Obstructed = false
				list.AddToListOfElevators(list.ListOfElevators, initialize.ElevatorGlob)

			}
		case a := <-drv_stop:
			if a{
				elevio.SetMotorDirection(elevio.MDStop)
			}else{
				elevio.SetMotorDirection(elevio.MotorDirection(initialize.ElevatorGlob.Dirn))
			}
		case <-time.After(time.Millisecond * time.Duration(initialize.ElevatorGlob.DoorOpenDuration*1000)):
			fsm.FsmOnDoorTimeout()
			UpdateGlobalHallCalls(list.ListOfElevators)
		}
		fsm.SetAllLights(initialize.ElevatorGlob)
		for floor_num, floor := range OneElevRequests {
			for btn_num := range floor {
				if OneElevRequests[floor_num][btn_num] {
					fsm.FsmOnRequestButtonPress(floor_num, elevio.ButtonType(btn_num))
				}
			}
		}
		AssignHallRequest()
		
	}
}
