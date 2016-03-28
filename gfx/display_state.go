package gfx

import (
	"github.com/go-gl/mathgl/mgl32/matstack"
	"github.com/goxjs/gl"
)

type displayState struct {
	color                  Color
	background_color       Color
	blend_mode             BlendMode
	line_width             float32
	line_style             LineStyle
	line_join              LineJoin
	pointSize              float32
	scissor                bool
	scissorBox             []int32
	stencilCompare         CompareMode
	stencilTestValue       int32
	font                   *Font
	shader                 *Shader
	colorMask              ColorMask
	wireframe              bool
	pixelSize              float32
	canvases               []*Canvas
	defaultFilter          Filter
	defaultMipmapFilter    FilterMode
	defaultMipmapSharpness float32
}

type glState struct {
	initialized            bool
	active                 bool
	boundTextures          []gl.Texture
	curTextureUnit         int
	viewport               []int32
	framebufferSRGBEnabled bool
	defaultTexture         gl.Texture
	defaultFBO             gl.Framebuffer
	projectionStack        *matstack.MatStack
	viewStack              *matstack.MatStack
	currentCanvas          *Canvas
	currentShader          *Shader
	textureCounters        []int
	writingToStencil       bool
}

func newDisplayState() displayState {
	return displayState{
		blend_mode:             BLENDMODE_ALPHA,
		pointSize:              1,
		pixelSize:              1,
		stencilCompare:         COMPARE_ALWAYS,
		line_width:             1,
		line_join:              LINE_JOIN_MITER,
		line_style:             LINE_SMOOTH,
		shader:                 defaultShader,
		defaultFilter:          newFilter(),
		defaultMipmapFilter:    FILTER_NEAREST,
		defaultMipmapSharpness: 0.0,
		colorMask:              ColorMask{r: true, g: true, b: true, a: true},
		scissorBox:             make([]int32, 4),
	}
}

type displayStateStack []displayState

func (stack *displayStateStack) push(state displayState) {
	*stack = append(*stack, state)
}

func (stack *displayStateStack) pop() displayState {
	var state displayState
	state, *stack = (*stack)[len(*stack)-1], (*stack)[:len(*stack)-1]
	return state
}

func (stack *displayStateStack) back() *displayState {
	return &(*stack)[len(*stack)-1]
}
