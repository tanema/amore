package touch

import (
	"github.com/veandco/go-sdl2/sdl"
)

type touchCB func(x, y, dx, dy, pressure float32)

var (
	press_default   touchCB = func(x, y, dx, dy, pressure float32) {}
	release_default touchCB = func(x, y, dx, dy, pressure float32) {}
	move_default    touchCB = func(x, y, dx, dy, pressure float32) {}

	touch_press_cb   = press_default
	touch_release_cb = release_default
	touch_move_cb    = move_default
)

// Delegate is used by amore/event to pass events to the touch package. It may
// also be useful to stub or fake events
func Delegate(event *sdl.TouchFingerEvent) {
	switch event.Type {
	case sdl.FINGERMOTION:
		touches[int64(event.TouchID)] = event
		touch_move_cb(event.X, event.Y, event.DX, event.DY, event.Pressure)
	case sdl.FINGERDOWN:
		touches[int64(event.TouchID)] = event
		touch_press_cb(event.X, event.Y, event.DX, event.DY, event.Pressure)
	case sdl.FINGERUP:
		delete(touches, int64(event.TouchID))
		touch_release_cb(event.X, event.Y, event.DX, event.DY, event.Pressure)
	}
}

// SetTouchPressCB set a callback to be called when there is a press event
func SetTouchPressCB(cb touchCB) {
	if cb == nil {
		touch_press_cb = press_default
	} else {
		touch_press_cb = cb
	}
}

// SetTouchReleaseCB set a callback to be called when there is a release event
func SetTouchReleaseCB(cb touchCB) {
	if cb == nil {
		touch_release_cb = release_default
	} else {
		touch_release_cb = cb
	}
}

// SetTouchMoveCB set a callback to be called when there is a move event
func SetTouchMoveCB(cb touchCB) {
	if cb == nil {
		touch_move_cb = move_default
	} else {
		touch_move_cb = cb
	}
}
