package joystick

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Vibration struct {
	Left, Right float32
	Effect      sdl.HapticEffect
	Data        [4]uint16
	ID          int
	Endtime     uint32
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
	joystick.vibration = &Vibration{
		ID:      -1,
		Left:    0.0,
		Right:   0.0,
		Endtime: sdl.HAPTIC_INFINITY,
		Effect:  sdl.HapticEffect{},
	}

	return joystick.haptic != nil
}

func (joystick *Joystick) runVibrationEffect() bool {
	if joystick.vibration.ID != -1 {
		if joystick.haptic.UpdateEffect(joystick.vibration.ID, &joystick.vibration.Effect) == 0 {
			if joystick.haptic.RunEffect(joystick.vibration.ID, 1) == 0 {
				return true
			}
		}

		// If the effect fails to update, we should destroy and re-create it.
		joystick.haptic.DestroyEffect(joystick.vibration.ID)
		joystick.vibration.ID = -1
	}

	joystick.vibration.ID = joystick.haptic.NewEffect(&joystick.vibration.Effect)

	if joystick.vibration.ID != -1 && joystick.haptic.RunEffect(joystick.vibration.ID, 1) == 0 {
		return true
	}

	return false
}

func (joystick *Joystick) SetVibration(args ...float32) bool {
	if !joystick.checkCreateHaptic() {
		return false
	}

	// no arguments given means stop the vibration
	if len(args) == 0 {
		return joystick.stopVibration()
	}

	left := float32(math.Min(math.Max(float64(args[0]), 0.0), 1.0))
	right := left
	if len(args) > 1 {
		right = float32(math.Min(math.Max(float64(right), 0.0), 1.0))
	}

	if left == 0.0 && right == 0.0 {
		return joystick.stopVibration()
	}

	length := sdl.HAPTIC_INFINITY
	if len(args) > 2 && args[2] >= 0.0 {
		maxduration := math.MaxUint32 / 1000.0
		length = int(math.Min(float64(args[2]), float64(maxduration)) * 1000)
	}

	success := false
	features := joystick.haptic.Query()
	axes := joystick.haptic.NumAxes()

	if (features & sdl.HAPTIC_LEFTRIGHT) != 0 {
		//joystick.vibration.Effect.Type = sdl.HAPTIC_LEFTRIGHT
		lr := joystick.vibration.Effect.LeftRight()
		lr.Length = uint32(length)
		lr.LargeMagnitude = uint16(left * math.MaxUint16)
		lr.SmallMagnitude = uint16(right * math.MaxUint16)
		success = joystick.runVibrationEffect()
	}

	// Some gamepad drivers only give support for controlling individual motors
	// through a custom FF effect.
	if !success && joystick.IsGamepad() && (features&sdl.HAPTIC_CUSTOM) != 0 && axes == 2 {
		// NOTE: this may cause issues with drivers which support custom effects
		// but aren't similar to https://github.com/d235j/360Controller .

		// Custom effect data is clamped to 0x7FFF in SDL.
		joystick.vibration.Data[0] = uint16(left * 0x7FFF)
		joystick.vibration.Data[2] = uint16(left * 0x7FFF)
		joystick.vibration.Data[1] = uint16(right * 0x7FFF)
		joystick.vibration.Data[3] = uint16(right * 0x7FFF)

		//joystick.vibration.Effect.Type = sdl.HAPTIC_CUSTOM
		custom := joystick.vibration.Effect.Custom()
		custom.Length = uint32(length)
		custom.Channels = 2
		custom.Period = 10
		custom.Samples = 2
		custom.Data = &joystick.vibration.Data[0]

		success = joystick.runVibrationEffect()
	}

	//// Fall back to a simple sine wave if all else fails. This only supports a
	//// single strength value.
	if !success && (features&sdl.HAPTIC_SINE) != 0 {
		//joystick.vibration.Effect.Type = sdl.HAPTIC_SINE
		periodic := joystick.vibration.Effect.Periodic()
		periodic.Length = uint32(length)
		periodic.Period = 10
		strength := math.Max(float64(left), float64(right))
		periodic.Magnitude = int16(strength * 0x7FFF)
		success = joystick.runVibrationEffect()
	}

	if success {
		joystick.vibration.Left = left
		joystick.vibration.Right = right
		if length == sdl.HAPTIC_INFINITY {
			joystick.vibration.Endtime = sdl.HAPTIC_INFINITY
		} else {
			joystick.vibration.Endtime = sdl.GetTicks() + uint32(length)
		}
	} else {
		joystick.vibration.Left = 0.0
		joystick.vibration.Right = 0.0
		joystick.vibration.Endtime = sdl.HAPTIC_INFINITY
	}

	return success
}

func (joystick *Joystick) stopVibration() bool {
	success := true

	if sdl.WasInit(sdl.INIT_HAPTIC) == 0 && joystick.haptic != nil && sdl.HapticIndex(joystick.haptic) != -1 {
		success = (joystick.haptic.StopEffect(joystick.vibration.ID) == 0)
	}

	if success {
		joystick.vibration.Left = 0.0
		joystick.vibration.Right = 0.0
	}

	return success
}

func (joystick *Joystick) GetVibration() (float32, float32) {
	if joystick.vibration.Endtime != sdl.HAPTIC_INFINITY {
		// With some drivers, the effect physically stops at the right time, but
		// SDL_HapticGetEffectStatus still thinks it's playing. So we explicitly
		// stop it once it's done, just to be sure.
		if joystick.vibration.Endtime-sdl.GetTicks() <= 0 {
			joystick.stopVibration()
			joystick.vibration.Endtime = sdl.HAPTIC_INFINITY
		}
	}

	// Check if the haptic effect has stopped playing.
	if joystick.haptic == nil || joystick.vibration.ID == -1 || joystick.haptic.GetEffectStatus(joystick.vibration.ID) != 1 {
		joystick.vibration.Left = 0.0
		joystick.vibration.Right = 0.0
	}

	return joystick.vibration.Left, joystick.vibration.Right
}
