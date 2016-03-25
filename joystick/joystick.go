// The joystick Package handles joystick, and gamepad events on the gl context
package joystick

import (
	"github.com/tanema/amore/window/ui"
)

var (
	activeSticks = []*ui.Joystick{}
)

func Init() {
	ui.InitJoyStickAndGamePad()
	for i := 0; i < ui.NumJoysticks(); i++ {
		addJoystick(i)
	}
}

func addJoystick(idx int) *ui.Joystick {
	if idx < 0 || idx >= ui.NumJoysticks() {
		return nil
	}
	joystick := ui.NewJoystick(idx)
	if !joystick.Open() {
		return nil
	}
	activeSticks = append(activeSticks, joystick)
	return joystick
}

func removeJoystick(joystick *ui.Joystick) {
	if joystick == nil {
		return
	}

	for i, stick := range activeSticks {
		if stick == joystick {
			activeSticks = append(activeSticks[:i], activeSticks[i+1:]...)
			break
		}
	}
}

func getDeviceGUID(idx int) string {
	if idx < 0 || idx >= ui.NumJoysticks() {
		return ""
	}
	return ui.GetJoystickName(idx)
}

//Gets the number of connected joysticks.
func GetJoystickCount() int {
	return len(activeSticks)
}

//Gets a list of connected Joysticks.
func GetJoysticks() []*ui.Joystick {
	return activeSticks
}

func getJoystickFromID(id int) *ui.Joystick {
	for _, stick := range activeSticks {
		if stick.GetID() == id {
			return stick
		}
	}
	return nil
}
