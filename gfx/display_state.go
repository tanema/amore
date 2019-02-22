package gfx

import (
	"github.com/go-gl/mathgl/mgl32/matstack"

	"github.com/goxjs/gl"
)

// displayState track a certain point in transformations
type displayState struct {
	color            []float32
	backgroundColor  []float32
	blendMode        string
	lineWidth        float32
	lineJoin         string
	pointSize        float32
	scissor          bool
	scissorBox       []int32
	stencilCompare   CompareMode
	stencilTestValue int32
	font             *Font
	shader           *Shader
	colorMask        ColorMask
	canvas           *Canvas
	defaultFilter    Filter
}

// glState keeps track of the context attributes
type openglState struct {
	initialized            bool
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

// newDisplayState initializes a display states default values
func newDisplayState() displayState {
	return displayState{
		blendMode:      "alpha",
		pointSize:      5,
		stencilCompare: CompareAlways,
		lineWidth:      1,
		lineJoin:       "miter",
		shader:         defaultShader,
		font:           defaultFont,
		defaultFilter:  newFilter(),
		color:          []float32{1, 1, 1, 1},
		colorMask:      ColorMask{r: true, g: true, b: true, a: true},
		scissorBox:     make([]int32, 4),
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
