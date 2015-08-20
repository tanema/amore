package gfx

import (
	"github.com/go-gl/gl/v2.1/gl"

	"github.com/tanema/amore/gfx/util"
)

const (
	ATTRIB_POS = iota
	ATTRIB_TEXCOORD
	ATTRIB_COLOR
	ATTRIB_MAX_ENUM
)

type viewport [4]int32

type stats struct {
	textureMemory    int
	drawCalls        int
	framebufferBinds int
}

type state struct {
	Color *Color

	// Texture unit state (currently bound texture for each texture unit.)
	BoundTextures []uint32

	// Currently active texture unit.
	CurTextureUnit int

	Viewport viewport
	Scissor  viewport

	PointSize              float32
	FramebufferSRGBEnabled bool
	DefaultTexture         uint32
}

type openGL struct {
	vendor     string
	transform  []*util.Matrix
	projection []*util.Matrix
	viewport   viewport
	stats      *stats
	state      *state

	maxAnisotropy          float32
	maxTextureSize         int32
	maxRenderTargets       int32
	maxRenderbufferSamples int32
	maxTextureUnits        int32
}

var (
	opengl *openGL
)

func InitContext() {
	if opengl != nil {
		return
	}

	// Okay, setup OpenGL.
	gl.Init()

	opengl = &openGL{
		vendor:     gl.GoStr(gl.GetString(gl.VENDOR)),
		transform:  []*util.Matrix{util.NewEmptyMatrix()},
		projection: []*util.Matrix{util.NewEmptyMatrix()},
		stats:      &stats{},
		state: &state{
			Color:    &Color{255, 255, 255, 255},
			Viewport: viewport{},
			Scissor:  viewport{},
		},
	}

	opengl.initMaxValues()

	gl.VertexAttrib4f(ATTRIB_COLOR, 1.0, 1.0, 1.0, 1.0)

	// Get the current viewport.
	gl.GetIntegerv(gl.VIEWPORT, &opengl.state.Viewport[0])

	// And the current scissor - but we need to compensate for GL scissors
	// starting at the bottom left instead of top left.
	gl.GetIntegerv(gl.SCISSOR_BOX, &opengl.state.Scissor[0])
	opengl.state.Scissor[1] = opengl.state.Viewport[3] - (opengl.state.Scissor[1] + opengl.state.Scissor[3])

	gl.GetFloatv(gl.POINT_SIZE, &opengl.state.PointSize)
	opengl.state.FramebufferSRGBEnabled = gl.IsEnabled(gl.FRAMEBUFFER_SRGB)

	var curgltextureunit int32
	gl.GetIntegerv(gl.ACTIVE_TEXTURE, &curgltextureunit)
	opengl.state.CurTextureUnit = int(curgltextureunit - gl.TEXTURE0)

	// Retrieve currently bound textures for each texture unit.
	var boundTex int32
	for i := 0; i < int(opengl.maxTextureUnits); i++ {
		gl.ActiveTexture(uint32(gl.TEXTURE0 + i))
		gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &boundTex)
		opengl.state.BoundTextures = append(opengl.state.BoundTextures, uint32(boundTex))
	}

	gl.ActiveTexture(uint32(curgltextureunit))

	opengl.createDefaultTexture()
}

func (g *openGL) initMaxValues() {
	// We'll need this value to clamp anisotropy.
	gl.GetFloatv(gl.MAX_TEXTURE_MAX_ANISOTROPY_EXT, &g.maxAnisotropy)
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &g.maxTextureSize)
	gl.GetIntegerv(gl.MAX_DRAW_BUFFERS, &g.maxRenderTargets)

	maxattachments := int32(1)
	gl.GetIntegerv(gl.MAX_COLOR_ATTACHMENTS, &maxattachments)
	if maxattachments < g.maxRenderTargets {
		g.maxRenderTargets = maxattachments
	}

	gl.GetIntegerv(gl.MAX_SAMPLES, &g.maxRenderbufferSamples)
	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &g.maxTextureUnits)
}

func (g *openGL) createDefaultTexture() {
	// Set the 'default' texture (id 0) as a repeating white pixel. Otherwise,
	// texture2D calls inside a shader would return black when drawing graphics
	// primitives, which would create the need to use different "passthrough"
	// shaders for untextured primitives vs images.
	curtexture := g.state.BoundTextures[g.state.CurTextureUnit]

	gl.GenTextures(1, &g.state.DefaultTexture)
	g.BindTexture(g.state.DefaultTexture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	pix := []uint8{255, 255, 255, 255}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pix))

	g.BindTexture(curtexture)
}

func (g *openGL) PrepareDraw() {
}

func (g *openGL) BindTexture(texture uint32) {
	if texture != g.state.BoundTextures[g.state.CurTextureUnit] {
		g.state.BoundTextures[g.state.CurTextureUnit] = texture
		gl.BindTexture(gl.TEXTURE_2D, texture)
	}
}

func (g *openGL) GetDefaultTexture() uint32 {
	return g.state.DefaultTexture
}

func (g *openGL) DeInit() {
	gl.DeleteTextures(1, &g.state.DefaultTexture)
	g.state.DefaultTexture = 0
}

func (g *openGL) pushTransform() {
	new_transform := *g.GetTransform()
	g.transform = append(g.transform, &new_transform)
}

func (g *openGL) popTransform() {
	g.transform = g.transform[:len(g.transform)-1]
}

func (g *openGL) GetTransform() *util.Matrix {
	return g.transform[len(g.transform)-1]
}
