package utils

import "time"

var (
	timerEndTime float64
	timerActive  bool
)

func Get_wall_time() float64 {
	seconds := time.Now().Unix()
	millisec := time.Now().UnixMilli()

	return float64(seconds) + float64(millisec)
}

func Timer_start(duration float64) {
	timerEndTime = Get_wall_time() + duration
	timerActive = true
}

func Timer_stop() {
	timerActive = false
}

func Timer_timedOut() bool {
	return (timerActive && Get_wall_time() > timerEndTime)
}
