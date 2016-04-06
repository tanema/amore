// The timer Package manages game timing by calling step so that the user can get
// FPS, and delta time from this pacakge
package timer

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	fps_update_frequency = 1 //update 1 per second
)

var (
	fps                 int     // frames per second
	frames              int     // frames since last update freq
	current_time        float32 //current frame time
	previous_time       float32 // last frame time
	previous_fps_update float32 // last time fps was updated
	dt                  float32 // change in time since last step
	average_delta       float32 // average change in time over update frequency
)

// Step should be called every game loop if rolling your own to keep track of time
// for the update function. It calculates fps and average fps.
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

// GetTime get the current time from the length that the application has been running.
func GetTime() float32 {
	return float32(sdl.GetTicks()) / 1000.0
}

// GetDelta returns the difference of time between the current frame and the last frame.
func GetDelta() float32 {
	return dt
}

// GetFPS returns the number of frames per second.
func GetFPS() int {
	return fps
}

// GetAverageDelta returns the average time of each frame
func GetAverageDelta() float32 {
	return average_delta
}
