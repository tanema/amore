package gfx

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/tanema/amore/gfx/gl"
)

type (
	polyLine struct {
		style      LineStyle
		join       LineJoin
		overdraw   bool
		halfwidth  float32
		pixel_size float32
		coord      []float32
		normals    []mgl32.Vec2
		vertices   []mgl32.Vec2
	}
)

func determinant(vec1, vec2 mgl32.Vec2) float32 {
	return vec1[0]*vec2[1] - vec1[1]*vec2[0]
}

func getNormal(v1 mgl32.Vec2, scale float32) mgl32.Vec2 {
	return mgl32.Vec2{-v1[1] * scale, v1[0] * scale}
}

func normalize(v1 mgl32.Vec2, length float32) mgl32.Vec2 {
	length_current := v1.Len()

	if length_current > 0 {
		v1 = v1.Mul(length / length_current)
	}

	return v1
}

func newPolyLine(join LineJoin, style LineStyle, line_width, pixel_size float32) polyLine {
	new_polyline := polyLine{
		style:      style,
		join:       join,
		overdraw:   style == LINE_SMOOTH,
		halfwidth:  line_width * 0.5,
		pixel_size: pixel_size,
	}
	return new_polyline
}

func (polyline *polyLine) render(coords []float32) {
	var sleeve, current, next mgl32.Vec2
	polyline.vertices = []mgl32.Vec2{}
	polyline.normals = []mgl32.Vec2{}

	coords_count := len(coords)
	is_looping := (coords[0] == coords[coords_count-2]) && (coords[1] == coords[coords_count-1])
	if !is_looping { // virtual starting point at second point mirrored on first point
		sleeve = mgl32.Vec2{coords[2] - coords[0], coords[3] - coords[1]}
	} else { // virtual starting point at last vertex
		sleeve = mgl32.Vec2{coords[0] - coords[coords_count-4], coords[1] - coords[coords_count-3]}
	}

	for i := 0; i+3 < coords_count; i += 2 {
		current = mgl32.Vec2{coords[i], coords[i+1]}
		next = mgl32.Vec2{coords[i+2], coords[i+3]}
		polyline.renderEdge(sleeve, current, next)
		sleeve = next.Sub(current)
	}

	if is_looping {
		polyline.renderEdge(sleeve, next, mgl32.Vec2{coords[2], coords[3]})
	} else {
		polyline.renderEdge(sleeve, next, next.Add(sleeve))
	}

	if polyline.join == LINE_JOIN_NONE {
		polyline.vertices = polyline.vertices[2 : len(polyline.vertices)-2]
	}

	polyline.draw(is_looping)
}

func (polyline *polyLine) renderEdge(sleeve, current, next mgl32.Vec2) {
	switch polyline.join {
	case LINE_JOIN_MITER:
		polyline.renderMiterEdge(sleeve, current, next)
	case LINE_JOIN_BEVEL:
		polyline.renderBevelEdge(sleeve, current, next)
	case LINE_JOIN_NONE:
		fallthrough
	default:
		polyline.renderNoEdge(sleeve, current, next)
	}
}

func (polyline *polyLine) generateEdges(current mgl32.Vec2, count int) {
	normal_count := len(polyline.normals)
	for i := count; i > 0; i-- {
		polyline.vertices = append(polyline.vertices, current.Add(polyline.normals[normal_count-i]))
	}
}

func (polyline *polyLine) renderNoEdge(sleeve, current, next mgl32.Vec2) {
	sleeve_normal := getNormal(sleeve, polyline.halfwidth/sleeve.Len())

	polyline.normals = append(polyline.normals, sleeve_normal)
	polyline.normals = append(polyline.normals, sleeve_normal.Mul(-1))

	sleeve = next.Sub(current)
	sleeve_normal = getNormal(sleeve, polyline.halfwidth/sleeve.Len())

	polyline.normals = append(polyline.normals, sleeve_normal.Mul(-1))
	polyline.normals = append(polyline.normals, sleeve_normal)

	polyline.generateEdges(current, 4)
}

