package main

import (
	"fmt"
	"math"

	"github.com/tanema/amore"
	"github.com/tanema/amore/audio"
	"github.com/tanema/amore/file"
	"github.com/tanema/amore/gfx"
	_ "github.com/tanema/amore/joystick"
	"github.com/tanema/amore/keyboard"
	"github.com/tanema/amore/mouse"
	"github.com/tanema/amore/timer"
	"github.com/tanema/amore/window"
)

var (
	tree          *gfx.Image
	ttf           *gfx.Font
	image_font    *gfx.Font
	mx, my        float32
	shader        *gfx.Shader
	bomb          *audio.Source
	use_shader    = false
	canvas        *gfx.Canvas
	quad          *gfx.Quad
	psystem       *gfx.ParticleSystem
	triangle_mesh *gfx.Mesh
	batch         *gfx.SpriteBatch
	text          *gfx.Text
	amore_text    *gfx.Text
	star          = []float32{
		133, 30, 198, 82, 259, 86, 197, 163,
		235, 243, 140, 184, 60, 243, 85, 158,
		34, 78, 113, 77, 132, 30,
	}
	mouseColor = gfx.NewColor(255, 255, 255, 255)
)

func main() {
	window.SetMouseVisible(false)
	keyboard.SetKeyReleaseCB(keyUp)
	mouse.SetButtonReleaseCB(mouseButtonUp)

	canvas = gfx.NewCanvas(800, 600)
	tree, _ = gfx.NewImage("images/palm_tree.png")
	quad = gfx.NewQuad(0, 0, 200, 200, tree.GetWidth(), tree.GetHeight())
	ttf = gfx.NewFont("fonts/arialbd.ttf", 20)
	image_font = gfx.NewImageFont("fonts/image_font.png", " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,!?-+/():;%&`'*#=[]\"")
	image_font.SetFallbacks(ttf)
	shader, _ = gfx.NewShader("shaders/blackandwhite.glsl")
	bomb, _ = audio.NewSource("audio/bomb.wav")
	bomb.SetLooping(true)
	text, _ = gfx.NewColorTextExt(ttf,
		[]string{file.ReadString("text/lorem.txt"), file.ReadString("text/lorem.txt")},
		[]gfx.Color{gfx.NewColor(255, 255, 255, 255), gfx.NewColor(255, 0, 255, 255)},
		500, gfx.ALIGN_CENTER)
	amore_text, _ = gfx.NewColorText(ttf, []string{"a", "m", "o", "r", "e"},
		[]gfx.Color{
			gfx.NewColor(0, 255, 0, 255),
			gfx.NewColor(255, 0, 255, 255),
			gfx.NewColor(255, 255, 0, 255),
			gfx.NewColor(0, 0, 255, 255),
			gfx.NewColor(255, 255, 255, 255),
		})

	particle, _ := gfx.NewImage("images/particle.png")
	psystem, _ = gfx.NewParticleSystem(particle, 10)
	psystem.SetParticleLifetime(2, 5) // Particles live at least 2s and at most 5s.
	psystem.SetEmissionRate(5)
	psystem.SetSizeVariation(1)
	psystem.SetLinearAcceleration(-20, -20, 20, 20) // Random movement in all directions.
	psystem.SetSpeed(3, 5)
	psystem.SetSpin(0.1, 0.5)
	psystem.SetSpinVariation(1)

	triangle_mesh, _ = gfx.NewMesh([]float32{
		25, 0,
		0, 0,
		255, 0, 0, 255,

		50, 50,
		1, 0,
		255, 0, 0, 255,

		0, 50,
		1, 1,
		255, 0, 0, 255,

		0, 0,
		0, 1,
		255, 0, 0, 255,
	}, 4)
	triangle_mesh.SetTexture(tree)

	q := gfx.NewQuad(50, 50, 50, 50, tree.GetWidth(), tree.GetHeight())
	q2 := gfx.NewQuad(100, 50, 50, 50, tree.GetWidth(), tree.GetHeight())
	batch = gfx.NewSpriteBatch(tree, 4)
	batch.Addq(q, 0, 0)
	batch.Addq(q2, 50, 0)
	batch.Addq(q, 50, 50)
	batch.Addq(q2, 0, 50)
	batch.Addq(q, 100, 50)

	amore.Start(update, draw)
}

func mouseButtonUp(x, y float32, button mouse.MouseButton) {
	if button == mouse.LeftButton {
		mouseColor = gfx.NewColor(255, 0, 0, 255)
	} else if button == mouse.RightButton {
		mouseColor = gfx.NewColor(0, 255, 0, 255)
	}
}

func keyUp(key keyboard.Key) {
	switch key {
	case keyboard.KeyEscape:
		amore.Quit()
	case keyboard.Key1:
		use_shader = !use_shader
	case keyboard.Key2:
		if bomb.IsPlaying() {
			bomb.Stop()
		} else {
			bomb.Play()
		}
	case keyboard.Key3:
		println(window.Confirm("test", "alert"))
	case keyboard.Key4:
		triangle_mesh.SetVertexMap([]uint16{0, 1, 2})
	case keyboard.Key5:
		triangle_mesh.ClearVertexMap()
	case keyboard.Key6:
		batch.SetDrawRange(2, 4)
	case keyboard.Key7:
		batch.ClearDrawRange()
	}
}

func update(deltaTime float32) {
	mx, my = mouse.GetPosition()
	mx, my = window.PixelToWindowCoords(mx, my)
	psystem.Update(deltaTime)
}

func draw() {
	if use_shader {
		gfx.SetShader(shader)
	} else {
		gfx.SetShader(nil)
	}

	gfx.Point(200, 20)

	gfx.SetLineWidth(1)
	//text
	gfx.SetColor(255, 255, 255, 255)
	gfx.Draw(text, 0, 300)
	gfx.Rect("line", 0, 300, 500, text.GetHeight())

	gfx.SetLineWidth(10)

	// line
	gfx.SetColor(0, 0, 170, 255)
	gfx.Line(star...)

	gfx.SetColor(255, 255, 255, 255)
	gfx.Draw(batch, 50, 150)

	//stencil
	gfx.Stencil(func() { gfx.Rect("fill", 426, 240, 426, 240) })
	gfx.SetStencilTest(gfx.COMPARE_EQUAL, 0)
	gfx.SetColor(239, 96, 17, 255)
	gfx.Rect("fill", 400, 200, 826, 440)
	gfx.ClearStencilTest()

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

	// image
	gfx.SetColor(255, 255, 255, 255)
	//x, y, rotate radians, scale x, y, offset x, y, shear x, y
	gfx.Draw(tree, 500, 50, -0.4, 0.5, 0.8, -100, -200, -0.2, 0.4)
	gfx.Drawq(tree, quad, 400, 50)

	// image font
	gfx.SetFont(image_font)
	gfx.Printf("test one @ two", 150, gfx.ALIGN_JUSTIFY, 0, 0)
	// ttf font
	gfx.SetFont(ttf)
	gfx.Print("test one two", 200, 100, math.Pi/2, 2, 2)

	gfx.Draw(psystem, 200, 200)

	gfx.SetColor(255, 255, 255, 255)
	gfx.Draw(triangle_mesh, 200, 200)

	//mouse position
	gfx.SetColorC(mouseColor)
	gfx.Circle("fill", mx, my, 20.0)

	//FPS
	gfx.SetColor(0, 170, 170, 255)
	gfx.Print(fmt.Sprintf("fps: %v", timer.GetFPS()), 720, 10)

	gfx.Draw(amore_text, 500, 400, 0, 3, 3)
}
