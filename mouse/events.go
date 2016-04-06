package mouse

import (
	"github.com/veandco/go-sdl2/sdl"
)

var (
	button_press_default   = func(x, y float32, button Button) {}
	button_release_default = func(x, y float32, button Button) {}
	move_default           = func(x, y, dx, dy float32) {}
	focus_default          = func(has_focus bool) {}

	button_press_cb   = button_press_default
	button_release_cb = button_release_default
	move_cb           = move_default
	focus_cb          = focus_default
)

// Delegate is used by amore/event to pass events to the mouse package. It may
// also be useful to stub or fake events
func Delegate(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.MouseMotionEvent:
		move_cb(float32(e.X), float32(e.Y), float32(e.XRel), float32(e.YRel))
	case *sdl.MouseButtonEvent:
		switch e.Type {
		case sdl.MOUSEBUTTONDOWN:
			button_press_cb(float32(e.X), float32(e.Y), Button(e.Button))
		case sdl.MOUSEBUTTONUP:
			button_release_cb(float32(e.X), float32(e.Y), Button(e.Button))
		}
	case *sdl.WindowEvent:
		focus_cb(e.Type == sdl.WINDOWEVENT_ENTER)
	case *sdl.MouseWheelEvent:
	}
}

// SetButtonPressCB will set a callbac to call when the a button on the mouse is
// pressed down.
func SetButtonPressCB(cb func(x, y float32, button Button)) {
	if cb == nil {
		button_press_cb = button_press_default
	} else {
		button_press_cb = cb
	}
}

// SetButtonReleaseCB will set a callback to call when the a button on the mouse is
// released.
func SetButtonReleaseCB(cb func(x, y float32, button Button)) {
	if cb == nil {
		button_release_cb = button_release_default
	} else {
		button_release_cb = cb
	}
}

// SetMoveCB will set a callback to call when the mouse is moved
func SetMoveCB(cb func(x, y, dx, dy float32)) {
	if cb == nil {
		move_cb = move_default
	} else {
		move_cb = cb
	}
}

// SetFocusCB will set a callback to call when the program has mouse focus or loses
// mouse focus.
func SetFocusCB(cb func(has_focus bool)) {
	if cb == nil {
		focus_cb = focus_default
	} else {
		focus_cb = cb
	}
}
