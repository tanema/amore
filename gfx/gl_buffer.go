package gfx

import (
	"math"

	"github.com/go-gl/gl/v2.1/gl"
)

type (
	glBuffer struct {
		is_bound               bool     // Whether the buffer is currently bound.
		is_mapped              bool     // Whether the buffer is currently mapped to main memory.
		size                   int      // The size of the buffer, in bytes.
		target                 uint32   // The target buffer object. (GL_ARRAY_BUFFER, GL_ELEMENT_ARRAY_BUFFER).
		usage                  Usage    // Usage hint. GL_[DYNAMIC, STATIC, STREAM]_DRAW.
		vbo                    uint32   // The VBO identifier. Assigned by OpenGL.
		memory_map             []uint32 // A pointer to mapped memory.
		modified_offset        int
		modified_size          int
		mapExplicitRangeModify bool
	}
	quadIndices struct {
		size        int
		indexBuffer *glBuffer
		indices     []uint32
	}
)

/**
 * QuadIndices manages one shared GLBuffer that stores the indices for an
 * element array. Vertex arrays using the vertex structure (or anything else
 * that can use the pattern below) can request a size and use it for the
 * drawElements call.
 *
 *  indices[i*6 + 0] = i*4 + 0;
 *  indices[i*6 + 1] = i*4 + 1;
 *  indices[i*6 + 2] = i*4 + 2;
 *
 *  indices[i*6 + 3] = i*4 + 2;
 *  indices[i*6 + 4] = i*4 + 1;
 *  indices[i*6 + 5] = i*4 + 3;
 *
 * There will always be a large enough GLBuffer around until all
 * QuadIndices instances have been deleted.
 *
 * Q: Why have something like QuadIndices?
 * A: The indices for the SpriteBatch do not change, only the array size
 * varies. Using one GLBuffer for all element arrays removes this
 * duplicated data and saves some memory.
 */
func newQuadIndices(size int) *quadIndices {
	new_qi := &quadIndices{
		size:        size,
		indexBuffer: newGlBuffer(size*6, gl.ELEMENT_ARRAY_BUFFER, gl.STATIC_DRAW, false),
		indices:     make([]uint32, size*6),
	}

	// 0----2
	// |  / |
	// | /  |
	// 1----3
	for i := 0; i < len(new_qi.indices); i++ {
		new_qi.indices[i*6+0] = uint32(i*4 + 0)
		new_qi.indices[i*6+1] = uint32(i*4 + 1)
		new_qi.indices[i*6+2] = uint32(i*4 + 2)

		new_qi.indices[i*6+3] = uint32(i*4 + 2)
		new_qi.indices[i*6+4] = uint32(i*4 + 1)
		new_qi.indices[i*6+5] = uint32(i*4 + 3)
	}

	new_qi.indexBuffer.bind()
	defer new_qi.indexBuffer.unbind()

	new_qi.indexBuffer.fill(0, new_qi.indexBuffer.size, new_qi.indices)

	return new_qi
}

func newGlBuffer(size int, target uint32, usage Usage, mapExplicitRangeModify bool) *glBuffer {
	new_buffer := &glBuffer{
		size:                   size,
		target:                 target,
		usage:                  usage,
		memory_map:             make([]uint32, size),
		mapExplicitRangeModify: mapExplicitRangeModify,
	}
	registerVolatile(new_buffer)
	return new_buffer
}

func (buffer *glBuffer) mapp() []uint32 {
	if buffer.is_mapped {
		return buffer.memory_map
	}
	buffer.is_mapped = true
	buffer.modified_offset = 0
	buffer.modified_size = 0
	return buffer.memory_map
}

func (buffer *glBuffer) unmapStatic() {
	if buffer.modified_size == 0 {
		return
	}
	// Upload the mapped data to the buffer.
	gl.BufferSubData(buffer.target, buffer.modified_offset, buffer.modified_size, gl.Ptr(buffer.memory_map[buffer.modified_offset]))
}

func (buffer *glBuffer) unmapStream() {
	// "orphan" current buffer to avoid implicit synchronisation on the GPU:
	// http://www.seas.upenn.edu/~pcozzi/OpenGLInsights/OpenGLInsights-AsynchronousBufferTransfers.pdf
	gl.BufferData(buffer.target, buffer.size, gl.Ptr(nil), uint32(buffer.usage))
	gl.BufferData(buffer.target, buffer.size, gl.Ptr(buffer.memory_map), uint32(buffer.usage))
}

func (buffer *glBuffer) unmap() {
	if !buffer.is_mapped {
		return
	}

	if buffer.mapExplicitRangeModify {
		buffer.modified_offset = int(math.Min(float64(buffer.modified_offset), float64(buffer.size-1)))
		buffer.modified_size = int(math.Min(float64(buffer.modified_size), float64(buffer.size-buffer.modified_offset)))
	} else {
		buffer.modified_offset = 0
		buffer.modified_size = buffer.size
	}

	// VBO::bind is a no-op when the VBO is mapped, so we have to make sure it's bound here.
	if !buffer.is_bound {
		gl.BindBuffer(buffer.target, buffer.vbo)
		buffer.is_bound = true
	}

	if buffer.modified_size > 0 {
		switch buffer.usage {
		case USAGE_STATIC:
			buffer.unmapStatic()
		case USAGE_STREAM:
			buffer.unmapStream()
		case USAGE_DYNAMIC:
			// It's probably more efficient to treat it like a streaming buffer if
			// at least a third of its contents have been modified during the map().
			if buffer.modified_size >= buffer.size/3 {
				buffer.unmapStream()
			} else {
				buffer.unmapStatic()
			}
		}
	}
	buffer.modified_offset = 0
	buffer.modified_size = 0
	buffer.is_mapped = false
}

func (buffer *glBuffer) setMappedRangeModified(offset, modifiedsize int) {
	if !buffer.is_mapped || buffer.mapExplicitRangeModify {
		return
	}

	// We're being conservative right now by internally marking the whole range
	// from the start of section a to the end of section b as modified if both
	// a and b are marked as modified.
	old_range_end := buffer.modified_offset + buffer.modified_size
	buffer.modified_offset = int(math.Min(float64(buffer.modified_offset), float64(offset)))

	new_range_end := int(math.Max(float64(offset+modifiedsize), float64(old_range_end)))
	buffer.modified_size = new_range_end - buffer.modified_offset
}

func (buffer *glBuffer) bind() {
	if !buffer.is_mapped {
		gl.BindBuffer(buffer.target, buffer.vbo)
		buffer.is_bound = true
	}
}

func (buffer *glBuffer) unbind() {
	if buffer.is_bound {
		gl.BindBuffer(buffer.target, 0)
	}
	buffer.is_bound = false
}

func (buffer *glBuffer) fill(offset, size int, data []uint32) {
	copy(buffer.memory_map[offset:], data[:size-1])
	if buffer.is_mapped {
		buffer.setMappedRangeModified(offset, size)
	} else {
		gl.BufferSubData(buffer.target, offset, size, gl.Ptr(data))
	}
}

func (buffer *glBuffer) loadVolatile() bool {
	gl.GenBuffers(1, &buffer.vbo)

	buffer.bind()
	defer buffer.unbind()

	// Note that if 'src' is '0', no data will be copied.
	gl.BufferData(buffer.target, buffer.size, gl.Ptr(buffer.memory_map), uint32(buffer.usage))

	return true
}

func (buffer *glBuffer) unloadVolatile() {
	buffer.is_mapped = false
	gl.DeleteBuffers(1, &buffer.vbo)
	buffer.vbo = 0
}
