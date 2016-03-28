package gfx

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
)

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

func NewSpriteBatch(text iTexture, size int) *SpriteBatch {
	return NewSpriteBatchExt(text, size, USAGE_DYNAMIC)
}

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

func (sprite_batch *SpriteBatch) Release() {
	sprite_batch.array_buf.Release()
	sprite_batch.quad_indices.Release()
}

func (sprite_batch *SpriteBatch) Add(args ...float32) error {
	return sprite_batch.addv(sprite_batch.texture.getVerticies(), generateModelMatFromArgs(args), -1)
}

func (sprite_batch *SpriteBatch) Addq(quad *Quad, args ...float32) error {
	return sprite_batch.addv(quad.getVertices(), generateModelMatFromArgs(args), -1)
}

func (sprite_batch *SpriteBatch) Set(index int, args ...float32) error {
	return sprite_batch.addv(sprite_batch.texture.getVerticies(), generateModelMatFromArgs(args), index)
}

func (sprite_batch *SpriteBatch) Setq(index int, quad *Quad, args ...float32) error {
	return sprite_batch.addv(quad.getVertices(), generateModelMatFromArgs(args), index)
}

func (sprite_batch *SpriteBatch) Clear() {
	if sprite_batch.array_buf != nil {
		releaseVolatile(sprite_batch.array_buf)
	}
	sprite_batch.array_buf = newVertexBuffer(sprite_batch.size*4*8, []float32{}, sprite_batch.usage)
	sprite_batch.count = 0
}

func (sprite_batch *SpriteBatch) flush() {
	sprite_batch.array_buf.bufferData()
}

func (sprite_batch *SpriteBatch) SetTexture(newtexture iTexture) {
	sprite_batch.texture = newtexture
}

func (sprite_batch *SpriteBatch) GetTexture() iTexture {
	return sprite_batch.texture
}

func (sprite_batch *SpriteBatch) SetColor(color *Color) {
	sprite_batch.color = color
}

func (sprite_batch *SpriteBatch) ClearColor() {
	sprite_batch.color = &Color{1, 1, 1, 1}
}

func (sprite_batch *SpriteBatch) GetColor() *Color {
	return sprite_batch.color
}

func (sprite_batch *SpriteBatch) GetCount() int {
	return sprite_batch.count
}

func (sprite_batch *SpriteBatch) SetBufferSize(newsize int) error {
	if newsize <= 0 {
		fmt.Errorf("Invalid SpriteBatch size.")
	} else if newsize == sprite_batch.size {
		return nil
	}

	sprite_batch.array_buf.Release()
	sprite_batch.array_buf = newVertexBuffer(newsize*4*8, sprite_batch.array_buf.getData(), sprite_batch.usage)

	sprite_batch.quad_indices.Release()
	sprite_batch.quad_indices = newQuadIndices(newsize)

	sprite_batch.size = newsize
	return nil
}

func (sprite_batch *SpriteBatch) GetBufferSize() int {
	return sprite_batch.size
}

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

func (sprite_batch *SpriteBatch) SetDrawRange(min, max int) error {
	if min < 0 || max < 0 || min > max {
		return fmt.Errorf("Invalid draw range.")
	}
	sprite_batch.rangeMin = min
	sprite_batch.rangeMax = max
	return nil
}

func (sprite_batch *SpriteBatch) ClearDrawRange() {
	sprite_batch.rangeMin = -1
	sprite_batch.rangeMax = -1
}

func (sprite_batch *SpriteBatch) GetDrawRange() (int, int) {
	min := 0
	max := sprite_batch.count - 1
	if sprite_batch.rangeMax >= 0 {
		max = int(math.Min(float64(sprite_batch.rangeMax), float64(max)))
	}
	if sprite_batch.rangeMin >= 0 {
		min = int(math.Min(float64(sprite_batch.rangeMin), float64(max)))
	}
	return min, max
}

func (sprite_batch *SpriteBatch) Draw(args ...float32) {
	if sprite_batch.count == 0 {
		return
	}

	prepareDraw(generateModelMatFromArgs(args))
	bindTexture(sprite_batch.texture.GetHandle())
	enableVertexAttribArrays(ATTRIB_POS, ATTRIB_TEXCOORD, ATTRIB_COLOR)
	defer disableVertexAttribArrays(ATTRIB_POS, ATTRIB_TEXCOORD, ATTRIB_COLOR)

	sprite_batch.array_buf.bind()
	defer sprite_batch.array_buf.unbind()

	gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 8*4, 0)
	gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, 8*4, 2*4)
	gl.VertexAttribPointer(ATTRIB_COLOR, 4, gl.FLOAT, false, 8*4, 4*4)

	min, max := sprite_batch.GetDrawRange()
	sprite_batch.quad_indices.drawElements(gl.TRIANGLES, min, max-min+1)
}
