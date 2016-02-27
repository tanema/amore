package gfx

import (
	"fmt"
	"strings"
)

type (
	Text struct {
		font      *Font
		strings   []string
		colors    []*Color
		wrapLimit float32
		align     AlignMode
		length    int
		batches   map[rasterizer]*SpriteBatch
		width     float32
		height    float32
		spaceSize float32
	}
)

func Print(fs string, argv ...float32) {
	text, err := NewText(GetFont(), fs)
	if err != nil {
		return
	}
	text.Draw(argv...)
	text.Release()
}

func Printc(strs []string, colors []*Color, argv ...float32) {
	text, err := NewColorText(GetFont(), strs, colors)
	if err != nil {
		return
	}
	text.Draw(argv...)
	text.Release()
}

func Printf(fs string, wrapLimit float32, align AlignMode, argv ...float32) {
	text, err := NewTextExt(GetFont(), fs, wrapLimit, align)
	if err != nil {
		return
	}
	text.Draw(argv...)
	text.Release()
}

func Printfc(strs []string, colors []*Color, wrapLimit float32, align AlignMode, argv ...float32) {
	text, err := NewColorTextExt(GetFont(), strs, colors, wrapLimit, align)
	if err != nil {
		return
	}
	text.Draw(argv...)
	text.Release()
}

func NewText(font *Font, text string) (*Text, error) {
	return NewTextExt(font, text, -1, ALIGN_LEFT)
}

func NewTextExt(font *Font, text string, wrap_limit float32, align AlignMode) (*Text, error) {
	if text == "" {
		return nil, fmt.Errorf("Cannot create an text object with blank string")
	}
	return NewColorTextExt(font, []string{text}, []*Color{NewColor(255, 255, 255, 255)}, wrap_limit, align)
}

func NewColorText(font *Font, strs []string, colors []*Color) (*Text, error) {
	return NewColorTextExt(font, strs, colors, -1, ALIGN_LEFT)
}

func NewColorTextExt(font *Font, strs []string, colors []*Color, wrap_limit float32, align AlignMode) (*Text, error) {
	if len(strs) == 0 {
		return nil, fmt.Errorf("Nothing to print")
	}

	if len(strs) != len(colors) {
		return nil, fmt.Errorf("Improper countof strings to colors")
	}

	new_text := &Text{
		font:      font,
		strings:   strs,
		colors:    colors,
		wrapLimit: wrap_limit,
		align:     align,
		length:    len(strings.Join(strs, "")),
		batches:   make(map[rasterizer]*SpriteBatch),
	}

	registerVolatile(new_text)

	return new_text, nil
}

func (text *Text) loadVolatile() bool {
	for _, rast := range text.font.rasterizers {
		text.batches[rast] = NewSpriteBatch(rast.getTexture(), text.length)
	}
	spaceGlyph := text.getSpaceGlyph()
	text.spaceSize = float32(spaceGlyph.advanceWidth)
	text.generate()
	return true
}

func (text *Text) unloadVolatile() {
	text.Release()
}

func (text *Text) generate() {
	for _, batch := range text.batches {
		batch.Clear()
	}
	if text.wrapLimit > 0 {
		text.generateFormatted()
	} else {
		text.generateUnformatted()
	}
}

func (text *Text) generateUnformatted() {
	var gx, gy float32
	for i, st := range text.strings {
		var prevChar rune
		for _, char := range st {
			for _, rast := range text.font.rasterizers {
				if rast.hasGlyph(char) {
					text.batches[rast].SetColor(text.colors[i])
					glyph := rast.getGlyphData(char)
					if prevChar != 0 {
						gx += text.font.getKerning(char, prevChar)
					}
					text.batches[rast].Addq(glyph.rec, gx+float32(glyph.leftSideBearing-rast.getOffset()), gy+float32(glyph.topSideBearing-rast.getOffset()))
					gx = gx + float32(glyph.advanceWidth)
					break
				}
			}
			prevChar = char
		}
	}
	text.compressBatches()
	text.width = text.font.GetWidth(strings.Join(text.strings, ""))
	text.height = text.font.GetHeight()
}

func (text *Text) generateFormatted() {
	var currentWidth, gy float32
	var currentLine []*word
	text.width = 0
	for _, w := range text.generateWords() {
		if (currentWidth + text.spaceSize + w.width) > text.wrapLimit {
			text.drawLine(currentLine, currentWidth, gy)
			currentLine = []*word{w}
			currentWidth = w.width
			gy += text.font.GetLineHeight()
		} else {
			currentLine = append(currentLine, w)
			currentWidth += (text.spaceSize + w.width)
		}
	}
	text.drawLine(currentLine, currentWidth, gy)
	text.compressBatches()
	text.width = text.wrapLimit
	text.height = gy + text.font.GetLineHeight()
}

