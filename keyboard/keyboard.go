package keyboard

import (
	"github.com/veandco/go-sdl2/sdl"
)

var (
	key_repeat = true
	text_input = true
)

//Checks whether a certain key is down.
func IsDown(key Key) bool {
	state := sdl.GetKeyboardState()
	scancode := GetScancodeFromKey(key)
	for _, code := range state {
		if code == uint8(scancode) {
			return true
		}
	}
	return false
}

func GetKeyFromScancode(code Scancode) Key {
	return Key(sdl.GetKeyFromScancode(sdl.Scancode(code)))
}

func GetScancodeFromKey(key Key) Scancode {
	return Scancode(sdl.GetScancodeFromKey(sdl.Keycode(key)))
}

////Enables or disables key repeat for love.keypressed.
func SetKeyRepeat(enabled bool) {
	key_repeat = enabled
}

////Gets whether key repeat is enabled.
func HasKeyRepeat() bool {
	return key_repeat
}

////Gets whether text input events are enabled.
func HasTextInput() bool {
	return sdl.IsTextInputActive() != false
}

////Enables or disables text input events.
func SetTextInput(enabled bool) {
	if enabled {
		sdl.StartTextInput()
	} else {
		sdl.StopTextInput()
	}
}
