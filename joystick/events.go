package joystick

import (
	"github.com/tanema/amore/window/ui"
)

func Delegate(event ui.Event) {
	switch e := event.(type) {
	case *ui.JoyAxisEvent:
	case *ui.JoyHatEvent:
	case *ui.ControllerAxisEvent:
	case *ui.ControllerButtonEvent:
		switch e.Type {
		case ui.CONTROLLERBUTTONDOWN:
		case ui.CONTROLLERBUTTONUP:
		}
	case *ui.JoyDeviceEvent:
		switch e.Type {
		case ui.JOYDEVICEADDED:
			addJoystick(int(e.Which))
		case ui.JOYDEVICEREMOVED:
			removeJoystick(getJoystickFromID(int(e.Which)))
		}
	case *ui.ControllerDeviceEvent:
		switch e.Type {
		case ui.CONTROLLERDEVICEADDED:
			addJoystick(int(e.Which))
		case ui.CONTROLLERDEVICEREMOVED:
			removeJoystick(getJoystickFromID(int(e.Which)))
		case ui.CONTROLLERDEVICEREMAPPED:
			println("joystick event: controller device remapped")
		}
	}
}
