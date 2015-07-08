package gfx

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
)

type BindCB func()

var filters = map[string]int32{"linear": gl.LINEAR, "nearest": gl.NEAREST}
var wraps = map[string]int32{"clamp": gl.CLAMP_TO_EDGE, "repeat": gl.REPEAT}

type Texture struct {
	textureId     uint32
	Width, Height float32
}

func NewTexture(file string) (*Texture, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	return NewImageTexture(img)
}

func NewImageTexture(img image.Image) (*Texture, error) {
	var texture_id uint32
	gl.GenTextures(1, &texture_id)

	bounds := img.Bounds()
	new_texture := &Texture{
		textureId: texture_id,
		Width:     float32(bounds.Dx()),
		Height:    float32(bounds.Dy()),
	}

	new_texture.Bind(func() {
		//generate a uniform image and upload to vram
		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, bounds, img, image.Point{0, 0}, draw.Src)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(bounds.Dx()), int32(bounds.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	})

	new_texture.SetFilter("nearest", "nearest")
	new_texture.SetWrap("clamp", "clamp")

	return new_texture, nil
}

func (texture *Texture) Bind(draw_cb BindCB) {
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, texture.textureId)
	draw_cb()
	gl.Disable(gl.TEXTURE_2D)
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
