package gfx

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"runtime"

	"github.com/tanema/amore/gfx/gl"
	"github.com/tanema/amore/gfx/mat"
)

type (
	// Filter is a representation of texture filtering that contains both min and
	// mag filter modes
	Filter struct {
		min, mag, mipmap FilterMode
		anisotropy       float32
	}
	// Wrap is a representation of texture rapping containing both s and t wrap
	Wrap struct {
		s, t WrapMode
	}
	// Texture is a struct to wrap the opengl texture object
	Texture struct {
		textureId     gl.Texture
		Width, Height int32
		vertices      []float32
		filter        Filter
		wrap          Wrap
		mipmaps       bool
	}
	// iTexture is an interface for any object that can be used like a texture.
	iTexture interface {
		getHandle() gl.Texture
		GetWidth() int32
		GetHeight() int32
		getVerticies() []float32
		Release()
	}
)

// newFilter will create a Filter with default values
func newFilter() Filter {
	return Filter{
		min:        FILTER_LINEAR,
		mag:        FILTER_LINEAR,
		mipmap:     FILTER_NONE,
		anisotropy: 1.0,
	}
}

// newTexture will return a new generated texture will not data uploaded to it.
func newTexture(width, height int32, mipmaps bool) *Texture {
	new_texture := &Texture{
		textureId: gl.CreateTexture(),
		Width:     width,
		Height:    height,
		wrap:      Wrap{s: WRAP_CLAMP, t: WRAP_CLAMP},
		filter:    newFilter(),
		mipmaps:   mipmaps,
	}

	new_texture.SetFilter(FILTER_NEAREST, FILTER_NEAREST)
	new_texture.SetWrap(WRAP_CLAMP, WRAP_CLAMP)

	if new_texture.mipmaps {
		new_texture.filter.mipmap = states.back().defaultMipmapFilter
		new_texture.SetMipmapSharpness(states.back().defaultMipmapSharpness)
	}

	if new_texture.mipmaps {
		// Auto-generate mipmaps every time the texture is modified, if glGenerateMipmap isn't supported.
		setTexMipMap()
	}

	new_texture.generateVerticies()

	return new_texture
}

// newImageTexture will generate a texture from an image. It will automatically
// upload the image data to the texture.
func newImageTexture(img image.Image, mipmaps bool) (*Texture, error) {
	bounds := img.Bounds()
	new_texture := newTexture(int32(bounds.Dx()), int32(bounds.Dy()), mipmaps)
	//generate a uniform image and upload to vram
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, bounds, img, image.Point{0, 0}, draw.Src)
	bindTexture(new_texture.getHandle())
	gl.TexImage2D(gl.TEXTURE_2D, 0, bounds.Dx(), bounds.Dy(), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	if new_texture.mipmaps {
		new_texture.generateMipmaps()
	}
	return new_texture, nil
}

// getHandle will return the gl texutre handle
func (texture *Texture) getHandle() gl.Texture {
	return texture.textureId
}

// generate both the x, y coords at origin and the uv coords.
func (texture *Texture) generateVerticies() {
	w := float32(texture.GetWidth())
	h := float32(texture.GetHeight())
	texture.vertices = []float32{
		0, 0, 0, 0,
		0, h, 0, 1,
		w, 0, 1, 0,
		w, h, 1, 1,
	}
}

// getVerticies will return the verticies generated when this texture was created.
func (texture *Texture) getVerticies() []float32 {
	return texture.vertices
}

// generateMipmaps will generate mipmaps for the gl texture
func (texture *Texture) generateMipmaps() {
	// The GL_GENERATE_MIPMAP texparameter is set in loadVolatile if we don't
	// have support for glGenerateMipmap.
	if texture.mipmaps {
		// Driver bug: http://www.opengl.org/wiki/Common_Mistakes#Automatic_mipmap_generation
		if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
			gl.Enable(gl.TEXTURE_2D)
		}

		gl.GenerateMipmap(gl.TEXTURE_2D)
	}
}

// SetWrap will set how the texture behaves when applies to a plane that is larger
// than itself.
func (texture *Texture) SetWrap(wrap_s, wrap_t WrapMode) {
	texture.wrap.s = wrap_s
	texture.wrap.t = wrap_t
	bindTexture(texture.getHandle())
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int(wrap_s))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int(wrap_t))
}

// GetWrap will return the wrapping for how the texture behaves on a plane that
// is larger than itself
func (texture *Texture) GetWrap() Wrap {
	return texture.wrap
}

// SetFilter will set the min, mag filters for the texture filtering.
func (texture *Texture) SetFilter(min, mag FilterMode) error {
	if !texture.validateFilter() {
		if texture.filter.mipmap != FILTER_NONE && !texture.mipmaps {
			return fmt.Errorf("Non-mipmapped image cannot have mipmap filtering.")
		} else {
			return fmt.Errorf("Invalid texture filter.")
		}
	}
	texture.filter.min = min
	texture.filter.mag = mag
	texture.setTextureFilter()
	return nil
}

