package gfx

import (
	"math"
)

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

// NewHSLColor will create a color from HSL values, all between 0 and 1
// This code is taken from github.com/gerow/go-color, credit goes to Mike Gerow gerow@mgerow.com
func NewHSLColor(h, s, l, a float32) *Color {
	if s == 0 { // it's gray
		return &Color{l, l, l, a}
	}

	var v1, v2 float32
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - (s * l)
	}

	v1 = 2*l - v2

	r := hueToRGB(v1, v2, h+(1.0/3.0))
	g := hueToRGB(v1, v2, h)
	b := hueToRGB(v1, v2, h-(1.0/3.0))

	return &Color{r, g, b, a}
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
	newColor := &Color{
		color[0] - (color[0] * percent),
		color[1] - (color[1] * percent),
		color[2] - (color[2] * percent),
		color[3],
	}

	//clamp colors between 0 and 1
	for i, c := range newColor {
		newColor[i] = float32(math.Max(0, math.Min(1, float64(c))))
	}

	return newColor
}

func hueToRGB(v1, v2, h float32) float32 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	switch {
	case 6*h < 1:
		return (v1 + (v2-v1)*6*h)
	case 2*h < 1:
		return v2
	case 3*h < 2:
		return v1 + (v2-v1)*((2.0/3.0)-h)*6
	}
	return v1
}
