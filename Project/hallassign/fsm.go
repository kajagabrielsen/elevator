package hallassign

import (
	"Elevator/driver-go-master/elevio"
	network "Elevator/networkcom"
	"Elevator/networkcom/network/listUpdate"
	"Elevator/utils"
	"fmt"
	"time"
)

func FSM(HelloRx  chan network.HelloMsg, 
	drv_buttons   chan elevio.ButtonEvent, 
	drv_floors    chan int, drv_obstr chan bool, 
	drv_stop      chan bool) {
	for {
		select {
		case E := <-drv_buttons:
			utils.ElevatorGlob.Requests = OneElevRequests
			utils.ElevatorGlob.Requests[E.Floor][E.Button] = true
			UpdateCabCalls(utils.ElevatorGlob.Requests)
			UpdateGlobalHallCalls(network.ListOfElevators)
		case F := <-drv_floors:
			utils.FsmOnFloorArrival(F, network.ListOfElevators)
			UpdateCabCalls(utils.ElevatorGlob.Requests)
			UpdateGlobalHallCalls(network.ListOfElevators)
		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				utils.ElevatorGlob.Obstructed = true
				listUpdate.RemoveFromListOfElevators(network.ListOfElevators, utils.ElevatorGlob)

			} else {
				utils.ElevatorGlob.Obstructed = false
				listUpdate.AddToListOfElevators(network.ListOfElevators, utils.ElevatorGlob)

			}
		case a := <-drv_stop:
			if a{
				elevio.SetMotorDirection(elevio.MDStop)
			}else{
				elevio.SetMotorDirection(elevio.MotorDirection(utils.ElevatorGlob.Dirn))
			}
		case <-time.After(time.Millisecond * time.Duration(utils.DoorOpenDuration*1000)):
			utils.FsmOnDoorTimeout()
			UpdateGlobalHallCalls(network.ListOfElevators)
		}
		utils.SetAllLights(utils.ElevatorGlob)
		for floor_num, floor := range OneElevRequests {
			for btn_num := range floor {
				if OneElevRequests[floor_num][btn_num] {
					utils.FsmOnRequestButtonPress(floor_num, elevio.ButtonType(btn_num))
				}
			}
		}
		AssignHallRequest()
		
	}
}
