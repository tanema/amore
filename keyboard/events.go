package keyboard

import (
	"github.com/veandco/go-sdl2/sdl"
)

var (
	key_press_default   = func(key Key, is_repeat bool) {}
	key_release_default = func(key Key) {}
	text_input_default  = func(str string) {}
	text_edit_default   = func(str string, start, length int32) {}

	key_press_cb   = key_press_default
	key_release_cb = key_release_default
	text_input_cb  = text_input_default
	text_edit_cb   = text_edit_default
)

// Delegate is used by amore/event to pass events to the keyboard package. It may
// also be useful to stub or fake events
func Delegate(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		is_repeat := (e.Repeat == 1)

		if is_repeat && !key_repeat {
			return
		}

		key := GetKeyFromScancode(Scancode(e.Keysym.Scancode))
		key_press_cb(key, is_repeat)
	case *sdl.KeyUpEvent:
		key := GetKeyFromScancode(Scancode(e.Keysym.Scancode))
		key_release_cb(key)
	case *sdl.TextEditingEvent:
		text_edit_cb(string(e.Text[:]), e.Start, e.Length)
	case *sdl.TextInputEvent:
		text_input_cb(string(e.Text[:]))
	}
}

// SetKeyPressCB will set a callback to call when the a key on the keyboard is
// pressed down.
func SetKeyPressCB(cb func(key Key, is_repeat bool)) {
	if cb == nil {
		key_press_cb = key_press_default
	} else {
		key_press_cb = cb
	}
}

// SetKeyReleaseCB will set a callback to call when the a key on the keyboard is
// released.
func SetKeyReleaseCB(cb func(key Key)) {
	if cb == nil {
		key_release_cb = key_release_default
	} else {
		key_release_cb = cb
	}
}

// SetTextInputCB is called when text has been entered by the user. For example if
// shift-2 is pressed on an American keyboard layout, the text "@" will be generated.
func SetTextInputCB(cb func(str string)) {
	if cb == nil {
		text_input_cb = text_input_default
	} else {
		text_input_cb = cb
	}
}

// SetTextEditCB is called when the candidate text for an IME (Input Method Editor)
// has changed.
func SetTextEditCB(cb func(str string, start, length int32)) {
	if cb == nil {
		text_edit_cb = text_edit_default
	} else {
		text_edit_cb = cb
	}
}
