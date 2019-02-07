package main

import (
	"fmt"

	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
)

func main() {
	amore.Start(update, draw)
}

func update(dt float32) {
	fmt.Println(dt)
}

func draw() {
	gfx.SetColor(255, 0, 0, 255)
	gfx.Rect("fill", 50, 50, 100, 100)

	gfx.Print(fmt.Sprintf("fps: %v", amore.GetFPS()), 0, 0)
}
