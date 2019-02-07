package main

import (
	"fmt"

	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/timer"
)

func main() {
	amore.Start(update, draw)
}

func update(deltaTime float32) {
}

func draw() {
	gfx.SetColor(255, 0, 0, 255)
	gfx.Rect("line", 50, 50, 100, 100)

	gfx.Print(fmt.Sprintf("fps: %v", timer.GetFPS()), 0, 0)
}
