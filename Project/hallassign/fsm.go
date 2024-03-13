package hallassign

import (
	"Elevator/driver-go-master/elevio"
	"Elevator/networkcom"
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
		case F := <-drv_floors:
			utils.FsmOnFloorArrival(F, network.ListOfElevators)
			UpdateCabCalls(utils.ElevatorGlob.Requests)
		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				utils.Obstructed = true
			} else {
				utils.Obstructed = false
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < utils.N_FLOORS; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		case <-time.After(time.Millisecond * time.Duration(utils.DoorOpenDuration*1000)):
			utils.FsmOnDoorTimeout()
		}
		AssignHallRequest()
		for floor_num, floor := range OneElevRequests {
			for btn_num := range floor {
				if OneElevRequests[floor_num][btn_num] {
					utils.FsmOnRequestButtonPress(floor_num, utils.Button(btn_num))
				}
			}
		}
	}
}
