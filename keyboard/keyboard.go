package keyboard

import (
	"github.com/tanema/go-sdl2/sdl"

	//"github.com/tanema/amore/window"
)

type KeyPressCB func(key Key, is_repeat bool)
type KeyReleaseCB func(key Key)
type TextInputCB func(str string)

var (
	key_repeat = true
	text_input = true

	key_press_default   KeyPressCB   = func(key Key, is_repeat bool) {}
	key_release_default KeyReleaseCB = func(key Key) {}
	text_input_default  TextInputCB  = func(str string) {}

	key_press_cb   = key_press_default
	key_release_cb = key_release_default
	text_input_cb  = text_input_default
)

func SetKeyReleaseCB(cb KeyReleaseCB) {
	if cb == nil {
		key_release_cb = key_release_default
	} else {
		key_release_cb = cb
	}
}

func SetTextInputCB(cb TextInputCB) {
	if cb == nil {
		text_input_cb = text_input_default
	} else {
		text_input_cb = cb
	}
}

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
