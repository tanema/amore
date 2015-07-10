package joystick

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type JoystickId int

const (
	JoyStick1    JoystickId = JoystickId(glfw.Joystick1)
	JoyStick2    JoystickId = JoystickId(glfw.Joystick2)
	JoyStick3    JoystickId = JoystickId(glfw.Joystick3)
	JoyStick4    JoystickId = JoystickId(glfw.Joystick4)
	JoyStick5    JoystickId = JoystickId(glfw.Joystick5)
	JoyStick6    JoystickId = JoystickId(glfw.Joystick6)
	JoyStick7    JoystickId = JoystickId(glfw.Joystick7)
	JoyStick8    JoystickId = JoystickId(glfw.Joystick8)
	JoyStick9    JoystickId = JoystickId(glfw.Joystick9)
	JoyStick10   JoystickId = JoystickId(glfw.Joystick10)
	JoyStick11   JoystickId = JoystickId(glfw.Joystick11)
	JoyStick12   JoystickId = JoystickId(glfw.Joystick12)
	JoyStick13   JoystickId = JoystickId(glfw.Joystick13)
	JoyStick14   JoystickId = JoystickId(glfw.Joystick14)
	JoyStick15   JoystickId = JoystickId(glfw.Joystick15)
	JoyStick16   JoystickId = JoystickId(glfw.Joystick16)
	JoyStickLast            = JoyStick16
)
