package gfx

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
)

type Image struct {
	texture *Texture
	img     image.Image
}

func NewImage(path string) (*Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	new_image := &Image{
		img: img,
	}

	Register(new_image)

	return new_image, nil
}

func (img *Image) LoadVolatile() bool {
	var err error
	img.texture, err = LoadImageTexture(img.img)
	return err == nil
}

func (image *Image) UnloadVolatile() {
	if image.texture != nil {
		return
	}

	image.texture.Release()
	image.texture = nil
}

func (image *Image) Draw(x, y, angle, sx, sy, ox, oy, kx, ky float64) {
	PrepareDraw()
	BindTexture(image.texture.GetHandle())
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
}
