// +build darwin linux
// +build arm arm64

package gfx

import (
	"github.com/goxjs/gl"

	"github.com/tanema/amore/window/ui"
)

func initGLContext(ctx ui.Context) {
	gl.ContextWatcher.OnMakeCurrent(ctx)
}

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
