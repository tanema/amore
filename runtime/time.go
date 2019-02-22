package runtime

import (
	"time"
)

var (
	fps               int       // frames per second
	frames            int       // frames since last update freq
	previousTime      time.Time // last frame time
	previousFPSUpdate time.Time // last time fps was updated
)

func step() float32 {
	frames++
	dt := float32(time.Since(previousTime).Seconds())
	previousTime = time.Now()
	timeSinceLast := float32(time.Since(previousFPSUpdate).Seconds())
	if timeSinceLast > 1 { //update 1 per second
		fps = int((float32(frames) / timeSinceLast) + 0.5)
		previousFPSUpdate = previousTime
		frames = 0
	}
	return dt
}
