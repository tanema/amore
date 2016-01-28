package amore

import (
	"errors"
	"runtime"

	"github.com/tanema/amore/event"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/timer"
	"github.com/tanema/amore/window"
)

var (
	current_window *window.Window
)

func Start(update func(float32), draw func()) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if current_window = window.GetCurrent(); current_window == nil {
		return errors.New("Cound not get a window")
	}
	defer current_window.Destroy()

	for !current_window.ShouldClose() {
		// update
		timer.Step()

		update(timer.GetDelta())

		// draw
		gfx.ClearC(gfx.GetBackgroundColorC())
		gfx.Origin()
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
