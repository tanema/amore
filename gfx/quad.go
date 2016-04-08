package gfx

type (
	// Quad is essentially a crop of an image/texture
	Quad struct {
		vertices []float32
		x        float32
		y        float32
		w        float32
		h        float32
		sw       float32
		sh       float32
	}
)

// New Quad will generate a new *Quad with the dimensions given
// x, y are position on the texture
// w, h are the size of the quad
// sw, sh are references on how large the texture is. image.GetWidth(), image.GetHeight()
func NewQuad(x, y, w, h, sw, sh int32) *Quad {
	new_quad := &Quad{
		x:  float32(x),
		y:  float32(y),
		w:  float32(w),
		h:  float32(h),
		sw: float32(sw),
		sh: float32(sh),
	}
	new_quad.generateVertices()
	return new_quad
}

// generateVertices generates an array of data for drawing the quad.
func (quad *Quad) generateVertices() {
	quad.vertices = []float32{
		0, 0, quad.x / quad.sw, quad.y / quad.sh,
		0, quad.h, quad.x / quad.sw, (quad.y + quad.h) / quad.sh,
		quad.w, 0, (quad.x + quad.w) / quad.sw, quad.y / quad.sh,
		quad.w, quad.h, (quad.x + quad.w) / quad.sw, (quad.y + quad.h) / quad.sh,
	}
}

// getVertices will return the generated verticies
func (quad *Quad) getVertices() []float32 {
	return quad.vertices
}

// SetViewport sets the texture coordinates according to a viewport.
func (quad *Quad) SetViewport(x, y, w, h float32) {
	quad.x = x
	quad.y = y
	quad.w = w
	quad.h = h
	quad.generateVertices()
}

// GetViewport gets the current viewport of this Quad.
func (quad *Quad) GetViewport() (x, y, w, h float32) {
	return quad.x, quad.y, quad.w, quad.h
}
