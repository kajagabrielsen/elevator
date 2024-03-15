package utils

import (
	"Elevator/driver-go-master/elevio"
)

func Init(addr string, numFloors int) {
	elevio.InitDriver(addr, numFloors)
}

func WrapRequestButton(floor int, button elevio.ButtonType) bool {
	return elevio.GetButton(button, floor)
}

func WrapRequestButtonLight(floor int, button elevio.ButtonType, value bool) {
	elevio.SetButtonLamp(button, floor, value)
}

func WrapMotorDirection(direction elevio.MotorDirection) {
	elevio.SetMotorDirection(direction)
}

// GetInputDevice returns the elevator's input device.
func GetInputDevice() ElevInputDevice {
	return ElevInputDevice{
		FloorSensor:   elevio.PollFloorSensor,
		RequestButton: WrapRequestButton,
		StopButton:    elevio.PollStopButton,
		Obstruction:   elevio.PollObstructionSwitch,
	}
}

// GetOutputDevice returns the elevator's output device.
func GetOutputDevice() ElevOutputDevice {
	return ElevOutputDevice{
		FloorIndicator:     elevio.SetFloorIndicator,
		RequestButtonLight: WrapRequestButtonLight,
		DoorLight:          elevio.SetDoorOpenLamp,
		StopButtonLight:    elevio.SetStopLamp,
		MotorDirection:     WrapMotorDirection,
	}
}

// DirnToString converts Direction to a string.
func DirnToString(d elevio.MotorDirection) string {
	switch d {
	case elevio.MDUp:
		return "D_Up"
	case elevio.MDDown:
		return "D_Down"
	case elevio.MDStop:
		return "D_Stop"
	default:
		return "D_UNDEFINED"
	}
}

// ButtonToString converts Button to a string.
func ButtonToString(b elevio.ButtonType) string {
	switch b {
	case elevio.BTHallUp:
		return "B_HallUp"
	case elevio.BTHallDown:
		return "B_HallDown"
	case elevio.BTCab:
		return "B_Cab"
	default:
		return "B_UNDEFINED"
	}
}