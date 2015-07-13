package window

import (
	"fmt"

	//"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"

	//"github.com/tanema/amore/keyboard"
	//"github.com/tanema/amore/mouse"
)

var (
	default_title  = "Amore"
	current_window *Window
	default_width  = 1280
	default_height = 720
)

type Window struct {
	sdl_window    *sdl.Window
	context       sdl.GLContext
	width, height int
	title         string
	should_close  bool
}

func New() (*Window, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(default_title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, default_width, default_height, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}

	context, err := sdl.GL_CreateContext(window)
	if err != nil {
		panic(err)
	}

	current_window = &Window{
		sdl_window:   window,
		context:      context,
		width:        default_width,
		height:       default_height,
		title:        default_title,
		should_close: false,
	}

	return current_window, nil
}

func GetCurrent() *Window {
	return current_window
}

//func onResize(win *glfw.Window, w, h int) {
//Height = h
//Width = w
//gl.Viewport(0, 0, int32(Width), int32(Height))
//}

func (window *Window) SetTitle(title string) {
	window.sdl_window.SetTitle(title)
	window.title = title
}

func (window *Window) GetTitle() string {
	return window.title
}

func (window *Window) GetWidth() int {
	return window.width
}

func (window *Window) GetHeight() int {
	return window.height
}

func (window *Window) ShouldClose() bool {
	return window.should_close
}

func (window *Window) SetShouldClose(should_close bool) {
	window.should_close = should_close
}

func (window *Window) SwapBuffers() {
	sdl.GL_SwapWindow(window.sdl_window)
}

func (window *Window) PollEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			window.SetShouldClose(true)
		case *sdl.MouseMotionEvent:
			fmt.Printf("[%d ms] MouseMotion\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n", t.Timestamp, t.Which, t.X, t.Y, t.XRel, t.YRel)
		}
	}
}

func (window *Window) Destroy() {
	sdl.GL_DeleteContext(window.context)
	window.sdl_window.Destroy()
	sdl.Quit()
}
