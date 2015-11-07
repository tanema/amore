package joystick

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Vibration struct {
	Left, Right float32
	Effect      sdl.HapticEffect
	Data        [4]uint16
	ID          int
	Endtime     uint32
}

func defaultVibration() *Vibration {
	return &Vibration{
		ID:      -1,
		Left:    0.0,
		Right:   0.0,
		Endtime: sdl.HAPTIC_INFINITY,
	}
}

type Joystick struct {
	id         int
	stick      *sdl.Joystick
	controller *sdl.GameController
	haptic     *sdl.Haptic
	vibration  *Vibration
}

func (joystick *Joystick) Open() bool {
	joystick.Close()

	joystick.stick = sdl.JoystickOpen(joystick.id)

	if joystick.stick != nil && sdl.IsGameController(joystick.id) {
		joystick.controller = sdl.GameControllerOpen(joystick.id)
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

func (joystick *Joystick) GetGUID() string {
	return sdl.JoystickGetGUIDString(joystick.stick.GetGUID())
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

func (joystick *Joystick) IsDown(button int) bool {
	if joystick.IsConnected() == false {
		return false
	}

	return joystick.stick.GetButton(button) == 1
}

func (joystick *Joystick) GetAxisCount() int {
	if joystick.IsConnected() == false {
		return 0
	}
	return joystick.stick.NumAxes()
}

func (joystick *Joystick) GetButtonCount() int {
	if joystick.IsConnected() == false {
		return 0
	}
	return joystick.stick.NumButtons()
}

func (joystick *Joystick) GetHatCount() int {
	if joystick.IsConnected() == false {
		return 0
	}
	return joystick.stick.NumHats()
}

func (joystick *Joystick) GetAxis(axisindex int) float32 {
	if joystick.IsConnected() == false || axisindex < 0 || axisindex >= joystick.GetAxisCount() {
		return 0.0
	}

	return clampval(float32(joystick.stick.GetAxis(axisindex)) / 32768.0)
}

func (joystick *Joystick) GetAxes() []float32 {
	count := joystick.GetAxisCount()
	axes := []float32{}

	if joystick.IsConnected() == false || count <= 0 {
		return axes
	}

	for i := 0; i < count; i++ {
		axes = append(axes, clampval(float32(joystick.stick.GetAxis(i))/32768.0))
	}

	return axes
}

func (joystick *Joystick) GetHat(hatindex int) byte {
	if joystick.IsConnected() == false || hatindex < 0 || hatindex >= joystick.GetHatCount() {
		return 0
	}

	return joystick.stick.GetHat(hatindex)
}

func (joystick *Joystick) GetGamepadAxis(axis GameControllerAxis) float32 {
	if joystick.IsConnected() == false || joystick.IsGamepad() == false {
		return 0.0
	}

	value := joystick.controller.GetAxis(sdl.GameControllerAxis(axis))

	return clampval(float32(value) / 32768.0)
}

func (joystick *Joystick) IsGamepadDown(button GameControllerButton) bool {
	if joystick.IsConnected() == false || joystick.IsGamepad() == false {
		return false
	}

	return joystick.controller.GetButton(sdl.GameControllerButton(button)) == 1
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
	joystick.vibration = defaultVibration()

	return joystick.haptic != nil
}

// TODO
func (joystick *Joystick) SetVibration(left, right, duration float32) bool {
	panic("set vibration not implemented yet")
	return false
}

func (joystick *Joystick) StopVibration() bool {
	panic("vibration not implemented yet")
	return false
}

func (joystick *Joystick) GetVibration() (float32, float32) {
	panic("vibration not implemented yet")
	return 0.0, 0.0
}
