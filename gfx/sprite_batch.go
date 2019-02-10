package gfx

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/goxjs/gl"
)

// SpriteBatch is a collection of images/quads/textures all drawn with a single draw call
type SpriteBatch struct {
	size        int
	count       int
	color       []float32 // Current color. This color, if present, will be applied to the next added sprite.
	arrayBuf    *vertexBuffer
	quadIndices *quadIndices
	usage       Usage
	texture     ITexture
	rangeMin    int
	rangeMax    int
}

// NewSpriteBatch will generate a new batch with the size provided and bind the texture to it.
func NewSpriteBatch(text ITexture, size int) *SpriteBatch {
	return NewSpriteBatchExt(text, size, UsageDynamic)
}

// NewSpriteBatchExt is like NewSpriteBatch but allows you to set the usage.
func NewSpriteBatchExt(texture ITexture, size int, usage Usage) *SpriteBatch {
	return &SpriteBatch{
		size:        size,
		texture:     texture,
		usage:       usage,
		color:       []float32{1, 1, 1, 1},
		arrayBuf:    newVertexBuffer(size*4*8, []float32{}, usage),
		quadIndices: newQuadIndices(size),
		rangeMin:    -1,
		rangeMax:    -1,
	}
}

// Add adds a sprite to the batch. Sprites are drawn in the order they are added.
// x, y The position to draw the object
// r rotation of the object
// sx, sy scale of the object
// ox, oy offset of the object
// kx, ky shear of the object
func (spriteBatch *SpriteBatch) Add(args ...float32) error {
	return spriteBatch.addv(spriteBatch.texture.getVerticies(), generateModelMatFromArgs(args), -1)
}

// Addq adds a Quad to the batch. This is very useful for something like a tilemap.
func (spriteBatch *SpriteBatch) Addq(quad *Quad, args ...float32) error {
	return spriteBatch.addv(quad.getVertices(), generateModelMatFromArgs(args), -1)
}

// Set changes a sprite in the batch with the same arguments as add
func (spriteBatch *SpriteBatch) Set(index int, args ...float32) error {
	return spriteBatch.addv(spriteBatch.texture.getVerticies(), generateModelMatFromArgs(args), index)
}

// Setq changes a sprite in the batch with the same arguments as addq
func (spriteBatch *SpriteBatch) Setq(index int, quad *Quad, args ...float32) error {
	return spriteBatch.addv(quad.getVertices(), generateModelMatFromArgs(args), index)
}

// Clear will remove all the sprites from the batch
func (spriteBatch *SpriteBatch) Clear() {
	spriteBatch.arrayBuf = newVertexBuffer(spriteBatch.size*4*8, []float32{}, spriteBatch.usage)
	spriteBatch.count = 0
}

// flush will ensure the data is uploaded to the buffer
func (spriteBatch *SpriteBatch) flush() {
	spriteBatch.arrayBuf.bufferData()
}

// SetTexture will change the texture of the batch to a new one
func (spriteBatch *SpriteBatch) SetTexture(newtexture ITexture) {
	spriteBatch.texture = newtexture
}

// GetTexture will return the currently bound texture of this sprite batch.
func (spriteBatch *SpriteBatch) GetTexture() ITexture {
	return spriteBatch.texture
}

// SetColor will set the color that will be used for the next add or set operations.
func (spriteBatch *SpriteBatch) SetColor(vals ...float32) {
	spriteBatch.color = vals
}

// ClearColor will reset the color back to white
func (spriteBatch *SpriteBatch) ClearColor() {
	spriteBatch.color = []float32{1, 1, 1, 1}
}

// GetColor will return the currently used color.
func (spriteBatch *SpriteBatch) GetColor() []float32 {
	return spriteBatch.color
}

// GetCount will return the amount of sprites already added to the batch
func (spriteBatch *SpriteBatch) GetCount() int {
	return spriteBatch.count
}

