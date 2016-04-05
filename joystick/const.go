package joystick

import (
	"github.com/veandco/go-sdl2/sdl"
)

type GameControllerAxis sdl.GameControllerAxis
type GameControllerButton sdl.GameControllerButton

const (
	AxisInvalid      = GameControllerAxis(sdl.CONTROLLER_AXIS_INVALID)
	AxisLeftx        = GameControllerAxis(sdl.CONTROLLER_AXIS_LEFTX)
	AxisLefty        = GameControllerAxis(sdl.CONTROLLER_AXIS_LEFTY)
	AxisRightx       = GameControllerAxis(sdl.CONTROLLER_AXIS_RIGHTX)
	AxisRighty       = GameControllerAxis(sdl.CONTROLLER_AXIS_RIGHTY)
	AxisTriggerleft  = GameControllerAxis(sdl.CONTROLLER_AXIS_TRIGGERLEFT)
	AxisTriggerright = GameControllerAxis(sdl.CONTROLLER_AXIS_TRIGGERRIGHT)
	AxisMax          = GameControllerAxis(sdl.CONTROLLER_AXIS_MAX)
)

const (
	ButtonInvalid       = GameControllerButton(sdl.CONTROLLER_BUTTON_INVALID)
	ButtonA             = GameControllerButton(sdl.CONTROLLER_BUTTON_A)
	ButtonB             = GameControllerButton(sdl.CONTROLLER_BUTTON_B)
	ButtonX             = GameControllerButton(sdl.CONTROLLER_BUTTON_X)
	ButtonY             = GameControllerButton(sdl.CONTROLLER_BUTTON_Y)
	ButtonBack          = GameControllerButton(sdl.CONTROLLER_BUTTON_BACK)
	ButtonGuide         = GameControllerButton(sdl.CONTROLLER_BUTTON_GUIDE)
	ButtonStart         = GameControllerButton(sdl.CONTROLLER_BUTTON_START)
	ButtonLeftstick     = GameControllerButton(sdl.CONTROLLER_BUTTON_LEFTSTICK)
	ButtonRightstick    = GameControllerButton(sdl.CONTROLLER_BUTTON_RIGHTSTICK)
	ButtonLeftshoulder  = GameControllerButton(sdl.CONTROLLER_BUTTON_LEFTSHOULDER)
	ButtonRightshoulder = GameControllerButton(sdl.CONTROLLER_BUTTON_RIGHTSHOULDER)
	ButtonDpadUp        = GameControllerButton(sdl.CONTROLLER_BUTTON_DPAD_UP)
	ButtonDpadDown      = GameControllerButton(sdl.CONTROLLER_BUTTON_DPAD_DOWN)
	ButtonDpadLeft      = GameControllerButton(sdl.CONTROLLER_BUTTON_DPAD_LEFT)
	ButtonDpadRight     = GameControllerButton(sdl.CONTROLLER_BUTTON_DPAD_RIGHT)
	ButtonMax           = GameControllerButton(sdl.CONTROLLER_BUTTON_MAX)
)

type HatDirection int

const (
	HatInvalid  HatDirection = -1
	HatCentered HatDirection = iota
	HatUp
	HatRight
	HatDown
	HatLeft
	HatRightUp
	HatRightDown
	HatLeftup
	HatLeftdown
)
