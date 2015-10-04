package gfx

import (
	"math"
)

type Vector struct {
	X, Y float64
}

func NewZeroVector() *Vector {
	return &Vector{}
}

func NewVector(x, y float64) *Vector {
	return &Vector{
		X: x,
		Y: y,
	}
}

func (v *Vector) GetLength() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector) GetNormal() *Vector {
	return &Vector{
		X: -v.Y,
		Y: v.X,
	}
}

func (v *Vector) GetNormals(scale float64) *Vector {
	return &Vector{
		X: -v.Y * scale,
		Y: v.X * scale,
	}
}

func (v *Vector) Normalize(length float64) float64 {
	length_current := v.GetLength()

	if length_current > 0 {
		v.Mul(length / length_current)
	}

	return length_current
}

func (v *Vector) Add(other *Vector) *Vector {
	return NewVector(v.X+other.X, v.Y+other.Y)
}

func (v *Vector) Sub(other *Vector) *Vector {
	return NewVector(v.X-other.X, v.Y-other.Y)
}

func (v *Vector) Mul(s float64) *Vector {
	return NewVector(v.X*s, v.Y*s)
}

func (v *Vector) Div(s float64) *Vector {
	return NewVector(v.X/s, v.Y/s)
}

func (v *Vector) Negate() *Vector {
	return NewVector(-v.X, -v.Y)
}

func (v *Vector) Dot(other *Vector) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v *Vector) Cross(other *Vector) float64 {
	return v.X*other.X - v.Y*other.Y
}

func (v *Vector) Eq(other *Vector) bool {
	return v.GetLength() == other.GetLength()
}

func (v *Vector) LessThan(other *Vector) bool {
	return v.GetLength() < other.GetLength()
}
