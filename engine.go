package amore

import (
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/window"
)

type LoadCb func()
type UpdateCb func(float64)
type DrawCb func()

var (
	current_window *glfw.Window
	current_time   = float64(0)
)

func Start(title string, load LoadCb, update UpdateCb, draw DrawCb) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()

	if current_window, err = window.New(); err != nil {
		return err
	}

	if err = gl.Init(); err != nil {
		return err
	}
	gl.Ortho(0, float64(window.GetWidth()), float64(window.GetHeight()), 0, -1, 1)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	load()

	defer glfw.Terminate()
	for !current_window.ShouldClose() {
		// update
		time := glfw.GetTime()
		update(time - current_time)
		current_time = time

		// draw
		gfx.Clear(gfx.Color{0.0, 0.0, 0.0, 0.0})
		draw()
		current_window.SwapBuffers()

		// get user interactions
		glfw.PollEvents()
	}

	return
}

func Quit() {
	current_window.SetShouldClose(true)
}
