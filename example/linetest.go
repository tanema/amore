package main

import (
	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/mouse"
)

var (
	mx, my float32
)

func main() { amore.Start(load, update, draw) }
func load() {}

func update(deltaTime float32) {
	mx, my = mouse.GetPosition()
}

func draw() {
	gfx.SetLineWidth(32)
	gfx.SetLineStyle(gfx.LINE_SMOOTH)
	gfx.SetLineJoin(gfx.LINE_JOIN_BEVEL)
	gfx.SetColor(0, 0, 170, 255)
	gfx.Line(100.0, 100.0, 200.0, 100.0, 200.0, 200.0, mx, my)
	gfx.Line(500.0, 100.0, 600.0, 100.0, 600.0, 200.0, 500.0, 200.0, 500.0, 100.0)
}
