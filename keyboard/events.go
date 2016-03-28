package keyboard

import (
	"github.com/tanema/amore/window/ui"
)

type KeyPressCB func(key Key, is_repeat bool)
type KeyReleaseCB func(key Key)
type TextInputCB func(str string)
type TextEditCB func(str string, start, length int32)

var (
	key_press_default   KeyPressCB   = func(key Key, is_repeat bool) {}
	key_release_default KeyReleaseCB = func(key Key) {}
	text_input_default  TextInputCB  = func(str string) {}
	text_edit_default   TextEditCB   = func(str string, start, length int32) {}

	key_press_cb   = key_press_default
	key_release_cb = key_release_default
	text_input_cb  = text_input_default
	text_edit_cb   = text_edit_default
)

func Delegate(event ui.Event) {
	switch e := event.(type) {
	case *ui.KeyDownEvent:
		is_repeat := (e.Repeat == 1)

		if is_repeat && !key_repeat {
			return
		}

		key := GetKeyFromScancode(Scancode(e.Keysym.Scancode))
		key_press_cb(key, is_repeat)
	case *ui.KeyUpEvent:
		key := GetKeyFromScancode(Scancode(e.Keysym.Scancode))
		key_release_cb(key)
	case *ui.TextEditingEvent:
		text_edit_cb(string(e.Text[:]), int32(e.Start), int32(e.Length))
	case *ui.TextInputEvent:
		text_input_cb(string(e.Text[:]))
	}
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

func SetTextEditCB(cb TextEditCB) {
	if cb == nil {
		text_edit_cb = text_edit_default
	} else {
		text_edit_cb = cb
	}
}
