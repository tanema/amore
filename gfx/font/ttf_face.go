package font

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/tanema/amore/file"
)

// Face just an alias so you don't ahve to import multople packages named font
type Face font.Face

// NewTTFFace will load up a ttf font face for creating a font in graphics
func NewTTFFace(filepath string, size float32) (font.Face, error) {
	fontBytes, err := file.Read(filepath)
	if err != nil {
		return nil, err
	}

	return ttfFromBytes(fontBytes, size)
}

// Bold will return an bold font face with the request font size
func Bold(fontSize float32) (font.Face, error) {
	return ttfFromBytes(gobold.TTF, fontSize)
}

// Default will return an regular font face with the request font size
func Default(fontSize float32) (font.Face, error) {
	return ttfFromBytes(goregular.TTF, fontSize)
}

// Italic will return an italic font face with the request font size
func Italic(fontSize float32) (font.Face, error) {
	return ttfFromBytes(goitalic.TTF, fontSize)
}

func ttfFromBytes(fontBytes []byte, size float32) (font.Face, error) {
	ttf, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	return truetype.NewFace(ttf, &truetype.Options{
		Size:              float64(size),
		GlyphCacheEntries: 1,
	}), nil
}
