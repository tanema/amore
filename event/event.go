// The event pacakge manages and delegates all the events in the gl context to
// thier respective handlers.
package event

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/joystick"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/mouse"
	"github.com/tanema/amore/touch"
	"github.com/tanema/amore/window"
)

// Poll is used by the game loop to gather events and delegate to each pacakge.
// Generally you should not have to use this method however if you are doing your
// own game loop this should be called at the end.
func Poll() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.WindowEvent:
			delegateWindowEvent(event, e)
		case *sdl.KeyDownEvent, *sdl.KeyUpEvent, *sdl.TextEditingEvent, *sdl.TextInputEvent:
			keyboard.Delegate(e)
		case *sdl.MouseMotionEvent, *sdl.MouseButtonEvent, *sdl.MouseWheelEvent:
			mouse.Delegate(e)
		case *sdl.JoyAxisEvent, *sdl.JoyBallEvent, *sdl.JoyHatEvent,
			*sdl.JoyButtonEvent, *sdl.JoyDeviceEvent, *sdl.ControllerAxisEvent,
			*sdl.ControllerButtonEvent, *sdl.ControllerDeviceEvent:
			joystick.Delegate(e)
		case *sdl.TouchFingerEvent:
			touch.Delegate(e)
		case *sdl.QuitEvent:
			window.Close(true)
		case *sdl.DropEvent, *sdl.RenderEvent, *sdl.UserEvent,
			*sdl.ClipboardEvent, *sdl.OSEvent, *sdl.CommonEvent:
			//discard not used in amore yet
		}
	}
}

// delegateWindowEvent handles window events and delegates them to the pacakges
// that handle those events.
func delegateWindowEvent(event sdl.Event, e *sdl.WindowEvent) {
	switch e.Type {
	case sdl.WINDOWEVENT_NONE:
		return
	case sdl.WINDOWEVENT_ENTER, sdl.WINDOWEVENT_LEAVE:
		mouse.Delegate(event)
	default:
		switch e.Event {
		case sdl.WINDOWEVENT_NONE:
			return
		case sdl.WINDOWEVENT_ENTER, sdl.WINDOWEVENT_LEAVE:
			mouse.Delegate(event)
		case sdl.WINDOWEVENT_SHOWN, sdl.WINDOWEVENT_FOCUS_GAINED:
			gfx.SetActive(true)
		case sdl.WINDOWEVENT_HIDDEN, sdl.WINDOWEVENT_FOCUS_LOST:
			gfx.SetActive(false)
		case sdl.WINDOWEVENT_RESIZED, sdl.WINDOWEVENT_SIZE_CHANGED:
			w, h := window.GetDrawableSize()
			gfx.SetViewportSize(w, h)
		case sdl.WINDOWEVENT_CLOSE:
			window.Close(true)
		}
		window.Delegate(e)
	}
}
