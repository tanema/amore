package gfx

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/tanema/amore/file"
)

type Image struct {
	*Texture
	filePath string
	mipmaps  bool
}

func NewImage(path string) (*Image, error) {
	//we do this first time to check the image before volitile load
	imgFile, new_err := file.NewFile(path)
	defer imgFile.Close()
	if new_err != nil {
		return nil, new_err
	}

	decoded_img, _, img_err := image.Decode(imgFile)
	if img_err != nil {
		return nil, img_err
	}

	bounds := decoded_img.Bounds()
	new_image := &Image{
		filePath: path,
		mipmaps:  false,
		Texture: &Texture{
			Width:  int32(bounds.Dx()),
			Height: int32(bounds.Dy()),
		},
	}

	registerVolatile(new_image)
	return new_image, nil
}

func NewMipmappedImage(path string) *Image {
	new_image := &Image{
		filePath: path,
		mipmaps:  true,
	}
	registerVolatile(new_image)
	return new_image
}

func (img *Image) loadVolatile() bool {
	imgFile, new_err := file.NewFile(img.filePath)
	defer imgFile.Close()
	if new_err != nil {
		return false
	}

	decoded_img, _, img_err := image.Decode(imgFile)
	if img_err != nil {
		return false
	}

	img.Texture, img_err = newImageTexture(decoded_img, img.mipmaps)
	if img_err != nil {
		return false
	}

	return true
}
