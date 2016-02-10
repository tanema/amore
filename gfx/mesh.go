package gfx

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v2.1/gl"
)

type Mesh struct {
	coords    []float32
	texcoords []float32
	colors    []float32
	mode      MeshDrawMode
	texture   iTexture
}

func NewMesh(components ...[]float32) *Mesh {
	return NewMeshExt(DRAWMODE_FAN, USAGE_DYNAMIC, components...)
}

func NewMeshExt(mode MeshDrawMode, usage Usage, components ...[]float32) *Mesh {
	if components == nil {
		panic("Cannot create mesh without verticies")
	}

	new_mesh := &Mesh{
		mode: mode,
	}

	switch len(components) {
	case 3:
		new_mesh.colors = components[2]
		fallthrough
	case 2:
		new_mesh.texcoords = components[1]
		fallthrough
	case 1:
		new_mesh.coords = components[0]
	}

	return &Mesh{}
}

func (mesh *Mesh) loadVolatile() {
}

func (mesh *Mesh) unloadVolatile() {
}

func (mesh *Mesh) SetMode(mode MeshDrawMode) {
	mesh.mode = mode
}

func (mesh *Mesh) GetMode() MeshDrawMode {
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

func (mesh *Mesh) drawv(model *mgl32.Mat4, coords, texcoords []float32) {
	prepareDraw(model)
	if mesh.texture != nil {
		bindTexture(mesh.texture.GetHandle())
	} else {
		bindTexture(gl_state.defaultTexture)
	}

	gl.EnableVertexAttribArray(ATTRIB_POS)
	gl.EnableVertexAttribArray(ATTRIB_TEXCOORD)
	gl.EnableVertexAttribArray(ATTRIB_COLOR)

	gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 0, gl.Ptr(coords))
	gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, 0, gl.Ptr(texcoords))
	gl.VertexAttribPointer(ATTRIB_COLOR, 4, gl.FLOAT, false, 0, gl.Ptr(mesh.colors))

	gl.DrawArrays(uint32(mesh.mode), 0, 4)

	gl.DisableVertexAttribArray(ATTRIB_COLOR)
	gl.DisableVertexAttribArray(ATTRIB_TEXCOORD)
	gl.DisableVertexAttribArray(ATTRIB_POS)
}

func (mesh *Mesh) Draw(args ...float32) {
	mesh.drawv(generateModelMatFromArgs(args), mesh.coords, mesh.texcoords)
}

func (mesh *Mesh) Drawq(quad *Quad, args ...float32) {
	mesh.drawv(generateModelMatFromArgs(args), quad.coords, quad.texcoords)
}
