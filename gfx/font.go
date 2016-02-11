package gfx

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/go-gl/gl/v2.1/gl"

	"github.com/tanema/amore/file"
	"github.com/tanema/freetype-go/freetype"
)

type (
	Font interface {
		Printf(x, y float32, fs string, argv ...interface{})
	}
	FontBase struct {
		*Texture
		filepath string
		Offset   int
		Glyphs   map[rune]Glyph
	}
	ImageFont struct {
		FontBase
		glyph_hints string
	}
	TTFont struct {
		FontBase
		font_size float32
	}

	textRec struct {
		X1, Y1, X2, Y2 float32
	}
	Glyph struct {
		TextureRec    textRec
		Width, Height int
		Advance       int
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

func NewFont(filename string, font_size float32) *TTFont {
	new_font := &TTFont{FontBase: FontBase{filepath: filename}, font_size: font_size}
	registerVolatile(new_font)
	return new_font
}

func NewImageFont(filename, glyph_hints string) *ImageFont {
	new_font := &ImageFont{FontBase: FontBase{filepath: filename}, glyph_hints: glyph_hints}
	registerVolatile(new_font)
	return new_font
}

func (font *TTFont) loadVolatile() bool {
	fontBytes, err := file.Read(font.filepath)
	ttf, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return false
	}

	glyphs := ttf.ListRunes()
	font.Glyphs = make(map[rune]Glyph)
	glyphsPerRow := int32(16)
	glyphsPerCol := (int32(len(glyphs)) / glyphsPerRow) + 1

	context := freetype.NewContext()
	context.SetDPI(72)
	context.SetFont(ttf)
	context.SetFontSize(float64(font.font_size))

	font_bounds := ttf.Bounds()
	glyph_width := context.FUnitToPixelRU(int(font_bounds.XMax - font_bounds.XMin))
	glyph_height := context.FUnitToPixelRU(int(font_bounds.YMax-font_bounds.YMin) + 5)
	image_width := pow2(uint32(int32(glyph_width) * glyphsPerRow))
	image_height := pow2(uint32(int32(glyph_height) * glyphsPerCol))

	rgba := image.NewRGBA(image.Rect(0, 0, int(image_width), int(image_height)))
	context.SetClip(rgba.Bounds())
	context.SetDst(rgba)
	context.SetSrc(image.White)

	var gx, gy int
	font.Offset = context.FUnitToPixelRU(ttf.UnitsPerEm())
	for i, ch := range glyphs {
		pt := freetype.Pt(gx+font.Offset, gy+font.Offset)
		context.DrawString(string(ch), pt)

		tx1 := float32(gx) / float32(image_width)
		ty1 := float32(gy) / float32(image_height)
		tx2 := (float32(gx) + float32(glyph_width)) / float32(image_width)
		ty2 := (float32(gy) + float32(glyph_height)) / float32(image_height)

		index := ttf.Index(ch)
		metric := ttf.HMetric(index)
		font.Glyphs[ch] = Glyph{
			TextureRec: textRec{tx1, ty1, tx2, ty2},
			Width:      glyph_width,
			Height:     glyph_height,
			Advance:    int(context.FUnitToPixelRU(int(metric.AdvanceWidth))),
		}

		if i%16 == 0 {
			gx = 0
			gy += glyph_height
		} else {
			gx += glyph_width
		}
	}

	font.Texture, err = newImageTexture(rgba, false)
	return err == nil
}

func (font *ImageFont) loadVolatile() bool {
	glyph_rune_hints := []rune(font.glyph_hints)

	imgFile, err := file.NewFile(font.filepath)
	if err != nil {
		return false
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return false
	}

	font.Glyphs = make(map[rune]Glyph)
	glyphsPerRow := int32(16)
	glyphsPerCol := (int32(len(glyph_rune_hints)) / glyphsPerRow) + 1
	glyph_width := img.Bounds().Dx() / len(glyph_rune_hints)
	glyph_height := img.Bounds().Dy()
	image_width := pow2(uint32(int32(glyph_width) * glyphsPerRow))
	image_height := pow2(uint32(int32(glyph_height) * glyphsPerCol))

	rgba := image.NewRGBA(image.Rect(0, 0, int(image_width), int(image_height)))

	var gx, gy int
	for i := 0; i < len(glyph_rune_hints); i++ {
		dst_rec := image.Rect(gx, gy, gx+glyph_width, gy+glyph_height)
		draw.Draw(rgba, dst_rec, img, image.Pt(i*glyph_width, 0), draw.Src)

		tx1 := float32(gx) / float32(image_width)
		ty1 := float32(gy) / float32(image_height)
		tx2 := (float32(gx) + float32(glyph_width)) / float32(image_width)
		ty2 := (float32(gy) + float32(glyph_height)) / float32(image_height)

		font.Glyphs[glyph_rune_hints[i]] = Glyph{
			TextureRec: textRec{tx1, ty1, tx2, ty2},
			Width:      glyph_width,
			Height:     glyph_height,
			Advance:    glyph_width,
		}

		if i%16 == 0 {
			gx = 0
			gy += glyph_height
		} else {
			gx += glyph_width
		}
	}

	font.Texture, err = newImageTexture(rgba, false)
	return err == nil
}

func (f *FontBase) Index(ch rune) {
}

func (f *FontBase) Release() {
	f.unloadVolatile()
}

func Printf(x, y float32, fs string, argv ...interface{}) {
	current_font := states.back().font
	if current_font == nil {
		return
	}
	current_font.Printf(x, y, fs, argv...)
}

func (font *FontBase) Printf(x, y float32, fs string, argv ...interface{}) {
	formatted_string := fmt.Sprintf(fs, argv...)
	if len(formatted_string) == 0 {
		return
	}

	x = x - float32(font.Offset)
	y = y - (float32(font.Offset) / 4.0)

	prepareDraw(nil)
	bindTexture(font.GetHandle())
	useVertexAttribArrays(ATTRIBFLAG_POS | ATTRIBFLAG_TEXCOORD)

	for _, ch := range formatted_string {
		if glyph, ok := font.Glyphs[ch]; ok {
			gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, 0, gl.Ptr([]float32{
				x, y,
				x, y + float32(glyph.Height),
				x + float32(glyph.Width), y + float32(glyph.Height),
				x + float32(glyph.Width), y,
			}))
			gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, 0, gl.Ptr([]float32{
				glyph.TextureRec.X1, glyph.TextureRec.Y1,
				glyph.TextureRec.X1, glyph.TextureRec.Y2,
				glyph.TextureRec.X2, glyph.TextureRec.Y2,
				glyph.TextureRec.X2, glyph.TextureRec.Y1,
			}))
			gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

			x = x + float32(glyph.Advance)
		}
	}

}
