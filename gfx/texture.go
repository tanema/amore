package gfx

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
)

var filters = map[string]int32{"linear": gl.LINEAR, "nearest": gl.NEAREST}
var wraps = map[string]int32{"clamp": gl.CLAMP_TO_EDGE, "repeat": gl.REPEAT}

type (
	Filter struct {
		min, mag, mipmap FilterMode
		anisotropy       float32
	}
	Wrap struct {
		s, t WrapMode
	}
	Texture struct {
		textureId     gl.Texture
		Width, Height int32
		vertices      []float32
		filter        Filter
		wrap          Wrap
		mipmaps       bool
	}
	iTexture interface {
		GetHandle() gl.Texture
		GetWidth() int32
		GetHeight() int32
		getVerticies() []float32
		Release()
	}
)

func newFilter() Filter {
	return Filter{
		min:        FILTER_LINEAR,
		mag:        FILTER_LINEAR,
		mipmap:     FILTER_NONE,
		anisotropy: 1.0,
	}
}

func newWrap() Wrap {
	return Wrap{s: WRAP_CLAMP, t: WRAP_CLAMP}
}

func newTexture(width, height int32, mipmaps bool) *Texture {
	new_texture := &Texture{
		textureId: gl.CreateTexture(),
		Width:     width,
		Height:    height,
		wrap:      newWrap(),
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
		//TODO transfer to non es build file
		// Auto-generate mipmaps every time the texture is modified, if
		// glGenerateMipmap isn't supported.
		//gl.TexParameteri(gl.TEXTURE_2D, gl.GENERATE_MIPMAP, gl.TRUE)
	}

	new_texture.generateVerticies()

	return new_texture
}

func newImageTexture(img image.Image, mipmaps bool) (*Texture, error) {
	bounds := img.Bounds()
	new_texture := newTexture(int32(bounds.Dx()), int32(bounds.Dy()), mipmaps)
	//generate a uniform image and upload to vram
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, bounds, img, image.Point{0, 0}, draw.Src)
	bindTexture(new_texture.GetHandle())
	gl.TexImage2D(gl.TEXTURE_2D, 0, bounds.Dx(), bounds.Dy(), gl.RGBA, gl.UNSIGNED_BYTE, rgba.Pix)

	if new_texture.mipmaps {
		new_texture.generateMipmaps()
	}
	return new_texture, nil
}

func (texture *Texture) GetHandle() gl.Texture {
	return texture.textureId
}

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

func (texture *Texture) getVerticies() []float32 {
	return texture.vertices
}

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

func (texture *Texture) SetMipmapSharpness(sharpness float32) {
	//TODO transfer to non es build file
	//var maxMipmapSharpness float32
	//gl.GetFloatv(gl.MAX_TEXTURE_LOD_BIAS, &maxMipmapSharpness)
	//mipmapSharpness := math.Min(math.Max(float64(sharpness), -float64(maxMipmapSharpness+0.01)), float64(maxMipmapSharpness-0.01))
	//bindTexture(texture.GetHandle())
	// negative bias is sharper
	//gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_LOD_BIAS, -float32(mipmapSharpness))
}

func (texture *Texture) SetWrap(wrap_s, wrap_t WrapMode) {
	texture.wrap.s = wrap_s
	texture.wrap.t = wrap_t
	bindTexture(texture.GetHandle())
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int(wrap_s))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int(wrap_t))
}

func (texture *Texture) GetWrap() Wrap {
	return texture.wrap
}

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

func (texture *Texture) setTextureFilter() {
	var gmin, gmag uint32

	bindTexture(texture.GetHandle())

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

func (texture *Texture) GetFilter() Filter {
	return texture.filter
}

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

func (texture *Texture) GetWidth() int32 {
	return texture.Width
}

func (texture *Texture) GetHeight() int32 {
	return texture.Height
}

func (texture *Texture) GetDimensions() (int32, int32) {
	return texture.Width, texture.Height
}

func (texture *Texture) Release() {
	releaseVolatile(texture)
}

func (texture *Texture) loadVolatile() bool {
	return false
}

func (texture *Texture) unloadVolatile() {
	if texture != nil {
		return
	}
	deleteTexture(texture.textureId)
	texture = nil
}

func (texture *Texture) drawv(model *mgl32.Mat4, vertices []float32) {
	prepareDraw(model)
	bindTexture(texture.GetHandle())
	useVertexAttribArrays(ATTRIBFLAG_POS | ATTRIBFLAG_TEXCOORD)
	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, f32Bytes(vertices...), gl.STATIC_DRAW)
	gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 4*4, 0)
	gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, 4*4, 2*4)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{})
	gl.DeleteBuffer(vbo)
}

func (texture *Texture) Draw(args ...float32) {
	texture.drawv(generateModelMatFromArgs(args), texture.vertices)
}

func (texture *Texture) Drawq(quad *Quad, args ...float32) {
	texture.drawv(generateModelMatFromArgs(args), quad.getVertices())
}
