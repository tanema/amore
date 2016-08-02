package gfx

import (
	"github.com/go-gl/gl/v2.1/gl"
)

// SetWireframe sets whether wireframe lines will be used when drawing.
func SetWireframe(enable bool) {
	if enable {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
	states.back().wireframe = enable
}

func enableMultisample() {
	gl.Enable(gl.MULTISAMPLE)
}

func initMaxValues() {
	gl_state.framebufferSRGBEnabled = gl.IsEnabled(gl.FRAMEBUFFER_SRGB)
	gl.GetFloatv(gl.POINT_SIZE, &states.back().pointSize)
	gl.GetFloatv(gl.MAX_TEXTURE_MAX_ANISOTROPY_EXT, &maxAnisotropy)
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &maxTextureSize)
	gl.GetIntegerv(gl.MAX_SAMPLES, &maxRenderbufferSamples)
	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &maxTextureUnits)
	gl_state.textureCounters = make([]int, maxTextureUnits)
	gl.GetIntegerv(gl.MAX_DRAW_BUFFERS, &maxRenderTargets)
	var maxattachments int32
	gl.GetIntegerv(gl.MAX_COLOR_ATTACHMENTS, &maxattachments)
	if maxattachments < maxRenderTargets {
		maxRenderTargets = maxattachments
	}
}

func setTexMipMap() {
	gl.TexParameteri(gl.TEXTURE_2D, gl.GENERATE_MIPMAP, gl.TRUE)
}

func (texture *Texture) SetMipmapSharpness(sharpness float32) {
	var maxMipmapSharpness float32
	gl.GetFloatv(gl.MAX_TEXTURE_LOD_BIAS, &maxMipmapSharpness)
	mipmapSharpness := Min(Max(sharpness, -(maxMipmapSharpness+0.01)), maxMipmapSharpness-0.01)
	bindTexture(texture.getHandle())
	//negative bias is sharper
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_LOD_BIAS, -float32(mipmapSharpness))
}

func (canvas *Canvas) attacheExtra(canvases []*Canvas) {
	if canvases != nil && len(canvases) > 0 {
		// Attach the canvas textures to the active FBO and set up MRTs.
		drawbuffers := []uint32{gl.COLOR_ATTACHMENT0}
		// Attach the canvas textures to the currently bound framebuffer.
		for i := 0; i < len(canvases); i++ {
			buf := gl.COLOR_ATTACHMENT1 + uint32(i)
			gl.FramebufferTexture2D(gl.FRAMEBUFFER, buf, gl.TEXTURE_2D, canvases[i].getHandle().Value, 0)
			drawbuffers = append(drawbuffers, buf)
		}
		// set up multiple render targets
		gl.DrawBuffers(int32(len(drawbuffers)), &drawbuffers[0])
	} else {
		// Make sure the FBO is only using a single draw buffer.
		gl.DrawBuffer(gl.COLOR_ATTACHMENT0)
	}
}
