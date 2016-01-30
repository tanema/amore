package amore

import (
	"runtime"

	"github.com/tanema/amore/event"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/timer"
	"github.com/tanema/amore/window"
)

var (
	current_window *window.Window
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread() //important SDL and OpenGl Demand it and stamp thier feet if you dont
}

func Start(update func(float32), draw func()) {
	current_window = window.GetCurrent()
	defer current_window.Destroy()
	gfx.InitContext(current_window.GetWidth(), current_window.GetHeight())
	defer gfx.DeInit()
	for !current_window.ShouldClose() {
		timer.Step()
		update(timer.GetDelta())
		gfx.ClearC(gfx.GetBackgroundColorC())
		gfx.Origin()
		draw()
		gfx.Present()
		event.Poll()
	}
}

func Quit() {
	current_window.SetShouldClose(true)
}
