package call

import (
	"Elevator/driver_go_master/elevio"
	"Elevator/elevator/initialize"
	"Elevator/elevator/fsm_func"
	"encoding/json"
	"fmt"
	"os"
)

type HRAElevState struct {
	ElevID      string      				   `json:"id"`
	Behaviour   string 						   `json:"behaviour"`
	Floor       int        					   `json:"floor"`
	Direction   string      				   `json:"direction"`
	CabRequests [initialize.N_FLOORS]bool      `json:"cabRequests"`
}

type HRAElevStatetemp struct {
	ElevID      string                         `json:"id"`
	Behaviour   initialize.ElevatorBehaviour   `json:"behaviour"`
	Floor       int                            `json:"floor"`
	Direction   elevio.MotorDirection          `json:"direction"`
	CabRequests [initialize.N_FLOORS]bool      `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [initialize.N_FLOORS][2]bool  `json:"hallRequests"`
	States       map[string]HRAElevState       `json:"states"`
}

func UpdateGlobalHallCalls(elevators []initialize.Elevator) {
	var n_elevators int = len(elevators)

	for floor := 0; floor < initialize.N_FLOORS; floor++ {
		up := elevators[0].Requests[floor][0]
		down := elevators[0].Requests[floor][1]
		for i := 0; i < n_elevators; i++ {
			up = up || elevators[i].Requests[floor][0]
			down = down || elevators[i].Requests[floor][1]
		}

		fsm.GlobalHallCalls[floor] = [2]bool{up, down}
	}
	
}


func GetCabCalls(elevator initialize.Elevator) ([initialize.N_FLOORS]bool, error) {
    data, err := ReadFromJSON("CabCallFile.json")
    if err != nil {
        return [initialize.N_FLOORS]bool{}, err
    }

    CabCalls, ok := data[elevator.ID]
    if !ok {
        return [initialize.N_FLOORS]bool{}, fmt.Errorf("key %s not found in the map", elevator.ID)
    }

    return CabCalls, nil
}


func UpdateCabCalls(Requests [initialize.N_FLOORS][initialize.N_BUTTONS]bool) error {
    // Check if the file already exists
    _, err := os.Stat("CabCallFile.json")
    var ExistingCabCallMap map[string][initialize.N_FLOORS]bool

    // If the file exists, read the existing data from it
    if !os.IsNotExist(err) {
        ExistingCabCallMap, err = ReadFromJSON("CabCallFile.json")
        if err != nil {
            return err
        }
    } else {
        // If the file does not exist, create a new map
        ExistingCabCallMap = make(map[string][initialize.N_FLOORS]bool)
    }

    // Update the existing map with the new data
	var CabCalls = [initialize.N_FLOORS]bool {}

	for i, floor := range Requests{
		CabCalls[i] = floor[2]
	}

    ExistingCabCallMap[initialize.ElevatorGlob.ID] = CabCalls

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


func ReadFromJSON(fileName string) (map[string][initialize.N_FLOORS]bool, error) {
    file, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    var data map[string][initialize.N_FLOORS]bool
    if err := decoder.Decode(&data); err != nil {
        return nil, err
    }
    return data, nil
}

func GetMyStates(elevators []initialize.Elevator) []HRAElevStatetemp {
	var n_elevators int = len(elevators)
	myStates := []HRAElevStatetemp{}
	for i := 0; i < n_elevators; i++ {
		CabCalls,_ := GetCabCalls(elevators[i])
		elevastate := HRAElevStatetemp{elevators[i].ID, elevators[i].Behaviour, elevators[i].Floor, elevators[i].Dirn, CabCalls}
		myStates = append(myStates, elevastate)
	}
	return myStates

}
