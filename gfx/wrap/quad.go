package wrap

import (
	"github.com/yuin/gopher-lua"

	"github.com/tanema/amore/gfx"
)

func toQuad(ls *lua.LState, offset int) *gfx.Quad {
	img := ls.CheckUserData(offset)
	if v, ok := img.Value.(*gfx.Quad); ok {
		return v
	}
	ls.ArgError(offset, "quad expected")
	return nil
}

func gfxNewQuad(ls *lua.LState) int {
	offset := 1
	args := []int32{}
	for x := ls.Get(offset); x != nil; offset++ {
		val := ls.Get(offset)
		if lv, ok := val.(lua.LNumber); ok {
			args = append(args, int32(lv))
		} else if val.Type() == lua.LTNil {
			break
		} else {
			ls.ArgError(offset, "argument wrong type, should be number")
		}
	}
	if offset < 6 {
		ls.ArgError(len(args)-1, "not enough arguments")
	}
	return returnUD(ls, "Quad", gfx.NewQuad(args[0], args[1], args[2], args[3], args[4], args[5]))
}

func gfxQuadGetWidth(ls *lua.LState) int {
	ls.Push(lua.LNumber(toQuad(ls, 1).GetWidth()))
	return 1
}

func gfxQuadGetHeight(ls *lua.LState) int {
	ls.Push(lua.LNumber(toQuad(ls, 1).GetHeight()))
	return 1
}

func gfxQuadGetViewport(ls *lua.LState) int {
	x, y, w, h := toQuad(ls, 1).GetViewport()
	ls.Push(lua.LNumber(x))
	ls.Push(lua.LNumber(y))
	ls.Push(lua.LNumber(w))
	ls.Push(lua.LNumber(h))
	return 4
}

func gfxQuadSetViewport(ls *lua.LState) int {
	toQuad(ls, 1).SetViewport(int32(toInt(ls, 2)), int32(toInt(ls, 3)), int32(toInt(ls, 4)), int32(toInt(ls, 5)))
	return 0
}