// setTextureFilter will set the texture filter on the actual gl texture. It will
// not reach this state if the filter is not valid.
func (texture *Texture) setTextureFilter() {
	var gmin, gmag uint32

	bindTexture(texture.getHandle())

	if texture.filter.mipmap == FILTER_NONE {
		if texture.filter.min == FILTER_NEAREST {
			gmin = gl.NEAREST
		} else { // f.min == FILTER_LINEAR
			gmin = gl.LINEAR
		}
	} else {
		if texture.filter.min == FILTER_NEAREST && texture.filter.mipmap == FILTER_NEAREST {
			gmin = gl.NEAREST_MIPMAP_NEAREST
		} else if texture.filter.min == FILTER_NEAREST && texture.filter.mipmap == FILTER_LINEAR {
			gmin = gl.NEAREST_MIPMAP_LINEAR
		} else if texture.filter.min == FILTER_LINEAR && texture.filter.mipmap == FILTER_NEAREST {
			gmin = gl.LINEAR_MIPMAP_NEAREST
		} else if texture.filter.min == FILTER_LINEAR && texture.filter.mipmap == FILTER_LINEAR {
			gmin = gl.LINEAR_MIPMAP_LINEAR
		} else {
			gmin = gl.LINEAR
		}
	}

	switch texture.filter.mag {
	case FILTER_NEAREST:
		gmag = gl.NEAREST
	case FILTER_LINEAR:
		fallthrough
	default:
		gmag = gl.LINEAR
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int(gmin))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int(gmag))
	//gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAX_ANISOTROPY_EXT, texture.filter.anisotropy)
}

// GetFilter will return the filter set on this texture.
func (texture *Texture) GetFilter() Filter {
	return texture.filter
}

// validateFilter will the the near and far filters and makes sure that it is possible
func (texture *Texture) validateFilter() bool {
	if !texture.mipmaps && texture.filter.mipmap != FILTER_NONE {
		return false
	}

	if texture.filter.mag != FILTER_LINEAR && texture.filter.mag != FILTER_NEAREST {
		return false
	}

	if texture.filter.min != FILTER_LINEAR && texture.filter.min != FILTER_NEAREST {
		return false
	}

	if texture.filter.mipmap != FILTER_LINEAR && texture.filter.mipmap != FILTER_NEAREST && texture.filter.mipmap != FILTER_NONE {
		return false
	}

	return true
}

// GetWidth will return the width of the texture.
func (texture *Texture) GetWidth() int32 {
	return texture.Width
}

// GetHeight will return the height of the texture.
func (texture *Texture) GetHeight() int32 {
	return texture.Height
}

// GetDimensions will return the width and height of the texture.
func (texture *Texture) GetDimensions() (int32, int32) {
	return texture.Width, texture.Height
}

// loadVolatile satisfies the volatile interface, so that it can be unloaded
func (texture *Texture) loadVolatile() bool {
	return false
}

// unloadVolatile release the texture data
func (texture *Texture) unloadVolatile() {
	if texture != nil {
		return
	}
	deleteTexture(texture.textureId)
	texture = nil
}

// Release will release the texture data and clean up the memory
func (texture *Texture) Release() {
	releaseVolatile(texture)
}

// drawv will take in verticies from the public draw calls and draw the texture
// with the verticies and the model matrix
func (texture *Texture) drawv(model *mat.Mat4, vertices []float32) {
	prepareDraw(model)
	bindTexture(texture.getHandle())
	useVertexAttribArrays(attribflag_pos | attribflag_texcoord)

	gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 4*4, gl.Ptr(vertices))
	gl.VertexAttribPointer(attrib_texcoord, 2, gl.FLOAT, false, 4*4, gl.Ptr(&vertices[2]))

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

// Draw satisfies the Drawable interface. Inputs are as follows
// x, y, r, sx, sy, ox, oy, kx, ky
// x, y are position
// r is rotation
// sx, sy is the scale, if sy is not given sy will equal sx
// ox, oy are offset
// kx, ky are the shear. If ky is not given ky will equal kx
func (texture *Texture) Draw(args ...float32) {
	texture.drawv(generateModelMatFromArgs(args), texture.vertices)
}

// Drawq satisfies the QuadDrawable interface.
// Inputs are as follows
// quad is the quad to crop the texture
// x, y, r, sx, sy, ox, oy, kx, ky
// x, y are position
// r is rotation
// sx, sy is the scale, if sy is not given sy will equal sx
// ox, oy are offset
// kx, ky are the shear. If ky is not given ky will equal kx
func (texture *Texture) Drawq(quad *Quad, args ...float32) {
	texture.drawv(generateModelMatFromArgs(args), quad.getVertices())
}