/** Calculate line boundary points.
 *
 * Sketch:
 *
 *              u1
 * -------------+---...___
 *              |         ```'''--  ---
 * p- - - - - - q- - . _ _           | w/2
 *              |          ` ' ' r   +
 * -------------+---...___           | w/2
 *              u2         ```'''-- ---
 *
 * u1 and u2 depend on four things:
 *   - the half line width w/2
 *   - the previous line vertex p
 *   - the current line vertex q
 *   - the next line vertex r
 *
 * u1/u2 are the intersection points of the parallel lines to p-q and q-r,
 * i.e. the point where
 *
 *    (q + w/2 * ns) + lambda * (q - p) = (q + w/2 * nt) + mu * (r - q)   (u1)
 *    (q - w/2 * ns) + lambda * (q - p) = (q - w/2 * nt) + mu * (r - q)   (u2)
 *
 * with nt,nt being the normals on the segments s = p-q and t = q-r,
 *
 *    ns = perp(s) / |s|
 *    nt = perp(t) / |t|.
 *
 * Using the linear equation system (similar for u2)
 *
 *         q + w/2 * ns + lambda * s - (q + w/2 * nt + mu * t) = 0                 (u1)
 *    <=>  q-q + lambda * s - mu * t                          = (nt - ns) * w/2
 *    <=>  lambda * s   - mu * t                              = (nt - ns) * w/2
 *
 * the intersection points can be efficiently calculated using Cramer's rule.
 */
func (polyline *polyLine) renderMiterEdge(sleeve, current, next mgl32.Vec2) {
	sleeve_normal := getNormal(sleeve, polyline.halfwidth/sleeve.Len())
	t := next.Sub(current)
	len_t := t.Len()

	det := determinant(sleeve, t)
	// lines parallel, compute as u1 = q + ns * w/2, u2 = q - ns * w/2
	if mgl32.Abs(det)/(sleeve.Len()*len_t) < LINES_PARALLEL_EPS && sleeve.Dot(t) > 0 {
		polyline.normals = append(polyline.normals, sleeve_normal)
		polyline.normals = append(polyline.normals, sleeve_normal.Mul(-1))
	} else {
		// cramers rule
		nt := getNormal(t, polyline.halfwidth/len_t)
		lambda := determinant(nt.Sub(sleeve_normal), t) / det
		d := sleeve_normal.Add(sleeve.Mul(lambda))

		polyline.normals = append(polyline.normals, d)
		polyline.normals = append(polyline.normals, d.Mul(-1))
	}
	polyline.generateEdges(current, 2)
}

/** Calculate line boundary points.
 *
 * Sketch:
 *
 *     uh1___uh2
 *      .'   '.
 *    .'   q   '.
 *  .'   '   '   '.
 *.'   '  .'.  '   '.
 *   '  .' ul'.  '
 * p  .'       '.  r
 *
 *
 * ul can be found as above, uh1 and uh2 are much simpler:
 *
 * uh1 = q + ns * w/2, uh2 = q + nt * w/2
 */
func (polyline *polyLine) renderBevelEdge(sleeve, current, next mgl32.Vec2) {
	t := next.Sub(current)
	len_t := t.Len()

	det := determinant(sleeve, t)
	if mgl32.Abs(det)/(sleeve.Len()*len_t) < LINES_PARALLEL_EPS && sleeve.Dot(t) > 0 {
		// lines parallel, compute as u1 = q + ns * w/2, u2 = q - ns * w/2
		n := getNormal(t, polyline.halfwidth/len_t)
		polyline.normals = append(polyline.normals, n)
		polyline.normals = append(polyline.normals, n.Mul(-1))
		polyline.generateEdges(current, 2)
		return // early out
	}

	// cramers rule
	sleeve_normal := getNormal(sleeve, polyline.halfwidth/sleeve.Len())
	nt := getNormal(t, polyline.halfwidth/len_t)
	lambda := determinant(nt.Sub(sleeve_normal), t) / det
	d := sleeve_normal.Add(sleeve.Mul(lambda))

	if det > 0 { // 'left' turn -> intersection on the top
		polyline.normals = append(polyline.normals, d)
		polyline.normals = append(polyline.normals, sleeve_normal.Mul(-1))
		polyline.normals = append(polyline.normals, d)
		polyline.normals = append(polyline.normals, nt.Mul(-1))
	} else {
		polyline.normals = append(polyline.normals, sleeve_normal)
		polyline.normals = append(polyline.normals, d.Mul(-1))
		polyline.normals = append(polyline.normals, nt)
		polyline.normals = append(polyline.normals, d.Mul(-1))
	}
	polyline.generateEdges(current, 4)
}

