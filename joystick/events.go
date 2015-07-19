package joystick

import (
	"github.com/veandco/go-sdl2/sdl"
)

func Delegate(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.JoyAxisEvent:
	case *sdl.JoyBallEvent:
	case *sdl.JoyHatEvent:
	case *sdl.JoyButtonEvent:
		switch e.Type {
		case sdl.JOYBUTTONDOWN:
		case sdl.JOYBUTTONUP:
		}
	case *sdl.JoyDeviceEvent:
		switch e.Type {
		case sdl.JOYDEVICEADDED:
		case sdl.JOYDEVICEREMOVED:
		}
	case *sdl.ControllerAxisEvent:
	case *sdl.ControllerButtonEvent:
		switch e.Type {
		case sdl.CONTROLLERBUTTONDOWN:
		case sdl.CONTROLLERBUTTONUP:
		}
	case *sdl.ControllerDeviceEvent:
		switch e.Type {
		case sdl.CONTROLLERDEVICEADDED:
		case sdl.CONTROLLERDEVICEREMOVED:
		case sdl.CONTROLLERDEVICEREMAPPED:
		}
	case *sdl.TouchFingerEvent:
	}
}
