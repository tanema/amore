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
		batches   []*SpriteBatch
		width     float32
		height    float32
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
	return NewColorTextExt(font, []string{text}, []*Color{NewColor(255, 255, 255, 255)}, -1, ALIGN_LEFT)
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
		batches:   []*SpriteBatch{},
	}

	new_text.generate()
	return new_text, nil
}

func (text *Text) generate() {
	if text.wrapLimit > 0 {
		text.generateFormatted()
	} else {
		text.generateUnformatted()
	}
}

func (text *Text) generateUnformatted() {
	batches := make(map[rasterizer]*SpriteBatch)
	for _, rast := range text.font.rasterizers {
		batches[rast] = NewSpriteBatch(rast.getTexture(), text.length)
	}

	var gx, gy float32
	for i, st := range text.strings {
		for _, char := range st {
			for _, rast := range text.font.rasterizers {
				if rast.hasGlyph(char) {
					batches[rast].SetColor(text.colors[i])
					glyph := rast.getGlyphData(char)
					batches[rast].Addq(glyph.rec, gx+float32(glyph.leftSideBearing-rast.getOffset()), gy+float32(glyph.topSideBearing-rast.getOffset()))
					gx = gx + float32(glyph.advanceWidth)
					break
				}
			}
		}
	}

	text.batches = []*SpriteBatch{}
	for _, rast := range text.font.rasterizers {
		if batches[rast].GetCount() > 0 {
			batches[rast].SetBufferSize(batches[rast].GetCount())
			text.batches = append(text.batches, batches[rast])
		}
	}

	text.width = text.font.GetWidth(strings.Join(text.strings, ""))
	text.height = text.font.GetHeight()
}

func (text *Text) generateFormatted() {
	batches := make(map[rasterizer]*SpriteBatch)
	for _, rast := range text.font.rasterizers {
		batches[rast] = NewSpriteBatch(rast.getTexture(), text.length)
	}

	var gx, gy float32
	for i, st := range text.strings {
		for _, char := range st {
			for _, rast := range text.font.rasterizers {
				if rast.hasGlyph(char) {
					batches[rast].SetColor(text.colors[i])
					glyph := rast.getGlyphData(char)
					batches[rast].Addq(glyph.rec, gx+float32(glyph.leftSideBearing-rast.getOffset()), gy+float32(glyph.topSideBearing-rast.getOffset()))
					gx = gx + float32(glyph.advanceWidth)
					break
				}
			}
		}
	}

	text.batches = []*SpriteBatch{}
	for _, rast := range text.font.rasterizers {
		if batches[rast].GetCount() > 0 {
			batches[rast].SetBufferSize(batches[rast].GetCount())
			text.batches = append(text.batches, batches[rast])
		}
	}

	text.width = text.font.GetWidth(strings.Join(text.strings, ""))
	text.height = text.font.GetHeight()
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
	text.generate()
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
