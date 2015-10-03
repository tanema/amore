package shader

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

func getUniformBaseType(t uint32) UniformType {
	return UNIFORM_UNKNOWN
}
