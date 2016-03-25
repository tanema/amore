// +build !js

package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type (
	Context sdl.GLContext
	Window  struct {
		*sdl.Window
	}
	Cursor               *sdl.Cursor
	MouseButton          uint32
	Keycode              sdl.Keycode
	Keymod               sdl.Keymod
	Scancode             uint32
	GameControllerAxis   sdl.GameControllerAxis
	GameControllerButton sdl.GameControllerButton
	HapticEffect         sdl.HapticEffect
	Haptic               sdl.Haptic
	Vibration            struct {
		Left, Right float32
		Effect      sdl.HapticEffect
		Data        [4]uint16
		ID          int
		Endtime     uint32
	}
	Joystick struct {
		id         int
		stick      *sdl.Joystick
		controller *sdl.GameController
		haptic     *sdl.Haptic
		vibration  *Vibration
	}

	Event                 sdl.Event
	JoyAxisEvent          sdl.JoyAxisEvent
	JoyBallEvent          sdl.JoyBallEvent
	JoyButtonEvent        sdl.JoyButtonEvent
	JoyHatEvent           sdl.JoyHatEvent
	ControllerAxisEvent   sdl.ControllerAxisEvent
	ControllerButtonEvent sdl.ControllerButtonEvent
	JoyDeviceEvent        sdl.JoyDeviceEvent
	ControllerDeviceEvent sdl.ControllerDeviceEvent
	TouchFingerEvent      sdl.TouchFingerEvent
	KeyDownEvent          sdl.KeyDownEvent
	KeyUpEvent            sdl.KeyUpEvent
	TextEditingEvent      sdl.TextEditingEvent
	TextInputEvent        sdl.TextInputEvent
	MouseMotionEvent      sdl.MouseMotionEvent
	MouseButtonEvent      sdl.MouseButtonEvent
	MouseWheelEvent       sdl.MouseWheelEvent
	WindowEvent           sdl.WindowEvent
	QuitEvent             sdl.QuitEvent
	DropEvent             sdl.DropEvent
	RenderEvent           sdl.RenderEvent
	UserEvent             sdl.UserEvent
	ClipboardEvent        sdl.ClipboardEvent
	OSEvent               sdl.OSEvent
	CommonEvent           sdl.CommonEvent
)
