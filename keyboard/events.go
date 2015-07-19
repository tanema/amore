package keyboard

import (
	"github.com/veandco/go-sdl2/sdl"
)

type KeyPressCB func(key Key, is_repeat bool)
type KeyReleaseCB func(key Key)
type TextInputCB func(str string)

var (
	key_press_default   KeyPressCB   = func(key Key, is_repeat bool) {}
	key_release_default KeyReleaseCB = func(key Key) {}
	text_input_default  TextInputCB  = func(str string) {}

	key_press_cb   = key_press_default
	key_release_cb = key_release_default
	text_input_cb  = text_input_default
)

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
		//not done right now
	case *sdl.TextInputEvent:
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
