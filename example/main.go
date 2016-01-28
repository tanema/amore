package main

import (
	"fmt"
	"image/png"
	"math"
	"os"

	"github.com/tanema/amore"
	"github.com/tanema/amore/audio"
	"github.com/tanema/amore/gfx"
	_ "github.com/tanema/amore/joystick"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/mouse"
	"github.com/tanema/amore/timer"
	"github.com/tanema/amore/window"
)

var (
	tree       *gfx.Image
	ttf        gfx.Font
	image_font gfx.Font
	mx, my     float32
	shader     *gfx.Shader
	bomb       *audio.Source
	use_shader = false
	vibrating  = false
	canvas     *gfx.Canvas
	quad       *gfx.Quad
)

func main() {
	window.GetCurrent().SetMouseVisible(false)
	keyboard.SetKeyReleaseCB(keyUp)

	canvas = gfx.NewCanvas(800, 600)
	tree = gfx.NewImage("images/palm_tree.png")
	quad = gfx.NewQuad(0, 0, 200, 200, tree.GetWidth(), tree.GetHeight())
	ttf = gfx.NewFont("fonts/arial.ttf", 20)
	image_font = gfx.NewImageFont("fonts/image_font.png", " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,!?-+/():;%&`'*#=[]\"")
	shader = gfx.NewShader("shaders/blackandwhite.glsl")
	bomb, _ = audio.NewStreamSource("audio/bomb.wav")

	if err := amore.Start(update, draw); err != nil {
		fmt.Println("Error starting engine: %v", err)
	}
}

func keyUp(key keyboard.Key) {
	switch key {
	case keyboard.KeyEscape:
		amore.Quit()
	case keyboard.Key1:
		use_shader = !use_shader
	case keyboard.Key2:
		bomb.Play()
	case keyboard.Key3:
		img := gfx.NewScreenshot()
		out, _ := os.Create("./output.png")
		png.Encode(out, img)
	}
}

func update(deltaTime float32) {
	mx, my = mouse.GetPosition()
	mx, my = window.GetCurrent().PixelToWindowCoords(mx, my)
}

func draw() {
	gfx.SetLineWidth(10)
	if use_shader {
		gfx.SetShader(shader)
	} else {
		gfx.SetShader(nil)
	}

	// rectangle
	gfx.SetCanvas(canvas)
	gfx.SetColor(0, 170, 0, 155)
	gfx.Rect("fill", 20.0, 20.0, 200.0, 200.0)
	gfx.Rect("line", 250.0, 20.0, 200.0, 200.0)
	gfx.SetCanvas()
	gfx.SetColor(255, 255, 255, 255)
	gfx.Draw(canvas, 300, 100)

	// circle
	gfx.SetColor(170, 0, 0, 255)
	gfx.Circle("fill", 100.0, 500.0, 50.0)
	gfx.Arc("fill", 200.0, 500.0, 50.0, 0, math.Pi)
	gfx.Ellipse("fill", 300.0, 500.0, 50.0, 20.0)
	gfx.Circle("line", 100.0, 600.0, 50.0)
	gfx.Arc("line", 200.0, 550.0, 50.0, 0, math.Pi)
	gfx.Ellipse("line", 300.0, 550.0, 50.0, 20.0)

	//// line
	gfx.SetColor(0, 0, 170, 255)
	gfx.Line(
		800.0, 100.0, 850.0, 100.0,
		900.0, 20.0, 950.0, 100.0,
		1030.0, 100.0, 950.0, 180.0,
	)

	// image
	gfx.SetColor(255, 255, 255, 255)
	//x, y, rotate radians, scale x, y, offset x, y, shear x, y
	gfx.Draw(tree, 500, 50, -0.4, 0.5, 0.8, -100, -200, -0.2, 0.4)
	gfx.Drawq(tree, quad, 1000, 500)

	// image font
	gfx.SetFont(image_font)
	gfx.Printf(20, 20, "test one two")
	// ttf font
	gfx.SetFont(ttf)
	gfx.Rotate(0.5)
	gfx.Scale(1.5)
	gfx.Printf(200, 100, "test one two")
	gfx.Origin()

	//FPS
	gfx.SetColor(0, 170, 170, 255)
	gfx.Printf(1200, 10, "fps: %v", timer.GetFPS())

	//mouse position
	gfx.Circle("fill", mx, my, 20.0)
}
