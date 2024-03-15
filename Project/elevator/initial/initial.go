package initial

import (
	"Elevator/driver_go_master/elevio"
)

// init function runs before main
func init() {
	ElevatorGlob = ElevatorInitialized()
	OutputDevice = GetOutputDevice()
}

const (
	N_FLOORS = 4
	N_BUTTONS = 3
)

type ElevatorBehaviour int

const (
	EBIdle ElevatorBehaviour = iota
	EBDoorOpen
	EBMoving
)

type ClearRequestVariantInt int

const (
	CV_All ClearRequestVariantInt = iota
	CV_InDirn
)

var ElevatorGlob Elevator

type Elevator struct {
	Floor                int
	Dirn                 elevio.MotorDirection
	Requests             [N_FLOORS][N_BUTTONS]bool
	Behaviour            ElevatorBehaviour
	ClearRequestVariant  ClearRequestVariantInt
	DoorOpenDuration   	 float64
	ID                   string
	Obstructed           bool

}

var OutputDevice ElevOutputDevice

type ElevOutputDevice struct {
	FloorIndicator     func(int)
	RequestButtonLight func(int, elevio.ButtonType, bool)
	DoorLight          func(bool)
	StopButtonLight     func(bool)
	MotorDirection      func(elevio.MotorDirection)
}

//Initializing the elevator with starting values
func ElevatorInitialized() Elevator {
	return Elevator{
		Floor:     1,
		Dirn:      elevio.MDStop,
		Behaviour: EBIdle,
		ClearRequestVariant: CV_All,
		DoorOpenDuration:   3.0,
		ID: "0",
	}
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