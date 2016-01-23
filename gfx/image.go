package gfx

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/go-gl/gl/v2.1/gl"

	"github.com/tanema/amore/file"
)

type Image struct {
	texture   *Texture
	img       image.Image
	coords    []float32
	texcoords []float32
}

func NewImage(path string) (*Image, error) {
	imgFile, err := file.NewFile(path)
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
	img.coords = []float32{
		0, 0,
		0, img.texture.Height,
		img.texture.Width, img.texture.Height,
		img.texture.Width, 0,
	}
	img.texcoords = []float32{
		0, 0,
		0, 1,
		1, 1,
		1, 0,
	}
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
	prepareDraw(generateModelMatFromArgs(args))
	bindTexture(image.texture.GetHandle())

	gl.EnableVertexAttribArray(ATTRIB_POS)
	gl.EnableVertexAttribArray(ATTRIB_TEXCOORD)

	gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 0, gl.Ptr(image.coords))
	gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, 0, gl.Ptr(image.texcoords))

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

	gl.DisableVertexAttribArray(ATTRIB_TEXCOORD)
	gl.DisableVertexAttribArray(ATTRIB_POS)
}
