package main

import "time"

var (
	timerEndTime float64
	timerActive  bool
)

func get_wall_time() float64 {
	seconds := time.Now().Unix()
	millisec := time.Now().UnixMilli()

	return float64(seconds) + float64(millisec)
}

func timer_start(duration float64) {
	timerEndTime = get_wall_time() + duration
	timerActive = true
}

func timer_stop() {
	timerActive = false
}

func timer_timedOut() bool {
	return (timerActive && get_wall_time() > timerEndTime)
}
