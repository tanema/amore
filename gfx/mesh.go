package gfx

import (
	"fmt"
	"math"

	"github.com/goxjs/gl"
)

type (
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

func NewMesh(verticies []float32, size int) (*Mesh, error) {
	return NewMeshExt(verticies, size, DRAWMODE_FAN, USAGE_DYNAMIC)
}

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

func (mesh *Mesh) generateFlags() error {
	switch mesh.vertexStride {
	case 8:
		mesh.enabledattribs = ATTRIBFLAG_POS | ATTRIBFLAG_TEXCOORD | ATTRIBFLAG_COLOR
	case 6:
		mesh.enabledattribs = ATTRIBFLAG_POS | ATTRIBFLAG_COLOR
	case 4:
		mesh.enabledattribs = ATTRIBFLAG_POS | ATTRIBFLAG_TEXCOORD
	case 2:
		mesh.enabledattribs = ATTRIBFLAG_POS
	default:
		return fmt.Errorf("invalid mesh verticies format, vertext stride was calculated as %v", mesh.vertexStride)
	}
	return nil
}

func (mesh *Mesh) SetDrawMode(mode MeshDrawMode) {
	mesh.mode = mode
}

func (mesh *Mesh) GetDrawMode() MeshDrawMode {
	return mesh.mode
}

func (mesh *Mesh) SetTexture(text iTexture) {
	mesh.texture = text
}

func (mesh *Mesh) ClearTexture() {
	mesh.texture = nil
}

func (mesh *Mesh) GetTexture() iTexture {
	return mesh.texture
}

func (mesh *Mesh) SetDrawRange(min, max int) error {
	if min < 0 || max < 0 || min > max {
		return fmt.Errorf("Invalid draw range.")
	}
	mesh.rangeMin = min
	mesh.rangeMax = max
	return nil
}

func (mesh *Mesh) ClearDrawRange() {
	mesh.rangeMin = -1
	mesh.rangeMax = -1
}

func (mesh *Mesh) GetDrawRange() (int, int) {
	var min, max int
	if mesh.ibo != nil && mesh.elementCount > 0 {
		max = mesh.elementCount - 1
	} else {
		max = mesh.vertexCount - 1
	}
	if mesh.rangeMax >= 0 {
		max = int(math.Min(float64(mesh.rangeMax), float64(max)))
	}
	if mesh.rangeMin >= 0 {
		min = int(math.Min(float64(mesh.rangeMin), float64(max)))
	}
	return min, max
}

func (mesh *Mesh) SetVertex(vertindex int, data []float32) error {
	if vertindex >= mesh.vertexCount {
		return fmt.Errorf("Invalid vertex index: %v", vertindex+1)
	} else if (len(data) / mesh.vertexStride) != 0 {
		return fmt.Errorf("Invalid vertex data, data given was of len %v and not divisible by the meshes vertex stride of %v", len(data), mesh.vertexStride)
	}
	mesh.vbo.fill(vertindex*mesh.vertexStride, data)
	return nil
}

// set vertext handles both because of how free form it is and its general checking
func (mesh *Mesh) SetVertices(startindex int, data []float32) error {
	return mesh.SetVertex(startindex, data)
}

func (mesh *Mesh) GetVertex(vertindex int) ([]float32, error) {
	if vertindex >= mesh.vertexCount {
		return []float32{}, fmt.Errorf("Invalid vertex index: %v", vertindex+1)
	}
	data := mesh.vbo.getData()
	return data[vertindex*mesh.vertexStride : mesh.vertexStride], nil
}

func (mesh *Mesh) GetVertexCount() int {
	return mesh.vertexCount
}

func (mesh *Mesh) GetVertexStride() int {
	return mesh.vertexStride
}

func (mesh *Mesh) GetVertexFormat() (vertex, text, color bool) {
	return mesh.enabledattribs&ATTRIBFLAG_POS > 0, mesh.enabledattribs&ATTRIBFLAG_TEXCOORD > 0, mesh.enabledattribs&ATTRIBFLAG_COLOR > 0
}

func (mesh *Mesh) SetVertexMap(vertex_map []uint32) {
	if len(vertex_map) > 0 {
		mesh.ibo = newIndexBuffer(len(vertex_map), vertex_map, mesh.vbo.usage)
		mesh.elementCount = len(vertex_map)
	}
}

func (mesh *Mesh) ClearVertexMap() {
	mesh.ibo = nil
	mesh.elementCount = 0
}

func (mesh *Mesh) GetVertexMap() []uint32 {
	if mesh.ibo == nil {
		return []uint32{}
	}
	return mesh.ibo.getData()
}

func (mesh *Mesh) Flush() {
	mesh.vbo.bufferData()
	if mesh.ibo != nil {
		mesh.ibo.bufferData()
	}
}

func (mesh *Mesh) bindEnabledAttributes() {
	useVertexAttribArrays(mesh.enabledattribs)

	mesh.vbo.bind()
	defer mesh.vbo.unbind()

	offset := 0
	if (mesh.enabledattribs & ATTRIBFLAG_POS) > 0 {
		gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, mesh.vertexStride*4, offset)
		offset += 2 * 4
	}
	if (mesh.enabledattribs & ATTRIBFLAG_TEXCOORD) > 0 {
		gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, mesh.vertexStride*4, offset)
		offset += 2 * 4
	}
	if (mesh.enabledattribs & ATTRIBFLAG_COLOR) > 0 {
		gl.VertexAttribPointer(ATTRIB_COLOR, 4, gl.FLOAT, false, mesh.vertexStride*4, offset)
	}
}

func (mesh *Mesh) bindTexture() {
	if mesh.texture != nil {
		bindTexture(mesh.texture.GetHandle())
	} else {
		bindTexture(gl_state.defaultTexture)
	}
}

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
