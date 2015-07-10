package main

import (
	"math"

	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/mouse"
)

var (
	tree        *gfx.Image
	ttf         *gfx.Font
	image_font  *gfx.Font
	fps, mx, my float64
)

func main() {
	keyboard.SetKeyReleaseCB(keyUp)
	amore.Start(load, update, draw)
}

func keyUp(key keyboard.Key) {
	if key == keyboard.KeyEscape {
		amore.Quit()
	}
}

func load() {
	var err error
	tree, err = gfx.NewImage("assets/palm_tree.png")
	if err != nil {
		panic(err)
	}
	ttf, _ = gfx.NewFont("assets/fonts/arial.ttf", 20)
	image_font, _ = gfx.NewImageFont("assets/fonts/image_font.png", " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,!?-+/():;%&`'*#=[]\"")
}

func update(deltaTime float64) {
	fps = 1 / deltaTime
	mx, my = mouse.GetPosition()
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

	// font
	gfx.SetFont(image_font)
	gfx.Printf(20, 20, "test one two")
	gfx.SetFont(ttf)
	gfx.Printf(20, 100, "test one two")

	//FPS
	gfx.SetColor(0.0, 170.0, 170.0, 255.0)
	gfx.Printf(1200, 10, "fps: %v", int(fps))

	//mouse position
	gfx.Circle("fill", float32(mx), float32(my), 20.0)
}
