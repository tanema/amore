package font

import (
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/tanema/amore/file"
)

// BitmapFace holds data for a bitmap font face to satisfy the font.Face interface
type BitmapFace struct {
	file.File
	img     image.Image
	glyphs  map[rune]glyphData
	advance fixed.Int26_6
	metrics font.Metrics
}

type glyphData struct {
	rect image.Rectangle
	pt   image.Point
}

// NewBitmapFace will load up an image font face for creating a font in graphics
func NewBitmapFace(filepath, glyphHints string) (font.Face, error) {
	imgFile, err := file.NewFile(filepath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	glyphRuneHints := []rune(glyphHints)
	advance := img.Bounds().Dx() / len(glyphRuneHints)
	newFace := BitmapFace{
		File:    imgFile,
		img:     img,
		glyphs:  make(map[rune]glyphData),
		advance: fixed.I(advance),
		metrics: font.Metrics{
			Height:  fixed.I(img.Bounds().Dy()),
			Ascent:  fixed.I(img.Bounds().Dy()),
			Descent: 0,
		},
	}

	var gx int
	for i, r := range glyphRuneHints {
		newFace.glyphs[r] = glyphData{
			rect: image.Rect(gx, 0, gx+advance, newFace.metrics.Height.Ceil()),
			pt:   image.Pt(i*advance, 0),
		}
		gx += advance
	}

	return newFace, err
}

// Glyph returns the draw.DrawMask parameters (dr, mask, maskp) to draw r's
// glyph at the sub-pixel destination location dot, and that glyph's
// advance width.
//
// It returns !ok if the face does not contain a glyph for r.
//
// The contents of the mask image returned by one Glyph call may change
// after the next Glyph call. Callers that want to cache the mask must make
// a copy.
func (face BitmapFace) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	var glyph glyphData
	glyph, ok = face.glyphs[r]
	if !ok {
		return
	}
	dr.Min = image.Point{
		X: dot.X.Floor(),
		Y: dot.Y.Floor() - face.metrics.Height.Floor(),
	}
	dr.Max = image.Point{
		X: dr.Min.X + face.advance.Floor(),
		Y: dr.Min.Y + face.metrics.Height.Floor(),
	}
	return dr, face.img, glyph.pt, face.advance, ok
}

// GlyphBounds returns the bounding box of r's glyph, drawn at a dot equal
// to the origin, and that glyph's advance width.
//
// It returns !ok if the face does not contain a glyph for r.
//
// The glyph's ascent and descent equal -bounds.Min.Y and +bounds.Max.Y. A
// visual depiction of what these metrics are is at
// https://developer.apple.com/library/mac/documentation/TextFonts/Conceptual/CocoaTextArchitecture/Art/glyph_metrics_2x.png
func (face BitmapFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	_, ok = face.glyphs[r]
	return fixed.R(0, 0, face.advance.Ceil(), face.metrics.Height.Ceil()), face.advance, ok
}

// GlyphAdvance returns the advance width of r's glyph.
//
// It returns !ok if the face does not contain a glyph for r.
func (face BitmapFace) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	return face.advance, true
}

// Kern returns the horizontal adjustment for the kerning pair (r0, r1). A
// positive kern means to move the glyphs further apart.
func (face BitmapFace) Kern(r0, r1 rune) fixed.Int26_6 {
	return 0.0
}

// Metrics returns the metrics for this Face.
func (face BitmapFace) Metrics() font.Metrics {
	return face.metrics
}
