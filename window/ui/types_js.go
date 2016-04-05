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
		focused           bool
		grabbed           bool
		x                 int
		y                 int
		width             int
		height            int
		minWidth          int
		minHeight         int
		windowListeners   map[string]func(*js.Object)
		documentListeners map[string]func(*js.Object)
		canvasListeners   map[string]func(*js.Object)
	}
	Cursor struct {
		icon string
		hx   int
		hy   int
	}
	MouseButton          int
	Keycode              string
	Keymod               string
	Scancode             int
	GameControllerAxis   int
	GameControllerButton int
	Joystick             struct {
		*js.Object
		id   int
		name string
	}
	Event        interface{}
	JoyAxisEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Which     int
		Axis      int
		Value     int
	}
	JoyBallEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Which     int
		Ball      int
		XRel      int
		YRel      int
	}
	JoyButtonEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Which     int
		Button    int
		State     int
	}
	JoyHatEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Which     int
		Hat       int
		Value     int
	}
	ControllerAxisEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Which     int
		Axis      int
		Value     int
	}
	ControllerButtonEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Which     int
		Button    int
		State     int
	}
	JoyDeviceEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Which     int
	}
	ControllerDeviceEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Which     int
	}
	TouchFingerEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		TouchID   int
		FingerID  int
		X         float32
		Y         float32
		DX        float32
		DY        float32
		Pressure  float32
	}
	Keysym struct {
		Scancode int
		Sym      string
		Mod      uint16
		Unicode  uint32
	}
	KeyDownEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		State     int
		Repeat    int
		Keysym    Keysym
	}
	KeyUpEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		State     int
		Repeat    int
		Keysym    Keysym
	}
	TextEditingEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Text      string
		Start     int
		Length    int
	}
	TextInputEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Text      string
	}
	MouseMotionEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Which     int
		State     int
		X         int
		Y         int
		XRel      int
		YRel      int
	}
	MouseButtonEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Which     int
		Button    int
		State     int
		X         int
		Y         int
	}
	MouseWheelEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Which     int
		X         int
		Y         int
	}
	WindowEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Event     int
		Data1     int
		Data2     int
	}
	QuitEvent struct {
		Type      int
		Timestamp int
		WindowID  int
	}
	DropEvent struct {
		Type      int
		Timestamp int
		WindowID  int
	}
	RenderEvent struct {
		Type      int
		Timestamp int
		WindowID  int
	}
	UserEvent struct {
		Type      int
		Timestamp int
		WindowID  int
		Code      int
		Data1     int
		Data2     int
	}
	ClipboardEvent struct {
		Type      int
		Timestamp int
		WindowID  int
	}
	OSEvent struct {
		Type      int
		Timestamp int
		WindowID  int
	}
	CommonEvent struct {
		Type      int
		Timestamp int
		WindowID  int
	}
)
