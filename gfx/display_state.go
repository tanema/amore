package gfx

import (
	"github.com/go-gl/mathgl/mgl32/matstack"

	"github.com/tanema/amore/gfx/gl"
)

// displayState track a certain point in transformations
type displayState struct {
	color                  *Color
	backgroundColor        *Color
	blendMode              BlendMode
	lineWidth              float32
	lineStyle              LineStyle
	lineJoin               LineJoin
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
type openglState struct {
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
	enabledAttribArrays    uint32
}

// newDisplayState initializes a display states default values
func newDisplayState() displayState {
	return displayState{
		blendMode:              BlendModeAlpha,
		pointSize:              5,
		pixelSize:              1,
		stencilCompare:         CompareAlways,
		lineWidth:              1,
		lineJoin:               LineJoinMiter,
		lineStyle:              LineSmooth,
		shader:                 defaultShader,
		defaultFilter:          newFilter(),
		defaultMipmapFilter:    FilterNearest,
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
