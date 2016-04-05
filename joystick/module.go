package joystick

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	joysticks    = []*Joystick{}
	activeSticks = []*Joystick{}
)

func init() {
	if err := sdl.InitSubSystem(sdl.INIT_JOYSTICK | sdl.INIT_GAMECONTROLLER); err != nil {
		panic(err)
	}

	for i := 0; i < sdl.NumJoysticks(); i++ {
		addJoystick(i)
	}

	sdl.JoystickEventState(sdl.ENABLE)
	sdl.GameControllerEventState(sdl.ENABLE)
}

func addJoystick(idx int) *Joystick {
	if idx < 0 || idx >= sdl.NumJoysticks() {
		return nil
	}

	guidstr := getDeviceGUID(idx)
	var joystick *Joystick
	reused := false

	for _, stick := range joysticks {
		// Try to re-use a disconnected Joystick with the same GUID.
		if stick.GetGUID() == guidstr {
			joystick = stick
			reused = true
			break
		}
	}

	if joystick == nil {
		joystick = &Joystick{id: idx}
		joysticks = append(joysticks, joystick)
	}

	// Make sure the Joystick object isn't in the active list already.
	removeJoystick(joystick)

	// Make sure multiple instances of the same physical joystick aren't added
	// to the active list.
	for _, stick := range activeSticks {
		if joystick.GetHandle() == stick.GetHandle() {
			// If we just created the stick, remove it since it's a duplicate.
			if !reused {
				joysticks = joysticks[:len(joysticks)-1]
			}
			return stick
		}
	}

	if !joystick.Open() {
		return nil
	}

	activeSticks = append(activeSticks, joystick)
	return joystick
}

func removeJoystick(joystick *Joystick) {
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
	if idx < 0 || idx >= sdl.NumJoysticks() {
		return ""
	}

	return sdl.JoystickGetGUIDString(sdl.JoystickGetDeviceGUID(idx))
}

//Gets the number of connected joysticks.
func GetJoystickCount() int {
	return len(activeSticks)
}

//Gets a list of connected Joysticks.
func GetJoysticks() []*Joystick {
	return activeSticks
}

func getJoystickFromID(id int) *Joystick {
	for _, stick := range activeSticks {
		if stick.GetID() == id {
			return stick
		}
	}
	return nil
}

func clampval(x float32) float32 {
	if math.Abs(float64(x)) < 0.01 {
		return 0.0
	}

	if x < -0.99 {
		return -1.0
	} else if x > 0.99 {
		return 1.0
	}
	return x
}
