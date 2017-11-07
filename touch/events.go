package touch

import (
	"github.com/veandco/go-sdl2/sdl"
)

type touchCB func(x, y, dx, dy, pressure float32)

var (
	// OnTouchPress is a callback that will be called when the press starts.
	OnTouchPress touchCB
	// OnTouchRelease is a callback that will be called when the press ends.
	OnTouchRelease touchCB
	// OnTouchMove is a callback that will be called when the press moves position.
	OnTouchMove touchCB
)

// Delegate is used by amore/event to pass events to the touch package. It may
// also be useful to stub or fake events
func Delegate(event *sdl.TouchFingerEvent) {
	switch event.Type {
	case sdl.FINGERMOTION:
		touches[int64(event.TouchID)] = event
		if OnTouchMove != nil {
			OnTouchMove(event.X, event.Y, event.DX, event.DY, event.Pressure)
		}
	case sdl.FINGERDOWN:
		touches[int64(event.TouchID)] = event
		if OnTouchPress != nil {
			OnTouchPress(event.X, event.Y, event.DX, event.DY, event.Pressure)
		}
	case sdl.FINGERUP:
		delete(touches, int64(event.TouchID))
		if OnTouchRelease != nil {
			OnTouchRelease(event.X, event.Y, event.DX, event.DY, event.Pressure)
		}
	}
}
