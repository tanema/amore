package gfx

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/tanema/amore/file"
)

type Image struct {
	*Texture
	filePath string
}

func NewImage(path string) (*Image, error) {
	new_image := &Image{filePath: path}
	Register(new_image)
	return new_image, nil
}

func (img *Image) LoadVolatile() bool {
	imgFile, new_err := file.NewFile(img.filePath)
	defer imgFile.Close()
	if new_err != nil {
		return false
	}

	decoded_img, _, img_err := image.Decode(imgFile)
	if img_err != nil {
		return false
	}

	img.Texture, img_err = LoadImageTexture(decoded_img)
	if img_err != nil {
		return false
	}

	img.generateVerticies()
	return true
}
