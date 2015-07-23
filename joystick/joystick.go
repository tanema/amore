package joystick

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type Joystick struct {
	id         int
	stick      *sdl.Joystick
	controller *sdl.GameController
	haptic     *sdl.Haptic
}

func (joystick *Joystick) Open() bool {
	joystick.Close()

	joystick.stick = sdl.JoystickOpen(joystick.id)

	if joystick.stick != nil {
		if sdl.IsGameController(joystick.id) {
			joystick.controller = sdl.GameControllerOpen(joystick.id)
		}

	}

	return joystick.IsConnected()
}

func (joystick *Joystick) IsGamepad() bool {
	return joystick.controller != nil
}

func (joystick *Joystick) GetID() int {
	return joystick.id
}

func (joystick *Joystick) GetName() string {
	// Prefer the Joystick name for consistency.
	name := joystick.stick.Name()
	if name == "" && joystick.controller != nil {
		name = joystick.controller.Name()
	}

	return name
}

func (joystick *Joystick) IsConnected() bool {
	return joystick.stick != nil && joystick.stick.GetAttached()
}

//TODO I dont think guid works
func (joystick *Joystick) GetGUID() string {
	guid := ""
	sdlguid := joystick.stick.GetGUID()
	sdl.JoystickGetGUIDString(sdlguid, guid, int(unsafe.Sizeof(guid)))
	return guid
}

func (joystick *Joystick) GetHandle() *sdl.Joystick {
	return joystick.stick
}

func (joystick *Joystick) Close() {
	if joystick.controller != nil {
		joystick.controller.Close()
	}
	if joystick.stick != nil {
		joystick.stick.Close()
	}
	joystick.stick = nil
	joystick.controller = nil
}

//TODO WTF is byte for get button
func (joystick *Joystick) IsDown(button int) bool {
	if joystick.IsConnected() == false {
		return false
	}

	println(joystick.stick.GetButton(button))
	return true
}

func (joystick *Joystick) GetNumAxes() int {
	if joystick.IsConnected() == false {
		return 0
	}
	return joystick.stick.NumAxes()
}

func (joystick *Joystick) GetNumButton() int {
	if joystick.IsConnected() == false {
		return 0
	}
	return joystick.stick.NumButtons()
}

func (joystick *Joystick) GetNumHats() int {
	if joystick.IsConnected() == false {
		return 0
	}
	return joystick.stick.NumHats()
}

func (joystick *Joystick) GetAxis(axisindex int) float32 {
	if joystick.IsConnected() == false || axisindex < 0 || axisindex >= joystick.GetNumAxes() {
		return 0.0
	}

	return float32(clampval(float64(joystick.stick.GetAxis(axisindex)) / 32768.0))
}

func (joystick *Joystick) GetAxes() []float32 {
	count := joystick.GetNumAxes()
	axes := []float32{}

	if joystick.IsConnected() == false || count <= 0 {
		return axes
	}

	for i := 0; i < count; i++ {
		axes = append(axes, float32(clampval(float64(joystick.stick.GetAxis(i))/32768.0)))
	}

	return axes
}

// TODO WFT byte still?
func (joystick *Joystick) GetHat(hatindex int) byte {
	if joystick.IsConnected() == false || hatindex < 0 || hatindex >= joystick.GetNumHats() {
		var a byte
		return a
	}

	return joystick.stick.GetHat(hatindex)
}

func (joystick *Joystick) GetGamepadAxis(axis GameControllerAxis) float32 {
	if joystick.IsConnected() == false || joystick.IsGamepad() == false {
		return 0.0
	}

	value := joystick.controller.GetAxis(sdl.GameControllerAxis(axis))

	return float32(clampval(float64(value) / 32768.0))
}

// TODO wtf is byte for get button?
func (joystick *Joystick) IsGamepadDown(button GameControllerButton) bool {
	if joystick.IsConnected() == false || joystick.IsGamepad() == false {
		return false
	}

	joystick.controller.GetButton(sdl.GameControllerButton(button))

	return false
}

func (joystick *Joystick) IsVibrationSupported() bool {
	if joystick.checkCreateHaptic() == false {
		return false
	}

	features := joystick.haptic.Query()

	if (features & sdl.HAPTIC_LEFTRIGHT) != 0 {
		return true
	}

	// Some gamepad drivers only support left/right motors via a custom effect.
	if joystick.IsGamepad() && (features&sdl.HAPTIC_CUSTOM) != 0 {
		return true
	}

	// Test for simple sine wave support as a last resort.
	if (features & sdl.HAPTIC_SINE) != 0 {
		return true
	}

	return false
}

func (joystick *Joystick) checkCreateHaptic() bool {
	if joystick.IsConnected() == false {
		return false
	}

	if sdl.WasInit(sdl.INIT_HAPTIC) == 0 && sdl.InitSubSystem(sdl.INIT_HAPTIC) != nil {
		return false
	}

	if joystick.haptic != nil && sdl.HapticIndex(joystick.haptic) != -1 {
		return true
	}

	if joystick.haptic != nil {
		joystick.haptic.Close()
	}

	joystick.haptic = sdl.HapticOpenFromJoystick(joystick.stick)

	return joystick.haptic != nil
}

//TODO
func (joystick *Joystick) SetVibration(left, right, duration float32) bool {
	return false
}

func (joystick *Joystick) StopVibration() bool {
	return false
}

func (joystick *Joystick) GetVibration() (float32, float32) {
	return 0.0, 0.0
}
