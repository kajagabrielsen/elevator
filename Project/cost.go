package main

import (
	//"Elevator/driver-go-master/elevio"
	"Elevator/utils"
	//"encoding/json"
	//"fmt"
	//"os/exec"
	"runtime"
	"strconv"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
	ElevID 		int         `json:"id"`
    Behaviour   string      `json:"behaviour"`
    Floor       int         `json:"floor"` 
    Direction   string      `json:"direction"`
    CabRequests []bool      `json:"cabRequests"`
}

type HRAElevStatetemp struct {
	ElevID 		int         `json:"id"`
    Behaviour   utils.ElevatorBehaviour      `json:"behaviour"`
    Floor       int         `json:"floor"` 
    Direction   utils.Dirn      `json:"direction"`
    CabRequests []bool      `json:"cabRequests"`
}

type HRAInput struct {
    HallRequests    [][2]bool                   `json:"hallRequests"`
    States          map[string]HRAElevState     `json:"states"`
}

func CalculateCostFunc(elevators []utils.Elevator) HRAInput{

    hraExecutable := ""
    switch runtime.GOOS {
        case "linux":   hraExecutable  = "hall_request_assigner"
        case "windows": hraExecutable  = "hall_request_assigner.exe"
        default:        panic("OS not supported")
    }

	var n_elevators int
	n_elevators =  len(elevators)
	myStates := []HRAElevStatetemp{}
	GlobalHallCalls := [][2]bool{}
	for i := 0; i < n_elevators; i++ {
		CabCalls := []bool{} 
		HallCalls := [2]bool{}
		for floor := 0; floor < utils.N_FLOORS; floor++{
			CabCalls[floor] = elevators[i].Requests[floor][2]
			HallCalls[0] = elevators[i].Requests[floor][0]
			HallCalls[1] = elevators[i].Requests[floor][1]
			GlobalHallCalls[floor] = HallCalls
		}
		elevastate := HRAElevStatetemp{i, elevators[i].Behaviour, elevators[i].Floor, elevators[i].Dirn, CabCalls}
		myStates[i] = elevastate
	}
	
	input := HRAInput{
		HallRequests: GlobalHallCalls,
		States: make(map[string]HRAElevState),
	}

	for _, elevatorStatus := range myStates {
		input.States[strconv.Itoa(elevatorStatus.ElevID)] = HRAElevState{
			Behaviour : func() string {
				if elevatorStatus.Behaviour == 0 {
					return "idle"
				} else if elevatorStatus.Behaviour == 1 {
					return "door open"
				} else {
					return "moving"
				}
			}(),	
			Floor : elevatorStatus.Floor,
			Direction : func() string {
				if elevatorStatus.Direction == -1 {
					return "down"
				} else if elevatorStatus.Direction == 0 {
					return "stop"
				} else {
					return "up"
				}
			}(),
			CabRequests : elevatorStatus.CabRequests,	
		}
	}



/*
    jsonBytes, err := json.Marshal(input)
    if err != nil {
        fmt.Println("json.Marshal error: ", err)
        return
    }
    
    ret, err := exec.Command(hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
    if err != nil {
        fmt.Println("exec.Command error: ", err)
        fmt.Println(string(ret))
        return
    }
    
    output := new(map[string][][2]bool)
    err = json.Unmarshal(ret, &output)
    if err != nil {
        fmt.Println("json.Unmarshal error: ", err)
        return
    }
        
    fmt.Printf("output: \n")
    for k, v := range *output {
        fmt.Printf("%6v :  %+v\n", k, v)
    }
	*/
	return input
}

