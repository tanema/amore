package gfx

import (
	"github.com/tanema/amore/gfx/gl"
)

type (
	WrapMode             int
	FilterMode           int
	BlendMode            int
	StencilAction        uint32
	CompareMode          uint32
	LineStyle            int
	LineJoin             int
	UniformType          int
	MeshDrawMode         uint32
	Usage                uint32
	ParticleDistribution int
	ParticleInsertion    int
	AlignMode            int
)

var (
	//opengl attribute variables
	ATTRIB_POS           = gl.Attrib{Value: 0}
	ATTRIB_TEXCOORD      = gl.Attrib{Value: 1}
	ATTRIB_COLOR         = gl.Attrib{Value: 2}
	ATTRIB_CONSTANTCOLOR = gl.Attrib{Value: 3}

	ATTRIBFLAG_POS           = uint32(1 << ATTRIB_POS.Value)
	ATTRIBFLAG_TEXCOORD      = uint32(1 << ATTRIB_TEXCOORD.Value)
	ATTRIBFLAG_COLOR         = uint32(1 << ATTRIB_COLOR.Value)
	ATTRIBFLAG_CONSTANTCOLOR = uint32(1 << ATTRIB_CONSTANTCOLOR.Value)
)

const (
	//texture wrap
	WRAP_CLAMP           WrapMode = WrapMode(gl.CLAMP_TO_EDGE)
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

	UNIFORM_BASE UniformType = iota
	UNIFORM_VEC
	UNIFORM_MAT

	//mesh draw modes
	DRAWMODE_FAN       MeshDrawMode = MeshDrawMode(gl.TRIANGLE_FAN)
	DRAWMODE_STRIP     MeshDrawMode = MeshDrawMode(gl.TRIANGLE_STRIP)
	DRAWMODE_TRIANGLES MeshDrawMode = MeshDrawMode(gl.TRIANGLES)
	DRAWMODE_POINTS    MeshDrawMode = MeshDrawMode(gl.POINTS)

	//mesh and spritebatch usage
	USAGE_STREAM  Usage = Usage(gl.STREAM_DRAW)
	USAGE_DYNAMIC Usage = Usage(gl.DYNAMIC_DRAW)
	USAGE_STATIC  Usage = Usage(gl.STATIC_DRAW)

	//particle distrobution
	DISTRIBUTION_NONE ParticleDistribution = iota
	DISTRIBUTION_UNIFORM
	DISTRIBUTION_NORMAL

	//particle insertion
	INSERT_MODE_TOP ParticleInsertion = iota
	INSERT_MODE_BOTTOM
	INSERT_MODE_RANDOM

	ALIGN_CENTER AlignMode = iota
	ALIGN_LEFT
	ALIGN_RIGHT
	ALIGN_JUSTIFY
)
