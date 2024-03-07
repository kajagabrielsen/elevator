package utils

import (
	"Elevator/driver-go-master/elevio"
)

// Configuration values
var (
	ClearRequestVariantString = "CV_InDirn"
	DoorOpenDuration    = 3.0
	InputPollRate       = 25
)


const (
	N_FLOORS = 4
	N_BUTTONS = 3
)

type Dirn int

const (
	D_Down Dirn = -1
	D_Stop      = 0
	D_Up        = 1
)

type Button int

const (
	B_HallUp Button = iota
	B_HallDown
	B_Cab
)

type ElevInputDevice struct {
	FloorSensor    func(chan<- int)
	RequestButton  func(int, elevio.ButtonType) bool
	StopButton     func(chan<- bool)
	Obstruction    func(chan<- bool)
}


type ElevOutputDevice struct {
	FloorIndicator     func(int)
	RequestButtonLight func(int, elevio.ButtonType, bool)
	DoorLight          func(bool)
	StopButtonLight     func(bool)
	MotorDirection      func(elevio.MotorDirection)
}

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type ClearRequestVariantInt int

const (
	CV_All ClearRequestVariantInt = iota
	CV_InDirn
)

type Elevator struct {
	Floor                int
	Dirn                 Dirn
	Requests             [N_FLOORS][N_BUTTONS]bool
	Behaviour            ElevatorBehaviour
	ClearRequestVariant  ClearRequestVariantInt
	DoorOpenDuration_s   float64
	ID                   string

}

type DirnBehaviourPair struct {
	Dirn      Dirn
	Behaviour ElevatorBehaviour
}
