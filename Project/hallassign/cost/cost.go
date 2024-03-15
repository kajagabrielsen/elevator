package cost

import (
	"Elevator/elevator/initialize"
	"Elevator/elevator/fsm_func"
	"Elevator/hallassign/call_handling"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

func CalculateCostFunc(elevators []initialize.Elevator) map[string][initialize.N_FLOORS][2]bool {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "./hallassign/hall_request_assigner.exe"
	case "darwin":
		hraExecutable = "./hallassign/hall_request_assigner_mac"
	default:
		panic("OS not supported")
	}

	input := call.HRAInput{
		HallRequests: fsm.GlobalHallCalls,
		States:       make(map[string]call.HRAElevState),
	}

	for _, elevatorStatus := range call.GetMyStates(elevators) {
		input.States[elevatorStatus.ElevID] = call.HRAElevState{
			Behaviour: func() string {
				if elevatorStatus.Behaviour == 0 {
					return "idle"
				} else if elevatorStatus.Behaviour == 1 {
					return "doorOpen"
				} else {
					return "moving"
				}
			}(),
			Floor: elevatorStatus.Floor,
			Direction: func() string {
				if elevatorStatus.Direction == -1 {
					return "down"
				} else if elevatorStatus.Direction == 0 {
					return "stop"
				} else {
					return "up"
				}
			}(),
			CabRequests: elevatorStatus.CabRequests,
		}
	}

	//Convert input to json format
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
	}

	//runds the hall_request_assigner file
	ret, err := exec.Command( hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
	}

	//convert the json received from hall_request_assigner to output
	output := make(map[string][initialize.N_FLOORS][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
	}

	fmt.Printf("output: \n")
	for k, v := range output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}

	return output
}
