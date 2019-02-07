// Copywrite 2016 Tim Anema. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be foind in the LICENSE file.

// Package amore  is simply for stopping and starting your game. It will
// also automatically lock the os thread and set the application to use all
// available cpus.
package amore

import (
	"runtime"

	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"

	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/timer"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread() //important OpenGl Demand it and stamp thier feet if you dont
}

// Start creates a window and context for the game to run on and runs the game
// loop. As such this function should be put as the last call in your main function.
// update and draw will be called synchronously because calls to OpenGL that are
// not on the main thread will crash your program.
func Start(update func(float32), draw func()) error {
	if err := glfw.Init(gl.ContextWatcher); err != nil {
		return err
	}
	defer glfw.Terminate()

	config, err := loadConfig()
	if err != nil {
		return err
	}

	window, err := glfw.CreateWindow(config.Width, config.Height, config.Title, nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()

	gfx.InitContext(int32(config.Width), int32(config.Height))
	defer gfx.DeInit()

	for !window.ShouldClose() {
		timer.Step()
		update(timer.GetDelta())

		gfx.ClearC(gfx.GetBackgroundColorC())
		gfx.Origin()
		draw()
		gfx.Present()

		window.SwapBuffers()
		glfw.PollEvents()
	}

	return nil
}

// Quit will prepare the window to close at the end of the next game loop. This
// will allow a nice clean destruction of all object that are allocated in OpenGL
func Quit() {
}
