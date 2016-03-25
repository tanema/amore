package joystick

import (
	"github.com/tanema/amore/window/ui"
)

type GameControllerAxis ui.GameControllerAxis
type GameControllerButton ui.GameControllerButton

const (
	AxisInvalid      = GameControllerAxis(ui.CONTROLLER_AXIS_INVALID)
	AxisLeftx        = GameControllerAxis(ui.CONTROLLER_AXIS_LEFTX)
	AxisLefty        = GameControllerAxis(ui.CONTROLLER_AXIS_LEFTY)
	AxisRightx       = GameControllerAxis(ui.CONTROLLER_AXIS_RIGHTX)
	AxisRighty       = GameControllerAxis(ui.CONTROLLER_AXIS_RIGHTY)
	AxisTriggerleft  = GameControllerAxis(ui.CONTROLLER_AXIS_TRIGGERLEFT)
	AxisTriggerright = GameControllerAxis(ui.CONTROLLER_AXIS_TRIGGERRIGHT)
	AxisMax          = GameControllerAxis(ui.CONTROLLER_AXIS_MAX)
)

const (
	ButtonInvalid       = GameControllerButton(ui.CONTROLLER_BUTTON_INVALID)
	ButtonA             = GameControllerButton(ui.CONTROLLER_BUTTON_A)
	ButtonB             = GameControllerButton(ui.CONTROLLER_BUTTON_B)
	ButtonX             = GameControllerButton(ui.CONTROLLER_BUTTON_X)
	ButtonY             = GameControllerButton(ui.CONTROLLER_BUTTON_Y)
	ButtonBack          = GameControllerButton(ui.CONTROLLER_BUTTON_BACK)
	ButtonGuide         = GameControllerButton(ui.CONTROLLER_BUTTON_GUIDE)
	ButtonStart         = GameControllerButton(ui.CONTROLLER_BUTTON_START)
	ButtonLeftstick     = GameControllerButton(ui.CONTROLLER_BUTTON_LEFTSTICK)
	ButtonRightstick    = GameControllerButton(ui.CONTROLLER_BUTTON_RIGHTSTICK)
	ButtonLeftshoulder  = GameControllerButton(ui.CONTROLLER_BUTTON_LEFTSHOULDER)
	ButtonRightshoulder = GameControllerButton(ui.CONTROLLER_BUTTON_RIGHTSHOULDER)
	ButtonDpadUp        = GameControllerButton(ui.CONTROLLER_BUTTON_DPAD_UP)
	ButtonDpadDown      = GameControllerButton(ui.CONTROLLER_BUTTON_DPAD_DOWN)
	ButtonDpadLeft      = GameControllerButton(ui.CONTROLLER_BUTTON_DPAD_LEFT)
	ButtonDpadRight     = GameControllerButton(ui.CONTROLLER_BUTTON_DPAD_RIGHT)
	ButtonMax           = GameControllerButton(ui.CONTROLLER_BUTTON_MAX)
)
