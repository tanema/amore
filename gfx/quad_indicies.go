package gfx

import (
	"github.com/tanema/amore/gfx/gl"
)

type indexBuffer struct {
	isBound bool      // Whether the buffer is currently bound.
	ibo     gl.Buffer // The IBO identifier. Assigned by OpenGL.
	data    []uint32  // A pointer to mapped memory.
}

func newIndexBuffer(data []uint32) *indexBuffer {
	newBuffer := &indexBuffer{data: data}
	registerVolatile(newBuffer)
	return newBuffer
}

func (buffer *indexBuffer) bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer.ibo)
	buffer.isBound = true
}

func (buffer *indexBuffer) unbind() {
	if buffer.isBound {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Buffer{Value: 0})
	}
	buffer.isBound = false
}

func (buffer *indexBuffer) drawElements(mode uint32, offset, size int) {
	buffer.bind()
	defer buffer.unbind()
	gl.DrawElements(gl.Enum(mode), size, gl.UNSIGNED_INT, gl.PtrOffset(offset*4))
}

func (buffer *indexBuffer) drawElementsLocal(mode uint32, offset, size int) {
	gl.DrawElements(gl.Enum(mode), size, gl.UNSIGNED_INT, gl.Ptr(&buffer.data[offset]))
}

func (buffer *indexBuffer) loadVolatile() bool {
	buffer.ibo = gl.CreateBuffer()
	buffer.bind()
	defer buffer.unbind()
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(buffer.data)*4, gl.Ptr(buffer.data), uint32(gl.STATIC_DRAW))
	return true
}

func (buffer *indexBuffer) unloadVolatile() {
	gl.DeleteBuffer(buffer.ibo)
	buffer.ibo.Value = 0
}

/**
 * QuadIndices manages one shared buffer that stores the indices for an
 * element array. Vertex arrays using the vertex structure (or anything else
 * that can use the pattern below) can request a size and use it for the
 * drawElements call.
 *
 *			0----2
 *			|  / |
 *			| /  |
 *			1----3
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
		indexBuffer: newIndexBuffer(indices),
	}
}

func newAltQuadIndices(size int) *quadIndices {
	indices := make([]uint32, size*6)
	for i := 0; i < size; i++ {
		indices[i*6+0] = uint32(i*4 + 0)
		indices[i*6+1] = uint32(i*4 + 1)
		indices[i*6+2] = uint32(i*4 + 2)

		indices[i*6+3] = uint32(i*4 + 2)
		indices[i*6+4] = uint32(i*4 + 3)
		indices[i*6+5] = uint32(i*4 + 1)
	}

	return &quadIndices{
		indexBuffer: newIndexBuffer(indices),
	}
}

func (qi *quadIndices) drawElements(mode uint32, offset, size int) {
	qi.indexBuffer.drawElements(mode, offset*6, size*6)
}

func (qi *quadIndices) drawElementsLocal(mode uint32, offset, size int) {
	qi.indexBuffer.drawElementsLocal(mode, offset*6, size*6)
}
