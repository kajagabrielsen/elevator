package hallassign

import (
	"Elevator/driver-go-master/elevio"
	"Elevator/networkcom"
	"Elevator/utils"
	"fmt"
	"time"
)

func FSM(HelloRx chan network.HelloMsg, drv_buttons chan elevio.ButtonEvent, drv_floors chan int, drv_obstr chan bool, drv_stop chan bool) {
	var d elevio.MotorDirection = elevio.MD_Up
	for {
		select {
		case E := <-drv_buttons:
			utils.ElevatorGlob.Requests = OneElevRequests
			utils.ElevatorGlob.Requests[E.Floor][E.Button] = true
		case F := <-drv_floors:
			utils.FsmOnFloorArrival(F)
			OneElevRequests[F] = [utils.N_BUTTONS]bool{false, false, false}
		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
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
