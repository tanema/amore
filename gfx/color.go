package gfx

type Color [4]float32

type ColorMask struct {
	r, g, b, a bool
}

func NewColor(r, g, b, a float32) Color {
	return Color{
		r / 255.0,
		g / 255.0,
		b / 255.0,
		a / 255.0,
	}
}

func (color *Color) RGB() (r, g, b float32) {
	return color[0], color[1], color[2]
}

func (color *Color) RGBA() (r, g, b, a float32) {
	return color[0], color[1], color[2], color[3]
}

func (color *Color) Add(c *Color) *Color {
	return &Color{
		color[0] + c[0],
		color[1] + c[1],
		color[2] + c[2],
		color[3] + c[3],
	}
}

func (color *Color) Mul(s float32) *Color {
	return &Color{
		color[0] * (s / 255.0),
		color[1] * (s / 255.0),
		color[2] * (s / 255.0),
		color[3] * (s / 255.0),
	}
}

func (color *Color) Div(s float32) *Color {
	return &Color{
		color[0] / (s / 255.0),
		color[1] / (s / 255.0),
		color[2] / (s / 255.0),
		color[3] / (s / 255.0),
	}
}
