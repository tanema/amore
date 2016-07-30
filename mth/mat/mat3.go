package mat

import (
	"github.com/tanema/amore/mth"
)

type Mat3 [9]float32

func New3(args ...float32) *Mat3 {
	var mat *Mat3
	mat.SetIdentity()
	if args != nil && len(args) > 0 {
		mat.SetTransformation(args...)
	}
	return mat
}

func (mat *Mat3) SetIdentity() {
	mat = new(Mat3)
	mat[0] = 1
	mat[4] = 1
	mat[8] = 1
}

func (mat *Mat3) Mul(other Mat3) *Mat3 {
	result := New3()

	result[0] = (mat[0] * other[0]) + (mat[3] * other[1]) + (mat[6] * other[2])
	result[3] = (mat[0] * other[3]) + (mat[3] * other[4]) + (mat[6] * other[5])
	result[6] = (mat[0] * other[6]) + (mat[3] * other[7]) + (mat[6] * other[8])

	result[1] = (mat[1] * other[0]) + (mat[4] * other[1]) + (mat[7] * other[2])
	result[4] = (mat[1] * other[3]) + (mat[4] * other[4]) + (mat[7] * other[5])
	result[7] = (mat[1] * other[6]) + (mat[4] * other[7]) + (mat[7] * other[8])

	result[2] = (mat[2] * other[0]) + (mat[5] * other[1]) + (mat[8] * other[2])
	result[5] = (mat[2] * other[3]) + (mat[5] * other[4]) + (mat[8] * other[5])
	result[8] = (mat[2] * other[6]) + (mat[5] * other[7]) + (mat[8] * other[8])

	return result
}

func (mat *Mat3) TransposedInverse() *Mat3 {
	det := mat[0]*(mat[4]*mat[8]-mat[7]*mat[5]) - mat[1]*(mat[3]*mat[8]-mat[5]*mat[6]) + mat[2]*(mat[3]*mat[7]-mat[4]*mat[6])
	invdet := 1.0 / det

	result := New3()
	result[0] = invdet * (mat[4]*mat[8] - mat[7]*mat[5])
	result[3] = -invdet * (mat[1]*mat[8] - mat[2]*mat[7])
	result[6] = invdet * (mat[1]*mat[5] - mat[2]*mat[4])
	result[1] = -invdet * (mat[3]*mat[8] - mat[5]*mat[6])
	result[4] = invdet * (mat[0]*mat[8] - mat[2]*mat[6])
	result[7] = -invdet * (mat[0]*mat[5] - mat[3]*mat[2])
	result[2] = invdet * (mat[3]*mat[7] - mat[6]*mat[4])
	result[5] = -invdet * (mat[0]*mat[7] - mat[6]*mat[1])
	result[8] = invdet * (mat[0]*mat[4] - mat[3]*mat[1])

	return result
}

func (mat *Mat3) SetTransformation(args ...float32) {
	x, y, angle, sx, sy, ox, oy, kx, ky := normalizeDrawCallArgs(args)
	mat.SetIdentity()
	c := mth.Cos(angle)
	s := mth.Sin(angle)

	mat[0] = c*sx - ky*s*sy // = a
	mat[1] = s*sx + ky*c*sy // = b
	mat[3] = kx*c*sx - s*sy // = c
	mat[4] = kx*s*sx + c*sy // = d
	mat[6] = x - ox*mat[0] - oy*mat[3]
	mat[7] = y - ox*mat[1] - oy*mat[4]

	mat[2] = 0.0
	mat[5] = 0.0
	mat[8] = 1.0
}

func (mat *Mat3) Transform(coords []float32) []float32 {
	output := make([]float32, len(coords))
	for i := 0; i < len(coords); i += 2 {
		x := coords[i]
		y := coords[i+1]

		output[i] = (mat[0] * x) + (mat[3] * y) + (mat[6])
		output[i+1] = (mat[1] * x) + (mat[4] * y) + (mat[7])
	}
	return output
}
