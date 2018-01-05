package gfx

import (
	"math"
)

type textLine struct {
	chars     []rune
	glyphs    []glyphData
	colors    []*Color
	rasts     []*rasterizer
	kern      []float32
	size      int
	lastBreak int
	width     float32
	y         float32
}

func generateLines(font *Font, text []string, color []*Color, wrapLimit float32) ([]*textLine, float32, float32) {
	var lines []*textLine
	var prevChar rune
	var width, gy float32
	currentLine := &textLine{}

	breakOff := func(y float32, immediate bool) {
		newLine := currentLine.breakOff(immediate)
		newLine.y = y
		width = float32(math.Max(float64(width), float64(currentLine.width)))
		lines = append(lines, currentLine)
		currentLine = newLine
		prevChar = 0
	}

	for i, st := range text {
		for _, char := range st {
			if char == '\n' {
				gy += font.GetLineHeight()
				breakOff(gy, true)
				continue
			} else if char == '\r' {
				breakOff(gy, true)
				continue
			} else if char == '\t' {
				if glyph, rast, ok := font.findGlyph(' '); ok {
					currentLine.add(' ', glyph, color[i], rast, font.Kern(prevChar, char))
					currentLine.add(' ', glyph, color[i], rast, font.Kern(' ', char))
					currentLine.add(' ', glyph, color[i], rast, font.Kern(' ', char))
					currentLine.add(' ', glyph, color[i], rast, font.Kern(' ', char))
				}
				continue
			}
			if glyph, rast, ok := font.findGlyph(char); ok {
				currentLine.add(char, glyph, color[i], rast, font.Kern(prevChar, char))
			}
			if wrapLimit > 0 && currentLine.width >= wrapLimit {
				gy += font.GetLineHeight()
				breakOff(gy, false)
			} else {
				prevChar = char
			}
		}
	}

	lines = append(lines, currentLine)

	return lines, width, gy + font.GetLineHeight()
}

func (l *textLine) add(char rune, g glyphData, color *Color, rast *rasterizer, kern float32) {
	if char == ' ' {
		l.lastBreak = l.size
	}
	l.chars = append(l.chars, char)
	l.glyphs = append(l.glyphs, g)
	l.colors = append(l.colors, color)
	l.rasts = append(l.rasts, rast)
	l.kern = append(l.kern, kern+g.lsb)
	l.size++
	l.width += g.advance + kern + g.lsb
}

func (l *textLine) breakOff(immediate bool) *textLine {
	breakPoint := l.lastBreak
	if l.lastBreak == -1 || immediate {
		breakPoint = l.size - 1
	}

	newLine := &textLine{}

	for i := l.size - 1; i > breakPoint; i-- {
		ch, g, cl, r, k := l.trimLastChar()
		newLine.chars = append([]rune{ch}, newLine.chars...)
		newLine.glyphs = append([]glyphData{g}, newLine.glyphs...)
		newLine.colors = append([]*Color{cl}, newLine.colors...)
		newLine.rasts = append([]*rasterizer{r}, newLine.rasts...)
		newLine.kern = append([]float32{k}, newLine.kern...)
		newLine.size++
		newLine.width += float32(g.advance) + k
	}

	for i := l.size - 1; i >= 0; i-- {
		if l.chars[i] == ' ' {
			l.trimLastChar()
			continue
		}
		break
	}

	return newLine
}

func (l *textLine) trimLastChar() (rune, glyphData, *Color, *rasterizer, float32) {
	i := l.size - 1
	ch, g, cl, r, k := l.chars[i], l.glyphs[i], l.colors[i], l.rasts[i], l.kern[i]
	l.chars = l.chars[:i]
	l.glyphs = l.glyphs[:i]
	l.colors = l.colors[:i]
	l.rasts = l.rasts[:i]
	l.kern = l.kern[:i]
	l.size--
	l.width -= float32(g.advance) + k
	return ch, g, cl, r, k
}
