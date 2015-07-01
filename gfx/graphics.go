package gfx

import (
	"github.com/go-gl/gl/v2.1/gl"
)

func Clear(c Color) {
	gl.ClearColor(c.R/255.0, c.G/255.0, c.B/255.0, c.A/255.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func SetColor(c Color) {
	SetColorf(c.R, c.G, c.B, c.A)
}

func SetColorf(r, g, b, a float32) {
	gl.Color4f(r/255.0, g/255.0, b/255.0, a/255.0)
}

func Line(x1, y1, x2, y2 float32) {
	PolyLine([]float32{x1, y1, x2, y2})
}

func PolyLine(coords []float32) {
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, gl.Ptr(coords))
	gl.DrawArrays(gl.LINE_STRIP, 0, int32(len(coords))/2)
	gl.DisableVertexAttribArray(0)
}

func Rect(mode string, x, y, w, h float32) {
	Polygon(mode, []float32{x, y, x, y + h, x + w, y + h, x + w, y, x, y})
}

func Polygon(mode string, coords []float32) {
	if mode == "line" {
		PolyLine(coords)
	} else {
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, gl.Ptr(coords))
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, int32(len(coords))/2-1)
		gl.DisableVertexAttribArray(0)
	}
}

func DrawS(drawable Drawable, x, y float32) {
	drawable.Draw(x, y, 0, 0, 0, 0, 0, 0, 0)
}

func Draw(drawable Drawable, x, y, angle, sx, sy, ox, oy, kx, ky float32) {
	drawable.Draw(x, y, angle, sx, sy, ox, oy, kx, ky)
}
