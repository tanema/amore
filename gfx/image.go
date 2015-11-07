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

func (image *Image) Draw(args ...float32) {
	x, y, angle, sx, sy, ox, oy, kx, ky := normalizeDrawCallArgs(args)

	BindTexture(image.texture.GetHandle())

	gl.EnableVertexAttribArray(ATTRIB_POS)
	gl.EnableVertexAttribArray(ATTRIB_TEXCOORD)

	gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 0, gl.Ptr([]float32{
		x, y,
		x, y + image.texture.Height,
		x + image.texture.Width, y + image.texture.Height,
		x + image.texture.Width, y,
	}))
	gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, 0, gl.Ptr([]float32{
		0, 0,
		0, 1,
		1, 1,
		1, 0,
	}))

	PrepareDraw()
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

	gl.DisableVertexAttribArray(ATTRIB_TEXCOORD)
	gl.DisableVertexAttribArray(ATTRIB_POS)
}