func (polyline *polyLine) renderOverdraw(is_looping bool) []mgl32.Vec2 {
	switch polyline.join {
	case LINE_JOIN_NONE:
		return polyline.renderTrianglesOverdraw()
	case LINE_JOIN_MITER:
		fallthrough
	case LINE_JOIN_BEVEL:
		fallthrough
	default:
		return polyline.renderTriangleStripOverdraw(is_looping)
	}
}

func (polyline *polyLine) renderTriangleStripOverdraw(is_looping bool) []mgl32.Vec2 {
	overdraw_vertex_count := 2 * len(polyline.vertices)
	if !is_looping {
		overdraw_vertex_count += 2
	}
	overdraw := make([]mgl32.Vec2, overdraw_vertex_count)
	for i := 0; i+1 < len(polyline.vertices); i += 2 {
		// upper segment
		overdraw[i] = polyline.vertices[i]
		overdraw[i+1] = polyline.vertices[i].Add(polyline.normals[i].Mul(polyline.pixel_size / polyline.normals[i].Len()))
		// lower segment
		k := len(polyline.vertices) - i - 1
		overdraw[len(polyline.vertices)+i] = polyline.vertices[k]
		overdraw[len(polyline.vertices)+i+1] = polyline.vertices[k].Add(polyline.normals[k].Mul(polyline.pixel_size / polyline.normals[i].Len()))
	}

	// if not looping, the outer overdraw vertices need to be displaced
	// to cover the line endings, i.e.:
	// +- - - - //- - +         +- - - - - //- - - +
	// +-------//-----+         : +-------//-----+ :
	// | core // line |   -->   : | core // line | :
	// +-----//-------+         : +-----//-------+ :
	// +- - //- - - - +         +- - - //- - - - - +
	if !is_looping {
		// left edge
		spacer := overdraw[1].Sub(overdraw[3])
		spacer = normalize(spacer, polyline.pixel_size)
		overdraw[1] = overdraw[1].Add(spacer)
		overdraw[overdraw_vertex_count-3] = overdraw[overdraw_vertex_count-3].Add(spacer)

		// right edge
		spacer = overdraw[len(polyline.vertices)-1].Sub(overdraw[len(polyline.vertices)-3])
		spacer = normalize(spacer, polyline.pixel_size)
		overdraw[len(polyline.vertices)-1] = overdraw[len(polyline.vertices)-1].Add(spacer)
		overdraw[len(polyline.vertices)+1] = overdraw[len(polyline.vertices)+1].Add(spacer)

		// we need to draw two more triangles to close the
		// overdraw at the line start.
		overdraw[overdraw_vertex_count-2] = overdraw[0]
		overdraw[overdraw_vertex_count-1] = overdraw[1]
	}
	return overdraw
}

