package main

import (
	"Elevator/driver-go-master/elevio"
	"encoding/json"
	"fmt"
	"os/exec"
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

type HRAInput struct {
    HallRequests    [][2]bool                   `json:"hallRequests"`
    States          map[string]HRAElevState     `json:"states"`
}

func CalculateCostFunc(elevators []Elevator){

    hraExecutable := ""
    switch runtime.GOOS {
        case "linux":   hraExecutable  = "hall_request_assigner"
        case "windows": hraExecutable  = "hall_request_assigner.exe"
        default:        panic("OS not supported")
    }


	n_elevators =  len(elevators)
	myStates := [n_elevators]HRAElevState{}
	for i := 0; i < n_elevators; i++ {
		CabCalls := [n_elevators]bool 
		for floor := 0; floor < N_FLOORS; floor++{
			CabCalls[floor] = elevators[i].Requests[floor][2]
		}
		elevastate := HRAElevState{i, elevators[i].Behaviour, elevators[i].Floor, elevators[i].Dirn, CabCalls}
		myStates[i] = elevastate
	}
	
	for 

	input := HRAInput{
		HallRequests: elevio.,
		States: make(map[string]HRAElevState),
	}

	for _, elevatorStatus := range myStates {
		input.States[strconv.Itoa(elevatorStatus.ElevID)] = HRAElevState{
			Behavior : func() string {
				if elevatorStatus.Behaviour == 0 {
					return "idle"
				}
				if elevatorStatus.Behaviour == 1 {
					return "door open"
				}
				if myStates[0].Behaviour == 2 {
					return "moving"
				}
			Floor : elevatorStatus.Floor
			Direction : func() string {
				if elevatorStatus.Dirn == -1 {
					return "down"
				}
				if elevatorStatus.Dirn == 0 {
					return "stop"
				}
				if elevatorStatus.Dirn == 1 {
					return "up"
				}
			}
			CabRequests : elevatorStatus.CabCalls
				
			}
		}
	}


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
}