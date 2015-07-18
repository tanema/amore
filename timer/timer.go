package timer

import (
	"github.com/tanema/go-sdl2/sdl"
)

const (
	fps_update_frequency = 1
)

var (
	fps                 int
	frames              int
	current_time        float64
	previous_time       float64
	previous_fps_update float64
	dt                  float64
	average_delta       float64
)

func Step() {
	frames++
	previous_time = current_time
	current_time = GetTime()
	dt = current_time - previous_time

	time_since_last := current_time - previous_fps_update
	if time_since_last > fps_update_frequency {
		fps = int((float64(frames) / time_since_last) + 0.5)
		average_delta = time_since_last / float64(frames)
		previous_fps_update = current_time
		frames = 0
	}
}

func GetTime() float64 {
	return float64(sdl.GetTicks()) / 1000.0
}

func GetDelta() float64 {
	return dt
}

func GetFPS() int {
	return fps
}

func GetAverageDelta() float64 {
	return average_delta
}
