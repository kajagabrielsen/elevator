package main

import "os/exec"
import "fmt"
import "encoding/json"
import "runtime"

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
    Behavior    string      `json:"behaviour"`
    Floor       int         `json:"floor"` 
    Direction   string      `json:"direction"`
    CabRequests []bool      `json:"cabRequests"`
}

type HRAInput struct {
    HallRequests    [][2]bool                   `json:"hallRequests"`
    States          map[string]HRAElevState     `json:"states"`
}



func CalculateCostFunc(){

    hraExecutable := ""
    switch runtime.GOOS {
        case "linux":   hraExecutable  = "hall_request_assigner"
        case "windows": hraExecutable  = "hall_request_assigner.exe"
        default:        panic("OS not supported")
    }

    input := HRAInput{
        HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
        States: map[string]HRAElevState{
            "1": HRAElevState{
				if elevator.Behavior == EB_Moving{
					Behavior: "moving"
				}
				if elevator.Behavior == EB_DoorOpen{
					Behavior: "door open"
				}
				if elevator.Behavior == EB_Idle{
					Behavior: "idle"
				},
                Floor:          elevator.Floor,
				if elevator.Dirn == D_Down{
					Direction: "down"
				}
				if elevator.Dirn == D_Stop{
					Direction: "stop"
				}
				if elevator.Dirn == D_Up{
					Direction: "up"
				},
                CabRequests:    []bool{false, false, false, true},
            },
            "2": HRAElevState{
                Behavior:       "idle",
                Floor:          0,
                Direction:      "stop",
                CabRequests:    []bool{false, false, false, false},
            },
        },
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