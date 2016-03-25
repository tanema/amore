package touch

import (
	"github.com/tanema/amore/window/ui"
)

type (
	touchCB func(x, y, dx, dy, pressure float32)
	Touch   struct {
		event *ui.TouchFingerEvent
	}
)

var (
	touches = make(map[int64]*ui.TouchFingerEvent)

	press_default   touchCB = func(x, y, dx, dy, pressure float32) {}
	release_default touchCB = func(x, y, dx, dy, pressure float32) {}
	move_default    touchCB = func(x, y, dx, dy, pressure float32) {}

	touch_press_cb   = press_default
	touch_release_cb = release_default
	touch_move_cb    = move_default
)

func Delegate(event *ui.TouchFingerEvent) {
	switch event.Type {
	case ui.FINGERMOTION:
		touches[int64(event.TouchID)] = event
		touch_move_cb(event.X, event.Y, event.DX, event.DY, event.Pressure)
	case ui.FINGERDOWN:
		touches[int64(event.TouchID)] = event
		touch_press_cb(event.X, event.Y, event.DX, event.DY, event.Pressure)
	case ui.FINGERUP:
		delete(touches, int64(event.TouchID))
		touch_release_cb(event.X, event.Y, event.DX, event.DY, event.Pressure)
	}
}

func SetTouchPressCB(cb touchCB) {
	if cb == nil {
		touch_press_cb = press_default
	} else {
		touch_press_cb = cb
	}
}

func SetTouchReleaseCB(cb touchCB) {
	if cb == nil {
		touch_release_cb = release_default
	} else {
		touch_release_cb = cb
	}
}

func SetTouchMoveCB(cb touchCB) {
	if cb == nil {
		touch_move_cb = move_default
	} else {
		touch_move_cb = cb
	}
}

func GetTouches() []Touch {
	fingers := []Touch{}
	for _, touch := range touches {
		fingers = append(fingers, Touch{event: touch})
	}
	return fingers
}

func (touch *Touch) GetPosition() (float32, float32) {
	return touch.event.X, touch.event.Y
}

func (touch *Touch) GetPressure(id int64) float32 {
	return touch.event.Pressure
}
