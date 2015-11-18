package system

import (
	"github.com/veandco/go-sdl2/sdl"
)

type PowerState int

const (
	POWERSTATE_UNKNOWN    PowerState = sdl.POWERSTATE_UNKNOWN
	POWERSTATE_ON_BATTERY PowerState = sdl.POWERSTATE_ON_BATTERY
	POWERSTATE_NO_BATTERY PowerState = sdl.POWERSTATE_NO_BATTERY
	POWERSTATE_CHARGING   PowerState = sdl.POWERSTATE_CHARGING
	POWERSTATE_CHARGED    PowerState = sdl.POWERSTATE_CHARGED
)
