package gfx

type (
	Quad struct {
		coords    []float32
		texcoords []float32
		x         float32
		y         float32
		w         float32
		h         float32
		sw        float32
		sh        float32
	}
)

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

func (quad *Quad) generateVertices() {
	quad.coords = []float32{0, 0, 0, quad.h, quad.w, 0, quad.w, quad.h}
	quad.texcoords = []float32{
		quad.x / quad.sw,
		quad.y / quad.sh,
		quad.x / quad.sw,
		(quad.y + quad.h) / quad.sh,
		(quad.x + quad.w) / quad.sw,
		quad.y / quad.sh,
		(quad.x + quad.w) / quad.sw,
		(quad.y + quad.h) / quad.sh,
	}
}

func (quad *Quad) SetViewport(x, y, w, h float32) {
	quad.x = x
	quad.y = y
	quad.w = w
	quad.h = h
	quad.generateVertices()
}

func (quad *Quad) GetViewport() (x, y, w, h float32) {
	return quad.x, quad.y, quad.w, quad.h
}
