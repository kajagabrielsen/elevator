package hallassign

import (
	"Elevator/driver-go-master/elevio"
	"Elevator/networkcom"
	"Elevator/utils"
	"fmt"
	"time"
)

func FSM(buttonPressCh chan elevio.ButtonEvent, drv_buttons chan elevio.ButtonEvent, drv_floors chan int, drv_obstr chan bool, drv_stop chan bool){
	var d elevio.MotorDirection = elevio.MD_Up
	for {
		select {
		case E := <-drv_buttons:
			utils.Elevator_glob.Requests[E.Floor][E.Button]=true
			AssignHallRequest()
			for floor_num, floor := range OneElevRequests{
				for btn_num, _ := range floor {
					if OneElevRequests[floor_num][btn_num]{
						utils.FsmOnRequestButtonPress(floor_num, utils.Button(btn_num))
						AssignHallRequest()
					}
				}
			}
		case F := <-drv_floors:
			AssignHallRequest()
			utils.FsmOnFloorArrival(F)
			
		case a := <-drv_obstr:
			AssignHallRequest()
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}
			

		case a := <-drv_stop:
			AssignHallRequest()
			fmt.Printf("%+v\n", a)
			for f := 0; f < utils.N_FLOORS; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		case <-time.After(time.Millisecond * time.Duration(utils.DoorOpenDuration*1000)):
			utils.FsmOnDoorTimeout()
		}
	}
}

func GetIndex(key string, list []string) int {
	for i, value := range list {
		if value == key {
			return i
		}
	}

	return 0
}

var OneElevRequests = [utils.N_FLOORS][utils.N_BUTTONS]bool{}

func AssignHallRequest() {
	ListOfElevators := network.ListOfElevators
    AssignedHallCalls := CalculateCostFunc(ListOfElevators)
    OneElevCabCalls := GetCabCalls(utils.Elevator_glob)
    OneElevHallCalls := AssignedHallCalls[utils.Elevator_glob.ID]

    for floor := 0; floor < utils.N_FLOORS; floor++ {
        OneElevRequests[floor][0] = OneElevHallCalls[floor][0]
        OneElevRequests[floor][1] = OneElevHallCalls[floor][1]
        OneElevRequests[floor][2] = OneElevCabCalls[floor]
    }
	fmt.Println(OneElevRequests)
   //utils.Elevator_glob.Requests = OneElevRequests



}


func HandleButtonPressUpdate( buttonPressCh chan elevio.ButtonEvent){
	for {
        select {
        case btn := <-buttonPressCh:
            utils.Elevator_glob.Requests[btn.Floor][btn.Button] = true
            flag := 0
            for i, element := range network.ListOfElevators {
                if element.ID == utils.Elevator_glob.ID {
                    network.ListOfElevators[i] = utils.Elevator_glob
                    flag = 1
                }
            }
            if flag == 0 {
                network.ListOfElevators = append(network.ListOfElevators, utils.Elevator_glob)
            }

			//AssignHallRequest()

        }
    }
}
