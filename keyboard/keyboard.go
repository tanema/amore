package keyboard

import (
	"github.com/go-gl/glfw/v3.1/glfw"
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

func OnKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		key_press_cb(Key(key), false)
	} else if action == glfw.Repeat {
		key_press_cb(Key(key), true)
	} else if action == glfw.Release {
		key_release_cb(Key(key))
	} else {
		println("OnKey: unknow keyboard action")
	}
}

func OnChar(window *glfw.Window, char rune) {
	text_input_cb(string(char))
}

func SetKeyPressCB(cb KeyPressCB) {
	if cb == nil {
		key_press_cb = key_press_default
	} else {
		key_press_cb = cb
	}
}

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
	action := glfw.GetCurrentContext().GetKey(glfw.Key(key))
	return action == glfw.Repeat || action == glfw.Press
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
	return text_input
}

//Enables or disables text input events.
func SetTextInput(enabled bool) {
	text_input = enabled
}
