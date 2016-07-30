package gfx

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/tanema/amore/gfx/gl"
	"github.com/tanema/amore/mth"
)

// SpriteBatch is a collection of images/quads/textures all drawn with a single draw call
type SpriteBatch struct {
	size         int
	count        int
	color        *Color // Current color. This color, if present, will be applied to the next added sprite.
	array_buf    *vertexBuffer
	quad_indices *quadIndices
	usage        Usage
	texture      iTexture
	rangeMin     int
	rangeMax     int
}

// NewSpriteBatch will generate a new batch with the size provided and bind the texture to it.
func NewSpriteBatch(text iTexture, size int) *SpriteBatch {
	return NewSpriteBatchExt(text, size, USAGE_DYNAMIC)
}

// NewSpriteBatchExt is like NewSpriteBatch but allows you to set the usage.
func NewSpriteBatchExt(texture iTexture, size int, usage Usage) *SpriteBatch {
	return &SpriteBatch{
		size:         size,
		texture:      texture,
		usage:        usage,
		color:        &Color{1, 1, 1, 1},
		array_buf:    newVertexBuffer(size*4*8, []float32{}, usage),
		quad_indices: newQuadIndices(size),
		rangeMin:     -1,
		rangeMax:     -1,
	}
}

// Release cleans up the gl object associates with the sprite batch. This should
// only be done when discarding this object.
func (sprite_batch *SpriteBatch) Release() {
	sprite_batch.array_buf.Release()
	sprite_batch.quad_indices.Release()
}

// Add adds a sprite to the batch. Sprites are drawn in the order they are added.
// x, y The position to draw the object
// r rotation of the object
// sx, sy scale of the object
// ox, oy offset of the object
// kx, ky shear of the object
func (sprite_batch *SpriteBatch) Add(args ...float32) error {
	return sprite_batch.addv(sprite_batch.texture.getVerticies(), generateModelMatFromArgs(args), -1)
}

// Adds a Quad to the batch. This is very useful for something like a tilemap.
func (sprite_batch *SpriteBatch) Addq(quad *Quad, args ...float32) error {
	return sprite_batch.addv(quad.getVertices(), generateModelMatFromArgs(args), -1)
}

// Set changes a sprite in the batch with the same arguments as add
func (sprite_batch *SpriteBatch) Set(index int, args ...float32) error {
	return sprite_batch.addv(sprite_batch.texture.getVerticies(), generateModelMatFromArgs(args), index)
}

// Set changes a sprite in the batch with the same arguments as addq
func (sprite_batch *SpriteBatch) Setq(index int, quad *Quad, args ...float32) error {
	return sprite_batch.addv(quad.getVertices(), generateModelMatFromArgs(args), index)
}

// Clear will remove all the sprites from the batch
func (sprite_batch *SpriteBatch) Clear() {
	if sprite_batch.array_buf != nil {
		releaseVolatile(sprite_batch.array_buf)
	}
	sprite_batch.array_buf = newVertexBuffer(sprite_batch.size*4*8, []float32{}, sprite_batch.usage)
	sprite_batch.count = 0
}

// flush will ensure the data is uploaded to the buffer
func (sprite_batch *SpriteBatch) flush() {
	sprite_batch.array_buf.bufferData()
}

// SetTexture will change the texture of the batch to a new one
func (sprite_batch *SpriteBatch) SetTexture(newtexture iTexture) {
	sprite_batch.texture = newtexture
}

// GetTexture will return the currently bound texture of this sprite batch.
func (sprite_batch *SpriteBatch) GetTexture() iTexture {
	return sprite_batch.texture
}

// SetColor will set the color that will be used for the next add or set operations.
func (sprite_batch *SpriteBatch) SetColor(color *Color) {
	sprite_batch.color = color
}

// ClearColor will reset the color back to white
func (sprite_batch *SpriteBatch) ClearColor() {
	sprite_batch.color = &Color{1, 1, 1, 1}
}

// GetColor will return the currently used color.
func (sprite_batch *SpriteBatch) GetColor() *Color {
	return sprite_batch.color
}

// GetCount will return the amount of sprites already added to the batch
func (sprite_batch *SpriteBatch) GetCount() int {
	return sprite_batch.count
}

