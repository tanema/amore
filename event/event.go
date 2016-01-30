package event

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/joystick"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/mouse"
	"github.com/tanema/amore/window"
)

func Poll() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.WindowEvent:
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
					w, h := window.GetCurrent().GetDrawableSize()
					gfx.SetViewportSize(w, h)
				case sdl.WINDOWEVENT_CLOSE:
					window.GetCurrent().SetShouldClose(true)
				}
				window.GetCurrent().Delegate(e)
			}
		case *sdl.KeyDownEvent, *sdl.KeyUpEvent, *sdl.TextEditingEvent, *sdl.TextInputEvent:
			go keyboard.Delegate(e)
		case *sdl.MouseMotionEvent, *sdl.MouseButtonEvent, *sdl.MouseWheelEvent:
			go mouse.Delegate(e)
		case *sdl.JoyAxisEvent, *sdl.JoyBallEvent, *sdl.JoyHatEvent,
			*sdl.JoyButtonEvent, *sdl.JoyDeviceEvent, *sdl.ControllerAxisEvent,
			*sdl.ControllerButtonEvent, *sdl.ControllerDeviceEvent, *sdl.TouchFingerEvent:
			go joystick.Delegate(e)
		case *sdl.DropEvent:
			println("drop event")
		case *sdl.RenderEvent:
			println("render event")
		case *sdl.UserEvent:
			println("user event")
		case *sdl.ClipboardEvent:
			println("clip event")
		case *sdl.OSEvent:
			println("OS event")
		case *sdl.QuitEvent:
			window.GetCurrent().SetShouldClose(true)
		case *sdl.CommonEvent:
			println("common")
		}
	}
}
