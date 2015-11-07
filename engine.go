package amore

import (
	"runtime"

	"github.com/tanema/amore/event"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/timer"
	"github.com/tanema/amore/window"
)

type LoadCb func()
type UpdateCb func(float32)
type DrawCb func()

var (
	current_window *window.Window
)

func Start(load LoadCb, update UpdateCb, draw DrawCb) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()

	if current_window, err = window.New(); err != nil {
		return err
	}
	defer current_window.Destroy()

	load()

	for !current_window.ShouldClose() {
		// update
		timer.Step()

		update(timer.GetDelta())

		// draw
		gfx.Reset()
		draw()
		current_window.SwapBuffers()

		// get user interactions
		event.Poll()
	}

	return
}

func Quit() {
	current_window.SetShouldClose(true)
}
