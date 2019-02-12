// Copywrite 2016 Tim Anema. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be foind in the LICENSE file.

// Package amore  is simply for stopping and starting your game. It will
// also automatically lock the os thread and set the application to use all
// available cpus.
package amore

import (
	"runtime"
	"time"

	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"

	"github.com/tanema/amore/gfx"
)

var (
	fps               int       // frames per second
	frames            int       // frames since last update freq
	previousTime      time.Time // last frame time
	previousFPSUpdate time.Time // last time fps was updated
	currentWindow     *glfw.Window
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

	currentWindow, err := glfw.CreateWindow(config.Width, config.Height, config.Title, nil, nil)
	if err != nil {
		return err
	}
	currentWindow.MakeContextCurrent()

	glfw.WindowHint(glfw.Samples, config.Msaa)
	gl.SampleCoverage(1, false)
	bufferW, bufferH := currentWindow.GetFramebufferSize()
	gfx.InitContext(int32(bufferW), int32(bufferH))
	defer gfx.DeInit()

	for !currentWindow.ShouldClose() {
		update(step())

		gfx.Clear(gfx.GetBackgroundColor()...)
		gfx.Origin()
		draw()
		gfx.Present()

		currentWindow.SwapBuffers()
		glfw.PollEvents()
	}

	return nil
}

func step() float32 {
	frames++
	dt := float32(time.Since(previousTime).Seconds())
	previousTime = time.Now()
	timeSinceLast := float32(time.Since(previousFPSUpdate).Seconds())
	if timeSinceLast > 1 { //update 1 per second
		fps = int((float32(frames) / timeSinceLast) + 0.5)
		previousFPSUpdate = previousTime
		frames = 0
	}
	return dt
}

// GetFPS returns the number of frames per second.
func GetFPS() int {
	return fps
}

// Quit will prepare the window to close at the end of the next game loop. This
// will allow a nice clean destruction of all object that are allocated in OpenGL
func Quit() {
	currentWindow.SetShouldClose(true)
}
