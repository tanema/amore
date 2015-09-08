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
	Location uint32
	Count    uint32
	Type     uint32
	BaseType UniformType
	Name     string
}
