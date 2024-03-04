package utils

import "time"

var (
	timerEndTime float64
	timerActive  bool
)

func GetWallTime() float64 {
	seconds := time.Now().Unix()
	millisec := time.Now().UnixMilli()

	return float64(seconds) + float64(millisec)
}

func TimerStart(duration float64) {
	timerEndTime = GetWallTime() + duration
	timerActive = true
}

func TimerStop() {
	timerActive = false
}

func TimerTimedOut() bool {
	return (timerActive && GetWallTime() > timerEndTime)
}
