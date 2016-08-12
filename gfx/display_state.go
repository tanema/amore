package gfx

import (
	"github.com/tanema/amore/gfx/gl"
	"github.com/tanema/amore/gfx/mat"
)

// displayState track a certain point in transformations
type displayState struct {
	color                  *Color
	background_color       *Color
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

// glState keeps track of the context attributes
type glState struct {
	initialized            bool
	active                 bool
	boundTextures          []gl.Texture
	curTextureUnit         int
	viewport               []int32
	framebufferSRGBEnabled bool
	defaultTexture         gl.Texture
	defaultFBO             gl.Framebuffer
	projectionStack        *mat.Stack
	viewStack              *mat.Stack
	currentCanvas          *Canvas
	currentShader          *Shader
	textureCounters        []int
	writingToStencil       bool
	enabledAttribArrays    uint32
}

// newDisplayState initializes a display states default values
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
		color:      NewColor(255, 255, 255, 255),
		colorMask:  ColorMask{r: true, g: true, b: true, a: true},
		scissorBox: make([]int32, 4),
	}
}

// displayStateStack is a simple stack for keeping track of display state.
type displayStateStack []displayState

// push a new element onto the top of the stack
func (stack *displayStateStack) push(state displayState) {
	*stack = append(*stack, state)
}

// take the top element off the stack
func (stack *displayStateStack) pop() displayState {
	var state displayState
	state, *stack = (*stack)[len(*stack)-1], (*stack)[:len(*stack)-1]
	return state
}

// get the top element in the stack
func (stack *displayStateStack) back() *displayState {
	return &(*stack)[len(*stack)-1]
}
