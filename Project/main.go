package main

import (
	"Elevator/driver-go-master/elevio"
	"fmt"
	"time"
)

func main() {
	fmt.Print("started")
	var inputPollRate_ms int = inputPollRate

	input := elevio.GetInputDevice()

	if input.FloorSensor() == -1 {
		fsmOnInitBetweenFloors()
	}

	for {
		// Request button
		prev := make([][N_BUTTONS]int, N_FLOORS)
		for f := 0; f < N_FLOORS; f++ {
			for b := 0; b < N_BUTTONS; b++ {
				v := input.RequestButton(f, elevio.ButtonType(b))
				if v != 0 && v != prev[f][b] {
					var BType Button = Button(elevio.ButtonType(b))
					fsmOnRequestButtonPress(f, BType)
				}
				prev[f][b] = v
			}
		}

		// Floor sensor
		var prev_floor int = -1
		f := input.FloorSensor()
		if f != -1 && f != prev_floor {
			fsmOnFloorArrival(f)
		}
		prev_floor = f

		// Timer
		if timer_timedOut() {
			timer_stop()
			fsmOnDoorTimeout()
		}

		time.Sleep(time.Duration(inputPollRate_ms) * time.Millisecond)
	}
}
