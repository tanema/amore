package gfx

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/tanema/amore/file"
	"github.com/tanema/freetype-go/freetype"
	"github.com/tanema/freetype-go/freetype/truetype"
)

type (
	rasterizer interface {
		getTexture() *Texture
		getOffset() int
		getHeight() int
		getAdvance() int
		getAscent() int
		getDescent() int
		getLineHeight() int
		getGlyphData(g rune) glyphData
		getGlyphCount() int
		hasGlyph(g rune) bool
		hasGlyphs(text string) bool
		getKerning(leftglyph, rightglyph rune) float32
		Release()
	}
	rasterizerBase struct {
		*Texture
		filepath string
		glyphs   map[rune]glyphData
		advance  int
		height   int
		ascent   int
		descent  int
		offset   int
	}
	imageFontRasterizer struct {
		rasterizerBase
		glyph_hints string
	}
	ttfFontRasterizer struct {
		rasterizerBase
		font_size float32
		ttf       *truetype.Font
		context   *freetype.Context
	}
	glyphData struct {
		rec             *Quad
		leftSideBearing int
		advanceWidth    int
		topSideBearing  int
		advanceHeight   int
		ascent          int
		descent         int
	}
)

func pow2(x uint32) uint32 {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return x + 1
}

func newTtfRasterizer(filename string, font_size float32) *ttfFontRasterizer {
	rasterizer := &ttfFontRasterizer{rasterizerBase: rasterizerBase{filepath: filename}, font_size: font_size}
	registerVolatile(rasterizer)
	return rasterizer
}

func newImageRasterizer(filename, glyph_hints string) *imageFontRasterizer {
	rasterizer := &imageFontRasterizer{rasterizerBase: rasterizerBase{filepath: filename}, glyph_hints: glyph_hints}
	registerVolatile(rasterizer)
	return rasterizer
}

func (rast *rasterizerBase) getTexture() *Texture {
	return rast.Texture
}

func (rast *rasterizerBase) getOffset() int {
	return rast.offset
}

func (rast *rasterizerBase) getHeight() int {
	return rast.height
}

func (rast *rasterizerBase) getAdvance() int {
	return rast.advance
}

func (rast *rasterizerBase) getAscent() int {
	return rast.ascent
}

func (rast *rasterizerBase) getDescent() int {
	return rast.descent
}

func (rast *rasterizerBase) getGlyphData(g rune) glyphData {
	return rast.glyphs[g]
}

func (rast *rasterizerBase) getGlyphCount() int {
	return len(rast.glyphs)
}

func (rast *rasterizerBase) hasGlyph(g rune) bool {
	_, ok := rast.glyphs[g]
	return ok
}

func (rast *rasterizerBase) hasGlyphs(text string) bool {
	for _, c := range text {
		if rast.hasGlyph(c) {
			return false
		}
	}
	return true
}

