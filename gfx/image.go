package gfx

import (
	"image"
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

// NewMipmappedImage is like NewImage but the image is mipmapped
func NewMipmappedImage(path string) *Image {
	new_image := &Image{
		filePath: path,
		mipmaps:  true,
	}
	registerVolatile(new_image)
	return new_image
}

// loadVolatile will create the volatile objects
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

// Drawv will take raw verticies so that you can draw the image on a polygon, specifying
// the image coords.
func (img *Image) Drawv(vertices, textCoords []float32) {
	prepareDraw(nil)
	bindTexture(img.Texture.getHandle())
	useVertexAttribArrays(attribflag_pos | attribflag_texcoord)

	gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 0, gl.Ptr(vertices))
	gl.VertexAttribPointer(attrib_texcoord, 2, gl.FLOAT, false, 0, gl.Ptr(textCoords))

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}
