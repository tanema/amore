package gfx

type Mesh struct {
}

func NewMesh(verticies []float32) *Mesh {
	return NewMeshExt(verticies, DRAWMODE_FAN, USAGE_DYNAMIC)
}

func NewMeshExt(verticies []float32, mode MeshDrawMode, usage Usage) *Mesh {
	return &Mesh{}
}
