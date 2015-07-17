package mouse

import (
	"github.com/veandco/go-sdl2/sdl"
)

func Delegate(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.MouseMotionEvent:
	case *sdl.MouseButtonEvent:
		switch e.Type {
		case sdl.MOUSEBUTTONDOWN:
		case sdl.MOUSEBUTTONUP:
		}
	case *sdl.MouseWheelEvent:
	}
}
