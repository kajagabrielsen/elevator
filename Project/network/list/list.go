package list

import (
	"Elevator/elevator/initialize"
)

var ListOfElevators []initialize.Elevator

//tar inn en liste av heiser og en heis, fjerner heisen fra lista og oppdaterer
func RemoveFromListOfElevators(list []initialize.Elevator, elevator initialize.Elevator) {
	var updatedList []initialize.Elevator

	for _, elev := range list {
		found := false
		if elev == elevator {
			found = true
			break
		}
		if !found {
			updatedList = append(updatedList, elev)
		}
	}
	ListOfElevators = updatedList
}

//tar inn en liste av heiser og en heis, legger til heisen i lista dersom den ikke finnes fra før og ikke er Obstructed
//oppdaterer den i lista dersom den finnes fra før
func AddToListOfElevators(list []initialize.Elevator, elevator initialize.Elevator) {

	flag := false
	for i, elev := range list {
		if elev.ID == elevator.ID {
			list[i] = elevator
			flag = true
		}
	}
	if !flag && !elevator.Obstructed {
		list = append(list, elevator)
	}

	ListOfElevators = list
}