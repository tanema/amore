package gfx

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v2.1/gl"
)

type Mesh struct {
	verticies []float32
	mode      MeshDrawMode
	texture   iTexture
}

func NewMesh(verticies []float32) *Mesh {
	return NewMeshExt(verticies, DRAWMODE_FAN, USAGE_DYNAMIC)
}

func NewMeshExt(verticies []float32, mode MeshDrawMode, usage Usage) *Mesh {
	new_mesh := &Mesh{
		mode:      mode,
		verticies: verticies,
	}

	return new_mesh
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

func (mesh *Mesh) drawv(model *mgl32.Mat4, verticies []float32) {
	prepareDraw(model)
	if mesh.texture != nil {
		bindTexture(mesh.texture.GetHandle())
	} else {
		bindTexture(gl_state.defaultTexture)
	}

	useVertexAttribArrays(ATTRIBFLAG_POS | ATTRIBFLAG_TEXCOORD | ATTRIBFLAG_COLOR)
	//gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	//gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(2*4))
	//gl.VertexAttribPointer(ATTRIB_COLOR, 4, gl.FLOAT, false, 8*4, gl.PtrOffset(4*4))
	gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 8*4, gl.Ptr(verticies))
	gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, 8*4, gl.Ptr(&verticies[2]))
	gl.VertexAttribPointer(ATTRIB_COLOR, 4, gl.FLOAT, false, 8*4, gl.Ptr(&verticies[4]))
	gl.DrawArrays(uint32(mesh.mode), 0, 4)
}

func (mesh *Mesh) Draw(args ...float32) {
	mesh.drawv(generateModelMatFromArgs(args), mesh.verticies)
}

func (mesh *Mesh) Drawq(quad *Quad, args ...float32) {
	mesh.drawv(generateModelMatFromArgs(args), quad.getVertices())
}
