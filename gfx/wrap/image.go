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

func gfxNewImage(ls *lua.LState) int {
	ls.Push(toUD(ls, "Image", gfx.NewImage(toString(ls, 1), ls.ToBool(2))))
	return 1
}

func gfxImageDraw(ls *lua.LState) int {
	image := toImage(ls, 1)
	drawArgs := extractFloatArray(ls, 2)
	image.Draw(drawArgs...)
	return 0
}

func gfxImageGetWidth(ls *lua.LState) int {
	image := toImage(ls, 1)
	ls.Push(lua.LNumber(int(image.Width)))
	return 1
}

func gfxImageGetHeight(ls *lua.LState) int {
	image := toImage(ls, 1)
	ls.Push(lua.LNumber(int(image.Height)))
	return 1
}

func gfxImageGetDimensions(ls *lua.LState) int {
	image := toImage(ls, 1)
	ls.Push(lua.LNumber(int(image.Width)))
	ls.Push(lua.LNumber(int(image.Height)))
	return 2
}
