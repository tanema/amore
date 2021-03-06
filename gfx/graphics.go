// Package gfx is used largly to simplify OpenGL calls and to manage state
// of transformations. Anything meant to be drawn to screen will come from this
// pacakge.
package gfx

import (
	"image"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/goxjs/gl"
)

// this is the default amount of points to allow a circle or arc to use when
// generating points
//const defaultPointCount = 30

// Circle will draw a circle at x, y with a radius as specified.
// points specifies how many points should be generated in the arc.
// If it is lower it will look jagged. If it is higher it will hit performace.
// The drawmode specifies either a fill or line draw
func Circle(mode string, x, y, radius float32, points int) {
	Ellipse(mode, x, y, radius, radius, points)
}

// Arc is like Arc except that you can define how many points you want to generate
// the arc.
// If it is lower it will look jagged. If it is higher it will hit performace.
// The drawmode specifies either a fill or line draw
func Arc(mode string, x, y, radius, angle1, angle2 float32, points int) {
	// Nothing to display with no points or equal angles. (Or is there with line mode?)
	if points <= 0 || angle1 == angle2 {
		return
	}

	// Oh, you want to draw a circle?
	if math.Abs(float64(angle1-angle2)) >= (2.0 * math.Pi) {
		Circle(mode, x, y, radius, points)
		return
	}

	angleShift := (angle2 - angle1) / float32(points)
	// Bail on precision issues.
	if angleShift == 0.0 {
		return
	}

	phi := angle1
	numCoords := (points + 3) * 2
	coords := make([]float32, numCoords)
	coords[0] = x
	coords[numCoords-2] = x
	coords[1] = y
	coords[numCoords-1] = y

	for i := 0; i <= points; i++ {
		phi = phi + angleShift
		coords[2*(i+1)] = x + radius*float32(math.Cos(float64(phi)))
		coords[2*(i+1)+1] = y + radius*float32(math.Sin(float64(phi)))
	}

	if mode == "line" {
		PolyLine(coords)
	} else {
		prepareDraw(nil)
		bindTexture(glState.defaultTexture)
		useVertexAttribArrays(shaderPos)

		buffer := newVertexBuffer(len(coords), coords, UsageStatic)
		buffer.bind()
		defer buffer.unbind()

		gl.VertexAttribPointer(shaderPos, 2, gl.FLOAT, false, 0, 0)
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, len(coords)/2-1)
	}
}

// Ellipse will draw a circle at x, y with a radius as specified.
// radiusx and radiusy will specify how much the width will be along those axis
// points specifies how many points should be generated in the arc.
// If it is lower it will look jagged. If it is higher it will hit performace.
// The drawmode specifies either a fill or line draw
func Ellipse(mode string, x, y, radiusx, radiusy float32, points int) {
	twoPi := math.Pi * 2.0
	if points <= 0 {
		points = 1
	}

	angleShift := float32(twoPi) / float32(points)
	phi := float32(0.0)

	coords := make([]float32, 2*(points+1))
	for i := 0; i < points; i++ {
		phi += angleShift
		coords[2*i+0] = x + radiusx*float32(math.Cos(float64(phi)))
		coords[2*i+1] = y + radiusy*float32(math.Sin(float64(phi)))
	}

	coords[2*points+0] = coords[0]
	coords[2*points+1] = coords[1]

	if mode == "line" {
		PolyLine(coords)
	} else {
		prepareDraw(nil)
		bindTexture(glState.defaultTexture)
		useVertexAttribArrays(shaderPos)

		buffer := newVertexBuffer(len(coords), coords, UsageStatic)
		buffer.bind()
		defer buffer.unbind()

		gl.VertexAttribPointer(shaderPos, 2, gl.FLOAT, false, 0, 0)
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, len(coords)/2-1)
	}
}

// Points will draw a point on the screen at x, y position. The size of the point
// is dependant on the point size set with SetPointSize.
func Points(coords []float32) {
	prepareDraw(nil)
	bindTexture(glState.defaultTexture)
	useVertexAttribArrays(shaderPos)

	buffer := newVertexBuffer(len(coords), coords, UsageStatic)
	buffer.bind()
	defer buffer.unbind()

	gl.VertexAttribPointer(shaderPos, 2, gl.FLOAT, false, 0, 0)
	gl.DrawArrays(gl.POINTS, 0, len(coords)/2)
}

