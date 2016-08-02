package gfx

import (
	"github.com/tanema/amore/gfx/gl"
)

type indexBuffer struct {
	is_bound        bool      // Whether the buffer is currently bound.
	usage           Usage     // Usage hint. GL_[DYNAMIC, STATIC, STREAM]_DRAW.
	ibo             gl.Buffer // The IBO identifier. Assigned by OpenGL.
	data            []uint32  // A pointer to mapped memory.
	modified_offset int
	modified_size   int
}

func newIndexBuffer(size int, data []uint32, usage Usage) *indexBuffer {
	new_buffer := &indexBuffer{
		usage: usage,
		data:  make([]uint32, size),
	}
	if len(data) > 0 {
		copy(new_buffer.data, data[:size])
	}
	registerVolatile(new_buffer)
	return new_buffer
}

func (buffer *indexBuffer) bufferStatic() {
	if buffer.modified_size == 0 {
		return
	}
	// Upload the mapped data to the buffer.
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, buffer.modified_offset*4, buffer.modified_size*4, gl.Ptr(&buffer.data[buffer.modified_offset]))
}

func (buffer *indexBuffer) bufferStream() {
	// "orphan" current buffer to avoid implicit synchronisation on the GPU:
	// http://www.seas.upenn.edu/~pcozzi/OpenGLInsights/OpenGLInsights-AsynchronousBufferTransfers.pdf
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(buffer.data)*4, gl.Ptr(nil), uint32(buffer.usage))
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(buffer.data)*4, gl.Ptr(buffer.data), uint32(buffer.usage))
}

func (buffer *indexBuffer) bufferData() {
	if buffer.modified_size != 0 { //if there is no modified size might as well do the whole buffer
		buffer.modified_offset = Mini(buffer.modified_offset, len(buffer.data)-1)
		buffer.modified_size = Mini(buffer.modified_size, len(buffer.data)-buffer.modified_offset)
	} else {
		buffer.modified_offset = 0
		buffer.modified_size = len(buffer.data)
	}

	buffer.bind()
	if buffer.modified_size > 0 {
		switch buffer.usage {
		case USAGE_STATIC:
			buffer.bufferStatic()
		case USAGE_STREAM:
			buffer.bufferStream()
		case USAGE_DYNAMIC:
			// It's probably more efficient to treat it like a streaming buffer if
			// at least a third of its contents have been modified during the map().
			if buffer.modified_size >= len(buffer.data)/3 {
				buffer.bufferStream()
			} else {
				buffer.bufferStatic()
			}
		}
	}

	buffer.modified_offset = 0
	buffer.modified_size = 0
}

func (buffer *indexBuffer) bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer.ibo)
	buffer.is_bound = true
}

func (buffer *indexBuffer) unbind() {
	if buffer.is_bound {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Buffer{Value: 0})
	}
	buffer.is_bound = false
}

func (buffer *indexBuffer) fill(offset int, data []uint32) {
	copy(buffer.data[offset:], data[:len(data)])
	if !buffer.ibo.Valid() {
		return
	}
	// We're being conservative right now by internally marking the whole range
	// from the start of section a to the end of section b as modified if both
	// a and b are marked as modified.
	old_range_end := buffer.modified_offset + buffer.modified_size
	buffer.modified_offset = Mini(buffer.modified_offset, offset)
	new_range_end := Maxi(offset+len(data), old_range_end)
	buffer.modified_size = new_range_end - buffer.modified_offset
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

func (buffer *indexBuffer) Release() {
	releaseVolatile(buffer)
}
