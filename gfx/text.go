package gfx

import (
	"strings"
)

type (
	// Text is a container of text, color and text formatting.
	Text struct {
		font      *Font
		strings   []string
		colors    [][]float32
		wrapLimit float32
		align     string
		batches   map[*rasterizer]*SpriteBatch
		width     float32
		height    float32
	}
)

// Print will print out a colored string. It accepts the normal drawable arguments
func Print(strs []string, colors [][]float32, argv ...float32) {
	NewText(GetFont(), strs, colors, -1, "left").Draw(argv...)
}

// Printf will print out a string with a wrap limit and alignment. It accepts the
// normal drawable arguments
func Printf(strs []string, colors [][]float32, wrapLimit float32, align string, argv ...float32) {
	NewText(GetFont(), strs, colors, wrapLimit, align).Draw(argv...)
}

// NewText will create a colored text object with the provided font and
// text. A wrap and alignment can be provided as well. If wrapLimit is < 0 it will
// not wrap
func NewText(font *Font, strs []string, colors [][]float32, wrapLimit float32, align string) *Text {
	newText := &Text{
		font:      font,
		strings:   strs,
		colors:    colors,
		wrapLimit: wrapLimit,
		align:     align,
		batches:   make(map[*rasterizer]*SpriteBatch),
	}
	registerVolatile(newText)
	return newText
}

func (text *Text) loadVolatile() bool {
	length := len(strings.Join(text.strings, ""))
	for _, rast := range text.font.rasterizers {
		text.batches[rast] = NewSpriteBatch(rast.texture, length, UsageDynamic)
	}
	text.generate()
	return true
}

func (text *Text) unloadVolatile() {}

func (text *Text) generate() {
	for _, batch := range text.batches {
		batch.Clear()
	}

	var lines []*textLine
	lines, text.width, text.height = generateLines(text.font, text.strings, text.colors, text.wrapLimit)

	for _, l := range lines {
		var gx, spacing float32

		if spaceGlyph, _, ok := text.font.findGlyph(' '); ok {
			spacing = spaceGlyph.advance
		} else {
			spacing = text.font.rasterizers[0].advance
		}

		switch text.align {
		case "left":
		case "right":
			gx = text.wrapLimit - l.width
		case "center":
			gx = (text.wrapLimit - l.width) / 2.0
		case "justify":
			amountOfSpace := float32(l.spaceCount-1) * spacing
			widthWithoutSpace := l.width - amountOfSpace
			spacing = (text.wrapLimit - widthWithoutSpace) / float32(l.spaceCount)
		}

		for i := 0; i < l.size; i++ {
			ch := l.chars[i]
			if ch == ' ' {
				gx += spacing
			} else {
				glyph := l.glyphs[i]
				rast := l.rasts[i]
				text.batches[rast].SetColor(l.colors[i]...)
				gx += l.kern[i]
				text.batches[rast].Addq(glyph.quad, gx, l.y+glyph.descent)
				gx += glyph.advance
			}
		}
	}

	for _, batch := range text.batches {
		if batch.GetCount() > 0 {
			batch.SetBufferSize(batch.GetCount())
		}
	}
}

// GetWidth will return the text obejcts set width which will be <= wrapLimit
func (text *Text) GetWidth() float32 {
	return text.width
}

// GetHeight will return the height of the text object after text wrap.
func (text *Text) GetHeight() float32 {
	return text.height
}

// GetDimensions will return the width and height of the text object
func (text *Text) GetDimensions() (float32, float32) {
	return text.width, text.height
}

// GetFont will return the font that this text object has been created with
func (text *Text) GetFont() *Font {
	return text.font
}

// SetFont will set the font in which this text object will use to render the
// string
func (text *Text) SetFont(f *Font) {
	text.font = f
	text.loadVolatile()
}

// Set will set the string and colors for this text object to be rendered.
func (text *Text) Set(strs []string, colors [][]float32) {
	text.strings = strs
	text.colors = colors
	text.generate()
}

// Draw satisfies the Drawable interface. Inputs are as follows
// x, y, r, sx, sy, ox, oy, kx, ky
// x, y are position
// r is rotation
// sx, sy is the scale, if sy is not given sy will equal sx
// ox, oy are offset
// kx, ky are the shear. If ky is not given ky will equal kx
func (text *Text) Draw(args ...float32) {
	for _, batch := range text.batches {
		if batch.GetCount() > 0 {
			batch.Draw(args...)
		}
	}
}
