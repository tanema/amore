// Package timer manages game timing by calling step so that the user can get
// FPS, and delta time from this pacakge
package timer

import (
	"time"
)

const (
	fpsUpdateFrequency = 1 //update 1 per second
)

var (
	fps               int       // frames per second
	frames            int       // frames since last update freq
	previousTime      time.Time // last frame time
	previousFPSUpdate time.Time // last time fps was updated
	dt                float32   // change in time since last step
	averageDelta      float32   // average change in time over update frequency
)

// Step should be called every game loop if rolling your own to keep track of time
// for the update function. It calculates fps and average fps.
func Step() {
	frames++
	dt = float32(time.Since(previousTime).Seconds())
	timeSinceLast := float32(time.Since(previousFPSUpdate).Seconds())
	if timeSinceLast > fpsUpdateFrequency {
		fps = int((float32(frames) / timeSinceLast) + 0.5)
		averageDelta = timeSinceLast / float32(frames)
		previousFPSUpdate = previousTime
		frames = 0
	}
	previousTime = time.Now()
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