func (rast *ttfFontRasterizer) loadVolatile() bool {
	fontBytes, err := file.Read(rast.filepath)
	rast.ttf, err = freetype.ParseFont(fontBytes)
	if err != nil {
		return false
	}

	glyphs := rast.ttf.ListRunes()
	rast.glyphs = make(map[rune]glyphData)
	glyphsPerRow := int32(16)
	glyphsPerCol := (int32(len(glyphs)) / glyphsPerRow) + 1

	rast.context = freetype.NewContext()
	rast.context.SetDPI(72)
	rast.context.SetFont(rast.ttf)
	rast.context.SetFontSize(float64(rast.font_size))

	font_bounds := rast.ttf.Bounds()
	rast.advance = rast.context.FUnitToPixelRU(int(font_bounds.XMax - font_bounds.XMin))
	rast.height = rast.context.FUnitToPixelRU(int(font_bounds.YMax-font_bounds.YMin) + 5)
	image_width := pow2(uint32(int32(rast.advance) * glyphsPerRow))
	image_height := pow2(uint32(int32(rast.height) * glyphsPerCol))

	rgba := image.NewRGBA(image.Rect(0, 0, int(image_width), int(image_height)))
	rast.context.SetClip(rgba.Bounds())
	rast.context.SetDst(rgba)
	rast.context.SetSrc(image.White)

	var gx, gy int
	rast.offset = rast.context.FUnitToPixelRU(rast.ttf.UnitsPerEm())
	for i, ch := range glyphs {
		pt := freetype.Pt(gx+rast.offset, gy+rast.offset)
		rast.context.DrawString(string(ch), pt)
		metric := rast.ttf.HMetric(rast.ttf.Index(ch))
		vmetric := rast.ttf.VMetric(rast.ttf.Index(ch))

		lsb := rast.context.FUnitToPixelRU(int(metric.LeftSideBearing))
		aw := rast.context.FUnitToPixelRU(int(metric.AdvanceWidth))
		tsb := rast.context.FUnitToPixelRU(int(vmetric.TopSideBearing))
		ah := rast.context.FUnitToPixelRU(int(vmetric.AdvanceHeight))
		descent := rast.height - ah
		ascent := rast.height - descent

		rast.glyphs[ch] = glyphData{
			rec:             NewQuad(int32(gx), int32(gy), int32(rast.advance), int32(rast.height), int32(image_width), int32(image_height)),
			leftSideBearing: lsb,
			advanceWidth:    aw,
			topSideBearing:  tsb,
			advanceHeight:   ah,
			ascent:          ascent,
			descent:         descent,
		}

		if i%16 == 0 {
			gx = 0
			gy += rast.height
		} else {
			gx += rast.advance
		}
	}

	for _, g := range rast.glyphs {
		rast.ascent = Maxi(rast.ascent, g.ascent)
		rast.descent = Maxi(rast.descent, g.descent)
	}

	rast.Texture, err = newImageTexture(rgba, false)
	return err == nil
}

func (rast *ttfFontRasterizer) getLineHeight() int {
	return int(float32(rast.getHeight()) * 1.25)
}

func (rast *ttfFontRasterizer) getKerning(leftglyph, rightglyph rune) float32 {
	return float32(rast.context.FUnitToPixelRU(int(rast.ttf.Kerning(rast.ttf.Index(leftglyph), rast.ttf.Index(rightglyph)))))
}

func (rast *imageFontRasterizer) loadVolatile() bool {
	glyph_rune_hints := []rune(rast.glyph_hints)

	imgFile, err := file.NewFile(rast.filepath)
	if err != nil {
		return false
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return false
	}

	rast.glyphs = make(map[rune]glyphData)
	glyphsPerRow := int32(16)
	glyphsPerCol := (int32(len(glyph_rune_hints)) / glyphsPerRow) + 1
	rast.advance = img.Bounds().Dx() / len(glyph_rune_hints)
	rast.height = img.Bounds().Dy()
	rast.ascent = rast.height
	rast.descent = 0
	image_width := pow2(uint32(int32(rast.advance) * glyphsPerRow))
	image_height := pow2(uint32(int32(rast.height) * glyphsPerCol))

	rgba := image.NewRGBA(image.Rect(0, 0, int(image_width), int(image_height)))

	var gx, gy int
	for i := 0; i < len(glyph_rune_hints); i++ {
		dst_rec := image.Rect(gx, gy, gx+rast.advance, gy+rast.height)
		draw.Draw(rgba, dst_rec, img, image.Pt(i*rast.advance, 0), draw.Src)

		rast.glyphs[glyph_rune_hints[i]] = glyphData{
			rec:             NewQuad(int32(gx), int32(gy), int32(rast.advance), int32(rast.height), int32(image_width), int32(image_height)),
			leftSideBearing: 0,
			advanceWidth:    rast.advance,
			topSideBearing:  0,
			advanceHeight:   rast.height + 5,
			ascent:          rast.height,
			descent:         0,
		}

		if i%16 == 0 {
			gx = 0
			gy += rast.height
		} else {
			gx += rast.advance
		}
	}

	rast.Texture, err = newImageTexture(rgba, false)
	return err == nil
}

func (rast *imageFontRasterizer) getLineHeight() int {
	return rast.getHeight()
}

func (rast *imageFontRasterizer) getKerning(leftglyph, rightglyph rune) float32 {
	return 0.0
}
