// Package timer manages game timing by calling step so that the user can get
// FPS, and delta time from this pacakge
package timer

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	fpsUpdateFrequency = 1 //update 1 per second
)

var (
	fps               int     // frames per second
	frames            int     // frames since last update freq
	currentTime       float32 //current frame time
	previousTime      float32 // last frame time
	previousFPSUpdate float32 // last time fps was updated
	dt                float32 // change in time since last step
	averageDelta      float32 // average change in time over update frequency
)

// Step should be called every game loop if rolling your own to keep track of time
// for the update function. It calculates fps and average fps.
func Step() {
	frames++
	previousTime = currentTime
	currentTime = GetTime()
	dt = currentTime - previousTime

	timeSinceLast := currentTime - previousFPSUpdate
	if timeSinceLast > fpsUpdateFrequency {
		fps = int((float32(frames) / timeSinceLast) + 0.5)
		averageDelta = timeSinceLast / float32(frames)
		previousFPSUpdate = currentTime
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
	return averageDelta
}
