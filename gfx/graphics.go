package gfx

import (
	"image"
	"math"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type (
	Drawable interface {
		Draw(args ...float32)
	}
	QuadDrawable interface {
		Drawq(quad *Quad, args ...float32)
	}
)

const defaultPointCount = 30

func Circle(mode string, x, y, radius float32) {
	Circlep(mode, x, y, radius, defaultPointCount)
}

func Circlep(mode string, x, y, radius float32, points int) {
	Ellipsep(mode, x, y, radius, radius, points)
}

func Arc(mode string, x, y, radius, angle1, angle2 float32) {
	Arcp(mode, x, y, radius, angle1, angle2, defaultPointCount)
}

func Arcp(mode string, x, y, radius, angle1, angle2 float32, points int) {
	// Nothing to display with no points or equal angles. (Or is there with line mode?)
	if points <= 0 || angle1 == angle2 {
		return
	}

	// Oh, you want to draw a circle?
	if math.Abs(float64(angle1-angle2)) >= (2.0 * math.Pi) {
		Circlep(mode, x, y, radius, points)
		return
	}

	angle_shift := (angle2 - angle1) / float32(points)
	// Bail on precision issues.
	if angle_shift == 0.0 {
		return
	}

	phi := angle1
	num_coords := (points + 3) * 2
	coords := make([]float32, num_coords)
	coords[0] = x
	coords[num_coords-2] = x
	coords[1] = y
	coords[num_coords-1] = y

	for i := 0; i <= points; i++ {
		phi = phi + angle_shift
		coords[2*(i+1)] = x + radius*float32(math.Cos(float64(phi)))
		coords[2*(i+1)+1] = y + radius*float32(math.Sin(float64(phi)))
	}

	Polygon(mode, coords)
}

func Ellipse(mode string, x, y, a, b float32) {
	Ellipsep(mode, x, y, a, b, defaultPointCount)
}

func Ellipsep(mode string, x, y, a, b float32, points int) {
	two_pi := math.Pi * 2.0
	if points <= 0 {
		points = 1
	}

	angle_shift := float32(two_pi) / float32(points)
	phi := float32(0.0)

	coords := make([]float32, 2*(points+1))
	for i := 0; i < points; i++ {
		phi += angle_shift
		coords[2*i+0] = x + a*float32(math.Cos(float64(phi)))
		coords[2*i+1] = y + b*float32(math.Sin(float64(phi)))
	}

	coords[2*points+0] = coords[0]
	coords[2*points+1] = coords[1]

	Polygon(mode, coords)
}

func Point(x, y float32) {
	prepareDraw(nil)
	bindTexture(gl_state.defaultTexture)
	useVertexAttribArrays(ATTRIBFLAG_POS)
	gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 0, gl.Ptr([]float32{x, y}))
	gl.DrawArrays(gl.POINTS, 0, 1)
}

func Line(args ...float32) {
	PolyLine(args)
}

func PolyLine(coords []float32) {
	polyline := newPolyLine(states.back().line_join, states.back().line_style, states.back().line_width, states.back().pixelSize)
	polyline.render(coords)
}

func Rect(mode string, x, y, w, h float32) {
	Polygon(mode, []float32{x, y, x, y + h, x + w, y + h, x + w, y, x, y})
}

func Polygon(mode string, coords []float32) {
	if mode == "line" {
		PolyLine(coords)
	} else {
		prepareDraw(nil)
		bindTexture(gl_state.defaultTexture)
		useVertexAttribArrays(ATTRIBFLAG_POS)
		gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 0, gl.Ptr(coords))
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, int32(len(coords))/2-1)
	}
}

func NewScreenshot() image.Image {
	// Temporarily unbind the currently active canvas (glReadPixels reads the active framebuffer, not the main one.)
	canvases := GetCanvas()
	SetCanvas()

	w, h := int32(1), int32(1)
	screenshot := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	stride := int32(screenshot.Stride)
	pixels := make([]byte, len(screenshot.Pix))
	gl.ReadPixels(0, 0, w, h, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	// OpenGL sucks and reads pixels from the lower-left. Let's fix that.
	for y := int32(0); y < h; y++ {
		i := (h - 1 - y) * stride
		copy(screenshot.Pix[y*stride:], pixels[i:i+w*4])
	}

	// Re-bind the active canvas, if necessary.
	SetCanvas(canvases...)

	return screenshot
}

func Draw(drawable Drawable, args ...float32) {
	drawable.Draw(args...)
}

func Drawq(drawable QuadDrawable, quad *Quad, args ...float32) {
	drawable.Drawq(quad, args...)
}

func normalizeDrawCallArgs(args []float32) (float32, float32, float32, float32, float32, float32, float32, float32, float32) {
	/**
	 * Normalized an array of floats into these params if they exist
	 * if they are not present then thier default values are returned
	 * x The position of the object along the x-axis.
	 * y The position of the object along the y-axis.
	 * angle The angle of the object (in radians).
	 * sx The scale factor along the x-axis.
	 * sy The scale factor along the y-axis.
	 * ox The origin offset along the x-axis.
	 * oy The origin offset along the y-axis.
	 * kx Shear along the x-axis.
	 * ky Shear along the y-axis.
	 **/
	var x, y, angle, sx, sy, ox, oy, kx, ky float32
	sx = 1
	sy = 1

	if args == nil || len(args) < 2 {
		return x, y, angle, sx, sy, ox, oy, kx, ky
	}

	args_length := len(args)

	switch args_length {
	case 9:
		ky = args[8]
		fallthrough
	case 8:
		kx = args[7]
		if args_length == 8 {
			ky = kx
		}
		fallthrough
	case 7:
		oy = args[6]
		fallthrough
	case 6:
		ox = args[5]
		if args_length == 6 {
			oy = ox
		}
		fallthrough
	case 5:
		sy = args[4]
		fallthrough
	case 4:
		sx = args[3]
		if args_length == 4 {
			sy = sx
		}
		fallthrough
	case 3:
		angle = args[2]
		fallthrough
	case 2:
		x = args[0]
		y = args[1]
	}

	return x, y, angle, sx, sy, ox, oy, kx, ky
}

//did manually to save computation
func generateModelMatFromArgs(args []float32) *mgl32.Mat4 {
	x, y, angle, sx, sy, ox, oy, kx, ky := normalizeDrawCallArgs(args)
	mat := mgl32.Ident4()
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))
	// matrix multiplication carried out on paper:
	// |1     x| |c -s    | |sx       | | 1 ky    | |1     -ox|
	// |  1   y| |s  c    | |   sy    | |kx  1    | |  1   -oy|
	// |    1  | |     1  | |      1  | |      1  | |    1    |
	// |      1| |       1| |        1| |        1| |       1 |
	//   move      rotate      scale       skew       origin
	mat[10] = 1
	mat[15] = 1
	mat[0] = c*sx - ky*s*sy // = a
	mat[1] = s*sx + ky*c*sy // = b
	mat[4] = kx*c*sx - s*sy // = c
	mat[5] = kx*s*sx + c*sy // = d
	mat[12] = x - ox*mat[0] - oy*mat[4]
	mat[13] = y - ox*mat[1] - oy*mat[5]

	return &mat
}
