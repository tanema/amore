package system

import (
	"github.com/veandco/go-sdl2/sdl"
)

// PowerState defines the battery connection to the PC
type PowerState int

// PowerStates are different status of the battery status
const (
	PowerStateUnknown   PowerState = sdl.POWERSTATE_UNKNOWN
	PowerStateOnBattery PowerState = sdl.POWERSTATE_ON_BATTERY
	PowerStateNoBattery PowerState = sdl.POWERSTATE_NO_BATTERY
	PowerStateCharging  PowerState = sdl.POWERSTATE_CHARGING
	PowerStateCharged   PowerState = sdl.POWERSTATE_CHARGED
)
