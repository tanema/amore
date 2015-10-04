package gfx

import (
	"github.com/go-gl/gl/v2.1/gl"
)

const (
	ATTRIB_POS = iota
	ATTRIB_TEXCOORD
	ATTRIB_COLOR
	ATTRIB_MAX_ENUM
)

type Viewport [4]int32 //The Viewport Values (X, Y, Width, Height)

type State struct {
	// Texture unit state (currently bound texture for each texture unit.)
	BoundTextures []uint32

	// Currently active texture unit.
	CurTextureUnit int

	Viewport Viewport
	Scissor  Viewport

	PointSize              float32
	FramebufferSRGBEnabled bool
	DefaultTexture         uint32

	LastProjectionMatrix *Matrix
	LastTransformMatrix  *Matrix
}

var (
	vendor     string
	transform  []*Matrix
	projection []*Matrix
	state      *State

	maxAnisotropy          float32
	maxTextureSize         int32
	maxRenderTargets       int32
	maxRenderbufferSamples int32
	maxTextureUnits        int32

	is_initialized = false
)

func InitContext() {
	if is_initialized {
		return
	}

	// Okay, setup OpenGL.
	gl.Init()

	vendor = gl.GoStr(gl.GetString(gl.VENDOR))
	transform = []*Matrix{NewEmptyMatrix()}
	projection = []*Matrix{NewEmptyMatrix()}
	state = &State{
		Viewport: Viewport{},
		Scissor:  Viewport{},
	}

	initMaxValues()

	gl.VertexAttrib4f(ATTRIB_COLOR, 1.0, 1.0, 1.0, 1.0)

	// Get the current viewport.
	gl.GetIntegerv(gl.VIEWPORT, &state.Viewport[0])

	// And the current scissor - but we need to compensate for GL scissors
	// starting at the bottom left instead of top left.
	gl.GetIntegerv(gl.SCISSOR_BOX, &state.Scissor[0])
	state.Scissor[1] = state.Viewport[3] - (state.Scissor[1] + state.Scissor[3])

	gl.GetFloatv(gl.POINT_SIZE, &state.PointSize)
	state.FramebufferSRGBEnabled = gl.IsEnabled(gl.FRAMEBUFFER_SRGB)

	var curgltextureunit int32
	gl.GetIntegerv(gl.ACTIVE_TEXTURE, &curgltextureunit)
	state.CurTextureUnit = int(curgltextureunit - gl.TEXTURE0)

	// Retrieve currently bound textures for each texture unit.
	var boundTex int32
	for i := 0; i < int(maxTextureUnits); i++ {
		gl.ActiveTexture(uint32(gl.TEXTURE0 + i))
		gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &boundTex)
		state.BoundTextures = append(state.BoundTextures, uint32(boundTex))
	}

	gl.ActiveTexture(uint32(curgltextureunit))

	createDefaultTexture()

	is_initialized = true
}

func initMaxValues() {
	// We'll need this value to clamp anisotropy.
	gl.GetFloatv(gl.MAX_TEXTURE_MAX_ANISOTROPY_EXT, &maxAnisotropy)
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &maxTextureSize)
	gl.GetIntegerv(gl.MAX_DRAW_BUFFERS, &maxRenderTargets)

	maxattachments := int32(1)
	gl.GetIntegerv(gl.MAX_COLOR_ATTACHMENTS, &maxattachments)
	if maxattachments < maxRenderTargets {
		maxRenderTargets = maxattachments
	}

	gl.GetIntegerv(gl.MAX_SAMPLES, &maxRenderbufferSamples)
	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &maxTextureUnits)
}

func createDefaultTexture() {
	// Set the 'default' texture (id 0) as a repeating white pixel. Otherwise,
	// texture2D calls inside a shader would return black when drawing graphics
	// primitives, which would create the need to use different "passthrough"
	// shaders for untextured primitives vs images.
	curtexture := state.BoundTextures[state.CurTextureUnit]

	gl.GenTextures(1, &state.DefaultTexture)
	BindTexture(state.DefaultTexture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	pix := []uint8{255, 255, 255, 255}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pix))

	BindTexture(curtexture)
}

func PrepareDraw() {
	shader := currentShader
	if shader != nil {
		// Make sure the active shader has the correct values for its
		// love-provided uniforms.
		shader.CheckSetScreenParams()
	}

	//curproj := projection[len(projection)-1]
	//curxform := transform[len(transform)-1]

	// We only need to re-upload the projection matrix if it's changed.
	//if state.LastProjectionMatrix != curproj {
	//gl.MatrixMode(gl.PROJECTION)
	//projection_elements := curproj.GetElements()
	//gl.LoadMatrixf(&projection_elements[0])
	//gl.MatrixMode(gl.MODELVIEW)

	//state.LastProjectionMatrix = curproj
	//}

	//if state.LastTransformMatrix != curxform {
	//transform_elements := curxform.GetElements()
	//gl.LoadMatrixf(&transform_elements[0])
	//state.LastTransformMatrix = curxform
	//}
}

func SetTextureUnit(textureunit int) {
	if textureunit < 0 || textureunit >= len(state.BoundTextures) {
		panic("Invalid texture unit index.")
	}

	if textureunit != state.CurTextureUnit {
		gl.ActiveTexture(gl.TEXTURE0 + (uint32)(textureunit))
	}

	state.CurTextureUnit = textureunit
}

func BindTexture(texture uint32) {
	if texture != state.BoundTextures[state.CurTextureUnit] {
		state.BoundTextures[state.CurTextureUnit] = texture
		gl.BindTexture(gl.TEXTURE_2D, texture)
	}
}

func BindTextureToUnit(texture uint32, textureunit int, restoreprev bool) {
	if textureunit < 0 || textureunit >= len(state.BoundTextures) {
		panic("Invalid texture unit index.")
	}

	if texture != state.BoundTextures[textureunit] {
		oldtextureunit := state.CurTextureUnit
		SetTextureUnit(textureunit)

		state.BoundTextures[textureunit] = texture
		gl.BindTexture(gl.TEXTURE_2D, texture)

		if restoreprev {
			SetTextureUnit(oldtextureunit)
		}
	}
}

func GetDefaultTexture() uint32 {
	return state.DefaultTexture
}

func DeInit() {
	gl.DeleteTextures(1, &state.DefaultTexture)
	state.DefaultTexture = 0
}

func pushTransform() {
	new_transform := *GetTransform()
	transform = append(transform, &new_transform)
}

func popTransform() {
	transform = transform[:len(transform)-1]
}

func GetTransform() *Matrix {
	return transform[len(transform)-1]
}

func GetViewport() Viewport {
	return state.Viewport
}
