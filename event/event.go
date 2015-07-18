package event

import (
	"github.com/tanema/go-sdl2/sdl"

	"github.com/tanema/amore/joystick"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/mouse"
	"github.com/tanema/amore/window"
)

func Poll() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.CommonEvent:
			println("common event")
		case *sdl.WindowEvent:
			window.GetCurrent().Delegate(e)
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
		default:
			break
		}
	}
}
