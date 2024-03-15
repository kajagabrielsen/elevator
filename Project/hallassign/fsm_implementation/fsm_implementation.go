package implfsm

import (
	"Elevator/driver_go_master/elevio"
	fsm "Elevator/elevator/fsm_func"
	"Elevator/elevator/initial"
	assign "Elevator/hallassign/assign_hall_request"
	call "Elevator/hallassign/call_handling"
	"Elevator/network/list"
	"fmt"
	"time"
)

func FSM(drv_buttons   chan elevio.ButtonEvent, 
		 drv_floors    chan int, 
		 drv_obstr 	   chan bool, 
		 drv_stop      chan bool) {
	for {
		select {
		case E := <-drv_buttons:
			initial.ElevatorGlob.Requests = assign.OneElevRequests
			initial.ElevatorGlob.Requests[E.Floor][E.Button] = true
			call.UpdateCabCalls(initial.ElevatorGlob.Requests)
			call.UpdateGlobalHallCalls(list.ListOfElevators)

		case F := <-drv_floors:
			fsm.FsmOnFloorArrival(F, list.ListOfElevators)
			call.UpdateCabCalls(initial.ElevatorGlob.Requests)
			call.UpdateGlobalHallCalls(list.ListOfElevators)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				initial.ElevatorGlob.Obstructed = true
				call.UpdateGlobalHallCalls(list.ListOfElevators)
				list.RemoveFromListOfElevators(list.ListOfElevators, initial.ElevatorGlob.ID)

			} else {
				initial.ElevatorGlob.Obstructed = false
				call.UpdateGlobalHallCalls(list.ListOfElevators)
				list.AddToListOfElevators(list.ListOfElevators, initial.ElevatorGlob)

			}
		case a := <-drv_stop:
			if a{
				elevio.SetMotorDirection(elevio.MDStop)
			}else{
				elevio.SetMotorDirection(elevio.MotorDirection(initial.ElevatorGlob.Dirn))
			}

		case <-time.After(time.Millisecond * time.Duration(initial.ElevatorGlob.DoorOpenDuration*1000)):
			fsm.FsmOnDoorTimeout()
			call.UpdateGlobalHallCalls(list.ListOfElevators)
		}
		
		fsm.SetAllLights(initial.ElevatorGlob)
		for floor_num, floor := range assign.OneElevRequests {
			for btn_num := range floor {
				if assign.OneElevRequests[floor_num][btn_num] {
					fsm.FsmOnRequestButtonPress(floor_num, elevio.ButtonType(btn_num))
				}
			}
		}
		assign.AssignHallRequest()
		
	}
}
