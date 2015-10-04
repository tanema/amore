package gfx

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"github.com/go-gl/gl/v2.1/gl"
)

var filters = map[string]int32{"linear": gl.LINEAR, "nearest": gl.NEAREST}
var wraps = map[string]int32{"clamp": gl.CLAMP_TO_EDGE, "repeat": gl.REPEAT}

type Texture struct {
	textureId     uint32
	Width, Height float64
}

func LoadImageTexture(img image.Image) (*Texture, error) {
	var texture_id uint32
	gl.GenTextures(1, &texture_id)

	bounds := img.Bounds()
	new_texture := &Texture{
		textureId: texture_id,
		Width:     float64(bounds.Dx()),
		Height:    float64(bounds.Dy()),
	}

	BindTexture(new_texture.GetHandle())
	//generate a uniform image and upload to vram
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, bounds, img, image.Point{0, 0}, draw.Src)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(bounds.Dx()), int32(bounds.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	new_texture.SetFilter("nearest", "nearest")
	new_texture.SetWrap("clamp", "clamp")

	return new_texture, nil
}

func (texture *Texture) GetHandle() uint32 {
	return texture.textureId
}

func (texture *Texture) SetWrap(wrap_s, wrap_t string) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, wraps[wrap_s])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, wraps[wrap_t])
}

func (texture *Texture) SetFilter(min, mag string) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, filters[min])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, filters[mag])
}

func (texture *Texture) Release() {
	gl.DeleteTextures(1, &texture.textureId)
}
