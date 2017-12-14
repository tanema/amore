package gfx

import (
	"image"
	"image/draw"
	// loading all image libs for loading
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"

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
		glyphHints string
	}
	ttfFontRasterizer struct {
		rasterizerBase
		fontSize float32
		ttf      *truetype.Font
		context  *freetype.Context
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

func newTtfRasterizer(filename string, fontSize float32) *ttfFontRasterizer {
	rasterizer := &ttfFontRasterizer{rasterizerBase: rasterizerBase{filepath: filename}, fontSize: fontSize}
	registerVolatile(rasterizer)
	return rasterizer
}

func newImageRasterizer(filename, glyphHints string) *imageFontRasterizer {
	rasterizer := &imageFontRasterizer{rasterizerBase: rasterizerBase{filepath: filename}, glyphHints: glyphHints}
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
	rast.context.SetFontSize(float64(rast.fontSize))

	fontBounds := rast.ttf.Bounds()
	rast.advance = rast.context.FUnitToPixelRU(int(fontBounds.XMax - fontBounds.XMin))
	rast.height = rast.context.FUnitToPixelRU(int(fontBounds.YMax-fontBounds.YMin) + 5)
	imageWidth := pow2(uint32(int32(rast.advance) * glyphsPerRow))
	imageHeight := pow2(uint32(int32(rast.height) * glyphsPerCol))

	rgba := image.NewRGBA(image.Rect(0, 0, int(imageWidth), int(imageHeight)))
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
			rec:             NewQuad(int32(gx), int32(gy), int32(rast.advance), int32(rast.height), int32(imageWidth), int32(imageHeight)),
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
		rast.ascent = int(math.Max(float64(rast.ascent), float64(g.ascent)))
		rast.descent = int(math.Max(float64(rast.descent), float64(g.descent)))
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
	glyphRuneHints := []rune(rast.glyphHints)

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
	glyphsPerCol := (int32(len(glyphRuneHints)) / glyphsPerRow) + 1
	rast.advance = img.Bounds().Dx() / len(glyphRuneHints)
	rast.height = img.Bounds().Dy()
	rast.ascent = rast.height
	rast.descent = 0
	imageWidth := pow2(uint32(int32(rast.advance) * glyphsPerRow))
	imageHeight := pow2(uint32(int32(rast.height) * glyphsPerCol))

	rgba := image.NewRGBA(image.Rect(0, 0, int(imageWidth), int(imageHeight)))

	var gx, gy int
	for i := 0; i < len(glyphRuneHints); i++ {
		dstRec := image.Rect(gx, gy, gx+rast.advance, gy+rast.height)
		draw.Draw(rgba, dstRec, img, image.Pt(i*rast.advance, 0), draw.Src)

		rast.glyphs[glyphRuneHints[i]] = glyphData{
			rec:             NewQuad(int32(gx), int32(gy), int32(rast.advance), int32(rast.height), int32(imageWidth), int32(imageHeight)),
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
