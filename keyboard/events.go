package keyboard

import (
	"github.com/veandco/go-sdl2/sdl"
)

var (
	// OnKeyDown is called when a key on the keyboard is pressed down.
	OnKeyDown func(key Key, isRepeat bool)
	// OnKeyUp is called when a key on the keyboard is released.
	OnKeyUp func(key Key)
	// OnTextInput is called when text has been entered by the user. For example if
	// shift-2 is pressed on an American keyboard layout, the text "@" will be generated.
	OnTextInput func(str string)
	// OnTextEdit is called when the candidate text for an IME (Input Method Editor) has changed.
	OnTextEdit func(str string, start, length int32)
)

// Delegate is used by amore/event to pass events to the keyboard package. It may
// also be useful to stub or fake events
func Delegate(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		isRepeat := (e.Repeat == 1)

		if isRepeat && !keyRepeat {
			return
		}

		key := GetKeyFromScancode(Scancode(e.Keysym.Scancode))
		if OnKeyDown != nil {
			OnKeyDown(key, isRepeat)
		}
	case *sdl.KeyUpEvent:
		key := GetKeyFromScancode(Scancode(e.Keysym.Scancode))
		if OnKeyUp != nil {
			OnKeyUp(key)
		}
	case *sdl.TextEditingEvent:
		if OnTextEdit != nil {
			OnTextEdit(string(e.Text[:]), e.Start, e.Length)
		}
	case *sdl.TextInputEvent:
		if OnTextInput != nil {
			OnTextInput(string(e.Text[:]))
		}
	}
}
