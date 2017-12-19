// Package surface is a simple package meant to load sdl surfaces without sdl_image.
// This enables removing the sdl_image dependancy
package surface

import (
	"image"
	"image/draw"
	// loading all the image libs to support loading
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"runtime"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/tanema/amore/file"
)

// Load creates a *sdl.Surface from an image file. Make sure to call surface.Free()
// on any surface you load. This method will call an error if the file does not exist
// or cannot decode it (not an image)
func Load(path string) (*sdl.Surface, error) {
	imgFile, newErr := file.NewFile(path)
	defer imgFile.Close()
	if newErr != nil {
		return &sdl.Surface{}, newErr
	}

	decodedImg, _, imgErr := image.Decode(imgFile)
	if imgErr != nil {
		return &sdl.Surface{}, imgErr
	}

	bounds := decodedImg.Bounds()
	rgba := image.NewRGBA(decodedImg.Bounds())
	draw.Draw(rgba, bounds, decodedImg, image.Point{0, 0}, draw.Src)

	var rmask, gmask, bmask, amask uint32
	switch runtime.GOARCH {
	case "mips64", "ppc64":
		rmask = 0xFF000000
		gmask = 0x00FF0000
		bmask = 0x0000FF00
		amask = 0x000000FF
	default:
		rmask = 0x000000FF
		gmask = 0x0000FF00
		bmask = 0x00FF0000
		amask = 0xFF000000
	}

	return sdl.CreateRGBSurfaceFrom(
		unsafe.Pointer(&rgba.Pix[0]),
		int32(bounds.Dx()),
		int32(bounds.Dy()),
		32,
		bounds.Dx()*4,
		rmask,
		gmask,
		bmask,
		amask,
	)
}
