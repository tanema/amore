package vec

import (
	"math"
)

type Vec2 [2]float32

func (v1 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v1[0] + v2[0], v1[1] + v2[1]}
}

func (v1 Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v1[0] - v2[0], v1[1] - v2[1]}
}

func (v1 Vec2) Mul(c float32) Vec2 {
	return Vec2{v1[0] * c, v1[1] * c}
}

func (v1 Vec2) Dot(v2 Vec2) float32 {
	return v1[0]*v2[0] + v1[1]*v2[1]
}

func (v1 Vec2) Len() float32 {
	return float32(math.Hypot(float64(v1[0]), float64(v1[1])))
}

func (v1 Vec2) Normalize() Vec2 {
	l := 1.0 / v1.Len()
	return Vec2{v1[0] * l, v1[1] * l}
}
