package window

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (window *Window) Delegate(event *sdl.WindowEvent) {
	switch event.Event {
	case sdl.WINDOWEVENT_NONE:
		return //handled already by event/event.go
	case sdl.WINDOWEVENT_SHOWN:
	case sdl.WINDOWEVENT_HIDDEN:
	case sdl.WINDOWEVENT_EXPOSED:
	case sdl.WINDOWEVENT_MOVED:
	case sdl.WINDOWEVENT_RESIZED, sdl.WINDOWEVENT_SIZE_CHANGED:
		window.OnSizeChanged(event.Data1, event.Data2)
	case sdl.WINDOWEVENT_MINIMIZED:
	case sdl.WINDOWEVENT_MAXIMIZED:
	case sdl.WINDOWEVENT_RESTORED:
	case sdl.WINDOWEVENT_ENTER, sdl.WINDOWEVENT_LEAVE:
		//handled by event/event.go and delegated to the mouse
	case sdl.WINDOWEVENT_FOCUS_GAINED:
		sdl.DisableScreenSaver()
	case sdl.WINDOWEVENT_FOCUS_LOST:
		sdl.EnableScreenSaver()
	case sdl.WINDOWEVENT_CLOSE:
	}
}
