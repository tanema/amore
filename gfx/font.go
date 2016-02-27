package gfx

import (
	"math"
	"strings"
)

type (
	Font struct {
		rasterizers []rasterizer
		lineHeight  float32
	}
)

func NewFont(filename string, font_size float32) *Font {
	rast := newTtfRasterizer(filename, font_size)
	return &Font{
		rasterizers: []rasterizer{rast},
	}
}

func NewImageFont(filename, glyph_hints string) *Font {
	rast := newImageRasterizer(filename, glyph_hints)
	return &Font{
		rasterizers: []rasterizer{rast},
	}
}

func (font *Font) Release() {
	font.rasterizers[0].Release() //only release the first because the rest are fallbacks
}

func (font *Font) setLineHeight(height float32) {
	font.lineHeight = height
}

func (font *Font) GetLineHeight() float32 {
	if font.lineHeight <= 0 {
		return float32(font.rasterizers[0].getLineHeight())
	} else {
		return font.lineHeight
	}
}

func (font *Font) SetFilter(min, mag FilterMode) error {
	for _, rasterizer := range font.rasterizers {
		if err := rasterizer.getTexture().SetFilter(min, mag); err != nil {
			return err
		}
	}
	return nil
}

func (font *Font) GetFilter() Filter {
	return font.rasterizers[0].getTexture().GetFilter()
}

func (font *Font) GetAscent() int {
	return font.rasterizers[0].getAscent()
}

func (font *Font) GetDescent() int {
	return font.rasterizers[0].getDescent()
}

func (font *Font) GetBaseline() int {
	return font.rasterizers[0].getLineHeight()
}

func (font *Font) HasGlyph(g rune) bool {
	for _, rasterizer := range font.rasterizers {
		if rasterizer.hasGlyph(g) {
			return true
		}
	}
	return false
}

func (font *Font) HasGlyphs(text string) bool {
	if len(text) == 0 {
		return false
	}

	for _, c := range text {
		if font.HasGlyph(c) {
			return false
		}
	}

	return true
}

func (font *Font) findGlyph(g rune) glyphData {
	for _, rasterizer := range font.rasterizers {
		if rasterizer.hasGlyph(g) {
			return rasterizer.getGlyphData(g)
		}
	}
	return font.rasterizers[0].getGlyphData(g)
}

func (font *Font) getKerning(first, second rune) float32 {
	k := font.rasterizers[0].getKerning(first, second)

	for _, r := range font.rasterizers {
		if r.hasGlyph(first) && r.hasGlyph(second) {
			k = r.getKerning(first, second)
			break
		}
	}

	return k
}

func (font *Font) SetFallbacks(fallbacks ...*Font) {
	if fallbacks == nil || len(fallbacks) == 0 {
		return
	}
	for _, fallback := range fallbacks {
		font.rasterizers = append(font.rasterizers, fallback.rasterizers...)
	}
}

func (font *Font) GetHeight() float32 {
	return float32(font.rasterizers[0].getHeight())
}

func (font *Font) GetWidth(text string) float32 {
	if len(text) == 0 {
		return 0
	}

	var max_width float32
	for _, line := range strings.Split(text, "\n") {
		var width float32
		var prevChar rune
		for i, char := range string(line[:]) {
			g := font.findGlyph(char)
			width += float32(g.advanceWidth)
			if i != 0 {
				width += font.getKerning(char, prevChar)
			}
			prevChar = char
		}
		max_width = float32(math.Max(float64(max_width), float64(width)))
	}

	return max_width
}

func (font *Font) GetWrap(text string, wrapLimit float32) (float32, []string) {
	var width, currentWidth float32
	var lines, currentLine []string

	for _, word := range strings.Split(text, " ") {
		wordWidth := font.GetWidth(word)
		if currentWidth+wordWidth > wrapLimit {
			if len(currentLine) > 0 {
				lines = append(lines, strings.Join(currentLine, " "))
				width = float32(math.Max(float64(currentWidth), float64(width)))
			}
			currentLine = []string{word}
			currentWidth = wordWidth
		} else {
			currentLine = append(currentLine, word)
			currentWidth += wordWidth
		}
	}

	lines = append(lines, strings.Join(currentLine, " "))
	width = float32(math.Max(float64(currentWidth), float64(width)))

	return width, lines
}
