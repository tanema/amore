package gfx

import (
	"github.com/go-gl/gl/v2.1/gl"
)

type Image struct {
	texture *Texture
}

func NewImage(path string) (*Image, error) {
	new_text, err := NewTexture(path)
	if err != nil {
		return nil, err
	}
	return &Image{texture: new_text}, nil
}

func (image *Image) Draw(x, y, angle, sx, sy, ox, oy, kx, ky float64) {
	image.texture.Bind(func() {
		gl.Begin(gl.QUADS)
		{
			gl.TexCoord2d(0, 0) // top-left
			gl.Vertex2d(x, y)
			gl.TexCoord2d(0, 1) // bottom-left
			gl.Vertex2d(x, y+image.texture.Height)
			gl.TexCoord2d(1, 1) // bottom-right
			gl.Vertex2d(x+image.texture.Width, y+image.texture.Height)
			gl.TexCoord2d(1, 0) // top-right
			gl.Vertex2d(x+image.texture.Width, y)
		}
		gl.End()
	})
}
