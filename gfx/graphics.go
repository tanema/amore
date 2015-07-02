package gfx

import (
	"github.com/go-gl/gl/v2.1/gl"
)

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
