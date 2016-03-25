package timer

import (
	"time"
)

const (
	fps_update_frequency = 1
)

var (
	fps                 int
	frames              int
	current_time        float32
	previous_time       float32
	previous_fps_update float32
	dt                  float32
	average_delta       float32
	ticks               int32
)

func init() {
	ticker := time.NewTicker(time.Millisecond)
	go func() {
		for range ticker.C {
			ticks++
		}
	}()
}

func Step() {
	frames++
	previous_time = current_time
	current_time = GetTime()
	dt = current_time - previous_time

	time_since_last := current_time - previous_fps_update
	if time_since_last > fps_update_frequency {
		fps = int((float32(frames) / time_since_last) + 0.5)
		average_delta = time_since_last / float32(frames)
		previous_fps_update = current_time
		frames = 0
	}
}

func GetTime() float32 {
	return float32(ticks) / 1000.0
}

func GetDelta() float32 {
	return dt
}

func GetFPS() int {
	return fps
}

func GetAverageDelta() float32 {
	return average_delta
}
