// The keyboard Pacakge handles the keyboard events on the gl context
package keyboard

import (
	"github.com/veandco/go-sdl2/sdl"
)

var (
	key_repeat = true
	text_input = true
)

// IsDown checks whether a certain key is down.
func IsDown(key Key) bool {
	return IsScancodeDown(GetScancodeFromKey(key))
}

// IsScancodeDown checks whether a certain scancode is down.
func IsScancodeDown(scancode Scancode) bool {
	state := sdl.GetKeyboardState()
	return state[int(scancode)] == 1
}

// GetKeyFromScancode translates a scancode to key.
func GetKeyFromScancode(code Scancode) Key {
	return Key(sdl.GetKeyFromScancode(sdl.Scancode(code)))
}

// GetScancodeFromKey translates a key to scancode.
func GetScancodeFromKey(key Key) Scancode {
	return Scancode(sdl.GetScancodeFromKey(sdl.Keycode(key)))
}

// SetKeyRepeat wnables or disables key repeat for love.keypressed.
func SetKeyRepeat(enabled bool) {
	key_repeat = enabled
}

// HasKeyRepeat gets whether key repeat is enabled.
func HasKeyRepeat() bool {
	return key_repeat
}

// HasTextInput gets whether text input events are enabled. For example if enabled, and
// shift-2 is pressed on an American keyboard layout, the text "@" will be generated.
// If disabled just a 2 will be sent to the keypress/keyrelease callbacks
func HasTextInput() bool {
	return sdl.IsTextInputActive() != false
}

// SetTextInput enables or disables text input events. For reference of what text input
// is, please refer to HasTextInput
func SetTextInput(enabled bool) {
	if enabled {
		sdl.StartTextInput()
	} else {
		sdl.StopTextInput()
	}
}
