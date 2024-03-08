package utils


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
	case D_Up:
		if RequestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{D_Down, EB_DoorOpen}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Down:
		if RequestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{D_Up, EB_DoorOpen}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Stop:
		if RequestsHere(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	default:
		return DirnBehaviourPair{D_Stop, EB_Idle}
	}
}

func RequestsShouldStop(e Elevator) bool {
	switch e.Dirn {
	case D_Down:
		return e.Requests[e.Floor][B_HallDown] ||
			e.Requests[e.Floor][B_Cab] ||
			!RequestsBelow(e)
	case D_Up:
		return e.Requests[e.Floor][B_HallUp] ||
			e.Requests[e.Floor][B_Cab] ||
			!RequestsAbove(e)
	case D_Stop:
	default:
		return true
	}
	return true
}

func RequestsShouldClearImmediately(e Elevator, btnFloor int, btnType Button) bool {
	switch e.ClearRequestVariant {
	case CV_All:
		return e.Floor == btnFloor
	case CV_InDirn:
		return e.Floor == btnFloor &&
			((e.Dirn == D_Up && btnType == B_HallUp) ||
				(e.Dirn == D_Down && btnType == B_HallDown) ||
				e.Dirn == D_Stop ||
				btnType == B_Cab)
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
		e.Requests[e.Floor][B_Cab] = false
		switch e.Dirn {
		case D_Up:
			if !RequestsAbove(e) && !e.Requests[e.Floor][B_HallUp] {
				e.Requests[e.Floor][B_HallDown] = false
			}
			e.Requests[e.Floor][B_HallUp] = false
		case D_Down:
			if !RequestsBelow(e) && !e.Requests[e.Floor][B_HallDown] {
				e.Requests[e.Floor][B_HallUp] = false
			}
			e.Requests[e.Floor][B_HallDown] = false
		case D_Stop:
		default:
			e.Requests[e.Floor][B_HallUp] = false
			e.Requests[e.Floor][B_HallDown] = false
		}
	default:
	}
	return e
}
