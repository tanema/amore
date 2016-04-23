package mouse

import (
	"github.com/veandco/go-sdl2/sdl"
)

var (
	// OnButtonDown is called when the a button on the mouse is pressed down.
	OnButtonDown func(x, y float32, button Button)
	// OnButtonUp is called when the a button on the mouse is released.
	OnButtonUp func(x, y float32, button Button)
	// OnMove is called when the mouse is moved.
	OnMove func(x, y, dx, dy float32)
	// OnFocus is called when the program has mouse focus or loses mouse focus.
	OnFocus func(has_focus bool)
	// OnWheelMove is called when the mouse wheel is changed
	OnWheelMove func(x, y float32)
)

// Delegate is used by amore/event to pass events to the mouse package. It may
// also be useful to stub or fake events
func Delegate(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.MouseMotionEvent:
		if OnMove != nil {
			OnMove(float32(e.X), float32(e.Y), float32(e.XRel), float32(e.YRel))
		}
	case *sdl.MouseButtonEvent:
		switch e.Type {
		case sdl.MOUSEBUTTONDOWN:
			if OnButtonDown != nil {
				OnButtonDown(float32(e.X), float32(e.Y), Button(e.Button))
			}
		case sdl.MOUSEBUTTONUP:
			if OnButtonUp != nil {
				OnButtonUp(float32(e.X), float32(e.Y), Button(e.Button))
			}
		}
	case *sdl.WindowEvent:
		if OnFocus != nil {
			OnFocus(e.Type == sdl.WINDOWEVENT_ENTER)
		}
	case *sdl.MouseWheelEvent:
		if OnWheelMove != nil {
			OnWheelMove(float32(e.X), float32(e.Y))
		}
	}
}
