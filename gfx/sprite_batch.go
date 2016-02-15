package gfx

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type SpriteBatch struct {
	size         int
	count        int
	color        *Color // Current color. This color, if present, will be applied to the next added sprite.
	array_buf    *vertexBuffer
	quad_indices *quadIndices
	usage        Usage
	texture      iTexture
}

func NewSpriteBatch(text iTexture, size int) *SpriteBatch {
	return NewSpriteBatchExt(text, size, USAGE_DYNAMIC)
}

func NewSpriteBatchExt(texture iTexture, size int, usage Usage) *SpriteBatch {
	return &SpriteBatch{
		texture:      texture,
		usage:        usage,
		color:        &Color{1, 1, 1, 1},
		array_buf:    newVertexBuffer(size*8, []float32{}, usage),
		quad_indices: newQuadIndices(size),
	}
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
	sprite_batch.array_buf = newVertexBuffer(sprite_batch.size*8, []float32{}, sprite_batch.usage)
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
	sprite_batch.array_buf = newVertexBuffer(newsize*8, sprite_batch.array_buf.data, sprite_batch.usage)
	sprite_batch.quad_indices = newQuadIndices(newsize)
	sprite_batch.size = newsize
	return nil
}

func (sprite_batch *SpriteBatch) GetBufferSize() int {
	return sprite_batch.size
}

func (sprite_batch *SpriteBatch) addv(verts []float32, mat *mgl32.Mat4, index int) error {
	if index == -1 && sprite_batch.count+1 == sprite_batch.size {
		return fmt.Errorf("Sprite Batch Buffer Full")
	}

	sprite_size := 8 * 4
	sprite := make([]float32, sprite_size)
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
		sprite_batch.array_buf.fill(sprite_batch.count*sprite_size, sprite)
		sprite_batch.count++
	} else {
		sprite_batch.array_buf.fill(index*sprite_size, sprite)
	}

	return nil
}

func (sprite_batch *SpriteBatch) Draw(args ...float32) {
	if sprite_batch.count == 0 {
		return
	}

	prepareDraw(generateModelMatFromArgs(args))
	bindTexture(sprite_batch.texture.GetHandle())

	sprite_batch.array_buf.bind()
	defer sprite_batch.array_buf.unbind()

	useVertexAttribArrays(ATTRIBFLAG_POS | ATTRIBFLAG_TEXCOORD | ATTRIBFLAG_COLOR)
	gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(2*4))
	gl.VertexAttribPointer(ATTRIB_COLOR, 4, gl.UNSIGNED_BYTE, true, 8*4, gl.PtrOffset(4*4))

	sprite_batch.quad_indices.drawElements(gl.TRIANGLES, 0, sprite_batch.count)
}
