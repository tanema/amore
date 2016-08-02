// The gfx Pacakge is used largly to simplify OpenGL calls and to manage state
// of transformations. Anything meant to be drawn to screen will come from this
// pacakge.
package gfx

import (
	"image"

	"github.com/tanema/amore/gfx/gl"
)

type (
	// Drawable interface defines all objects that can be drawn. Inputs are as follows
	// x, y, r, sx, sy, ox, oy, kx, ky
	// x, y are position
	// r is rotation
	// sx, sy is the scale, if sy is not given sy will equal sx
	// ox, oy are offset
	// kx, ky are the shear. If ky is not given ky will equal kx
	Drawable interface {
		Draw(args ...float32)
	}
	// QuadDrawable interface defines all objects that can be drawn with a quad.
	// Inputs are as follows
	// quad is the quad to crop the texture
	// x, y, r, sx, sy, ox, oy, kx, ky
	// x, y are position
	// r is rotation
	// sx, sy is the scale, if sy is not given sy will equal sx
	// ox, oy are offset
	// kx, ky are the shear. If ky is not given ky will equal kx
	QuadDrawable interface {
		Drawq(quad *Quad, args ...float32)
	}
)

// this is the default amount of points to allow a circle or arc to use when
// generating points
const defaultPointCount = 30

// Circle will draw a circle at x, y with a radius as specified.
// The drawmode specifies either a fill or line draw
func Circle(mode DrawMode, x, y, radius float32) {
	Circlep(mode, x, y, radius, defaultPointCount)
}

// Circlep will draw a circle at x, y with a radius as specified.
// points specifies how many points should be generated in the arc.
// If it is lower it will look jagged. If it is higher it will hit performace.
// The drawmode specifies either a fill or line draw
func Circlep(mode DrawMode, x, y, radius float32, points int) {
	Ellipsep(mode, x, y, radius, radius, points)
}

// Arc will draw a part of a circle at the point x, y with the radius provied.
// The arc will start at angle1 (radians) and end at angle2 (radians)
// The drawmode specifies either a fill or line draw
func Arc(mode DrawMode, x, y, radius, angle1, angle2 float32) {
	Arcp(mode, x, y, radius, angle1, angle2, defaultPointCount)
}

// Arcp is like Arc except that you can define how many points you want to generate
// the arc.
// If it is lower it will look jagged. If it is higher it will hit performace.
// The drawmode specifies either a fill or line draw
func Arcp(mode DrawMode, x, y, radius, angle1, angle2 float32, points int) {
	// Nothing to display with no points or equal angles. (Or is there with line mode?)
	if points <= 0 || angle1 == angle2 {
		return
	}

	// Oh, you want to draw a circle?
	if Abs(angle1-angle2) >= (2.0 * Pi) {
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
		coords[2*(i+1)] = x + radius*Cos(phi)
		coords[2*(i+1)+1] = y + radius*Sin(phi)
	}

	if mode == LINE {
		PolyLine(coords)
	} else {
		prepareDraw(nil)
		bindTexture(gl_state.defaultTexture)
		useVertexAttribArrays(attribflag_pos)
		gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 0, gl.Ptr(coords))
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, len(coords)/2-1)
	}
}

// Ellipse will draw a circle at x, y with a radius as specified.
// radiusx and radiusy will specify how much the width will be along those axis
// If it is lower it will look jagged. If it is higher it will hit performace.
// The drawmode specifies either a fill or line draw
func Ellipse(mode DrawMode, x, y, radiusx, radiusy float32) {
	Ellipsep(mode, x, y, radiusx, radiusy, defaultPointCount)
}

