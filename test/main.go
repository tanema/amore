package main

import (
	"fmt"

	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
)

//var img *gfx.Image

func main() {
	//img, _ = gfx.NewImage("icon.png")

	amore.On("keyboard", "escape", "release", func(device, button, action string, modifiers []string) {
		amore.Quit()
	})

	amore.Start(update, draw)
}

func update(dt float32) {
}

func draw() {
	gfx.SetColor(1, 1, 1, 1)
	gfx.Print(fmt.Sprintf("fps: %v", amore.GetFPS()), 0, 0)
	gfx.Rect("fill", 300, 300, 480, 440)
	gfx.Draw(img, 300, 300)

	gfx.SetLineWidth(2)
	gfx.SetLineJoin(gfx.LineJoinBevel)
	gfx.SetColor(1, 0, 0, 1)
	gfx.Rect("line", 50, 50, 100, 100)

	gfx.SetColor(1, 1, 1, 1)
	gfx.Line(0, 0, 100, 100, 200, 100)
}
