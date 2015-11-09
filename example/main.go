package main

import (
	"fmt"
	"math"

	"github.com/tanema/amore"
	_ "github.com/tanema/amore/audio"
	"github.com/tanema/amore/file"
	"github.com/tanema/amore/gfx"
	_ "github.com/tanema/amore/joystick"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/mouse"
	"github.com/tanema/amore/timer"
)

var (
	tree       *gfx.Image
	ttf        *gfx.Font
	image_font *gfx.Font
	mx, my     float32
	shader     *gfx.Shader
	use_shader = false
)

func main() {
	if err := amore.Start(load, update, draw); err != nil {
		fmt.Println("Error starting engine: %v", err)
	}
}

func load() {
	keyboard.SetKeyReleaseCB(keyUp)
	var err error
	tree, err = gfx.NewImage("assets/palm_tree.png")
	if err != nil {
		panic(err)
	}
	ttf, _ = gfx.NewFont("assets/fonts/arial.ttf", 20)
	image_font, _ = gfx.NewImageFont("assets/fonts/image_font.png", " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,!?-+/():;%&`'*#=[]\"")

	shader = gfx.NewShader(file.Read("shaders/blackandwhite.glsl"))
}

func keyUp(key keyboard.Key) {
	if key == keyboard.KeyEscape {
		amore.Quit()
	}
	if key == keyboard.Key1 {
		use_shader = !use_shader
	}
}

func update(deltaTime float32) {
	mx, my = mouse.GetPosition()
}

func draw() {
	if use_shader {
		gfx.SetShader(shader)
	} else {
		gfx.SetShader(nil)
	}

	// rectangle
	gfx.SetColor(0, 170, 0, 155)
	gfx.Rect("fill", 20.0, 20.0, 200.0, 200.0)
	gfx.Rect("line", 250.0, 20.0, 200.0, 200.0)

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

	// font
	gfx.SetFont(image_font)
	gfx.Printf(20, 20, "test one two")
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
