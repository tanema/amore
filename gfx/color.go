package gfx

// Color represents an rgba color
type Color [4]float32

// ColorMask contains an rgba color mask
type ColorMask struct {
	r, g, b, a bool
}

// NewColor will take the values r, g, b, a in the range from 0 to 255 and return
// a pointer to the new color
func NewColor(r, g, b, a float32) *Color {
	return &Color{
		r / 255.0,
		g / 255.0,
		b / 255.0,
		a / 255.0,
	}
}

// Sub will subtract one color from another and return the resulting color
func (color *Color) Sub(c *Color) *Color {
	return &Color{
		color[0] - c[0],
		color[1] - c[1],
		color[2] - c[2],
		color[3] - c[3],
	}
}

// Add will add one color to another and return the resulting color
func (color *Color) Add(c *Color) *Color {
	return &Color{
		color[0] + c[0],
		color[1] + c[1],
		color[2] + c[2],
		color[3] + c[3],
	}
}

// Mul will multiply one color to another and return the resulting color
func (color *Color) Mul(s float32) *Color {
	return &Color{
		color[0] * (s / 255.0),
		color[1] * (s / 255.0),
		color[2] * (s / 255.0),
		color[3] * (s / 255.0),
	}
}

// Div will divide one color from another and return the resulting color
func (color *Color) Div(s float32) *Color {
	return &Color{
		color[0] / (s / 255.0),
		color[1] / (s / 255.0),
		color[2] / (s / 255.0),
		color[3] / (s / 255.0),
	}
}

// Darken will darken the color by the percent given
func (color *Color) Darken(percent float32) *Color {
	return color.changeShade(percent)
}

// Lighten will lighten the color by the percent given
func (color *Color) Lighten(percent float32) *Color {
	return color.changeShade(-percent)
}

func (color *Color) changeShade(percent float32) *Color {
	if percent > 100 || percent < -100 {
		return color
	}

	percent = percent * 0.01
	new_color := &Color{
		color[0] - (color[0] * percent),
		color[1] - (color[1] * percent),
		color[2] - (color[2] * percent),
		color[3],
	}

	for i, c := range new_color {
		new_color[i] = clamp(c, 0, 1)
	}

	return new_color
}