// PolyLine will draw a line with an array in the form of x1, y1, x2, y2, x3, y3, ..... xn, yn
func PolyLine(coords []float32) {
	polyline := newPolyLine(states.back().lineJoin, states.back().lineWidth)
	polyline.render(coords)
}

// Rect draws a rectangle with the top left corner at x, y with the specified width
// and height
// The drawmode specifies either a fill or line draw
func Rect(mode string, x, y, width, height float32) {
	Polygon(mode, []float32{x, y, x, y + height, x + width, y + height, x + width, y})
}

// Polygon will draw a closed polygon with an array in the form of x1, y1, x2, y2, x3, y3, ..... xn, yn
// The drawmode specifies either a fill or line draw
func Polygon(mode string, coords []float32) {
	coords = append(coords, coords[0], coords[1])
	if mode == "line" {
		PolyLine(coords)
	} else {
		prepareDraw(nil)
		bindTexture(glState.defaultTexture)
		useVertexAttribArrays(shaderPos)

		buffer := newVertexBuffer(len(coords), coords, UsageStatic)
		buffer.bind()
		defer buffer.unbind()

		gl.VertexAttribPointer(shaderPos, 2, gl.FLOAT, false, 0, 0)
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, len(coords)/2-1)
	}
}

// NewScreenshot will take a screenshot of the screen and convert it to an image.Image
func NewScreenshot() *Image {
	// Temporarily unbind the currently active canvas (glReadPixels reads the active framebuffer, not the main one.)
	canvas := GetCanvas()
	SetCanvas(nil)

	w, h := int32(screenWidth), int32(screenHeight)
	screenshot := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	stride := int32(screenshot.Stride)
	pixels := make([]byte, len(screenshot.Pix))
	gl.ReadPixels(pixels, 0, 0, int(w), int(h), gl.RGBA, gl.UNSIGNED_BYTE)

	// OpenGL sucks and reads pixels from the lower-left. Let's fix that.
	for y := int32(0); y < h; y++ {
		i := (h - 1 - y) * stride
		copy(screenshot.Pix[y*stride:], pixels[i:i+w*4])
	}

	// Re-bind the active canvas, if necessary.
	SetCanvas(canvas)
	newImage := &Image{Texture: newImageTexture(screenshot, false)}
	registerVolatile(newImage)
	return newImage
}

// Normalized an array of floats into these params if they exist
// if they are not present then thier default values are returned
// x The position of the object along the x-axis.
// y The position of the object along the y-axis.
// angle The angle of the object (in radians).
// sx The scale factor along the x-axis.
// sy The scale factor along the y-axis.
// ox The origin offset along the x-axis.
// oy The origin offset along the y-axis.
// kx Shear along the x-axis.
// ky Shear along the y-axis.
func normalizeDrawCallArgs(args []float32) (float32, float32, float32, float32, float32, float32, float32, float32, float32) {
	var x, y, angle, sx, sy, ox, oy, kx, ky float32
	sx = 1
	sy = 1

	if args == nil || len(args) < 2 {
		return x, y, angle, sx, sy, ox, oy, kx, ky
	}

	argsLength := len(args)

	switch argsLength {
	case 9:
		ky = args[8]
		fallthrough
	case 8:
		kx = args[7]
		if argsLength == 8 {
			ky = kx
		}
		fallthrough
	case 7:
		oy = args[6]
		fallthrough
	case 6:
		ox = args[5]
		if argsLength == 6 {
			oy = ox
		}
		fallthrough
	case 5:
		sy = args[4]
		fallthrough
	case 4:
		sx = args[3]
		if argsLength == 4 {
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

// generateModelMatFromArgs will take in the arguments
// x, y, r, sx, sy, ox, oy, kx, ky
// and generate a matrix to be applied to the model transformation.
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

func f32Bytes(values []float32) []byte {
	b := make([]byte, 4*len(values))
	for i, v := range values {
		u := math.Float32bits(v)
		b[4*i+0] = byte(u >> 0)
		b[4*i+1] = byte(u >> 8)
		b[4*i+2] = byte(u >> 16)
		b[4*i+3] = byte(u >> 24)
	}
	return b
}

func ui32Bytes(values []uint32) []byte {
	b := make([]byte, 4*len(values))
	for i, v := range values {
		b[4*i+0] = byte(v)
		b[4*i+1] = byte(v >> 8)
		b[4*i+2] = byte(v >> 16)
		b[4*i+3] = byte(v >> 24)
	}
	return b
}
