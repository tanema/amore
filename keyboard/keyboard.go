// The keyboard Pacakge handles the keyboard events on the gl context
package keyboard

import (
	"github.com/tanema/amore/window/ui"
)

var (
	key_repeat = true
	text_input = true
)

//Checks whether a certain key is down.
func IsDown(key Key) bool {
	return ui.IsKeyDown(ui.Keycode(key))
}

//Checks whether a certain scancode is down.
func IsScancodeDown(scancode Scancode) bool {
	return ui.IsScancodeDown(ui.Scancode(scancode))
}

func GetKeyFromScancode(code Scancode) Key {
	return Key(ui.GetKeyFromScancode(ui.Scancode(code)))
}

func GetScancodeFromKey(key Key) Scancode {
	return Scancode(ui.GetScancodeFromKey(ui.Keycode(key)))
}

//Enables or disables key repeat for love.keypressed.
func SetKeyRepeat(enabled bool) {
	key_repeat = enabled
}

//Gets whether key repeat is enabled.
func HasKeyRepeat() bool {
	return key_repeat
}

//Gets whether text input events are enabled.
func HasTextInput() bool {
	return ui.HasTextInput()
}

//Enables or disables text input events.
func SetTextInput(enabled bool) {
	ui.SetTextInput(enabled)
}