func (polyline *polyLine) renderTrianglesOverdraw() []mgl32.Vec2 {
	overdraw_vertex_count := 4 * (len(polyline.vertices) - 2) // less than ideal
	overdraw := make([]mgl32.Vec2, overdraw_vertex_count)
	for i := 2; i+3 < len(polyline.vertices); i += 4 {
		s := normalize(polyline.vertices[i].Sub(polyline.vertices[i+3]), polyline.pixel_size)
		t := normalize(polyline.vertices[i].Sub(polyline.vertices[i+1]), polyline.pixel_size)

		k := 4 * (i - 2)
		overdraw[k] = polyline.vertices[i]
		overdraw[k+1] = polyline.vertices[i].Add(s.Add(t))
		overdraw[k+2] = polyline.vertices[i+1].Add(s.Sub(t))
		overdraw[k+3] = polyline.vertices[i+1]

		overdraw[k+4] = polyline.vertices[i+1]
		overdraw[k+5] = polyline.vertices[i+1].Add(s.Sub(t))
		overdraw[k+6] = polyline.vertices[i+2].Sub(s.Sub(t))
		overdraw[k+7] = polyline.vertices[i+2]

		overdraw[k+8] = polyline.vertices[i+2]
		overdraw[k+9] = polyline.vertices[i+2].Sub(s.Sub(t))
		overdraw[k+10] = polyline.vertices[i+3].Sub(s.Add(t))
		overdraw[k+11] = polyline.vertices[i+3]

		overdraw[k+12] = polyline.vertices[i+3]
		overdraw[k+13] = polyline.vertices[i+3].Sub(s.Add(t))
		overdraw[k+14] = polyline.vertices[i].Add(s.Add(t))
		overdraw[k+15] = polyline.vertices[i]
	}
	return overdraw
}

func (polyline *polyLine) generateColorArray(count int, c *Color) []Color {
	colors := make([]Color, count)
	for i := 0; i < count; i++ {
		colors[i] = *c
		if i%2 == 1 || (polyline.join == LINE_JOIN_NONE && (i%4 == 2 || i%4 == 1)) {
			colors[i][3] = 0
		}
	}
	return colors
}

func (polyline *polyLine) draw(is_looping bool) {
	switch polyline.join {
	case LINE_JOIN_NONE:
		polyline.drawTriangles(is_looping)
	case LINE_JOIN_MITER:
		fallthrough
	case LINE_JOIN_BEVEL:
		fallthrough
	default:
		polyline.drawTriangleStrip(is_looping)
	}
}

func (polyline *polyLine) drawTriangles(is_looping bool) {
	var overdraw []mgl32.Vec2
	if polyline.overdraw {
		overdraw = polyline.renderOverdraw(is_looping)
	}

	numindices := int(math.Max(float64(len(polyline.vertices)/4), float64(len(overdraw)/4)))
	indices := newAltQuadIndices(numindices)

	prepareDraw(nil)
	bindTexture(gl_state.defaultTexture)
	useVertexAttribArrays(attribflag_pos)
	gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 0, gl.Ptr(polyline.vertices))
	gl.DrawElements(gl.TRIANGLES, (len(polyline.vertices)/4)*6, gl.UNSIGNED_SHORT, gl.Ptr(indices))
	if polyline.overdraw {
		c := GetColor()
		colors := polyline.generateColorArray(len(overdraw), c)
		useVertexAttribArrays(attribflag_pos | attribflag_color)
		gl.VertexAttribPointer(attrib_color, 4, gl.UNSIGNED_BYTE, true, 0, gl.Ptr(colors))
		gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 0, gl.Ptr(overdraw))
		gl.DrawElements(gl.TRIANGLES, (len(overdraw)/4)*6, gl.UNSIGNED_SHORT, gl.Ptr(indices))
		SetColorC(c)
	}
}

func (polyline *polyLine) drawTriangleStrip(is_looping bool) {
	prepareDraw(nil)
	bindTexture(gl_state.defaultTexture)
	useVertexAttribArrays(attribflag_pos)
	gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 0, gl.Ptr(polyline.vertices))
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, len(polyline.vertices))
	if polyline.overdraw { // prepare colors:
		c := GetColor()
		overdraw := polyline.renderOverdraw(is_looping)
		colors := polyline.generateColorArray(len(overdraw), c)
		useVertexAttribArrays(attribflag_pos | attribflag_color)
		gl.VertexAttribPointer(attrib_color, 4, gl.UNSIGNED_BYTE, true, 0, gl.Ptr(colors))
		gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 0, gl.Ptr(overdraw))
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, len(overdraw))
		SetColorC(c)
	}
}
