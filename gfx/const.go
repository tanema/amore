package gfx

import "github.com/goxjs/gl"

type (
	// WrapMode is used for setting texture/image/canvas wrap
	WrapMode int
	// FilterMode is used for setting texture/image/canvas filters
	FilterMode int
	// StencilAction is how a stencil function modifies the stencil values of pixels it touches.
	StencilAction uint32
	// CompareMode defines different types of per-pixel stencil test comparisons.
	// The pixels of an object will be drawn if the comparison succeeds, for each
	// pixel that the object touches.
	CompareMode uint32
	// Usage is used for sprite batch usage, and specifies if it is static, dynamic, or stream
	Usage uint32
)

// ColorMask contains an rgba color mask
type ColorMask struct {
	r, g, b, a bool
}

var ( //opengl attribute variables
	shaderPos           = gl.Attrib{Value: 0}
	shaderTexCoord      = gl.Attrib{Value: 1}
	shaderColor         = gl.Attrib{Value: 2}
	shaderConstantColor = gl.Attrib{Value: 3}
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

// spritebatch usage
const (
	UsageStream  Usage = 0x88E0
	UsageStatic  Usage = 0x88E4
	UsageDynamic Usage = 0x88E8
)
