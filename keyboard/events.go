package keyboard

import (
	"github.com/veandco/go-sdl2/sdl"
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
