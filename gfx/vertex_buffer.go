package gfx

import (
	"github.com/tanema/amore/gfx/gl"
	"github.com/tanema/amore/mth"
)

type vertexBuffer struct {
	is_bound        bool      // Whether the buffer is currently bound.
	usage           Usage     // Usage hint. GL_[DYNAMIC, STATIC, STREAM]_DRAW.
	vbo             gl.Buffer // The VBO identifier. Assigned by OpenGL.
	data            []float32 // A pointer to mapped memory.
	modified_offset int
	modified_size   int
}

func newVertexBuffer(size int, data []float32, usage Usage) *vertexBuffer {
	new_buffer := &vertexBuffer{
		usage: usage,
		data:  make([]float32, size),
	}
	if len(data) > 0 {
		copy(new_buffer.data, data[:size])
	}
	registerVolatile(new_buffer)
	return new_buffer
}

func (buffer *vertexBuffer) bufferStatic() {
	if buffer.modified_size == 0 {
		return
	}
	// Upload the mapped data to the buffer.
	gl.BufferSubData(gl.ARRAY_BUFFER, buffer.modified_offset*4, buffer.modified_size*4, gl.Ptr(&buffer.data[buffer.modified_offset]))
}

func (buffer *vertexBuffer) bufferStream() {
	// "orphan" current buffer to avoid implicit synchronisation on the GPU:
	// http://www.seas.upenn.edu/~pcozzi/OpenGLInsights/OpenGLInsights-AsynchronousBufferTransfers.pdf
	gl.BufferData(gl.ARRAY_BUFFER, len(buffer.data)*4, gl.Ptr(nil), uint32(buffer.usage))
	gl.BufferData(gl.ARRAY_BUFFER, len(buffer.data)*4, gl.Ptr(buffer.data), uint32(buffer.usage))
}

func (buffer *vertexBuffer) bufferData() {
	if buffer.modified_size != 0 { //if there is no modified size might as well do the whole buffer
		buffer.modified_offset = mth.Mini(buffer.modified_offset, len(buffer.data)-1)
		buffer.modified_size = mth.Mini(buffer.modified_size, len(buffer.data)-buffer.modified_offset)
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

func (buffer *vertexBuffer) bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.vbo)
	buffer.is_bound = true
}

func (buffer *vertexBuffer) unbind() {
	if buffer.is_bound {
		gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{})
	}
	buffer.is_bound = false
}

func (buffer *vertexBuffer) fill(offset int, data []float32) {
	copy(buffer.data[offset:], data[:len(data)])
	if !buffer.vbo.Valid() {
		return
	}
	// We're being conservative right now by internally marking the whole range
	// from the start of section a to the end of section b as modified if both
	// a and b are marked as modified.
	old_range_end := buffer.modified_offset + buffer.modified_size
	buffer.modified_offset = mth.Mini(buffer.modified_offset, offset)
	new_range_end := mth.Maxi(offset+len(data), old_range_end)
	buffer.modified_size = new_range_end - buffer.modified_offset
	buffer.bufferData()
}

func (buffer *vertexBuffer) loadVolatile() bool {
	buffer.vbo = gl.CreateBuffer()
	buffer.bind()
	defer buffer.unbind()
	gl.BufferData(gl.ARRAY_BUFFER, len(buffer.data)*4, gl.Ptr(buffer.data), uint32(buffer.usage))
	return true
}

func (buffer *vertexBuffer) unloadVolatile() {
	gl.DeleteBuffer(buffer.vbo)
	buffer.vbo.Value = 0
}

func (buffer *vertexBuffer) Release() {
	releaseVolatile(buffer)
}
