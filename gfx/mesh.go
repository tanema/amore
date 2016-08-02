package gfx

import (
	"fmt"

	"github.com/tanema/amore/gfx/gl"
)

type (
	// Mesh is a collection of points that a texture can be applied to.
	Mesh struct {
		mode           MeshDrawMode
		texture        iTexture
		vbo            *vertexBuffer
		vertexStride   int
		enabledattribs uint32
		rangeMin       int
		rangeMax       int
		vertexCount    int
		ibo            *indexBuffer
		elementCount   int
	}
	meshAttribute struct {
		name string
		size int
	}
)

// NewMesh will generate a mesh with the verticies and it with a max size of the size
// provided. MeshDrawMode is default to DRAWMODE_FAN and Usage is USAGE_DYNAMIC
func NewMesh(verticies []float32, size int) (*Mesh, error) {
	return NewMeshExt(verticies, size, DRAWMODE_FAN, USAGE_DYNAMIC)
}

// NewMeshExt is like NewMesh but with access to setting the MeshDrawMode and Usage.
func NewMeshExt(vertices []float32, size int, mode MeshDrawMode, usage Usage) (*Mesh, error) {
	count := len(vertices)
	stride := count / size

	if count == 0 || stride == 0 {
		return nil, fmt.Errorf("Not enough data to establish a mesh")
	}

	new_mesh := &Mesh{
		rangeMin:     -1,
		rangeMax:     -1,
		mode:         mode,
		vertexStride: stride,
		vertexCount:  count / stride,
		vbo:          newVertexBuffer(size*stride, vertices, usage),
	}

	new_mesh.generateFlags()

	return new_mesh, nil
}

// generateFlags will generate a attribs flag setting with the vertex stride. It
// will guess which components it was provided and enable that functionality.
func (mesh *Mesh) generateFlags() error {
	switch mesh.vertexStride {
	case 8:
		mesh.enabledattribs = attribflag_pos | attribflag_texcoord | attribflag_color
	case 6:
		mesh.enabledattribs = attribflag_pos | attribflag_color
	case 4:
		mesh.enabledattribs = attribflag_pos | attribflag_texcoord
	case 2:
		mesh.enabledattribs = attribflag_pos
	default:
		return fmt.Errorf("invalid mesh verticies format, vertext stride was calculated as %v", mesh.vertexStride)
	}
	return nil
}

// SetDrawMode sets the MeshDrawMode. Please see the constant definitions for explanation.
func (mesh *Mesh) SetDrawMode(mode MeshDrawMode) {
	mesh.mode = mode
}

// GetDrawMode will return the current draw mode for the mesh
func (mesh *Mesh) GetDrawMode() MeshDrawMode {
	return mesh.mode
}

// SetTexture will apply a texture to the mesh. This will go a lot better if uv coords
// were provided to the mesh.
func (mesh *Mesh) SetTexture(text iTexture) {
	mesh.texture = text
}

// ClearTexture will unbind the texture of the mesh
func (mesh *Mesh) ClearTexture() {
	mesh.texture = nil
}

// GetTexture will return an interface iTexture if there is a texture bound, and
// nil if there isnt
func (mesh *Mesh) GetTexture() iTexture {
	return mesh.texture
}

// SetDrawRange will set a range in the points to draw. This is useful if you only
// need to render a portion of the mesh.
func (mesh *Mesh) SetDrawRange(min, max int) error {
	if min < 0 || max < 0 || min > max {
		return fmt.Errorf("Invalid draw range.")
	}
	mesh.rangeMin = min
	mesh.rangeMax = max
	return nil
}

// ClearDrawRange will reset the draw range if you want to draw the whole mesh again.
func (mesh *Mesh) ClearDrawRange() {
	mesh.rangeMin = -1
	mesh.rangeMax = -1
}

// GetDrawRange will return the min, max range set on the mesh. If no range is set
// the range will return -1, -1
func (mesh *Mesh) GetDrawRange() (int, int) {
	var min, max int
	if mesh.ibo != nil && mesh.elementCount > 0 {
		max = mesh.elementCount - 1
	} else {
		max = mesh.vertexCount - 1
	}
	if mesh.rangeMax >= 0 {
		max = Mini(mesh.rangeMax, max)
	}
	if mesh.rangeMin >= 0 {
		min = Mini(mesh.rangeMin, max)
	}
	return min, max
}

// SetVertex will replace vertex data at the specified index. The index specifies
// the point and not the index in the array provided to the mesh.
func (mesh *Mesh) SetVertex(vertindex int, data []float32) error {
	if vertindex >= mesh.vertexCount {
		return fmt.Errorf("Invalid vertex index: %v", vertindex+1)
	} else if (len(data) / mesh.vertexStride) != 0 {
		return fmt.Errorf("Invalid vertex data, data given was of len %v and not divisible by the meshes vertex stride of %v", len(data), mesh.vertexStride)
	}
	mesh.vbo.fill(vertindex*mesh.vertexStride, data)
	return nil
}

