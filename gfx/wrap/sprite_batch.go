package wrap

import (
	"github.com/yuin/gopher-lua"

	"github.com/tanema/amore/gfx"
)

func toSpriteBatch(ls *lua.LState, offset int) *gfx.SpriteBatch {
	img := ls.CheckUserData(offset)
	if v, ok := img.Value.(*gfx.SpriteBatch); ok {
		return v
	}
	ls.ArgError(offset, "sprite batch expected")
	return nil
}

func gfxNewSpriteBatch(ls *lua.LState) int {
	return returnUD(
		ls,
		"SpriteBatch",
		gfx.NewSpriteBatch(toTexture(ls, 1), toIntD(ls, 2, 1000), toUsage(ls, 3)),
	)
}

func gfxSpriteBatchAdd(ls *lua.LState) int {
	toSpriteBatch(ls, 1).Add(extractFloatArray(ls, 2)...)
	return 0
}

func gfxSpriteBatchAddq(ls *lua.LState) int {
	toSpriteBatch(ls, 1).Addq(toQuad(ls, 2), extractFloatArray(ls, 3)...)
	return 0
}

func gfxSpriteBatchSet(ls *lua.LState) int {
	toSpriteBatch(ls, 1).Set(toInt(ls, 2), extractFloatArray(ls, 3)...)
	return 0
}

func gfxSpriteBatchSetq(ls *lua.LState) int {
	toSpriteBatch(ls, 1).Setq(toInt(ls, 2), toQuad(ls, 3), extractFloatArray(ls, 4)...)
	return 0
}

func gfxSpriteBatchClear(ls *lua.LState) int {
	toSpriteBatch(ls, 1).Clear()
	return 0
}

func gfxSpriteBatchSetTexture(ls *lua.LState) int {
	toSpriteBatch(ls, 1).SetTexture(toTexture(ls, 2))
	return 0
}

func gfxSpriteBatchGetTexture(ls *lua.LState) int {
	return returnUD(ls, "Image", toSpriteBatch(ls, 1).GetTexture())
}

func gfxSpriteBatchSetColor(ls *lua.LState) int {
	batch := toSpriteBatch(ls, 1)
	if len(extractFloatArray(ls, 2)) == 0 {
		batch.ClearColor()
	} else {
		r, g, b, a := extractColor(ls, 2)
		batch.SetColor(r, g, b, a)
	}
	return 0
}

func gfxSpriteBatchGetColor(ls *lua.LState) int {
	batch := toSpriteBatch(ls, 1)
	for _, x := range batch.GetColor() {
		ls.Push(lua.LNumber(x))
	}
	return 4
}

func gfxSpriteBatchGetCount(ls *lua.LState) int {
	ls.Push(lua.LNumber(toSpriteBatch(ls, 1).GetCount()))
	return 1
}

func gfxSpriteBatchSetBufferSize(ls *lua.LState) int {
	toSpriteBatch(ls, 1).SetBufferSize(toInt(ls, 2))
	return 0
}

func gfxSpriteBatchGetBufferSize(ls *lua.LState) int {
	ls.Push(lua.LNumber(toSpriteBatch(ls, 1).GetBufferSize()))
	return 1
}

func gfxSpriteBatchSetDrawRange(ls *lua.LState) int {
	toSpriteBatch(ls, 1).SetDrawRange(toIntD(ls, 2, -1), toIntD(ls, 3, -1))
	return 0
}

func gfxSpriteBatchGetDrawRange(ls *lua.LState) int {
	min, max := toSpriteBatch(ls, 1).GetDrawRange()
	ls.Push(lua.LNumber(min))
	ls.Push(lua.LNumber(max))
	return 2
}

func gfxSpriteBatchDraw(ls *lua.LState) int {
	toSpriteBatch(ls, 1).Draw(extractFloatArray(ls, 2)...)
	return 0
}
