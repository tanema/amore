package gfx

import (
	"github.com/tanema/amore/gfx/gl"
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
	*indexBuffer
}

func newQuadIndices(size int) *quadIndices {
	indices := make([]uint32, size*6)
	for i := 0; i < size; i++ {
		indices[i*6+0] = uint32(i*4 + 0)
		indices[i*6+1] = uint32(i*4 + 1)
		indices[i*6+2] = uint32(i*4 + 2)
		indices[i*6+3] = uint32(i*4 + 2)
		indices[i*6+4] = uint32(i*4 + 1)
		indices[i*6+5] = uint32(i*4 + 3)
	}

	return &quadIndices{
		indexBuffer: newIndexBuffer(len(indices), indices, gl.STATIC_DRAW),
	}
}

func newAltQuadIndices(size int) *quadIndices {
	indices := make([]uint32, size*6)
	for i := 0; i < size; i++ {
		indices[i*6+0] = uint32(i*4 + 0)
		indices[i*6+1] = uint32(i*4 + 1)
		indices[i*6+2] = uint32(i*4 + 2)

		indices[i*6+3] = uint32(i*4 + 0)
		indices[i*6+4] = uint32(i*4 + 2)
		indices[i*6+5] = uint32(i*4 + 3)
	}

	return &quadIndices{
		indexBuffer: newIndexBuffer(len(indices), indices, gl.STATIC_DRAW),
	}
}

func (qi *quadIndices) drawElements(mode uint32, offset, size int) {
	qi.indexBuffer.drawElements(mode, offset*6, size*6)
}

func (qi *quadIndices) drawElementsLocal(mode uint32, offset, size int) {
	qi.indexBuffer.drawElementsLocal(mode, offset*6, size*6)
}
