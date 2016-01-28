package gfx

import (
	"github.com/go-gl/mathgl/mgl32"
)

type (
	Quadable interface {
		GetWidth() int32
		GetHeight() int32
		drawv(model *mgl32.Mat4, coords, texcoords []float32)
	}
	Quad struct {
		quadable  Quadable
		coords    []float32
		texcoords []float32
		x         float32
		y         float32
		w         float32
		h         float32
	}
)

func NewQuad(quadable Quadable, x, y, w, h float32) *Quad {
	new_quad := &Quad{quadable: quadable, x: x, y: y, w: w, h: h}
	new_quad.generateVertices()
	return new_quad
}

func (quad *Quad) generateVertices() {
	sw := float32(quad.quadable.GetWidth())
	sh := float32(quad.quadable.GetHeight())
	quad.coords = []float32{0, 0, 0, quad.h, quad.w, 0, quad.w, quad.h}
	quad.texcoords = []float32{
		quad.x / sw,
		quad.y / sh,
		quad.x / sw,
		(quad.y + quad.h) / sh,
		(quad.x + quad.w) / sw,
		quad.y / sh,
		(quad.x + quad.w) / sw,
		(quad.y + quad.h) / sh,
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

func (quad *Quad) getVertices() ([]float32, []float32) {
	return quad.coords, quad.texcoords
}

func (quad *Quad) Draw(args ...float32) {
	quad.quadable.drawv(generateModelMatFromArgs(args), quad.coords, quad.texcoords)
}
