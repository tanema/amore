package joystick

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Delegate is used by amore/event to pass events to the mouse package. It may
// also be useful to stub or fake events
func Delegate(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.JoyAxisEvent:
	case *sdl.JoyHatEvent:
	case *sdl.ControllerAxisEvent:
	case *sdl.ControllerButtonEvent:
		switch e.Type {
		case sdl.CONTROLLERBUTTONDOWN:
		case sdl.CONTROLLERBUTTONUP:
		}

	case *sdl.JoyDeviceAddedEvent:
		addJoystick(int(e.Which))
	case *sdl.JoyDeviceRemovedEvent:
		removeJoystick(getJoystickFromID(int(e.Which)))
	case *sdl.ControllerDeviceEvent:
		switch e.Type {
		case sdl.CONTROLLERDEVICEADDED:
			addJoystick(int(e.Which))
		case sdl.CONTROLLERDEVICEREMOVED:
			removeJoystick(getJoystickFromID(int(e.Which)))
		case sdl.CONTROLLERDEVICEREMAPPED:
			println("joystick event: controller device remapped")
		}
	}
}