func (text *Text) drawLine(currentLine []*word, lineWidth, gy float32) {
	if len(currentLine) == 0 {
		return
	}
	spaceing := text.spaceSize
	var gx float32
	switch text.align {
	case ALIGN_LEFT:
	case ALIGN_RIGHT:
		gx = text.wrapLimit - lineWidth
	case ALIGN_CENTER:
		gx = (text.wrapLimit - lineWidth) / 2.0
	case ALIGN_JUSTIFY:
		spaceing = (text.wrapLimit - lineWidth) / float32(len(currentLine)-1)
	}

	for _, w := range currentLine {
		for i := 0; i < w.size; i++ {
			glyph := w.glyphs[i]
			rast := w.rasts[i]
			text.batches[rast].SetColor(w.colors[i])
			gx += w.kern[i]
			text.batches[rast].Addq(glyph.rec, gx+float32(glyph.leftSideBearing-rast.getOffset()), gy+float32(glyph.topSideBearing-rast.getOffset()))
			gx += float32(glyph.advanceWidth)
		}
		gx += spaceing
	}
}

func (text *Text) compressBatches() {
	for _, batch := range text.batches {
		if batch.GetCount() > 0 {
			batch.SetBufferSize(batch.GetCount())
		}
	}
}

func (text *Text) getSpaceGlyph() glyphData {
	for _, rast := range text.font.rasterizers {
		if rast.hasGlyph(' ') {
			return rast.getGlyphData(' ')
		}
	}
	return text.font.rasterizers[0].getGlyphData(' ')
}

func (text *Text) generateWords() []*word {
	words := []*word{}
	currentWord := newWord()

	for i, st := range text.strings {
		var prevChar rune
		for _, char := range st {
			if char == ' ' {
				words = append(words, currentWord)
				currentWord = newWord()
				//prevChar = 0
				continue
			}
			for _, rast := range text.font.rasterizers {
				if rast.hasGlyph(char) {
					var kern float32
					if prevChar != 0 {
						kern = text.font.getKerning(char, prevChar)
					}
					currentWord.add(rast.getGlyphData(char), text.colors[i], rast, kern)
					break
				}
			}
			prevChar = char
		}
	}

	words = append(words, currentWord)

	return words
}

func (text *Text) GetWidth() float32 {
	return text.width
}

func (text *Text) GetHeight() float32 {
	return text.height
}

func (text *Text) GetDimensions() (float32, float32) {
	return text.width, text.height
}

func (text *Text) GetFont() *Font {
	return text.font
}

func (text *Text) SetFont(f *Font) {
	text.font = f
	text.loadVolatile()
}

func (text *Text) Set(t string) {
	text.Setc([]string{t}, []*Color{NewColor(255, 255, 255, 255)})
}

func (text *Text) Setc(strs []string, colors []*Color) {
	text.strings = strs
	text.colors = colors
	text.generate()
}

func (text *Text) Add(t string, args ...float32) {
	text.Addc([]string{t}, []*Color{NewColor(255, 255, 255, 255)}, args...)
}

func (text *Text) Addc(strs []string, colors []*Color, args ...float32) {
	text.Addfc(strs, colors, -1, ALIGN_LEFT, args...)
}

func (text *Text) Addf(t string, wrapLimit float32, align AlignMode, args ...float32) {
	text.Addfc([]string{t}, []*Color{NewColor(255, 255, 255, 255)}, wrapLimit, align, args...)
}

func (text *Text) Addfc(strs []string, colors []*Color, wrapLimit float32, align AlignMode, args ...float32) {
}

func (text *Text) Release() {
	for _, batch := range text.batches {
		batch.Release()
	}
}

func (text *Text) Draw(args ...float32) {
	for _, batch := range text.batches {
		batch.Draw(args...)
	}
}

type word struct {
	glyphs []glyphData
	colors []*Color
	rasts  []rasterizer
	kern   []float32
	size   int
	width  float32
}

func newWord() *word {
	return &word{
		glyphs: []glyphData{},
		colors: []*Color{},
		rasts:  []rasterizer{},
	}
}

func (w *word) add(g glyphData, color *Color, rast rasterizer, kern float32) {
	w.glyphs = append(w.glyphs, g)
	w.colors = append(w.colors, color)
	w.rasts = append(w.rasts, rast)
	w.kern = append(w.kern, kern)
	w.size++
	w.width += float32(g.advanceWidth) + kern
}
