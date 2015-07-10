package joystick

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

//Gets the number of connected joysticks.
func GetJoystickCount() int {
	return len(GetJoysticks())
}

//Gets a list of connected Joysticks.
func GetJoysticks() []*Joystick {
	joysticks := []*Joystick{}
	for i := Joystick1; i <= JoystickLast; i++ {
		if glfw.JoystickPresent(glfw.Joystick(i)) {
			append(joysticks, &Joystick{id: i})
		}
	}
}

func Refresh() {
	joysticks := GetJoysticks()
	for _, joystick := range joysticks {
		joystick.Refresh()
	}
}
