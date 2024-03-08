package network

import (
	"Elevator/utils"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//
//	will be received as zero-values.
type HelloMsg struct {
	Elevator utils.Elevator
	Iter     int
}

var AliveElevatorsID []string

var ListOfElevators []utils.Elevator

