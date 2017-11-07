package gfx

import (
	"image"
	// All image types have been imported for loading them
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/tanema/amore/file"
	"github.com/tanema/amore/gfx/gl"
)

// Image is an image that is drawable to the screen
type Image struct {
	*Texture
	filePath string
	mipmaps  bool
}

// NewImage will create a new texture for this image and return the *Image. If the
// file does not exist or cannot be decoded it will return an error.
func NewImage(path string) (*Image, error) {
	//we do this first time to check the image before volitile load
	imgFile, newErr := file.NewFile(path)
	defer imgFile.Close()
	if newErr != nil {
		return nil, newErr
	}

	decodedImg, _, imgErr := image.Decode(imgFile)
	if imgErr != nil {
		return nil, imgErr
	}

	bounds := decodedImg.Bounds()
	newImage := &Image{
		filePath: path,
		mipmaps:  false,
		Texture: &Texture{
			Width:  int32(bounds.Dx()),
			Height: int32(bounds.Dy()),
		},
	}

	registerVolatile(newImage)
	return newImage, nil
}

// NewMipmappedImage is like NewImage but the image is mipmapped
func NewMipmappedImage(path string) *Image {
	newImage := &Image{
		filePath: path,
		mipmaps:  true,
	}
	registerVolatile(newImage)
	return newImage
}

// loadVolatile will create the volatile objects
func (img *Image) loadVolatile() bool {
	imgFile, newErr := file.NewFile(img.filePath)
	defer imgFile.Close()
	if newErr != nil {
		return false
	}

	decodedImg, _, imgErr := image.Decode(imgFile)
	if imgErr != nil {
		return false
	}

	img.Texture, imgErr = newImageTexture(decodedImg, img.mipmaps)
	if imgErr != nil {
		return false
	}

	return true
}
