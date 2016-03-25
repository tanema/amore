package mouse

import (
	"github.com/tanema/amore/window/ui"
)

type MouseButton ui.MouseButton

const (
	LeftButton   MouseButton = MouseButton(ui.BUTTON_LEFT)
	MiddleButton MouseButton = MouseButton(ui.BUTTON_MIDDLE)
	RightButton  MouseButton = MouseButton(ui.BUTTON_RIGHT)
	X1Button     MouseButton = MouseButton(ui.BUTTON_X1)
	X2Button     MouseButton = MouseButton(ui.BUTTON_X2)
)
