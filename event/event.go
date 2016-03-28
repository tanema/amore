package event

import (
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/joystick"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/mouse"
	"github.com/tanema/amore/touch"
	"github.com/tanema/amore/window"
	"github.com/tanema/amore/window/ui"
)

func Poll() {
	for event := ui.PollEvent(); event != nil; event = ui.PollEvent() {
		switch e := event.(type) {
		case *ui.WindowEvent:
			delegateWindowEvent(event, e)
		case *ui.KeyDownEvent, *ui.KeyUpEvent, *ui.TextEditingEvent, *ui.TextInputEvent:
			keyboard.Delegate(e)
		case *ui.MouseMotionEvent, *ui.MouseButtonEvent, *ui.MouseWheelEvent:
			mouse.Delegate(e)
		case *ui.JoyAxisEvent, *ui.JoyBallEvent, *ui.JoyHatEvent,
			*ui.JoyButtonEvent, *ui.JoyDeviceEvent, *ui.ControllerAxisEvent,
			*ui.ControllerButtonEvent, *ui.ControllerDeviceEvent:
			joystick.Delegate(e)
		case *ui.TouchFingerEvent:
			touch.Delegate(e)
		case *ui.QuitEvent:
			window.SetShouldClose(true)
		case *ui.DropEvent, *ui.RenderEvent, *ui.UserEvent,
			*ui.ClipboardEvent, *ui.OSEvent, *ui.CommonEvent:
			//discard not used in amore yet
		}
	}
}

func delegateWindowEvent(event ui.Event, e *ui.WindowEvent) {
	switch e.Type {
	case ui.WINDOWEVENT_NONE:
		return
	case ui.WINDOWEVENT_ENTER, ui.WINDOWEVENT_LEAVE:
		mouse.Delegate(event)
	default:
		switch e.Event {
		case ui.WINDOWEVENT_NONE:
			return
		case ui.WINDOWEVENT_ENTER, ui.WINDOWEVENT_LEAVE:
			mouse.Delegate(event)
		case ui.WINDOWEVENT_SHOWN, ui.WINDOWEVENT_FOCUS_GAINED:
			gfx.SetActive(true)
			ui.DisableScreenSaver()
		case ui.WINDOWEVENT_HIDDEN, ui.WINDOWEVENT_FOCUS_LOST:
			gfx.SetActive(false)
			ui.EnableScreenSaver()
		case ui.WINDOWEVENT_RESIZED, ui.WINDOWEVENT_SIZE_CHANGED:
			w, h := window.GetDrawableSize()
			gfx.SetViewportSize(int32(w), int32(h))
			window.OnSizeChanged(int32(e.Data1), int32(e.Data2))
		case ui.WINDOWEVENT_CLOSE:
			window.SetShouldClose(true)
		}
	}
}
