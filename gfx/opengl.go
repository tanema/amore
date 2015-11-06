package gfx

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	ATTRIB_POS = iota
	ATTRIB_TEXCOORD
	ATTRIB_COLOR
	ATTRIB_MAX_ENUM
)

type Viewport [4]int32 //The Viewport Values (X, Y, Width, Height)

var (
	is_initialized         = false
	opengl_version         string
	opengl_vendor          string
	maxAnisotropy          float32
	maxTextureSize         int32
	maxRenderTargets       int32
	maxRenderbufferSamples int32
	maxTextureUnits        int32
	viewport               Viewport
	scissor                Viewport
	pointSize              float32
	framebufferSRGBEnabled bool
	defaultTexture         uint32
	projectionMat          mgl32.Mat4
	viewMat                mgl32.Mat4
	modelMat               mgl32.Mat4
	screen_width           = 0
	screen_height          = 0
)

func InitContext(w, h int) {
	if is_initialized {
		return
	}

	// Okay, setup OpenGL.
	gl.Init()

	//Get system info
	opengl_version = gl.GoStr(gl.GetString(gl.VERSION))
	opengl_vendor = gl.GoStr(gl.GetString(gl.VENDOR))
	framebufferSRGBEnabled = gl.IsEnabled(gl.FRAMEBUFFER_SRGB)
	gl.GetIntegerv(gl.VIEWPORT, &viewport[0])
	gl.GetIntegerv(gl.SCISSOR_BOX, &scissor[0])
	gl.GetFloatv(gl.POINT_SIZE, &pointSize)
	gl.GetFloatv(gl.MAX_TEXTURE_MAX_ANISOTROPY_EXT, &maxAnisotropy)
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &maxTextureSize)
	gl.GetIntegerv(gl.MAX_SAMPLES, &maxRenderbufferSamples)
	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &maxTextureUnits)
	gl.GetIntegerv(gl.MAX_DRAW_BUFFERS, &maxRenderTargets)
	var maxattachments int32
	gl.GetIntegerv(gl.MAX_COLOR_ATTACHMENTS, &maxattachments)
	if maxattachments < maxRenderTargets {
		maxRenderTargets = maxattachments
	}

	// Enable blending
	gl.Enable(gl.BLEND)
	// Auto-generated mipmaps should be the best quality possible
	gl.Hint(gl.GENERATE_MIPMAP_HINT, gl.NICEST)
	// Make sure antialiasing works when set elsewhere
	gl.Enable(gl.MULTISAMPLE)
	// Set pixel row alignment
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	SetViewportSize(w, h)
	SetBackgroundColor(0, 0, 0, 1)
	createDefaultTexture()

	//default matricies
	projectionMat = mgl32.Ortho(0, float32(screen_width), float32(screen_height), 0, -1, 1)
	viewMat = mgl32.Translate3D(0, 0, 0)
	modelMat = mgl32.Translate3D(0, 0, 0)

	// We always need a default shader.
	defaultShader = NewShader()
	SetShader(defaultShader)

	is_initialized = true
}

// Set the 'default' texture (id 0) as a repeating white pixel. Otherwise,
// texture2D calls inside a shader would return black when drawing graphics
// primitives, which would create the need to use different "passthrough"
// shaders for untextured primitives vs images.
func createDefaultTexture() {
	gl.GenTextures(1, &defaultTexture)
	BindTexture(defaultTexture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	pix := []uint8{255, 255, 255, 255}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pix))
}

func PrepareDraw() {
	currentShader.SendMat4("ProjectionMat", projectionMat)
	currentShader.SendMat4("ViewMat", viewMat)
	currentShader.SendMat4("ModelMat", modelMat)
	currentShader.SendFloat("ScreenSize", float32(screen_width), float32(screen_height), 0, 0)
	currentShader.SendFloat("PointSize", pointSize)
}

func BindTexture(texture uint32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
}

func DeInit() {
	UnloadAll()
	gl.DeleteTextures(1, &defaultTexture)
	defaultTexture = 0
}

func GetViewport() Viewport {
	return viewport
}

func SetViewportSize(w, h int) {
	screen_width = w
	screen_height = h
	// Set the viewport to top-left corner.
	gl.Viewport(0, 0, int32(screen_width), int32(screen_height))
}
