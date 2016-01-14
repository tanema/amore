package gfx

type Color [4]float32

func (color *Color) RGB() (r, g, b float32) {
	return color[0], color[1], color[2]
}

func (color *Color) RGBA() (r, g, b, a float32) {
	return color[0], color[1], color[2], color[3]
}

func (color *Color) Add(c *Color) {
	color[0] = color[0] + c[0]
	color[1] = color[1] + c[1]
	color[2] = color[2] + c[2]
	color[3] = color[3] + c[3]
}

func (color *Color) Mul(s int) {
	color[0] = color[0] * (float32(s) / 255.0)
	color[1] = color[1] * (float32(s) / 255.0)
	color[2] = color[2] * (float32(s) / 255.0)
	color[3] = color[3] * (float32(s) / 255.0)
}

func (color *Color) Div(s int) {
	color[0] = color[0] / (float32(s) / 255.0)
	color[1] = color[1] / (float32(s) / 255.0)
	color[2] = color[2] / (float32(s) / 255.0)
	color[3] = color[3] / (float32(s) / 255.0)
}
