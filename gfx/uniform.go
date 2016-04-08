package gfx

import (
	"github.com/tanema/amore/gfx/gl"
)

// uniform represents a uniform in the shaders
type uniform struct {
	Location   gl.Uniform
	Type       gl.Enum
	BaseType   UniformType
	SecondType UniformType
	Count      int
	TypeSize   int
	Name       string
}

func (u *uniform) CalculateTypeInfo() {
	u.BaseType = u.getBaseType()
	u.SecondType = u.getSecondType()
	u.TypeSize = u.getTypeSize()
}

func (u *uniform) getTypeSize() int {
	switch u.Type {
	case gl.INT, gl.FLOAT, gl.BOOL, gl.SAMPLER_2D, gl.SAMPLER_CUBE:
		return 1
	case gl.INT_VEC2, gl.FLOAT_VEC2, gl.FLOAT_MAT2, gl.BOOL_VEC2:
		return 2
	case gl.INT_VEC3, gl.FLOAT_VEC3, gl.FLOAT_MAT3, gl.BOOL_VEC3:
		return 3
	case gl.INT_VEC4, gl.FLOAT_VEC4, gl.FLOAT_MAT4, gl.BOOL_VEC4:
		return 4
	}
	return 1
}

func (u *uniform) getBaseType() UniformType {
	switch u.Type {
	case gl.INT, gl.INT_VEC2, gl.INT_VEC3, gl.INT_VEC4:
		return UNIFORM_INT
	case gl.FLOAT, gl.FLOAT_VEC2, gl.FLOAT_VEC3,
		gl.FLOAT_VEC4, gl.FLOAT_MAT2, gl.FLOAT_MAT3, gl.FLOAT_MAT4:
		return UNIFORM_FLOAT
	case gl.BOOL, gl.BOOL_VEC2, gl.BOOL_VEC3, gl.BOOL_VEC4:
		return UNIFORM_BOOL
	case gl.SAMPLER_2D, gl.SAMPLER_CUBE:
		return UNIFORM_SAMPLER
	}
	return UNIFORM_UNKNOWN
}

func (u uniform) getSecondType() UniformType {
	switch u.Type {
	case gl.INT_VEC2, gl.INT_VEC3, gl.INT_VEC4, gl.FLOAT_VEC2,
		gl.FLOAT_VEC3, gl.FLOAT_VEC4, gl.BOOL_VEC2, gl.BOOL_VEC3, gl.BOOL_VEC4:
		return UNIFORM_VEC
	case gl.FLOAT_MAT2, gl.FLOAT_MAT3, gl.FLOAT_MAT4:
		return UNIFORM_MAT
	}
	return UNIFORM_BASE
}

func translateUniformBaseType(t UniformType) string {
	switch t {
	case UNIFORM_FLOAT:
		return "float"
	case UNIFORM_INT:
		return "int"
	case UNIFORM_BOOL:
		return "bool"
	case UNIFORM_SAMPLER:
		return "sampler"
	}
	return "unknown"
}
