package gfx

import (
	"github.com/go-gl/gl/v2.1/gl"
)

func Translate(x, y float64) {
	gl.Translated(x, y, 0)
}

func Reset() {
	Origin()
	SetBlendMode("alpha")
	Clear(0.0, 0.0, 0.0, 0.0)
}

func Origin() {
	//reset transforms
	gl.LoadIdentity()
	//set our coord system to flow form top left
	gl.Ortho(0, float64(width), float64(height), 0, -1, 1)
}

func Clear(r, g, b, a float32) {
	gl.ClearColor(r/255.0, g/255.0, b/255.0, a/255.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func Rotate(angle float64) {
	gl.Rotated(angle, 0, 0, 0)
}

func RotateAround(angle, x, y float64) {
	gl.Rotated(angle, x, y, 0)
}

func Scale(sx float64) {
	gl.Scaled(sx, sx, 0)
}

func Scale2(sx, sy float64) {
	gl.Scaled(sx, sy, 0)
}

func Push() {
	gl.PushMatrix()
}

func Pop() {
	gl.PopMatrix()
}

func SetScissor(x, y, width, height int32) {
	gl.Scissor(x, y, width, height)
}

func SetShader(shader *Shader) {
	shader.Use()
}

func ClearShader() {
	gl.UseProgram(0)
}

func SetColor(r, g, b, a int) {
	gl.Color4d(float64(r)/255.0, float64(g)/255.0, float64(b)/255.0, float64(a)/255.0)
}

func SetLineWidth(width float32) {
	gl.LineWidth(width)
}

func SetBlendMode(mode string) {
	fn := gl.FUNC_ADD
	srcRGB := gl.ONE
	srcA := gl.ONE
	dstRGB := gl.ZERO
	dstA := gl.ZERO

	switch mode {
	case "alpha":
		srcRGB = gl.SRC_ALPHA
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_ALPHA
		dstA = gl.ONE_MINUS_SRC_ALPHA
	case "multiplicative":
		srcRGB = gl.DST_COLOR
		srcA = gl.DST_COLOR
		dstRGB = gl.ZERO
		dstA = gl.ZERO
	case "premultiplied":
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_ALPHA
		dstA = gl.ONE_MINUS_SRC_ALPHA
	case "subtractive":
		fn = gl.FUNC_REVERSE_SUBTRACT
	case "additive":
		srcRGB = gl.SRC_ALPHA
		srcA = gl.SRC_ALPHA
		dstRGB = gl.ONE
		dstA = gl.ONE
	case "screen":
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_COLOR
		dstA = gl.ONE_MINUS_SRC_COLOR
		break
	case "replace":
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ZERO
		dstA = gl.ZERO
	}

	gl.BlendEquation(uint32(fn))
	gl.BlendFuncSeparate(uint32(srcRGB), uint32(dstRGB), uint32(srcA), uint32(dstA))
}
