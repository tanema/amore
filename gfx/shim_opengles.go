// +build darwin linux
// +build arm arm64

package gfx

//Not supported
func enableMultisample()                               {}
func setTexMipMap()                                    {}
func (canvas *Canvas) attacheExtra(canvases []*Canvas) {}

// SetMipmapSharpness is disabled in ES
func (texture *Texture) SetMipmapSharpness(sharpness float32) {}

// SetWireframe is disabled in ES
func SetWireframe(enable bool) {}

func initMaxValues() {
	maxTextureSize = int32(gl.GetInteger(gl.MAX_TEXTURE_SIZE))
	maxTextureUnits = int32(gl.GetInteger(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS))
	glState.textureCounters = make([]int, maxTextureUnits)
}
