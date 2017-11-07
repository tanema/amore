package gfx

import (
	"math"

	"github.com/tanema/amore/gfx/gl"
)

type indexBuffer struct {
	isBound        bool      // Whether the buffer is currently bound.
	usage          Usage     // Usage hint. GL_[DYNAMIC, STATIC, STREAM]_DRAW.
	ibo            gl.Buffer // The IBO identifier. Assigned by OpenGL.
	data           []uint32  // A pointer to mapped memory.
	modifiedOffset int
	modifiedSize   int
}

func newIndexBuffer(size int, data []uint32, usage Usage) *indexBuffer {
	newBuffer := &indexBuffer{
		usage: usage,
		data:  make([]uint32, size),
	}
	if len(data) > 0 {
		copy(newBuffer.data, data[:size])
	}
	registerVolatile(newBuffer)
	return newBuffer
}

func (buffer *indexBuffer) bufferStatic() {
	if buffer.modifiedSize == 0 {
		return
	}
	// Upload the mapped data to the buffer.
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, buffer.modifiedOffset*4, buffer.modifiedSize*4, gl.Ptr(&buffer.data[buffer.modifiedOffset]))
}

func (buffer *indexBuffer) bufferStream() {
	// "orphan" current buffer to avoid implicit synchronisation on the GPU:
	// http://www.seas.upenn.edu/~pcozzi/OpenGLInsights/OpenGLInsights-AsynchronousBufferTransfers.pdf
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(buffer.data)*4, gl.Ptr(nil), uint32(buffer.usage))
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(buffer.data)*4, gl.Ptr(buffer.data), uint32(buffer.usage))
}

func (buffer *indexBuffer) bufferData() {
	if buffer.modifiedSize != 0 { //if there is no modified size might as well do the whole buffer
		buffer.modifiedOffset = int(math.Min(float64(buffer.modifiedOffset), float64(len(buffer.data)-1)))
		buffer.modifiedSize = int(math.Min(float64(buffer.modifiedSize), float64(len(buffer.data)-buffer.modifiedOffset)))
	} else {
		buffer.modifiedOffset = 0
		buffer.modifiedSize = len(buffer.data)
	}

	buffer.bind()
	if buffer.modifiedSize > 0 {
		switch buffer.usage {
		case UsageStatic:
			buffer.bufferStatic()
		case UsageStream:
			buffer.bufferStream()
		case UsageDynamic:
			// It's probably more efficient to treat it like a streaming buffer if
			// at least a third of its contents have been modified during the map().
			if buffer.modifiedSize >= len(buffer.data)/3 {
				buffer.bufferStream()
			} else {
				buffer.bufferStatic()
			}
		}
	}

	buffer.modifiedOffset = 0
	buffer.modifiedSize = 0
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

func (buffer *indexBuffer) fill(offset int, data []uint32) {
	copy(buffer.data[offset:], data[:])
	if !buffer.ibo.Valid() {
		return
	}
	// We're being conservative right now by internally marking the whole range
	// from the start of section a to the end of section b as modified if both
	// a and b are marked as modified.
	oldRangeEnd := buffer.modifiedOffset + buffer.modifiedSize
	buffer.modifiedOffset = int(math.Min(float64(buffer.modifiedOffset), float64(offset)))
	newRangeEnd := int(math.Max(float64(offset+len(data)), float64(oldRangeEnd)))
	buffer.modifiedSize = newRangeEnd - buffer.modifiedOffset
	buffer.bufferData()
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
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(buffer.data)*4, gl.Ptr(buffer.data), uint32(buffer.usage))
	return true
}

func (buffer *indexBuffer) unloadVolatile() {
	gl.DeleteBuffer(buffer.ibo)
	buffer.ibo.Value = 0
}
