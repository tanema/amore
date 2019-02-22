package runtime

import (
	"github.com/goxjs/glfw"
)

type config struct {
	Title      string
	Width      int
	Height     int
	Resizable  bool
	Fullscreen bool
	MouseShown bool
}

var (
	conf = config{
		Width:  800,
		Height: 600,
	}
)

type window struct {
	*glfw.Window
	active bool
}

func createWindow(conf config) (window, error) {
	newWin := window{active: true}

	var err error
	newWin.Window, err = glfw.CreateWindow(conf.Width, conf.Height, conf.Title, nil, nil)
	if err != nil {
		return window{}, err
	}
	newWin.MakeContextCurrent()
	if conf.Resizable {
		glfw.WindowHint(glfw.Resizable, 1)
	} else {
		glfw.WindowHint(glfw.Resizable, 0)
	}
	glfw.WindowHint(glfw.Samples, 4.0)
	newWin.SetFocusCallback(newWin.focus)
	newWin.SetIconifyCallback(newWin.iconify)
	return newWin, nil
}

func (win *window) focus(w *glfw.Window, focused bool)     { win.active = focused }
func (win *window) iconify(w *glfw.Window, iconified bool) { win.active = !iconified }
