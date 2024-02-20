package main

import (
	"Elevator/driver-go-master/elevio"
)

func Init(addr string, numFloors int) {
	elevio.Init(addr, numFloors)
}

func wrapRequestButton(floor int, button elevio.ButtonType) bool {
	return elevio.GetButton(button, floor)
}

func wrapRequestButtonLight(floor int, button elevio.ButtonType, value bool) {
	elevio.SetButtonLamp(button, floor, value)
}

func wrapMotorDirection(direction elevio.MotorDirection) {
	elevio.SetMotorDirection(direction)
}

// GetInputDevice returns the elevator's input device.
func GetInputDevice() ElevInputDevice {
	return ElevInputDevice{
		FloorSensor:   elevio.PollFloorSensor,//elevatorHardwareGetFloorSensorSignal,
		RequestButton: wrapRequestButton,
		StopButton:    elevio.PollStopButton,//elevatorHardwareGetStopSignal,
		Obstruction:   elevio.PollObstructionSwitch,//elevatorHardwareGetObstructionSignal,
	}
}

// GetOutputDevice returns the elevator's output device.
func GetOutputDevice() ElevOutputDevice {
	return ElevOutputDevice{
		FloorIndicator:     elevio.SetFloorIndicator,//elevatorHardwareSetFloorIndicator,
		RequestButtonLight: wrapRequestButtonLight,
		DoorLight:          elevio.SetDoorOpenLamp,//elevatorHardwareSetDoorOpenLamp,
		StopButtonLight:    elevio.SetStopLamp,//elevatorHardwareSetStopLamp,
		MotorDirection:     wrapMotorDirection,
	}
}

// DirnToString converts Direction to a string.
func DirnToString(d Dirn) string {
	switch d {
	case D_Up:
		return "D_Up"
	case D_Down:
		return "D_Down"
	case D_Stop:
		return "D_Stop"
	default:
		return "D_UNDEFINED"
	}
}

// ButtonToString converts Button to a string.
func ButtonToString(b Button) string {
	switch b {
	case B_HallUp:
		return "B_HallUp"
	case B_HallDown:
		return "B_HallDown"
	case B_Cab:
		return "B_Cab"
	default:
		return "B_UNDEFINED"
	}
}