package main

import (
	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
)

var (
	img *gfx.Image
)

func main() {
	amore.Start("Test Game", load, update, draw)
}

func load() {
	var err error
	img, err = gfx.NewImage("assets/palm_tree.png")
	if err != nil {
		panic(err)
	}
}

func update(deltaTime float64) {
}

func draw() {
	gfx.Translate(100, 20)
	// rectangle
	gfx.SetColor(0.0, 170.0, 0.0, 155.0)
	gfx.Rect("fill", 20.0, 20.0, 400.0, 200.0)
	// line
	gfx.SetColor(255.0, 170.0, 0.0, 255.0)
	gfx.Line(800.0, 100.0, 900.0, 100.0)
	// image
	gfx.SetColor(255.0, 255.0, 255.0, 255.0)
	gfx.DrawS(img, 500, 100)
}