// SetBufferSize will resize the buffer, change the limit of sprites you can add
// to this batch.
func (sprite_batch *SpriteBatch) SetBufferSize(newsize int) error {
	if newsize <= 0 {
		fmt.Errorf("Invalid SpriteBatch size.")
	} else if newsize == sprite_batch.size {
		return nil
	}

	sprite_batch.array_buf.Release()
	sprite_batch.array_buf = newVertexBuffer(newsize*4*8, sprite_batch.array_buf.data, sprite_batch.usage)

	sprite_batch.quad_indices.Release()
	sprite_batch.quad_indices = newQuadIndices(newsize)

	sprite_batch.size = newsize
	return nil
}

// GetBufferSize will return the limit of sprites you can add to this batch.
func (sprite_batch *SpriteBatch) GetBufferSize() int {
	return sprite_batch.size
}

// addv will add a sprite to the batch using the verts, a transform and an index to
// place it
func (sprite_batch *SpriteBatch) addv(verts []float32, mat *mgl32.Mat4, index int) error {
	if index == -1 && sprite_batch.count >= sprite_batch.size {
		return fmt.Errorf("Sprite Batch Buffer Full")
	}

	sprite := make([]float32, 8*4)
	for i := 0; i < 32; i += 8 {
		j := (i / 2)
		sprite[i+0] = (mat[0] * verts[j+0]) + (mat[4] * verts[j+1]) + mat[12]
		sprite[i+1] = (mat[1] * verts[j+0]) + (mat[5] * verts[j+1]) + mat[13]
		sprite[i+2] = verts[j+2]
		sprite[i+3] = verts[j+3]
		sprite[i+4] = sprite_batch.color[0]
		sprite[i+5] = sprite_batch.color[1]
		sprite[i+6] = sprite_batch.color[2]
		sprite[i+7] = sprite_batch.color[3]
	}

	if index == -1 {
		sprite_batch.array_buf.fill(sprite_batch.count*4*8, sprite)
		sprite_batch.count++
	} else {
		sprite_batch.array_buf.fill(index*4*8, sprite)
	}

	return nil
}

// SetDrawRange will set a range in the points to draw. This is useful if you only
// need to render a portion of the batch.
func (sprite_batch *SpriteBatch) SetDrawRange(min, max int) error {
	if min < 0 || max < 0 || min > max {
		return fmt.Errorf("Invalid draw range.")
	}
	sprite_batch.rangeMin = min
	sprite_batch.rangeMax = max
	return nil
}

// ClearDrawRange will reset the draw range if you want to draw the whole batch again.
func (sprite_batch *SpriteBatch) ClearDrawRange() {
	sprite_batch.rangeMin = -1
	sprite_batch.rangeMax = -1
}

// GetDrawRange will return the min, max range set on the batch. If no range is set
// the range will return -1, -1
func (sprite_batch *SpriteBatch) GetDrawRange() (int, int) {
	min := 0
	max := sprite_batch.count - 1
	if sprite_batch.rangeMax >= 0 {
		max = mth.Mini(sprite_batch.rangeMax, max)
	}
	if sprite_batch.rangeMin >= 0 {
		min = mth.Mini(sprite_batch.rangeMin, max)
	}
	return min, max
}

// Draw satisfies the Drawable interface. Inputs are as follows
// x, y, r, sx, sy, ox, oy, kx, ky
// x, y are position
// r is rotation
// sx, sy is the scale, if sy is not given sy will equal sx
// ox, oy are offset
// kx, ky are the shear. If ky is not given ky will equal kx
func (sprite_batch *SpriteBatch) Draw(args ...float32) {
	if sprite_batch.count == 0 {
		return
	}

	prepareDraw(generateModelMatFromArgs(args))
	bindTexture(sprite_batch.texture.getHandle())
	useVertexAttribArrays(attribflag_pos | attribflag_texcoord | attribflag_color)

	sprite_batch.array_buf.bind()
	defer sprite_batch.array_buf.unbind()

	gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(attrib_texcoord, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(2*4))
	gl.VertexAttribPointer(attrib_color, 4, gl.FLOAT, false, 8*4, gl.PtrOffset(4*4))

	min, max := sprite_batch.GetDrawRange()
	sprite_batch.quad_indices.drawElements(gl.TRIANGLES, min, max-min+1)
}
