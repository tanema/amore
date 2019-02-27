package wrap

import (
	"github.com/yuin/gopher-lua"

	"github.com/tanema/amore/gfx"
)

const luaImageType = "Image"

func toImage(ls *lua.LState, offset int) *gfx.Image {
	img := ls.CheckUserData(offset)
	if v, ok := img.Value.(*gfx.Image); ok {
		if v.Texture == nil {
			ls.ArgError(offset, "image not loaded")
		}
		return v
	}
	ls.ArgError(offset, "image expected")
	return nil
}

func toCanvas(ls *lua.LState, offset int) *gfx.Canvas {
	canvas := ls.CheckUserData(offset)
	if v, ok := canvas.Value.(*gfx.Canvas); ok {
		return v
	}
	ls.ArgError(offset, "canvas expected")
	return nil
}

func toTexture(ls *lua.LState, offset int) *gfx.Texture {
	text := ls.CheckUserData(offset)
	if v, ok := text.Value.(*gfx.Canvas); ok {
		return v.Texture
	} else if v, ok := text.Value.(*gfx.Image); ok {
		return v.Texture
	}
	ls.ArgError(offset, "texture expected")
	return nil
}

func gfxNewImage(ls *lua.LState) int {
	return returnUD(ls, "Image", gfx.NewImage(toString(ls, 1), ls.ToBool(2)))
}

func gfxNewCanvas(ls *lua.LState) int {
	w, h := gfx.GetDimensions()
	cw, ch := toIntD(ls, 1, int(w)), toIntD(ls, 2, int(h))
	return returnUD(ls, "Canvas", gfx.NewCanvas(int32(cw), int32(ch)))
}

func gfxCanvasNewImage(ls *lua.LState) int {
	canvas := toCanvas(ls, 1)
	cw, ch := canvas.GetDimensions()
	x, y, w, h := toIntD(ls, 2, 0), toIntD(ls, 3, 0), toIntD(ls, 4, int(cw)), toIntD(ls, 5, int(ch))
	img, err := canvas.NewImageData(int32(x), int32(y), int32(w), int32(h))
	if err == nil {
		return returnUD(ls, "Image", img)
	}
	ls.Push(lua.LNil)
	return 1
}

func gfxTextureDraw(ls *lua.LState) int {
	toTexture(ls, 1).Draw(extractFloatArray(ls, 2)...)
	return 0
}

func gfxTextureDrawq(ls *lua.LState) int {
	toTexture(ls, 1).Drawq(toQuad(ls, 2), extractFloatArray(ls, 3)...)
	return 0
}

func gfxTextureSetWrap(ls *lua.LState) int {
	toTexture(ls, 1).SetWrap(toWrap(ls, 2), toWrap(ls, 3))
	return 0
}

func gfxTextureSetFilter(ls *lua.LState) int {
	toTexture(ls, 1).SetFilter(toFilter(ls, 2), toFilter(ls, 3))
	return 0
}

func gfxTextureGetWidth(ls *lua.LState) int {
	ls.Push(lua.LNumber(int(toTexture(ls, 1).Width)))
	return 1
}

func gfxTextureGetHeight(ls *lua.LState) int {
	ls.Push(lua.LNumber(int(toTexture(ls, 1).Height)))
	return 1
}

func gfxTextureGetDimensions(ls *lua.LState) int {
	text := toTexture(ls, 1)
	ls.Push(lua.LNumber(int(text.Width)))
	ls.Push(lua.LNumber(int(text.Height)))
	return 2
}
