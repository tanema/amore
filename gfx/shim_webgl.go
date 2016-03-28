// +build js

package gfx

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/goxjs/gl"
)

//Not supported
func enableMultisample()                                      {}
func setTexMipMap()                                           {}
func (texture *Texture) SetMipmapSharpness(sharpness float32) {}
func (canvas *Canvas) attacheExtra(canvases []*Canvas)        {}
func SetWireframe(enable bool)                                {}

func initMaxValues() {
	maxTextureSize = int32(gl.GetInteger(gl.MAX_TEXTURE_SIZE))
	maxTextureUnits = int32(gl.GetInteger(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS))
	gl_state.textureCounters = make([]int, maxTextureUnits)
}

func getCurrentFBO() gl.Framebuffer {
	current_fbo := gl.GetInteger(gl.FRAMEBUFFER_BINDING)
	return gl.Framebuffer{js.MakeWrapper(current_fbo)}
}
