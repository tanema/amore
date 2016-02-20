package gfx

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
)

type (
	Font struct {
		rasterizers []rasterizer
		lineHeight  float32
	}
)

func NewFont(filename string, font_size float32) *Font {
	return &Font{
		rasterizers: []rasterizer{newTtfRasterizer(filename, font_size)},
		lineHeight:  1,
	}
}

func NewImageFont(filename, glyph_hints string) *Font {
	return &Font{
		rasterizers: []rasterizer{newImageRasterizer(filename, glyph_hints)},
		lineHeight:  1,
	}
}

func Printf(x, y float32, fs string, argv ...interface{}) {
	getCheckFont().Printf(x, y, fs, argv...)
}

func (font *Font) setLineHeight(height float32) {
	font.lineHeight = height
}

func (font *Font) GetLineHeight() float32 {
	return font.lineHeight
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

func (font *Font) SetFallbacks(fallbacks ...*Font) {
	if fallbacks == nil || len(fallbacks) == 0 {
		return
	}
	for _, fallback := range fallbacks {
		font.rasterizers = append(font.rasterizers, fallback.rasterizers...)
	}
}

func (font *Font) Printf(x, y float32, fs string, argv ...interface{}) {
	formatted_string := fmt.Sprintf(fs, argv...)
	if len(formatted_string) == 0 {
		return
	}

	rast := font.rasterizers[0]

	x = x - float32(rast.getOffset())
	y = y - float32(rast.getOffset())

	bindTexture(rast.getTexture().GetHandle())
	useVertexAttribArrays(ATTRIBFLAG_POS | ATTRIBFLAG_TEXCOORD)
	for _, ch := range formatted_string {
		if rast.hasGlyph(ch) {
			g := rast.getGlyphData(ch)
			prepareDraw(generateModelMatFromArgs([]float32{x + float32(g.leftSideBearing), y + float32(g.topSideBearing)}))
			gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 4*4, gl.Ptr(g.rec.vertices))
			gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, 4*4, gl.Ptr(&g.rec.vertices[2]))
			gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
			x = x + float32(g.advanceWidth)
		}
	}
}
