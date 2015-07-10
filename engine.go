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

func Start(load LoadCb, update UpdateCb, draw DrawCb) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()

	if current_window, err = window.New(); err != nil {
		return err
	}

	if err = setupGL(); err != nil {
		return err
	}

	load()

	defer glfw.Terminate()
	for !current_window.ShouldClose() {
		// update
		time := glfw.GetTime()
		update(time - current_time)
		current_time = time

		// draw
		gfx.Reset()
		draw()
		current_window.SwapBuffers()

		// get user interactions
		glfw.PollEvents()
	}

	return
}

func setupGL() error {
	if err := gl.Init(); err != nil {
		return err
	}
	gl.Enable(gl.BLEND)
	// Auto-generated mipmaps should be the best quality possible
	gl.Hint(gl.GENERATE_MIPMAP_HINT, gl.NICEST)
	// Set pixel row alignment
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	return nil
}

func Quit() {
	current_window.SetShouldClose(true)
}
