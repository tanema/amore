package gfx

import (
	"image"
	"image/draw"
	"math"
	"unicode"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type (
	rasterizer struct {
		face       font.Face
		atlasImg   *image.RGBA
		texture    *Texture
		mapping    map[rune]glyphData
		lineHeight float32
		ascent     float32
		descent    float32
	}
	glyphData struct {
		quad    *Quad
		advance float32
		decent  float32
		lsb     float32
	}
)

const glyphPadding int = 2

func newRasterizer(face font.Face, runeSets ...[]rune) *rasterizer {
	runes := uniqRunesForSets(runeSets...)
	imageWidth, imageHeight, advance, height := calcSquareMapping(face, runes)
	atlasImg := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	mapping := make(map[rune]glyphData)
	var gx, gy int
	for _, r := range runes {
		rect, srcImg, srcPoint, adv, ok := face.Glyph(fixed.P(gx, gy+height), r)
		if !ok {
			continue
		}
		draw.Draw(atlasImg, rect, srcImg, srcPoint, draw.Src)
		mapping[r] = glyphData{
			decent: float32(rect.Min.Y - gy),
			lsb:    float32(rect.Min.X - gx),
			quad: NewQuad(
				int32(rect.Min.X), int32(rect.Min.Y),
				int32(rect.Dx()), int32(rect.Dy()),
				int32(imageWidth), int32(imageHeight),
			),
			advance: i2f(adv),
		}

		gx += advance + glyphPadding
		if gx+advance >= imageWidth {
			gx = 0
			gy += height + glyphPadding
		}
	}

	newRast := &rasterizer{
		face:       face,
		atlasImg:   atlasImg,
		mapping:    mapping,
		ascent:     i2f(face.Metrics().Ascent),
		descent:    i2f(face.Metrics().Descent),
		lineHeight: i2f(face.Metrics().Height),
	}
	registerVolatile(newRast)
	return newRast
}

func (rast *rasterizer) loadVolatile() bool {
	var err error
	rast.texture, err = newImageTexture(rast.atlasImg, false)
	return err == nil
}

func (rast *rasterizer) unloadVolatile() {}

func uniqRunesForSets(runeSets ...[]rune) []rune {
	seen := make(map[rune]bool)
	runes := []rune{unicode.ReplacementChar}
	for _, set := range runeSets {
		for _, r := range set {
			if !seen[r] {
				runes = append(runes, r)
				seen[r] = true
			}
		}
	}
	return runes
}

func calcSquareMapping(face font.Face, runes []rune) (imageWidth, imageHeight, advance, height int) {
	var maxAdvance fixed.Int26_6
	for _, r := range runes {
		_, adv, ok := face.GlyphBounds(r)
		if !ok {
			continue
		}
		if adv > maxAdvance {
			maxAdvance = adv
		}
	}
	a := i2f(maxAdvance)
	h := i2f(face.Metrics().Ascent + face.Metrics().Descent)
	squareWidth := float32(math.Ceil(math.Sqrt(float64(len(runes)))))
	imageWidth = int(pow2(uint32(a * squareWidth)))
	imageHeight = int(pow2(uint32(h * squareWidth)))
	return imageWidth, imageHeight, ceil(a), ceil(h)
}

func i2f(i fixed.Int26_6) float32 {
	return float32(i) / (1 << 6)
}

func f2i(f float64) fixed.Int26_6 {
	return fixed.Int26_6(f * (1 << 6))
}

func pow2(x uint32) uint32 {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return x + 1
}

func ceil(x float32) int {
	return int(math.Ceil(float64(x)))
}
