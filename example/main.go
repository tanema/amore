package main

import (
	"math"

	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
)

var (
	tree *gfx.Image
)

func main() {
	amore.Start("Test Game", load, update, draw)
}

func load() {
	var err error
	tree, err = gfx.NewImage("assets/palm_tree.png")
	if err != nil {
		panic(err)
	}
}

func update(deltaTime float64) {
}

func draw() {
	// rectangle
	gfx.SetColor(0.0, 170.0, 0.0, 155.0)
	gfx.Rect("fill", 20.0, 20.0, 200.0, 200.0)
	gfx.Rect("line", 250.0, 20.0, 200.0, 200.0)

	// circle
	gfx.SetColor(170.0, 0.0, 0.0, 255.0)
	gfx.Circle("fill", 100.0, 500.0, 50.0)
	gfx.Arc("fill", 200.0, 500.0, 50.0, 0, math.Pi)
	gfx.Ellipse("fill", 300.0, 500.0, 50.0, 20.0)
	gfx.Circle("line", 100.0, 600.0, 50.0)
	gfx.Arc("line", 200.0, 550.0, 50.0, 0, math.Pi)
	gfx.Ellipse("line", 300.0, 550.0, 50.0, 20.0)

	// line
	gfx.SetColor(0.0, 0.0, 170.0, 255.0)
	gfx.Line(
		800.0, 100.0, 850.0, 100.0,
		900.0, 20.0, 950.0, 100.0,
		1030.0, 100.0, 950.0, 180.0,
	)

	// image
	gfx.SetColor(255.0, 255.0, 255.0, 255.0)
	gfx.DrawS(tree, 500, 100)
}
