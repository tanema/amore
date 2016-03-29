package main

import (
	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
)

func main() {
	err := amore.Start(update, draw)
	if err != nil {
		panic(err)
	}
}

func update(deltaTime float32) {
}

func draw() {
	gfx.SetColor(239, 96, 17, 255)
	gfx.Rect("fill", 100, 100, 446, 440)
}
