package gfx

import (
	"math"

	"github.com/goxjs/gl"
)

type glBuffer struct {
	is_bound        bool      // Whether the buffer is currently bound.
	usage           Usage     // Usage hint. GL_[DYNAMIC, STATIC, STREAM]_DRAW.
	vbo             gl.Buffer // The VBO identifier. Assigned by OpenGL.
	target          gl.Enum   // gl.ARRAY_BUFFER gl.ELEMENT_ARRAY_BUFFER
	data            []byte    // A pointer to mapped memory.
	modified_offset int
	modified_size   int
}

func newGlBuffer(target gl.Enum, size int, data []byte, usage Usage) *glBuffer {
	new_buffer := &glBuffer{
		target: target,
		usage:  usage,
		data:   make([]byte, size),
	}
	if len(data) > 0 {
		copy(new_buffer.data, data[:size])
	}
	registerVolatile(new_buffer)
	return new_buffer
}

func (buffer *glBuffer) bufferStatic() {
	if buffer.modified_size == 0 {
		return
	}
	// Upload the mapped data to the buffer.
	gl.BufferSubData(buffer.target, buffer.modified_offset, buffer.data[buffer.modified_offset:buffer.modified_size])
}

func (buffer *glBuffer) bufferStream() {
	// "orphan" current buffer to avoid implicit synchronisation on the GPU:
	// http://www.seas.upenn.edu/~pcozzi/OpenGLInsights/OpenGLInsights-AsynchronousBufferTransfers.pdf
	//gl.BufferData(buffer.target, nil, gl.Enum(buffer.usage))
	gl.BufferData(buffer.target, buffer.data, gl.Enum(buffer.usage))
}

func (buffer *glBuffer) bufferData() {
	if buffer.modified_size != 0 { //if there is no modified size might as well do the whole buffer
		buffer.modified_offset = int(math.Min(float64(buffer.modified_offset), float64(len(buffer.data)-1)))
		buffer.modified_size = int(math.Min(float64(buffer.modified_size), float64(len(buffer.data)-buffer.modified_offset)))
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

func (buffer *glBuffer) bind() {
	gl.BindBuffer(buffer.target, buffer.vbo)
	buffer.is_bound = true
}

func (buffer *glBuffer) unbind() {
	if buffer.is_bound {
		gl.BindBuffer(buffer.target, gl.Buffer{})
	}
	buffer.is_bound = false
}

func (buffer *glBuffer) fill(offset int, data []byte) {
	copy(buffer.data[offset:], data[:len(data)])
	if !buffer.vbo.Valid() {
		return
	}
	// We're being conservative right now by internally marking the whole range
	// from the start of section a to the end of section b as modified if both
	// a and b are marked as modified.
	old_range_end := buffer.modified_offset + buffer.modified_size
	buffer.modified_offset = int(math.Min(float64(buffer.modified_offset), float64(offset)))
	new_range_end := int(math.Max(float64(offset+len(data)), float64(old_range_end)))
	buffer.modified_size = new_range_end - buffer.modified_offset
	buffer.bufferData()
}

func (buffer *glBuffer) loadVolatile() bool {
	buffer.vbo = gl.CreateBuffer()
	buffer.bind()
	defer buffer.unbind()
	gl.BufferData(buffer.target, buffer.data, gl.Enum(buffer.usage))
	return true
}

func (buffer *glBuffer) unloadVolatile() {
	gl.DeleteBuffer(buffer.vbo)
	buffer.vbo = gl.Buffer{}
}

func (buffer *glBuffer) Release() {
	releaseVolatile(buffer)
}

// Vertex Buffer

type vertexBuffer struct {
	*glBuffer
}

func newVertexBuffer(size int, data []float32, usage Usage) *vertexBuffer {
	new_buffer := &vertexBuffer{
		glBuffer: newGlBuffer(gl.ARRAY_BUFFER, size*4, f32Bytes(data...), usage),
	}
	return new_buffer
}

func (buffer *vertexBuffer) fill(offset int, data []float32) {
	buffer.glBuffer.fill(offset*4, f32Bytes(data...))
}

func (buffer *vertexBuffer) getData() []float32 {
	return f32FromByte(buffer.glBuffer.data)
}

// Index Buffer

type indexBuffer struct {
	*glBuffer
}

func newIndexBuffer(size int, data []uint32, usage Usage) *indexBuffer {
	new_buffer := &indexBuffer{
		glBuffer: newGlBuffer(gl.ELEMENT_ARRAY_BUFFER, size*4, u32Bytes(data...), usage),
	}
	return new_buffer
}

func (buffer *indexBuffer) fill(offset int, data []uint32) {
	buffer.glBuffer.fill(offset*4, u32Bytes(data...))
}

func (buffer *indexBuffer) drawElements(mode uint32, offset, size int) {
	buffer.bind()
	defer buffer.unbind()
	gl.DrawElements(gl.Enum(mode), size*4, gl.UNSIGNED_INT, offset*4)
}

func (buffer *indexBuffer) getData() []uint32 {
	return u32FromByte(buffer.glBuffer.data)
}
