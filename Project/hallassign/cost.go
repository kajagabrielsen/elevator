package hallassign

import (
	"Elevator/utils"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"os"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
	ElevID      string    `json:"id"`
	Behaviour   string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests [utils.N_FLOORS]bool `json:"cabRequests"`
}

type HRAElevStatetemp struct {
	ElevID      string                     `json:"id"`
	Behaviour   utils.ElevatorBehaviour `json:"behaviour"`
	Floor       int                     `json:"floor"`
	Direction   utils.Dirn              `json:"direction"`
	CabRequests [utils.N_FLOORS]bool    `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [utils.N_FLOORS][2]bool `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func CalculateCostFunc(elevators []utils.Elevator) map[string][utils.N_FLOORS][2]bool {

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

	input := HRAInput{
		HallRequests: GetHallCalls(elevators),
		States:       make(map[string]HRAElevState),
	}

	for _, elevatorStatus := range GetMyStates(elevators) {
		input.States[elevatorStatus.ElevID] = HRAElevState{
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
		//return
	}

	//runds the hall_request_assigner file
	ret, err := exec.Command( hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		//return
	}

	//convert the json received from hall_request_assigner to output
	output := make(map[string][utils.N_FLOORS][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		//return
	}

	fmt.Printf("output: \n")
	for k, v := range output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}

	return output
}

func GetHallCalls(elevators []utils.Elevator) [utils.N_FLOORS][2]bool {
	var n_elevators int = len(elevators)
	//fmt.Println(elevators)

	for floor := 0; floor < utils.N_FLOORS; floor++ {
		up := elevators[0].Requests[floor][0]
		down := elevators[0].Requests[floor][1]
		for i := 0; i < n_elevators; i++ {
			up = up || elevators[i].Requests[floor][0]
			down = down || elevators[i].Requests[floor][1]
		}

		utils.GlobalHallCalls[floor] = [2]bool{up, down}
	}
	//fmt.Println(GlobalHallCalls)
	return utils.GlobalHallCalls
}


/**func GetCabCalls(elevator utils.Elevator) [utils.N_FLOORS]bool {
	CabCalls := [utils.N_FLOORS]bool{}
	for floor := 0; floor < utils.N_FLOORS; floor++ {
		CabCalls[floor] = elevator.Requests[floor][2]
	}
	return CabCalls
}**/

func GetCabCalls(elevator utils.Elevator) ([utils.N_FLOORS]bool, error) {
    data, err := ReadFromJSON("CabCallFile.json")
    if err != nil {
        return [utils.N_FLOORS]bool{}, err
    }

    CabCalls, ok := data[elevator.ID]
    if !ok {
        return [utils.N_FLOORS]bool{}, fmt.Errorf("key %s not found in the map", elevator.ID)
    }

    return CabCalls, nil
}


func UpdateCabCalls(Requests [utils.N_FLOORS][utils.N_BUTTONS]bool) error {
    // Check if the file already exists
    _, err := os.Stat("CabCallFile.json")
    var ExistingCabCallMap map[string][utils.N_FLOORS]bool

    // If the file exists, read the existing data from it
    if !os.IsNotExist(err) {
        ExistingCabCallMap, err = ReadFromJSON("CabCallFile.json")
        if err != nil {
            return err
        }
    } else {
        // If the file does not exist, create a new map
        ExistingCabCallMap = make(map[string][utils.N_FLOORS]bool)
    }

    // Update the existing map with the new data
	var CabCalls = [utils.N_FLOORS]bool {}

	for i, floor := range Requests{
		CabCalls[i] = floor[2]
	}

    ExistingCabCallMap[utils.ElevatorGlob.ID] = CabCalls

    // Write the updated map to the JSON file
    file, err := os.OpenFile("CabCallFile.json", os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(ExistingCabCallMap); err != nil {
        return err
    }
    return nil
}


func ReadFromJSON(fileName string) (map[string][utils.N_FLOORS]bool, error) {
    file, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    var data map[string][utils.N_FLOORS]bool
    if err := decoder.Decode(&data); err != nil {
        return nil, err
    }
    return data, nil
}

func GetMyStates(elevators []utils.Elevator) []HRAElevStatetemp {
	var n_elevators int = len(elevators)
	myStates := []HRAElevStatetemp{}
	for i := 0; i < n_elevators; i++ {
		CabCalls,_ := GetCabCalls(elevators[i])
		elevastate := HRAElevStatetemp{elevators[i].ID, elevators[i].Behaviour, elevators[i].Floor, elevators[i].Dirn, CabCalls}
		myStates = append(myStates, elevastate)
	}
	return myStates

}
