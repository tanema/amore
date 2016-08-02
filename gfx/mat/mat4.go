package mat

type Mat4 [16]float32

func New4(args ...float32) *Mat4 {
	mat := &Mat4{}
	mat.SetIdentity()
	if args != nil && len(args) > 0 {
		mat.SetTransformation(args...)
	}
	return mat
}

func Ortho(left, right, bottom, top float32) *Mat4 {
	m := New4()

	m[0] = 2.0 / (right - left)
	m[5] = 2.0 / (top - bottom)
	m[10] = -1.0

	m[12] = -(right + left) / (right - left)
	m[13] = -(top + bottom) / (top - bottom)

	return m
}

func (mat *Mat4) SetIdentity() {
	*mat = Mat4{}
	mat[0] = 1
	mat[5] = 1
	mat[10] = 1
	mat[15] = 1
}

func (mat *Mat4) SetTranslation(x, y float32) {
	mat.SetIdentity()
	mat[12] = x
	mat[13] = y
}

func (mat *Mat4) SetRotation(rad float32) {
	mat.SetIdentity()
	c := cos(rad)
	s := sin(rad)
	mat[0] = c
	mat[4] = -s
	mat[1] = s
	mat[5] = c
}

func (mat *Mat4) SetScale(sx, sy float32) {
	mat.SetIdentity()
	mat[0] = sx
	mat[5] = sy
}

func (mat *Mat4) SetShear(kx, ky float32) {
	mat.SetIdentity()
	mat[1] = ky
	mat[4] = kx
}

func (mat *Mat4) SetTransformation(args ...float32) {
	x, y, angle, sx, sy, ox, oy, kx, ky := normalizeDrawCallArgs(args)
	mat.SetIdentity()
	c := cos(angle)
	s := sin(angle)
	mat[10] = 1
	mat[15] = 1
	mat[0] = c*sx - ky*s*sy // = a
	mat[1] = s*sx + ky*c*sy // = b
	mat[4] = kx*c*sx - s*sy // = c
	mat[5] = kx*s*sx + c*sy // = d
	mat[12] = x - ox*mat[0] - oy*mat[4]
	mat[13] = y - ox*mat[1] - oy*mat[5]
}

func (mat *Mat4) Transform(coords []float32) []float32 {
	output := make([]float32, len(coords))
	for i := 0; i < len(coords); i += 2 {
		x := coords[i]
		y := coords[i+1]

		output[i] = (mat[0] * x) + (mat[4] * y) + (0) + (mat[12])
		output[i+1] = (mat[1] * x) + (mat[5] * y) + (0) + (mat[13])
	}
	return output
}

func (mat *Mat4) Translate(x, y float32) {
	t := New4()
	t.SetTranslation(x, y)
	mat = mat.Mul(t)
}

func (mat *Mat4) Rotate(rad float32) {
	t := New4()
	t.SetRotation(rad)
	mat = mat.Mul(t)
}

func (mat *Mat4) Scale(sx, sy float32) {
	t := New4()
	t.SetScale(sx, sy)
	mat = mat.Mul(t)
}

func (mat *Mat4) Shear(kx, ky float32) {
	t := New4()
	t.SetShear(kx, ky)
	mat = mat.Mul(t)
}

func (mat *Mat4) Mul(other *Mat4) *Mat4 {
	result := New4()

	result[0] = (mat[0] * other[0]) + (mat[4] * other[1]) + (mat[8] * other[2]) + (mat[12] * other[3])
	result[4] = (mat[0] * other[4]) + (mat[4] * other[5]) + (mat[8] * other[6]) + (mat[12] * other[7])
	result[8] = (mat[0] * other[8]) + (mat[4] * other[9]) + (mat[8] * other[10]) + (mat[12] * other[11])
	result[12] = (mat[0] * other[12]) + (mat[4] * other[13]) + (mat[8] * other[14]) + (mat[12] * other[15])

	result[1] = (mat[1] * other[0]) + (mat[5] * other[1]) + (mat[9] * other[2]) + (mat[13] * other[3])
	result[5] = (mat[1] * other[4]) + (mat[5] * other[5]) + (mat[9] * other[6]) + (mat[13] * other[7])
	result[9] = (mat[1] * other[8]) + (mat[5] * other[9]) + (mat[9] * other[10]) + (mat[13] * other[11])
	result[13] = (mat[1] * other[12]) + (mat[5] * other[13]) + (mat[9] * other[14]) + (mat[13] * other[15])

	result[2] = (mat[2] * other[0]) + (mat[6] * other[1]) + (mat[10] * other[2]) + (mat[14] * other[3])
	result[6] = (mat[2] * other[4]) + (mat[6] * other[5]) + (mat[10] * other[6]) + (mat[14] * other[7])
	result[10] = (mat[2] * other[8]) + (mat[6] * other[9]) + (mat[10] * other[10]) + (mat[14] * other[11])
	result[14] = (mat[2] * other[12]) + (mat[6] * other[13]) + (mat[10] * other[14]) + (mat[14] * other[15])

	result[3] = (mat[3] * other[0]) + (mat[7] * other[1]) + (mat[11] * other[2]) + (mat[15] * other[3])
	result[7] = (mat[3] * other[4]) + (mat[7] * other[5]) + (mat[11] * other[6]) + (mat[15] * other[7])
	result[11] = (mat[3] * other[8]) + (mat[7] * other[9]) + (mat[11] * other[10]) + (mat[15] * other[11])
	result[15] = (mat[3] * other[12]) + (mat[7] * other[13]) + (mat[11] * other[14]) + (mat[15] * other[15])

	return result
}
