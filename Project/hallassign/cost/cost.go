package cost

import (
	fsm "Elevator/elevator/fsm_func"
	"Elevator/elevator/initial"
	call "Elevator/hallassign/call_handling"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

func CalculateCostFunc(elevators []initial.Elevator) map[string][initial.NFloors][2]bool {

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

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
	}

	ret, err := exec.Command( hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
	}

	output := make(map[string][initial.NFloors][2]bool)
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
