package gfx

import (
	"image"
	// All image types have been imported for loading them
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/tanema/amore/file"
)

// Image is an image that is drawable to the screen
type Image struct {
	*Texture
	img      image.Image
	filePath string
	mipmaps  bool
}

// NewImage will create a new texture for this image and return the *Image. If the
// file does not exist or cannot be decoded it will return an error.
func NewImage(path string) (*Image, error) {
	//we do this first time to check the image before volitile load
	decodedImg, err := loadImageFromPath(path)
	if err != nil {
		return nil, err
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

// NewImageFrom will create a new texture for this image and return the *Image.
func NewImageFrom(img image.Image) *Image {
	bounds := img.Bounds()
	newImage := &Image{
		img:     img,
		mipmaps: false,
		Texture: &Texture{
			Width:  int32(bounds.Dx()),
			Height: int32(bounds.Dy()),
		},
	}
	registerVolatile(newImage)
	return newImage
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
	var err error
	decodedImg := img.img

	if img.filePath != "" {
		if decodedImg, err = loadImageFromPath(img.filePath); err != nil {
			return false
		}
	}

	if decodedImg == nil {
		return false
	}

	if img.Texture, err = newImageTexture(decodedImg, img.mipmaps); err != nil {
		return false
	}

	img.img = nil

	return true
}

func loadImageFromPath(path string) (image.Image, error) {
	imgFile, newErr := file.NewFile(path)
	defer imgFile.Close()
	if newErr != nil {
		return nil, newErr
	}

	decodedImg, _, imgErr := image.Decode(imgFile)
	return decodedImg, imgErr
}
