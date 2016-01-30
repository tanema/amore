package gfx

import (
	"github.com/go-gl/gl/v2.1/gl"
)

type (
	WrapMode      int
	FilterMode    int
	BlendMode     int
	StencilAction uint32
	CompareMode   uint32
	LineStyle     int
	LineJoin      int
	UniformType   int
)

const (
	//opengl attribute variables
	ATTRIB_POS = iota
	ATTRIB_TEXCOORD
	ATTRIB_COLOR
	ATTRIB_MAX_ENUM

	//texture wrap
	WRAP_CLAMP           WrapMode = WrapMode(gl.CLAMP)
	WRAP_REPEAT          WrapMode = WrapMode(gl.REPEAT)
	WRAP_MIRRORED_REPEAT WrapMode = WrapMode(gl.MIRRORED_REPEAT)

	//texture filter
	FILTER_NONE    FilterMode = FilterMode(gl.NONE)
	FILTER_LINEAR  FilterMode = FilterMode(gl.LINEAR)
	FILTER_NEAREST FilterMode = FilterMode(gl.NEAREST)

	//opengl blending constants
	BLENDMODE_ALPHA BlendMode = iota
	BLENDMODE_MULTIPLICATIVE
	BLENDMODE_PREMULTIPLIED
	BLENDMODE_SUBTRACTIVE
	BLENDMODE_ADDITIVE
	BLENDMODE_SCREEN
	BLENDMODE_REPLACE

	//stencil actions
	STENCIL_REPLACE        StencilAction = StencilAction(gl.REPLACE)
	STENCIL_INCREMENT      StencilAction = StencilAction(gl.INCR)
	STENCIL_DECREMENT      StencilAction = StencilAction(gl.DECR)
	STENCIL_INCREMENT_WRAP StencilAction = StencilAction(gl.INCR_WRAP)
	STENCIL_DECREMENT_WRAP StencilAction = StencilAction(gl.DECR_WRAP)
	STENCIL_INVERT         StencilAction = StencilAction(gl.INVERT)

	/**
	 * Q: Why are some of the compare modes inverted (e.g. COMPARE_LESS becomes
	 * GL_GREATER)?
	 *
	 * A: OpenGL / GPUs do the comparison in the opposite way that makes sense
	 * for this API. For example, if the compare function is GL_GREATER then the
	 * stencil test will pass if the reference value is greater than the value
	 * in the stencil buffer. With our API it's more intuitive to assume that
	 * setStencilTest(COMPARE_GREATER, 4) will make it pass if the stencil
	 * buffer has a value greater than 4.
	 **/
	COMPARE_LESS     CompareMode = CompareMode(gl.GREATER)
	COMPARE_LEQUAL   CompareMode = CompareMode(gl.GEQUAL)
	COMPARE_EQUAL    CompareMode = CompareMode(gl.EQUAL)
	COMPARE_GEQUAL   CompareMode = CompareMode(gl.LEQUAL)
	COMPARE_GREATER  CompareMode = CompareMode(gl.LESS)
	COMPARE_NOTEQUAL CompareMode = CompareMode(gl.NOTEQUAL)
	COMPARE_ALWAYS   CompareMode = CompareMode(gl.ALWAYS)

	// treat adjacent segments with angles between their directions <5 degree as straight
	LINES_PARALLEL_EPS float32 = 0.05
	//line styles for overdraw
	LINE_ROUGH LineStyle = iota
	LINE_SMOOTH
	//line joins for nicer corners
	LINE_JOIN_NONE LineJoin = iota
	LINE_JOIN_MITER
	LINE_JOIN_BEVEL

	//uniform types for shaders
	UNIFORM_FLOAT UniformType = iota
	UNIFORM_INT
	UNIFORM_BOOL
	UNIFORM_SAMPLER
	UNIFORM_UNKNOWN
	UNIFORM_MAX_ENUM
	UNIFORM_BASE UniformType = iota
	UNIFORM_VEC
	UNIFORM_MAT
)