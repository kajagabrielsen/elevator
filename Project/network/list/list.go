package list

import (
	"Elevator/elevator/initial"
)

var ListOfElevators []initial.Elevator

func RemoveFromListOfElevators(list []initial.Elevator, id string) {
	var updatedList []initial.Elevator
	for _, elev := range list {
		found := false
		if elev.ID == id {
			found = true
			break
		}
		if !found {
			updatedList = append(updatedList, elev)
		}
	}
	ListOfElevators = updatedList
}

func AddToListOfElevators(list []initial.Elevator, elevator initial.Elevator) {
	found := false
	for i, elev := range list {
		if elev.ID == elevator.ID {
			list[i] = elevator
			found = true
		}
	}
	if !found && !elevator.Obstructed {
		list = append(list, elevator)
	}
	ListOfElevators = list
}