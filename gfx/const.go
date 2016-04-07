package gfx

import (
	"github.com/tanema/amore/gfx/gl"
)

type (
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
	// LineStyle specifies if the line drawing is smooth or rough
	LineStyle int
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

var (
	//opengl attribute variables
	attrib_pos           = gl.Attrib{Value: 0}
	attrib_texcoord      = gl.Attrib{Value: 1}
	attrib_color         = gl.Attrib{Value: 2}
	attrib_constantcolor = gl.Attrib{Value: 3}

	attribflag_pos           = uint32(1 << attrib_pos.Value)
	attribflag_texcoord      = uint32(1 << attrib_texcoord.Value)
	attribflag_color         = uint32(1 << attrib_color.Value)
	attribflag_constantcolor = uint32(1 << attrib_constantcolor.Value)
)

const (
	//texture wrap
	WRAP_CLAMP           WrapMode = 0x812F
	WRAP_REPEAT          WrapMode = 0x2901
	WRAP_MIRRORED_REPEAT WrapMode = 0x8370

	//texture filter
	FILTER_NONE    FilterMode = 0
	FILTER_NEAREST FilterMode = 0x2600
	FILTER_LINEAR  FilterMode = 0x2601

	//opengl blending constants
	BLENDMODE_ALPHA BlendMode = iota
	BLENDMODE_MULTIPLICATIVE
	BLENDMODE_PREMULTIPLIED
	BLENDMODE_SUBTRACTIVE
	BLENDMODE_ADDITIVE
	BLENDMODE_SCREEN
	BLENDMODE_REPLACE

	//stencil actions
	STENCIL_REPLACE        StencilAction = 0x1E01
	STENCIL_INCREMENT      StencilAction = 0x1E02
	STENCIL_DECREMENT      StencilAction = 0x1E03
	STENCIL_INCREMENT_WRAP StencilAction = 0x8507
	STENCIL_DECREMENT_WRAP StencilAction = 0x8508
	STENCIL_INVERT         StencilAction = 0x150A

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
	COMPARE_GREATER  CompareMode = 0x0201
	COMPARE_EQUAL    CompareMode = 0x0202
	COMPARE_GEQUAL   CompareMode = 0x0203
	COMPARE_LESS     CompareMode = 0x0204
	COMPARE_NOTEQUAL CompareMode = 0x0205
	COMPARE_LEQUAL   CompareMode = 0x0206
	COMPARE_ALWAYS   CompareMode = 0x0207

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
	DRAWMODE_POINTS    MeshDrawMode = 0x0000
	DRAWMODE_TRIANGLES MeshDrawMode = 0x0004
	DRAWMODE_STRIP     MeshDrawMode = 0x0005
	DRAWMODE_FAN       MeshDrawMode = 0x0006

	//mesh and spritebatch usage
	USAGE_STREAM  Usage = 0x88E0
	USAGE_STATIC  Usage = 0x88E4
	USAGE_DYNAMIC Usage = 0x88E8

	//particle distrobution
	DISTRIBUTION_NONE ParticleDistribution = iota
	DISTRIBUTION_UNIFORM
	DISTRIBUTION_NORMAL

	//particle insertion
	INSERT_MODE_TOP ParticleInsertion = iota
	INSERT_MODE_BOTTOM
	INSERT_MODE_RANDOM

	// text align
	ALIGN_CENTER AlignMode = iota
	ALIGN_LEFT
	ALIGN_RIGHT
	ALIGN_JUSTIFY
)
