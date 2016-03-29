// +build js

package ui

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

type (
	Context struct {
		*js.Object
	}
	Window struct {
		*dom.HTMLCanvasElement
	}
	Cursor               struct{}
	MouseButton          int
	Keycode              string
	Keymod               string
	Scancode             int
	GameControllerAxis   int
	GameControllerButton int
	Joystick             struct {
		id int
	}
	Event     interface{}
	BaseEvent struct {
		*js.Object
		Type      int
		Timestamp int
		WindowID  int
	}
	JoyAxisEvent struct {
		*BaseEvent
		Which int
		Axis  int
		Value int
	}
	JoyBallEvent struct {
		*BaseEvent
		Which int
		Ball  int
		XRel  int
		YRel  int
	}
	JoyButtonEvent struct {
		*BaseEvent
		Which  int
		Button int
		State  int
	}
	JoyHatEvent struct {
		*BaseEvent
		Which int
		Hat   int
		Value int
	}
	ControllerAxisEvent struct {
		*BaseEvent
		Which int
		Axis  int
		Value int
	}
	ControllerButtonEvent struct {
		*BaseEvent
		Which  int
		Button int
		State  int
	}
	JoyDeviceEvent struct {
		*BaseEvent
		Which int
	}
	ControllerDeviceEvent struct {
		*BaseEvent
		Which int
	}
	TouchFingerEvent struct {
		*BaseEvent
		TouchID  int
		FingerID int
		X        float32
		Y        float32
		DX       float32
		DY       float32
		Pressure float32
	}
	Keysym struct {
		Scancode int
		Sym      string
		Mod      uint16
		Unicode  uint32
	}
	KeyDownEvent struct {
		*BaseEvent
		State  int
		Repeat int
		Keysym Keysym
	}
	KeyUpEvent struct {
		*BaseEvent
		State  int
		Repeat int
		Keysym Keysym
	}
	TextEditingEvent struct {
		*BaseEvent
		Text   string
		Start  int
		Length int
	}
	TextInputEvent struct {
		*BaseEvent
		Text string
	}
	MouseMotionEvent struct {
		*BaseEvent
		Which int
		State int
		X     int
		Y     int
		XRel  int
		YRel  int
	}
	MouseButtonEvent struct {
		*BaseEvent
		Which  int
		Button int
		State  int
		X      int
		Y      int
	}
	MouseWheelEvent struct {
		*BaseEvent
		Which int
		X     int
		Y     int
	}
	WindowEvent struct {
		*BaseEvent
		Event int
		Data1 int
		Data2 int
	}
	QuitEvent struct {
		*BaseEvent
	}
	DropEvent struct {
		*BaseEvent
	}
	RenderEvent struct {
		*BaseEvent
	}
	UserEvent struct {
		*BaseEvent
		WindowID int
		Code     int
		Data1    int
		Data2    int
	}
	ClipboardEvent struct {
		*BaseEvent
	}
	OSEvent struct {
		*BaseEvent
	}
	CommonEvent struct {
		*BaseEvent
	}
)
