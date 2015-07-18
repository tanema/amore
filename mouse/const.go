package mouse

import (
	"github.com/tanema/go-sdl2/sdl"
)

type Button uint32

const (
	LeftButton   Button = sdl.BUTTON_LEFT
	MiddleButton Button = sdl.BUTTON_MIDDLE
	RightButton  Button = sdl.BUTTON_RIGHT
	X1Button     Button = sdl.BUTTON_X1
	X2Button     Button = sdl.BUTTON_X2
)
