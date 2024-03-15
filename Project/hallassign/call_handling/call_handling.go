package call

import (
	"Elevator/driver_go_master/elevio"
	fsm "Elevator/elevator/fsm_func"
	"Elevator/elevator/initial"
	"encoding/json"
	"fmt"
	"os"
)

type HRAElevState struct {
	ElevID      string      				   `json:"id"`
	Behaviour   string 						   `json:"behaviour"`
	Floor       int        					   `json:"floor"`
	Direction   string      				   `json:"direction"`
	CabRequests [initial.NFloors]bool          `json:"cabRequests"`
}

type HRAElevStatetemp struct {
	ElevID      string                         `json:"id"`
	Behaviour   initial.ElevatorBehaviour      `json:"behaviour"`
	Floor       int                            `json:"floor"`
	Direction   elevio.MotorDirection          `json:"direction"`
	CabRequests [initial.NFloors]bool          `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [initial.NFloors][2]bool      `json:"hallRequests"`
	States       map[string]HRAElevState       `json:"states"`
}

func UpdateGlobalHallCalls(elevators []initial.Elevator) {
	var n_elevators int = len(elevators)

	for floor := 0; floor < initial.NFloors; floor++ {
		up := elevators[0].Requests[floor][0]
		down := elevators[0].Requests[floor][1]
		for i := 0; i < n_elevators; i++ {
			up = up || elevators[i].Requests[floor][0]
			down = down || elevators[i].Requests[floor][1]
		}

		fsm.GlobalHallCalls[floor] = [2]bool{up, down}
	}
	
}

func ReadFromJSON(fileName string) (map[string][initial.NFloors]bool, error) {
    file, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    var data map[string][initial.NFloors]bool
    if err := decoder.Decode(&data); err != nil {
        return nil, err
    }
    return data, nil
}

func GetCabCalls(elevator initial.Elevator) ([initial.NFloors]bool, error) {
    data, err := ReadFromJSON("CabCallFile.json")
    if err != nil {
        return [initial.NFloors]bool{}, err
    }

    CabCalls, ok := data[elevator.ID]
    if !ok {
        return [initial.NFloors]bool{}, fmt.Errorf("key %s not found in the map", elevator.ID)
    }

    return CabCalls, nil
}


func UpdateCabCalls(Requests [initial.NFloors][initial.NButtons]bool) error {
    _, err := os.Stat("CabCallFile.json")
    var ExistingCabCallMap map[string][initial.NFloors]bool

    if !os.IsNotExist(err) {
        ExistingCabCallMap, err = ReadFromJSON("CabCallFile.json")
        if err != nil {
            return err
        }
    } else {
        ExistingCabCallMap = make(map[string][initial.NFloors]bool)
    }

	var CabCalls = [initial.NFloors]bool {}

	for floor, btn := range Requests{
		CabCalls[floor] = btn[elevio.BTCab]
	}

    ExistingCabCallMap[initial.ElevatorGlob.ID] = CabCalls

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

func GetMyStates(elevators []initial.Elevator) []HRAElevStatetemp {
	var n_elevators int = len(elevators)
	myStates := []HRAElevStatetemp{}
	for i := 0; i < n_elevators; i++ {
		CabCalls,_ := GetCabCalls(elevators[i])
		elevastate := HRAElevStatetemp{elevators[i].ID, elevators[i].Behaviour, elevators[i].Floor, elevators[i].Dirn, CabCalls}
		myStates = append(myStates, elevastate)
	}
	return myStates

}
