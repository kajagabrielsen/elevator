package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	N_FLOORS  = 4
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
	FloorSensor   func() int
	RequestButton func(int, Button) int
	StopButton    func() int
	Obstruction   func() int
}

type ElevOutputDevice struct {
	FloorIndicator     func(int)
	RequestButtonLight func(int, Button, int)
	DoorLight          func(int)
	StopButtonLight    func(int)
	MotorDirection     func(Dirn)
}

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type ClearRequestVariant int

const (
	CV_All ClearRequestVariant = iota
	CV_InDirn
)

type Elevator struct {
	Floor     int
	Dirn      Dirn
	Requests  [N_FLOORS][N_BUTTONS]int
	Behaviour ElevatorBehaviour
	Config    struct {
		ClearRequestVariant ClearRequestVariant
		DoorOpenDuration_s  float64
	}
}

type DirnBehaviourPair struct {
	Dirn      Dirn
	Behaviour ElevatorBehaviour
}

func conLoad(file string, cases func(string, string)) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Unable to open config file %s\n", file)
		return
	}
	defer f.Close()

	var line string
	for {
		_, err := fmt.Fscanf(f, "%s", &line)
		if err != nil {
			break
		}
		if strings.HasPrefix(line, "--") {
			var key, val string
			fmt.Fscanf(f, "--%s %s", &key, &val)
			cases(key, val)
		}
	}
}

func conVal(key string, variable interface{}, format string) {
	if strings.EqualFold(key, _key) {
		fmt.Sscanf(_val, format, variable)
	}
}

func conEnum(key string, variable interface{}, matchCases func()) {
	if strings.EqualFold(key, _key) {
		var v interface{}
		matchCases()
		reflect.ValueOf(variable).Set(reflect.ValueOf(v))
	}
}

func conMatch(id string) {
	if strings.EqualFold(_val, id) {
		_v = id
	}
}
