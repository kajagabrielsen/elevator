package utils

import 	"Elevator/driver_go_master/elevio"

	
func RequestsAbove(e Elevator) bool {
	for f := e.Floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func RequestsBelow(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func RequestsHere(e Elevator) bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

func RequestsChooseDirection(e Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case elevio.MDUp:
		if RequestsAbove(e) {
			return DirnBehaviourPair{elevio.MDUp, EB_Moving}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{elevio.MDDown, EB_DoorOpen}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{elevio.MDDown, EB_Moving}
		} else {
			return DirnBehaviourPair{elevio.MDStop, EB_Idle}
		}
	case elevio.MDDown:
		if RequestsBelow(e) {
			return DirnBehaviourPair{elevio.MDDown, EB_Moving}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{elevio.MDUp, EB_DoorOpen}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{elevio.MDUp, EB_Moving}
		} else {
			return DirnBehaviourPair{elevio.MDStop, EB_Idle}
		}
	case elevio.MDStop:
		if RequestsHere(e) {
			return DirnBehaviourPair{elevio.MDStop, EB_DoorOpen}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{elevio.MDUp, EB_Moving}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{elevio.MDDown, EB_Moving}
		} else {
			return DirnBehaviourPair{elevio.MDStop, EB_Idle}
		}
	default:
		return DirnBehaviourPair{elevio.MDStop, EB_Idle}
	}
}

func RequestsShouldStop(e Elevator) bool {
	switch e.Dirn {
	case elevio.MDDown:
		return e.Requests[e.Floor][elevio.BTHallDown] ||
			e.Requests[e.Floor][elevio.BTCab] ||
			!RequestsBelow(e)
	case elevio.MDUp:
		return e.Requests[e.Floor][elevio.BTHallUp] ||
			e.Requests[e.Floor][elevio.BTCab] ||
			!RequestsAbove(e)
	case elevio.MDStop:
	default:
		return true
	}
	return true
}

func RequestsShouldClearImmediately(e Elevator, btnFloor int, btnType elevio.ButtonType) bool {
	switch e.ClearRequestVariant {
	case CV_All:
		return e.Floor == btnFloor
	case CV_InDirn:
		return e.Floor == btnFloor &&
			((e.Dirn == elevio.MDUp && btnType == elevio.BTHallUp) ||
				(e.Dirn == elevio.MDDown && btnType == elevio.BTHallDown) ||
				e.Dirn == elevio.MDStop ||
				btnType == elevio.BTCab)
	default:
		return false
	}
}

func RequestsClearAtCurrentFloor(e Elevator) Elevator {
	switch e.ClearRequestVariant {
	case CV_All:
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.Requests[e.Floor][btn] = false
		}
	case CV_InDirn:
		e.Requests[e.Floor][elevio.BTCab] = false
		switch e.Dirn {
		case elevio.MDUp:
			if !RequestsAbove(e) && !e.Requests[e.Floor][elevio.BTHallUp] {
				e.Requests[e.Floor][elevio.BTHallDown] = false
			}
			e.Requests[e.Floor][elevio.BTHallUp] = false
		case elevio.MDDown:
			if !RequestsBelow(e) && !e.Requests[e.Floor][elevio.BTHallDown] {
				e.Requests[e.Floor][elevio.BTHallUp] = false
			}
			e.Requests[e.Floor][elevio.BTHallDown] = false
		case elevio.MDStop:
		default:
			e.Requests[e.Floor][elevio.BTHallUp] = false
			e.Requests[e.Floor][elevio.BTHallDown] = false
		}
	default:
	}
	return e
}
