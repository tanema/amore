package window

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (window *Window) Delegate(event *sdl.WindowEvent) {
	switch event.Type {
	case sdl.WINDOWEVENT_NONE:
		return
	case sdl.WINDOWEVENT_SHOWN:
	case sdl.WINDOWEVENT_HIDDEN:
	case sdl.WINDOWEVENT_EXPOSED:
	case sdl.WINDOWEVENT_MOVED:
	case sdl.WINDOWEVENT_RESIZED:
	case sdl.WINDOWEVENT_SIZE_CHANGED:
	case sdl.WINDOWEVENT_MINIMIZED:
	case sdl.WINDOWEVENT_MAXIMIZED:
	case sdl.WINDOWEVENT_RESTORED:
	case sdl.WINDOWEVENT_ENTER:
	case sdl.WINDOWEVENT_LEAVE:
	case sdl.WINDOWEVENT_FOCUS_GAINED:
	case sdl.WINDOWEVENT_FOCUS_LOST:
	case sdl.WINDOWEVENT_CLOSE:
	}
}
