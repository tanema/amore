package mouse

import (
	"github.com/tanema/amore/window/ui"
)

type ButtonPressCB func(x, y float32, button MouseButton)
type ButtonReleaseCB func(x, y float32, button MouseButton)
type MoveCB func(x, y, dx, dy float32)
type FocusCb func(has_focus bool)

var (
	button_press_default   ButtonPressCB   = func(x, y float32, button MouseButton) {}
	button_release_default ButtonReleaseCB = func(x, y float32, button MouseButton) {}
	move_default           MoveCB          = func(x, y, dx, dy float32) {}
	focus_default          FocusCb         = func(has_focus bool) {}

	button_press_cb   = button_press_default
	button_release_cb = button_release_default
	move_cb           = move_default
	focus_cb          = focus_default
)

func Delegate(event ui.Event) {
	switch e := event.(type) {
	case *ui.MouseMotionEvent:
		move_cb(float32(e.X), float32(e.Y), float32(e.XRel), float32(e.YRel))
	case *ui.MouseButtonEvent:
		switch e.Type {
		case ui.MOUSEBUTTONDOWN:
			button_press_cb(float32(e.X), float32(e.Y), MouseButton(e.Button))
		case ui.MOUSEBUTTONUP:
			button_release_cb(float32(e.X), float32(e.Y), MouseButton(e.Button))
		}
	case *ui.WindowEvent:
		focus_cb(e.Type == ui.WINDOWEVENT_ENTER)
	case *ui.MouseWheelEvent:
	}
}

func SetButtonPressCB(cb ButtonPressCB) {
	if cb == nil {
		button_press_cb = button_press_default
	} else {
		button_press_cb = cb
	}
}

func SetButtonReleaseCB(cb ButtonReleaseCB) {
	if cb == nil {
		button_release_cb = button_release_default
	} else {
		button_release_cb = cb
	}
}

func SetMoveCB(cb MoveCB) {
	if cb == nil {
		move_cb = move_default
	} else {
		move_cb = cb
	}
}

func SetFocusCB(cb FocusCb) {
	if cb == nil {
		focus_cb = focus_default
	} else {
		focus_cb = cb
	}
}
