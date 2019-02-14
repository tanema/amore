package gfx

import (
	"github.com/goxjs/gl"
)

type (
	// DrawMode is used to specify line or fill draws on primitives
	DrawMode string
	// WrapMode is used for setting texture/image/canvas wrap
	WrapMode int
	// FilterMode is used for setting texture/image/canvas filters
	FilterMode int
	// BlendMode specifies different ways to do color blending.
	BlendMode int
	// StencilAction is how a stencil function modifies the stencil values of pixels it touches.
	StencilAction uint32
	// CompareMode defines different types of per-pixel stencil test comparisons.
	// The pixels of an object will be drawn if the comparison succeeds, for each
	// pixel that the object touches.
	CompareMode uint32
	// LineJoin specifies how each lines are joined together
	LineJoin int
	// UniformType is the data type of a uniform
	UniformType int
	// MeshDrawMode specifies the tesselation of the mesh points
	MeshDrawMode uint32
	// Usage is used for sprite batch usage, and specifies if it is static, dynamic, or stream
	Usage uint32
	// ParticleDistribution specifies which direction particle will be send in when spawned
	ParticleDistribution int
	// ParticleInsertion specifies which level a particle will be inserted on spawn.
	ParticleInsertion int
	// AlignMode is the normal text align, center, left and right
	AlignMode int
)

// ColorMask contains an rgba color mask
type ColorMask struct {
	r, g, b, a bool
}

var (
	//opengl attribute variables
	shaderPos           = gl.Attrib{Value: 0}
	shaderTexCoord      = gl.Attrib{Value: 1}
	shaderColor         = gl.Attrib{Value: 2}
	shaderConstantColor = gl.Attrib{Value: 3}
)

// Draw modes
const (
	LINE DrawMode = "line"
	FILL DrawMode = "fill"
)

//texture wrap
const (
	WrapClamp          WrapMode = 0x812F
	WrapRepeat         WrapMode = 0x2901
	WrapMirroredRepeat WrapMode = 0x8370
)

//texture filter
const (
	FilterNone    FilterMode = 0
	FilterNearest FilterMode = 0x2600
	FilterLinear  FilterMode = 0x2601
)

//opengl blending constants
const (
	// Alpha blending (normal). The alpha of what's drawn determines its opacity.
	BlendModeAlpha BlendMode = iota
	// The pixel colors of what's drawn are multiplied with the pixel colors
	// already on the screen (darkening them). The alpha of drawn objects is
	// multiplied with the alpha of the screen rather than determining how much
	// the colors on the screen are affected.
	BlendModeMultiplicative
	BlendModePremultiplied
	// The pixel colors of what's drawn are subtracted from the pixel colors
	// already on the screen. The alpha of the screen is not modified.
	BlendModeSubtractive
	// The pixel colors of what's drawn are added to the pixel colors already on
	// the screen. The alpha of the screen is not modified.
	BlendModeAdditive
	// screen blending
	BlendModeScreen
	// The colors of what's drawn completely replace what was on the screen, with
	// no additional blending.
	BlendModeReplace
)

//stencil actions
const (
	StencilReplace       StencilAction = 0x1E01
	StencilIncrement     StencilAction = 0x1E02
	StencilDecrement     StencilAction = 0x1E03
	StencilIncrementWrap StencilAction = 0x8507
	StencilDecrementWrap StencilAction = 0x8508
	StencilInvert        StencilAction = 0x150A
)

// stenicl test modes
const (
	CompareGreater  CompareMode = 0x0201
	CompareEqual    CompareMode = 0x0202
	CompareGequal   CompareMode = 0x0203
	CompareLess     CompareMode = 0x0204
	CompareNotequal CompareMode = 0x0205
	CompareLequal   CompareMode = 0x0206
	CompareAlways   CompareMode = 0x0207
)

//line joins for nicer corners
const (
	LineJoinMiter LineJoin = iota
	LineJoinBevel
)

//uniform types for shaders
const (
	UniformFloat UniformType = iota
	UniformInt
	UniformBool
	UniformSampler
	UniformUnknown
	UniformBase UniformType = iota
	UniformVec
	UniformMat
)

// spritebatch usage
const (
	UsageStream  Usage = 0x88E0
	UsageStatic  Usage = 0x88E4
	UsageDynamic Usage = 0x88E8
)

// text align
const (
	AlignCenter AlignMode = iota
	AlignLeft
	AlignRight
	AlignJustify
)
