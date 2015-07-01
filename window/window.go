package window

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/mouse"
)

var (
	Title          = "Engine Test"
	current_window *glfw.Window
	Width          = 1280
	Height         = 720
)

func New() (window *glfw.Window, err error) {
	if err = glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	if current_window, err = glfw.CreateWindow(Width, Height, Title, nil, nil); err != nil {
		return current_window, err
	}
	current_window.MakeContextCurrent()
	current_window.SetSizeCallback(onResize)

	current_window.SetKeyCallback(keyboard.OnKey)
	current_window.SetCharCallback(keyboard.OnChar)

	current_window.SetMouseButtonCallback(mouse.OnClick)
	current_window.SetScrollCallback(mouse.OnScroll)

	return current_window, nil
}

func GetCurrent() *glfw.Window {
	return current_window
}

func onResize(win *glfw.Window, w, h int) {
	Height = h
	Width = w
	gl.Viewport(0, 0, int32(Width), int32(Height))
}

func SetTitle(title string) {
	current_window.SetTitle(title)
	Title = title
}

func GetTitle() string {
	return Title
}

func GetWidth() int {
	return Width
}

func GetHeight() int {
	return Height
}
