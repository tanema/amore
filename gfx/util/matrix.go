package util

import (
	"math"
)

type Matrix struct {
	e [16]float64
}

func newEmptyMatrix() *Matrix {
	return &Matrix{}
}

func newOrthoMatrix(left, right, bottom, top float64) *Matrix {
	m := newEmptyMatrix()
	m.e[0] = 2.0 / (right - left)
	m.e[5] = 2.0 / (top - bottom)
	m.e[10] = -1.0
	m.e[12] = -(right + left) / (right - left)
	m.e[13] = -(top + bottom) / (top - bottom)
	return m
}

func newMatrix(x, y, angle, sx, sy, ox, oy, kx, ky float64) *Matrix {
	new_matrix := &Matrix{}
	new_matrix.SetTransformation(x, y, angle, sx, sy, ox, oy, kx, ky)
	return new_matrix
}

func (mat *Matrix) SetTransformation(x, y, angle, sx, sy, ox, oy, kx, ky float64) {
	c := math.Cos(angle)
	s := math.Sin(angle)
	// matrix multiplication carried out on paper:
	// |1     x| |c -s    | |sx       | | 1 ky    | |1     -ox|
	// |  1   y| |s  c    | |   sy    | |kx  1    | |  1   -oy|
	// |    1  | |     1  | |      1  | |      1  | |    1    |
	// |      1| |       1| |        1| |        1| |       1 |
	//   move      rotate      scale       skew       origin
	mat.e[10] = 1.0
	mat.e[15] = 1.0
	mat.e[0] = c*sx - ky*s*sy // = a
	mat.e[1] = s*sx + ky*c*sy // = b
	mat.e[4] = kx*c*sx - s*sy // = c
	mat.e[5] = kx*s*sx + c*sy // = d
	mat.e[12] = x - ox*mat.e[0] - oy*mat.e[4]
	mat.e[13] = y - ox*mat.e[1] - oy*mat.e[5]
}

func (mat *Matrix) SetIdentity() {
	mat.e[0] = 1.0
	mat.e[5] = 1.0
	mat.e[10] = 1.0
	mat.e[15] = 1.0
}

func (mat *Matrix) GetElements() [16]float64 {
	return mat.e
}

func (mat *Matrix) SetTranslation(x, y float64) {
	mat.SetIdentity()
	mat.e[12] = x
	mat.e[13] = y
}

func (mat *Matrix) SetRotation(rad float64) {
	mat.SetIdentity()
	c := math.Cos(rad)
	s := math.Sin(rad)
	mat.e[0] = c
	mat.e[4] = -s
	mat.e[1] = s
	mat.e[5] = c
}

func (mat *Matrix) SetScale(sx, sy float64) {
	mat.SetIdentity()
	mat.e[0] = sx
	mat.e[5] = sy
}

func (mat *Matrix) SetShear(kx, ky float64) {
	mat.SetIdentity()
	mat.e[1] = ky
	mat.e[4] = kx
}

//                 | e0 e4 e8  e12 |
//                 | e1 e5 e9  e13 |
//                 | e2 e6 e10 e14 |
//                 | e3 e7 e11 e15 |
// | e0 e4 e8  e12 |
// | e1 e5 e9  e13 |
// | e2 e6 e10 e14 |
// | e3 e7 e11 e15 |

func (mat *Matrix) Mul(m *Matrix) *Matrix {
	t := newEmptyMatrix()

	t.e[0] = (mat.e[0] * m.e[0]) + (mat.e[4] * m.e[1]) + (mat.e[8] * m.e[2]) + (mat.e[12] * m.e[3])
	t.e[4] = (mat.e[0] * m.e[4]) + (mat.e[4] * m.e[5]) + (mat.e[8] * m.e[6]) + (mat.e[12] * m.e[7])
	t.e[8] = (mat.e[0] * m.e[8]) + (mat.e[4] * m.e[9]) + (mat.e[8] * m.e[10]) + (mat.e[12] * m.e[11])
	t.e[12] = (mat.e[0] * m.e[12]) + (mat.e[4] * m.e[13]) + (mat.e[8] * m.e[14]) + (mat.e[12] * m.e[15])

	t.e[1] = (mat.e[1] * m.e[0]) + (mat.e[5] * m.e[1]) + (mat.e[9] * m.e[2]) + (mat.e[13] * m.e[3])
	t.e[5] = (mat.e[1] * m.e[4]) + (mat.e[5] * m.e[5]) + (mat.e[9] * m.e[6]) + (mat.e[13] * m.e[7])
	t.e[9] = (mat.e[1] * m.e[8]) + (mat.e[5] * m.e[9]) + (mat.e[9] * m.e[10]) + (mat.e[13] * m.e[11])
	t.e[13] = (mat.e[1] * m.e[12]) + (mat.e[5] * m.e[13]) + (mat.e[9] * m.e[14]) + (mat.e[13] * m.e[15])

	t.e[2] = (mat.e[2] * m.e[0]) + (mat.e[6] * m.e[1]) + (mat.e[10] * m.e[2]) + (mat.e[14] * m.e[3])
	t.e[6] = (mat.e[2] * m.e[4]) + (mat.e[6] * m.e[5]) + (mat.e[10] * m.e[6]) + (mat.e[14] * m.e[7])
	t.e[10] = (mat.e[2] * m.e[8]) + (mat.e[6] * m.e[9]) + (mat.e[10] * m.e[10]) + (mat.e[14] * m.e[11])
	t.e[14] = (mat.e[2] * m.e[12]) + (mat.e[6] * m.e[13]) + (mat.e[10] * m.e[14]) + (mat.e[14] * m.e[15])

	t.e[3] = (mat.e[3] * m.e[0]) + (mat.e[7] * m.e[1]) + (mat.e[11] * m.e[2]) + (mat.e[15] * m.e[3])
	t.e[7] = (mat.e[3] * m.e[4]) + (mat.e[7] * m.e[5]) + (mat.e[11] * m.e[6]) + (mat.e[15] * m.e[7])
	t.e[11] = (mat.e[3] * m.e[8]) + (mat.e[7] * m.e[9]) + (mat.e[11] * m.e[10]) + (mat.e[15] * m.e[11])
	t.e[15] = (mat.e[3] * m.e[12]) + (mat.e[7] * m.e[13]) + (mat.e[11] * m.e[14]) + (mat.e[15] * m.e[15])

	return t
}

func (mat *Matrix) Translate(x, y float64) {
	t := newEmptyMatrix()
	t.SetTranslation(x, y)
	mat.e = mat.Mul(m).GetElements()
}

func (mat *Matrix) Rotate(rad float64) {
	t := newEmptyMatrix()
	t.SetRotation(rad)
	mat.e = mat.Mul(m).GetElements()
}

func (mat *Matrix) Scale(sx, sy float64) {
	t := newEmptyMatrix()
	t.SetScale(sx, sy)
	mat.e = mat.Mul(m).GetElements()
}

func (mat *Matrix) Shear(kx, ky float64) {
	t := newEmptyMatrix()
	t.SetShear(kx, ky)
	mat.e = mat.Mul(m).GetElements()
}
