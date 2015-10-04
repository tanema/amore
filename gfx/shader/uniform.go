package shader

import (
	"github.com/go-gl/gl/v2.1/gl"
)

type UniformType int

const (
	UNIFORM_FLOAT UniformType = iota
	UNIFORM_INT
	UNIFORM_BOOL
	UNIFORM_SAMPLER
	UNIFORM_UNKNOWN
	UNIFORM_MAX_ENUM
)

type Uniform struct {
	Location int32
	Count    int32
	Type     uint32
	BaseType UniformType
	Name     string
}

func (u Uniform) SetBaseType() {
	u.BaseType = getUniformBaseType(u.Type)
}

func (u Uniform) GetTypeSize() int32 {
	return getUniformTypeSize(u.Type)
}

func getUniformTypeSize(t uint32) int32 {
	switch t {
	case gl.INT:
	case gl.FLOAT:
	case gl.BOOL:
	case gl.SAMPLER_1D:
	case gl.SAMPLER_2D:
	case gl.SAMPLER_3D:
		return 1
	case gl.INT_VEC2:
	case gl.FLOAT_VEC2:
	case gl.FLOAT_MAT2:
	case gl.BOOL_VEC2:
		return 2
	case gl.INT_VEC3:
	case gl.FLOAT_VEC3:
	case gl.FLOAT_MAT3:
	case gl.BOOL_VEC3:
		return 3
	case gl.INT_VEC4:
	case gl.FLOAT_VEC4:
	case gl.FLOAT_MAT4:
	case gl.BOOL_VEC4:
		return 4
	}
	return 1
}

func getUniformBaseType(t uint32) UniformType {
	switch t {
	case gl.INT:
	case gl.INT_VEC2:
	case gl.INT_VEC3:
	case gl.INT_VEC4:
		return UNIFORM_INT
	case gl.FLOAT:
	case gl.FLOAT_VEC2:
	case gl.FLOAT_VEC3:
	case gl.FLOAT_VEC4:
	case gl.FLOAT_MAT2:
	case gl.FLOAT_MAT3:
	case gl.FLOAT_MAT4:
		return UNIFORM_FLOAT
	case gl.BOOL:
	case gl.BOOL_VEC2:
	case gl.BOOL_VEC3:
	case gl.BOOL_VEC4:
		return UNIFORM_BOOL
	case gl.SAMPLER_1D:
	case gl.SAMPLER_2D:
	case gl.SAMPLER_3D:
		return UNIFORM_SAMPLER
	}
	return UNIFORM_UNKNOWN
}
