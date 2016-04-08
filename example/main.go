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
	_ "github.com/tanema/amore/touch"
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
)

func main() {
	window.SetMouseVisible(false)
	keyboard.SetKeyReleaseCB(keyUp)

	canvas = gfx.NewCanvas(800, 600)
	tree, _ = gfx.NewImage("images/palm_tree.png")
	quad = gfx.NewQuad(0, 0, 200, 200, tree.GetWidth(), tree.GetHeight())
	ttf = gfx.NewFont("fonts/arialbd.ttf", 20)
	image_font = gfx.NewImageFont("fonts/image_font.png", " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,!?-+/():;%&`'*#=[]\"")
	image_font.SetFallbacks(ttf)
	shader = gfx.NewShader("shaders/blackandwhite.glsl")
	var er error
	bomb, er = audio.NewSource("audio/bomb.wav", true)
	if er != nil {
		panic(er)
	}
	bomb.SetLooping(true)
	text, _ = gfx.NewColorTextExt(ttf,
		[]string{file.ReadString("text/lorem.txt"), file.ReadString("text/lorem.txt")},
		[]*gfx.Color{gfx.NewColor(255, 255, 255, 255), gfx.NewColor(255, 0, 255, 255)},
		500, gfx.ALIGN_CENTER)
	amore_text, _ = gfx.NewColorText(ttf, []string{"a", "m", "o", "r", "e"},
		[]*gfx.Color{
			gfx.NewColor(0, 255, 0, 255),
			gfx.NewColor(255, 0, 255, 255),
			gfx.NewColor(255, 255, 0, 255),
			gfx.NewColor(0, 0, 255, 255),
			gfx.NewColor(255, 255, 255, 255),
		})

	particle, _ := gfx.NewImage("images/particle.png")
	psystem = gfx.NewParticleSystem(particle, 10)
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
		window.ShowMessageBox("title", "message", []string{"okay", "cancel"}, window.MESSAGEBOX_INFO, true)
	case keyboard.Key4:
		triangle_mesh.SetVertexMap([]uint32{0, 1, 2})
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
	gfx.Rect(gfx.LINE, 0, 300, 500, text.GetHeight())

	gfx.SetLineWidth(10)

	gfx.SetColor(255, 255, 255, 255)
	gfx.Draw(batch, 50, 150)

	//stencil
	gfx.Stencil(func() { gfx.Rect(gfx.FILL, 426, 240, 426, 240) })
	gfx.SetStencilTest(gfx.COMPARE_EQUAL, 0)
	gfx.SetColor(239, 96, 17, 255)
	gfx.Rect(gfx.FILL, 400, 200, 826, 440)
	gfx.ClearStencilTest()

	// rectangle
	gfx.SetCanvas(canvas)
	gfx.SetColor(0, 170, 0, 155)
	gfx.Rect(gfx.FILL, 20.0, 20.0, 200.0, 200.0)
	gfx.Rect(gfx.LINE, 250.0, 20.0, 200.0, 200.0)
	gfx.SetCanvas()
	gfx.SetColor(255, 255, 255, 255)
	gfx.Draw(canvas, 300, 100)

	// circle
	gfx.SetColor(170, 0, 0, 255)
	gfx.Circle(gfx.FILL, 100.0, 500.0, 50.0)
	gfx.Arc(gfx.FILL, 200.0, 500.0, 50.0, 0, math.Pi)
	gfx.Ellipse(gfx.FILL, 300.0, 500.0, 50.0, 20.0)
	gfx.Circle(gfx.LINE, 100.0, 600.0, 50.0)
	gfx.Arc(gfx.LINE, 200.0, 550.0, 50.0, 0, math.Pi)
	gfx.Ellipse(gfx.LINE, 300.0, 550.0, 50.0, 20.0)

	// line
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
	gfx.Printf("test one @ two", 150, gfx.ALIGN_JUSTIFY, 0, 0)
	// ttf font
	gfx.SetFont(ttf)
	gfx.Print("test one two", 200, 100, math.Pi/2, 2, 2)

	//FPS
	gfx.SetColor(0, 170, 170, 255)
	gfx.Print(fmt.Sprintf("fps: %v", timer.GetFPS()), 1200, 10)

	gfx.Draw(psystem, 200, 200)

	gfx.SetColor(255, 255, 255, 255)
	gfx.Draw(triangle_mesh, 50, 50)

	//mouse position
	gfx.Circle(gfx.FILL, mx, my, 20.0)

	gfx.Draw(amore_text, 500, 400, 0, 3, 3)
}
