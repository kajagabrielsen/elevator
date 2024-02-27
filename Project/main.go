package main

import (
	"Elevator/driver-go-master/elevio"
	"Elevator/utils"
	"fmt"
	"time"
)

func main() {
	fmt.Print("started")
	var inputPollRate_ms int = utils.InputPollRate

	input := utils.GetInputDevice()

	floorCh := make(chan int)

	// Call the FloorSensor function with the channel
	input.FloorSensor(floorCh)

	// Receive the floor value from the channel
	currentFloor := <-floorCh

	if currentFloor == -1 {
		utils.FsmOnInitBetweenFloors()
	}

	for {
		// Request button
		prev := make([][utils.N_BUTTONS]bool, utils.N_FLOORS)
		for f := 0; f < utils.N_FLOORS; f++ {
			for b := 0; b < utils.N_BUTTONS; b++ {
				v := input.RequestButton(f, elevio.ButtonType(b))
				if v && v != prev[f][b] {
					var BType utils.Button = utils.Button(elevio.ButtonType(b))
					utils.FsmOnRequestButtonPress(f, BType)
				}
				prev[f][b] = v
			}
		}

		// Floor sensor
		var prev_floor int = -1
		f := currentFloor
		if f != -1 && f != prev_floor {
			utils.FsmOnFloorArrival(f)
		}
		prev_floor = f

		// Timer
		if utils.Timer_timedOut() {
			utils.Timer_stop()
			utils.FsmOnDoorTimeout()
		}

		time.Sleep(time.Duration(inputPollRate_ms) * time.Millisecond)
	}
}