// Circlep will draw a circle at x, y with a radius as specified.
// radiusx and radiusy will specify how much the width will be along those axis
// points specifies how many points should be generated in the arc.
// If it is lower it will look jagged. If it is higher it will hit performace.
// The drawmode specifies either a fill or line draw
func Ellipsep(mode DrawMode, x, y, radiusx, radiusy float32, points int) {
	two_pi := Pi * 2.0
	if points <= 0 {
		points = 1
	}

	angle_shift := float32(two_pi) / float32(points)
	phi := float32(0.0)

	coords := make([]float32, 2*(points+1))
	for i := 0; i < points; i++ {
		phi += angle_shift
		coords[2*i+0] = x + radiusx*Cos(phi)
		coords[2*i+1] = y + radiusy*Sin(phi)
	}

	coords[2*points+0] = coords[0]
	coords[2*points+1] = coords[1]

	if mode == LINE {
		PolyLine(coords)
	} else {
		prepareDraw(nil)
		bindTexture(gl_state.defaultTexture)
		useVertexAttribArrays(attribflag_pos)
		gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 0, gl.Ptr(coords))
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, len(coords)/2-1)
	}
}

// Point will draw a point on the screen at x, y position. The size of the point
// is dependant on the point size set with SetPointSize.
func Point(x, y float32) {
	prepareDraw(nil)
	bindTexture(gl_state.defaultTexture)
	useVertexAttribArrays(attribflag_pos)
	gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 0, gl.Ptr([]float32{x, y}))
	gl.DrawArrays(gl.POINTS, 0, 1)
}

// Line is a short form of Polyline so you can enter your params not in an array
func Line(args ...float32) {
	PolyLine(args)
}

// PolyLine will draw a line with an array in the form of x1, y1, x2, y2, x3, y3, ..... xn, yn
func PolyLine(coords []float32) {
	polyline := newPolyLine(states.back().line_join, states.back().line_style, states.back().line_width, states.back().pixelSize)
	polyline.render(coords)
}

// Rect draws a rectangle with the top left corner at x, y with the specified width
// and height
// The drawmode specifies either a fill or line draw
func Rect(mode DrawMode, x, y, width, height float32) {
	Polygon(mode, []float32{x, y, x, y + height, x + width, y + height, x + width, y})
}

// Polygon will draw a closed polygon with an array in the form of x1, y1, x2, y2, x3, y3, ..... xn, yn
// The drawmode specifies either a fill or line draw
func Polygon(mode DrawMode, coords []float32) {
	coords = append(coords, coords[0], coords[1])
	if mode == LINE {
		PolyLine(coords)
	} else {
		prepareDraw(nil)
		bindTexture(gl_state.defaultTexture)
		useVertexAttribArrays(attribflag_pos)
		gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 0, gl.Ptr(coords))
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, len(coords)/2-1)
	}
}

// NewScreenshot will take a screenshot of the screen and convert it to an image.Image
func NewScreenshot() image.Image {
	// Temporarily unbind the currently active canvas (glReadPixels reads the active framebuffer, not the main one.)
	canvases := GetCanvas()
	SetCanvas()

	w, h := int32(screen_width), int32(screen_height)
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
	SetCanvas(canvases...)

	return screenshot
}

// Draw calls draw on any drawable object with the inputs
// Inputs are as follows
// x, y, r, sx, sy, ox, oy, kx, ky
// x, y are position
// r is rotation
// sx, sy is the scale, if sy is not given sy will equal sx
// ox, oy are offset
// kx, ky are the shear. If ky is not given ky will equal kx
func Draw(drawable Drawable, args ...float32) {
	drawable.Draw(args...)
}

// Drawq will draw any object that can have a quad applied to it
// Inputs are as follows
// quad is the quad to crop the texture
// x, y, r, sx, sy, ox, oy, kx, ky
// x, y are position
// r is rotation
// sx, sy is the scale, if sy is not given sy will equal sx
// ox, oy are offset
// kx, ky are the shear. If ky is not given ky will equal kx
func Drawq(drawable QuadDrawable, quad *Quad, args ...float32) {
	drawable.Drawq(quad, args...)
}