// SetVertices is like SetVertex but will do more than one point
func (mesh *Mesh) SetVertices(startindex int, data []float32) error {
	// set vertext handles both because of how free form it is and its general checking
	return mesh.SetVertex(startindex, data)
}

// GetVertex will return all the data for that vertex at the given index.
func (mesh *Mesh) GetVertex(vertindex int) ([]float32, error) {
	if vertindex >= mesh.vertexCount {
		return []float32{}, fmt.Errorf("Invalid vertex index: %v", vertindex+1)
	}
	return mesh.vbo.data[vertindex*mesh.vertexStride : mesh.vertexStride], nil
}

// GetVertexCount will return how many vertexes there are in the provided data.
func (mesh *Mesh) GetVertexCount() int {
	return mesh.vertexCount
}

// GetVertexStride will return the number of components in each vertex. i.e.
// x, y, u, v, r, g, b, a
func (mesh *Mesh) GetVertexStride() int {
	return mesh.vertexStride
}

// GetVertexFormat will return true for each attribute that is enabled.
func (mesh *Mesh) GetVertexFormat() (vertex, text, color bool) {
	return mesh.enabledattribs&attribflag_pos > 0, mesh.enabledattribs&attribflag_texcoord > 0, mesh.enabledattribs&attribflag_color > 0
}

// SetVertexMap will allow to set indexes of verticies that should be drawn.
func (mesh *Mesh) SetVertexMap(vertex_map []uint32) {
	if len(vertex_map) > 0 {
		mesh.ibo = newIndexBuffer(len(vertex_map), vertex_map, mesh.vbo.usage)
		mesh.elementCount = len(vertex_map)
	}
}

// ClearVertexMap disabled the vertex map and re-enabled drawing the whole mesh again.
func (mesh *Mesh) ClearVertexMap() {
	mesh.ibo = nil
	mesh.elementCount = 0
}

// GetVertexMap returns the currently set vertex map and an empty slice if there
// is not one set.
func (mesh *Mesh) GetVertexMap() []uint32 {
	if mesh.ibo == nil {
		return []uint32{}
	}
	return mesh.ibo.data
}

// Flush immediately sends all modified vertex data in the Mesh to the graphics card.
// Normally it isn't necessary to call this method as Draw(mesh, ...) will do it
// automatically if needed, but explicitly using Flush gives more control over when
// the work happens. If this method is used, it generally shouldn't be called more
// than once (at most) between draw(mesh, ...) calls.
func (mesh *Mesh) Flush() {
	mesh.vbo.bufferData()
	if mesh.ibo != nil {
		mesh.ibo.bufferData()
	}
}

// bindEnabledAttributes will take the enabled attrib flags and use them to enable
// all the attributes that we need.
func (mesh *Mesh) bindEnabledAttributes() {
	useVertexAttribArrays(mesh.enabledattribs)

	mesh.vbo.bind()
	defer mesh.vbo.unbind()

	offset := 0
	if (mesh.enabledattribs & attribflag_pos) > 0 {
		gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, mesh.vertexStride*4, gl.PtrOffset(offset))
		offset += 2 * 4
	}
	if (mesh.enabledattribs & attribflag_texcoord) > 0 {
		gl.VertexAttribPointer(attrib_texcoord, 2, gl.FLOAT, false, mesh.vertexStride*4, gl.PtrOffset(offset))
		offset += 2 * 4
	}
	if (mesh.enabledattribs & attribflag_color) > 0 {
		gl.VertexAttribPointer(attrib_color, 4, gl.FLOAT, false, mesh.vertexStride*4, gl.PtrOffset(offset))
	}
}

// bindTexture will bind our current texture if we have one or the framework default
// if there wan't a texture provided.
func (mesh *Mesh) bindTexture() {
	if mesh.texture != nil {
		bindTexture(mesh.texture.getHandle())
	} else {
		bindTexture(gl_state.defaultTexture)
	}
}

// Draw satisfies the Drawable interface. Inputs are as follows
// x, y, r, sx, sy, ox, oy, kx, ky
// x, y are position
// r is rotation
// sx, sy is the scale, if sy is not given sy will equal sx
// ox, oy are offset
// kx, ky are the shear. If ky is not given ky will equal kx
func (mesh *Mesh) Draw(args ...float32) {
	prepareDraw(generateModelMatFromArgs(args))
	mesh.bindTexture()
	mesh.bindEnabledAttributes()
	min, max := mesh.GetDrawRange()
	if mesh.ibo != nil && mesh.elementCount > 0 {
		mesh.ibo.drawElements(uint32(mesh.mode), min, max-min+1)
	} else {
		gl.DrawArrays(gl.Enum(mesh.mode), min, max-min+1)
	}
}