// SetBufferSize will resize the buffer, change the limit of sprites you can add
// to this batch.
func (spriteBatch *SpriteBatch) SetBufferSize(newsize int) error {
	if newsize <= 0 {
		return fmt.Errorf("invalid SpriteBatch size")
	} else if newsize == spriteBatch.size {
		return nil
	}
	spriteBatch.arrayBuf = newVertexBuffer(newsize*4*8, spriteBatch.arrayBuf.data, spriteBatch.usage)
	spriteBatch.quadIndices = newQuadIndices(newsize)
	spriteBatch.size = newsize
	return nil
}

// GetBufferSize will return the limit of sprites you can add to this batch.
func (spriteBatch *SpriteBatch) GetBufferSize() int {
	return spriteBatch.size
}

// addv will add a sprite to the batch using the verts, a transform and an index to
// place it
func (spriteBatch *SpriteBatch) addv(verts []float32, mat *mgl32.Mat4, index int) error {
	if index == -1 && spriteBatch.count >= spriteBatch.size {
		return fmt.Errorf("Sprite Batch Buffer Full")
	}

	sprite := make([]float32, 8*4)
	for i := 0; i < 32; i += 8 {
		j := (i / 2)
		sprite[i+0] = (mat[0] * verts[j+0]) + (mat[4] * verts[j+1]) + mat[12]
		sprite[i+1] = (mat[1] * verts[j+0]) + (mat[5] * verts[j+1]) + mat[13]
		sprite[i+2] = verts[j+2]
		sprite[i+3] = verts[j+3]
		sprite[i+4] = spriteBatch.color[0]
		sprite[i+5] = spriteBatch.color[1]
		sprite[i+6] = spriteBatch.color[2]
		sprite[i+7] = spriteBatch.color[3]
	}

	if index == -1 {
		spriteBatch.arrayBuf.fill(spriteBatch.count*4*8, sprite)
		spriteBatch.count++
	} else {
		spriteBatch.arrayBuf.fill(index*4*8, sprite)
	}

	return nil
}

// SetDrawRange will set a range in the points to draw. This is useful if you only
// need to render a portion of the batch.
func (spriteBatch *SpriteBatch) SetDrawRange(min, max int) error {
	if min < 0 || max < 0 || min > max {
		return fmt.Errorf("invalid draw range")
	}
	spriteBatch.rangeMin = min
	spriteBatch.rangeMax = max
	return nil
}

// ClearDrawRange will reset the draw range if you want to draw the whole batch again.
func (spriteBatch *SpriteBatch) ClearDrawRange() {
	spriteBatch.rangeMin = -1
	spriteBatch.rangeMax = -1
}

// GetDrawRange will return the min, max range set on the batch. If no range is set
// the range will return -1, -1
func (spriteBatch *SpriteBatch) GetDrawRange() (int, int) {
	min := 0
	max := spriteBatch.count - 1
	if spriteBatch.rangeMax >= 0 {
		max = int(math.Min(float64(spriteBatch.rangeMax), float64(max)))
	}
	if spriteBatch.rangeMin >= 0 {
		min = int(math.Min(float64(spriteBatch.rangeMin), float64(max)))
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
func (spriteBatch *SpriteBatch) Draw(args ...float32) {
	if spriteBatch.count == 0 {
		return
	}

	prepareDraw(generateModelMatFromArgs(args))
	bindTexture(spriteBatch.texture.getHandle())
	useVertexAttribArrays(shaderPosFlag | shaderTexCoordFlag | shaderColorFlag)

	spriteBatch.arrayBuf.bind()
	defer spriteBatch.arrayBuf.unbind()

	gl.VertexAttribPointer(gl.Attrib{Value: 0}, 2, gl.FLOAT, false, 8*4, 0)
	gl.VertexAttribPointer(gl.Attrib{Value: 1}, 2, gl.FLOAT, false, 8*4, 2*4)
	gl.VertexAttribPointer(gl.Attrib{Value: 2}, 4, gl.FLOAT, false, 8*4, 4*4)

	min, max := spriteBatch.GetDrawRange()
	spriteBatch.quadIndices.drawElements(gl.TRIANGLES, min, max-min+1)
}
