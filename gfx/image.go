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
	filePath string
	mipmaps  bool
}

// NewImage will create a new texture for this image and return the *Image. If the
// file does not exist or cannot be decoded it will return an error.
func NewImage(path string, mipmapped bool) *Image {
	newImage := &Image{filePath: path, mipmaps: mipmapped}
	registerVolatile(newImage)
	return newImage
}

// loadVolatile will create the volatile objects
func (img *Image) loadVolatile() bool {
	if img.filePath == "" {
		return false
	}

	imgFile, err := file.NewFile(img.filePath)
	if err != nil {
		return false
	}
	defer imgFile.Close()

	decodedImg, _, err := image.Decode(imgFile)
	if err != nil || decodedImg == nil {
		return false
	}

	img.Texture = newImageTexture(decodedImg, img.mipmaps)
	return true
}
