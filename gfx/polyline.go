package gfx

import (
	"math"

	"github.com/goxjs/gl"
)

// treat adjacent segments with angles between their directions <5 degree as straight
const linesParallelEPS float32 = 0.05

type polyLine struct {
	join      LineJoin
	halfwidth float32
	pixelSize float32
}

func determinant(vec1, vec2 []float32) float32 {
	return vec1[0]*vec2[1] - vec1[1]*vec2[0]
}

func getNormal(v1 []float32, scale float32) []float32 {
	return []float32{-v1[1] * scale, v1[0] * scale}
}

func normalize(v1 []float32, length float32) []float32 {
	lengthCurrent := vecLen(v1)
	if lengthCurrent > 0 {
		scale := length / lengthCurrent
		return []float32{v1[0] * scale, v1[1] * scale}
	}
	return v1
}

func vecLen(v1 []float32) float32 {
	return float32(math.Hypot(float64(v1[0]), float64(v1[1])))
}

func abs(a float32) float32 {
	if a < 0 {
		return -a
	} else if a == 0 {
		return 0
	}
	return a
}

func newPolyLine(join LineJoin, lineWidth, pixelSize float32) polyLine {
	newPolyline := polyLine{
		join:      join,
		halfwidth: lineWidth * 0.5,
		pixelSize: pixelSize,
	}
	return newPolyline
}

func (polyline *polyLine) render(coords []float32) {
	var sleeve, current, next []float32
	vertices := []float32{}

	coordsCount := len(coords)
	isLooping := (coords[0] == coords[coordsCount-2]) && (coords[1] == coords[coordsCount-1])
	if !isLooping { // virtual starting point at second point mirrored on first point
		sleeve = []float32{coords[2] - coords[0], coords[3] - coords[1]}
	} else { // virtual starting point at last vertex
		sleeve = []float32{coords[0] - coords[coordsCount-4], coords[1] - coords[coordsCount-3]}
	}

	for i := 0; i+3 < coordsCount; i += 2 {
		current = []float32{coords[i], coords[i+1]}
		next = []float32{coords[i+2], coords[i+3]}
		vertices = append(vertices, polyline.renderEdge(sleeve, current, next)...)
		sleeve = []float32{next[0] - current[0], next[1] - current[1]}
	}

	if isLooping {
		vertices = append(vertices, polyline.renderEdge(sleeve, next, []float32{coords[2], coords[3]})...)
	} else {
		vertices = append(vertices, polyline.renderEdge(sleeve, next, []float32{next[0] + sleeve[0], next[1] + sleeve[1]})...)
	}

	prepareDraw(nil)
	bindTexture(glState.defaultTexture)
	useVertexAttribArrays(shaderPosFlag)

	buffer := newVertexBuffer(len(vertices), vertices, UsageStatic)
	buffer.bind()
	defer buffer.unbind()

	gl.VertexAttribPointer(gl.Attrib{Value: 0}, 2, gl.FLOAT, false, 0, 0)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, len(vertices)/2)
}

func (polyline *polyLine) renderEdge(sleeve, current, next []float32) []float32 {
	if polyline.join == LineJoinBevel {
		return polyline.renderBevelEdge(sleeve, current, next)
	}
	return polyline.renderMiterEdge(sleeve, current, next)
}

func (polyline *polyLine) generateEdges(current []float32, normals ...float32) []float32 {
	verts := make([]float32, len(normals))
	for i := 0; i < len(normals); i += 2 {
		verts[i] = current[0] + normals[i]
		verts[i+1] = current[1] + normals[i+1]
	}
	return verts
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
func (polyline *polyLine) renderMiterEdge(sleeve, current, next []float32) []float32 {
	sleeveNormal := getNormal(sleeve, polyline.halfwidth/vecLen(sleeve))
	t := []float32{next[0] - current[0], next[1] - current[1]}
	lenT := vecLen(t)

	det := determinant(sleeve, t)
	// lines parallel, compute as u1 = q + ns * w/2, u2 = q - ns * w/2
	if abs(det)/(vecLen(sleeve)*lenT) < linesParallelEPS && (sleeve[0]*t[0]+sleeve[1]*t[1]) > 0 {
		return polyline.generateEdges(current, sleeveNormal[0], sleeveNormal[1], sleeveNormal[0]*-1, sleeveNormal[1]*-1)
	}
	// cramers rule
	nt := getNormal(t, polyline.halfwidth/lenT)
	lambda := determinant([]float32{nt[0] - sleeveNormal[0], nt[1] - sleeveNormal[1]}, t) / det
	sleeveChange := []float32{sleeve[0] * lambda, sleeve[1] * lambda}
	d := []float32{sleeveNormal[0] + sleeveChange[0], sleeveNormal[1] + sleeveChange[1]}
	return polyline.generateEdges(current, d[0], d[1], d[0]*-1, d[1]*-1)
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
func (polyline *polyLine) renderBevelEdge(sleeve, current, next []float32) []float32 {
	t := []float32{next[0] - current[0], next[1] - current[1]}
	lenT := vecLen(t)

	det := determinant(sleeve, t)
	if abs(det)/(vecLen(sleeve)*lenT) < linesParallelEPS && (sleeve[0]*t[0]+sleeve[1]*t[1]) > 0 {
		// lines parallel, compute as u1 = q + ns * w/2, u2 = q - ns * w/2
		n := getNormal(t, polyline.halfwidth/lenT)
		return polyline.generateEdges(current, n[0], n[1], n[0]*-1, n[1]*-1)
	}

	// cramers rule
	sleeveNormal := getNormal(sleeve, polyline.halfwidth/vecLen(sleeve))
	nt := getNormal(t, polyline.halfwidth/lenT)
	lambda := determinant([]float32{nt[0] - sleeveNormal[0], nt[1] - sleeveNormal[1]}, t) / det
	sleeveChange := []float32{sleeve[0] * lambda, sleeve[1] * lambda}
	d := []float32{sleeveNormal[0] + sleeveChange[0], sleeveNormal[1] + sleeveChange[1]}

	if det > 0 { // 'left' turn -> intersection on the top
		return polyline.generateEdges(current,
			d[0], d[1],
			sleeveNormal[0]*-1, sleeveNormal[1]*-1,
			d[0], d[1],
			nt[0]*-1, nt[1]*-1,
		)
	}

	return polyline.generateEdges(current,
		sleeveNormal[0], sleeveNormal[1],
		d[0]*-1, d[1]*-1,
		nt[0], nt[1],
		d[0]*-1, d[1]*-1,
	)
}
