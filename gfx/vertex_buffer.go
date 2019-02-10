package gfx

import (
	"encoding/binary"
	"math"

	"github.com/goxjs/gl"
	"golang.org/x/mobile/exp/f32"
)

type vertexBuffer struct {
	isBound        bool      // Whether the buffer is currently bound.
	usage          Usage     // Usage hint. GL_[DYNAMIC, STATIC, STREAM]_DRAW.
	vbo            gl.Buffer // The VBO identifier. Assigned by OpenGL.
	data           []float32 // A pointer to mapped memory.
	modifiedOffset int
	modifiedSize   int
}

func newVertexBuffer(size int, data []float32, usage Usage) *vertexBuffer {
	newBuffer := &vertexBuffer{
		usage: usage,
		data:  make([]float32, size),
	}
	if len(data) > 0 {
		copy(newBuffer.data, data[:size])
	}
	registerVolatile(newBuffer)
	return newBuffer
}

func (buffer *vertexBuffer) bufferStatic() {
	if buffer.modifiedSize == 0 {
		return
	}
	// Upload the mapped data to the buffer.
	f32.Bytes(binary.LittleEndian, buffer.data[buffer.modifiedOffset:buffer.modifiedSize]...)
}

func (buffer *vertexBuffer) bufferStream() {
	gl.BufferData(gl.ARRAY_BUFFER, f32.Bytes(binary.LittleEndian, buffer.data...), gl.Enum(buffer.usage))
}

func (buffer *vertexBuffer) bufferData() {
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

func (buffer *vertexBuffer) bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.vbo)
	buffer.isBound = true
}

func (buffer *vertexBuffer) unbind() {
	if buffer.isBound {
		gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{})
	}
	buffer.isBound = false
}

func (buffer *vertexBuffer) fill(offset int, data []float32) {
	copy(buffer.data[offset:], data[:])
	if !buffer.vbo.Valid() {
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

func (buffer *vertexBuffer) loadVolatile() bool {
	buffer.vbo = gl.CreateBuffer()
	buffer.bind()
	defer buffer.unbind()
	gl.BufferData(gl.ARRAY_BUFFER, f32.Bytes(binary.LittleEndian, buffer.data...), gl.Enum(buffer.usage))
	return true
}

func (buffer *vertexBuffer) unloadVolatile() {
	gl.DeleteBuffer(buffer.vbo)
	buffer.vbo.Value = 0
}
