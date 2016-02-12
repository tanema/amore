package gfx

import (
	"github.com/go-gl/gl/v2.1/gl"
)

/**
 * QuadIndices manages one shared buffer that stores the indices for an
 * element array. Vertex arrays using the vertex structure (or anything else
 * that can use the pattern below) can request a size and use it for the
 * drawElements call.
 *
 * 			0----2
 * 			|  / |
 * 			| /  |
 * 			1----3
 *
 *  indices[i*6 + 0] = i*4 + 0;
 *  indices[i*6 + 1] = i*4 + 1;
 *  indices[i*6 + 2] = i*4 + 2;
 *
 *  indices[i*6 + 3] = i*4 + 2;
 *  indices[i*6 + 4] = i*4 + 1;
 *  indices[i*6 + 5] = i*4 + 3;
 *
 * There will always be a large enough buffer around until all
 * QuadIndices instances have been deleted.
 *
 * Q: Why have something like QuadIndices?
 * A: The indices for the SpriteBatch do not change, only the array size
 * varies. Using one buffer for all element arrays removes this
 * duplicated data and saves some memory.
 */
type quadIndices struct {
	vbo     uint32
	indices []uint32
}

func newQuadIndices(size int) *quadIndices {
	new_qi := &quadIndices{
		indices: make([]uint32, size*6),
	}

	for i := 0; i < size; i++ {
		new_qi.indices[i*6+0] = uint32(i*4 + 0)
		new_qi.indices[i*6+1] = uint32(i*4 + 1)
		new_qi.indices[i*6+2] = uint32(i*4 + 2)

		new_qi.indices[i*6+3] = uint32(i*4 + 2)
		new_qi.indices[i*6+4] = uint32(i*4 + 1)
		new_qi.indices[i*6+5] = uint32(i*4 + 3)
	}

	registerVolatile(new_qi)
	return new_qi
}

func (qi *quadIndices) bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, qi.vbo)
}

func (qi *quadIndices) unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (qi *quadIndices) drawElements(offset, size int) {
	gl.DrawElements(gl.TRIANGLES, int32(size*6), gl.UNSIGNED_INT, gl.PtrOffset(offset*6))
}

func (qi *quadIndices) drawElementsLocal(offset, size int) {
	gl.DrawElements(gl.TRIANGLES, int32(size*6), gl.UNSIGNED_INT, gl.Ptr(&qi.indices[offset*6]))
}

func (qi *quadIndices) loadVolatile() bool {
	gl.GenBuffers(1, &qi.vbo)
	qi.bind()
	defer qi.unbind()
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(qi.indices), gl.Ptr(qi.indices), gl.STATIC_DRAW)
	return true
}

func (qi *quadIndices) unloadVolatile() {}
