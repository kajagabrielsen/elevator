package main

import (
	"Elevator/driver-go-master/elevio"
	"Elevator/utils"
	"fmt"
	"time"
)
func main(){       

    numFloors := 4

    elevio.Init("localhost:15657", numFloors)
    
    var d elevio.MotorDirection = elevio.MD_Up
    //elevio.SetMotorDirection(d)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)    
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)

	utils.FsmOnInitBetweenFloors()


for{
    select{
    case E := <- drv_buttons:
        fmt.Printf("button")
        utils.FsmOnRequestButtonPress(E.Floor, utils.Button(E.Button))
    case F := <- drv_floors:
        fmt.Printf("floor")
        utils.FsmOnFloorArrival(F)
    case a := <- drv_obstr:
        fmt.Printf("obs")
        fmt.Printf("%+v\n", a)
        if a {
            elevio.SetMotorDirection(elevio.MD_Stop)
        } else {
            elevio.SetMotorDirection(d)
        }
        
    case a := <- drv_stop:
        fmt.Printf("stop")
        fmt.Printf("%+v\n", a)
        for f := 0; f < numFloors; f++ {
            for b := elevio.ButtonType(0); b < 3; b++ {
                elevio.SetButtonLamp(b, f, false)
            }
        }
    case <- time.After(time.Millisecond*time.Duration(utils.DoorOpenDuration*1000)):
        utils.FsmOnDoorTimeout()
    }
    }
}
