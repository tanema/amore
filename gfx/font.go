package gfx

import (
	"math"
	"strings"
)

// Font is a rasterized font data
type Font struct {
	rasterizers []rasterizer
	lineHeight  float32
}

// NewFont rasterizes a ttf font and returns a pointer to a new Font
func NewFont(filename string, fontSize float32) *Font {
	rast := newTtfRasterizer(filename, fontSize)
	return &Font{
		rasterizers: []rasterizer{rast},
	}
}

// NewImageFont rasterizes an image using the glyphHints. The glyphHints should
// list all characters in the image. The characters should all have equal width
// and height. Using the glyphHints, the image is split up into equal rectangles
// for rasterization. The function will return a pointer to a new Font
func NewImageFont(filename, glyphHints string) *Font {
	rast := newImageRasterizer(filename, glyphHints)
	return &Font{
		rasterizers: []rasterizer{rast},
	}
}

// SetLineHeight sets the height between lines
func (font *Font) SetLineHeight(height float32) {
	font.lineHeight = height
}

// GetLineHeight will return the current line height of the font
func (font *Font) GetLineHeight() float32 {
	if font.lineHeight <= 0 {
		return float32(font.rasterizers[0].getLineHeight())
	}
	return font.lineHeight
}

// SetFilter sets the filtering on the font.
func (font *Font) SetFilter(min, mag FilterMode) error {
	for _, rasterizer := range font.rasterizers {
		if err := rasterizer.getTexture().SetFilter(min, mag); err != nil {
			return err
		}
	}
	return nil
}

// GetFilter will return the filter of the font
func (font *Font) GetFilter() Filter {
	return font.rasterizers[0].getTexture().GetFilter()
}

// GetAscent gets the height of the font from the baseline
func (font *Font) GetAscent() int {
	return font.rasterizers[0].getAscent()
}

// GetDescent gets the height of the font below the base line
func (font *Font) GetDescent() int {
	return font.rasterizers[0].getDescent()
}

// GetBaseline returns the position of the base line.
func (font *Font) GetBaseline() int {
	return font.rasterizers[0].getLineHeight()
}

// HasGlyph checks if this font has a character for the given rune
func (font *Font) HasGlyph(g rune) bool {
	for _, rasterizer := range font.rasterizers {
		if rasterizer.hasGlyph(g) {
			return true
		}
	}
	return false
}

// HasGlyphs will check every glyph in the string and return if it has a character
// for each character
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

// findGlyph will fetch the glyphData for the given rune
func (font *Font) findGlyph(g rune) glyphData {
	for _, rasterizer := range font.rasterizers {
		if rasterizer.hasGlyph(g) {
			return rasterizer.getGlyphData(g)
		}
	}
	return font.rasterizers[0].getGlyphData(g)
}

// getKerning will return the space between two characters
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

// SetFallbacks will add extra fonts in case some characters are not available
// in this font. If the character is not available it will be rendered with one
// of the fallback characters
func (font *Font) SetFallbacks(fallbacks ...*Font) {
	if fallbacks == nil || len(fallbacks) == 0 {
		return
	}
	for _, fallback := range fallbacks {
		font.rasterizers = append(font.rasterizers, fallback.rasterizers...)
	}
}

// GetHeight will get the height of the font.
func (font *Font) GetHeight() float32 {
	return float32(font.rasterizers[0].getHeight())
}

// GetWidth will get the width of a given string after rendering.
func (font *Font) GetWidth(text string) float32 {
	if len(text) == 0 {
		return 0
	}

	var maxWidth float32
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
		maxWidth = float32(math.Max(float64(maxWidth), float64(width)))
	}

	return maxWidth
}

// GetWrap will split a string given a wrap limit. It will return the max width
// of the longest string and it will return the string split into the strings that
// are smaller than the wrap limit.
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
