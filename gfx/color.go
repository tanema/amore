package gfx

type Color struct {
	R, G, B, A int
}

func (color *Color) Add(c *Color) {
	color.R = color.R + c.R
	color.G = color.G + c.G
	color.B = color.B + c.B
	color.A = color.A + c.A
}

func (color *Color) Mul(s int) {
	color.R = color.R * s
	color.G = color.G * s
	color.B = color.B * s
	color.A = color.A * s
}

func (color *Color) Div(s int) {
	color.R = color.R / s
	color.G = color.G / s
	color.B = color.B / s
	color.A = color.A / s
}
