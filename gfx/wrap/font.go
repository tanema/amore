package wrap

import (
	"github.com/yuin/gopher-lua"

	"github.com/tanema/amore/gfx"
)

func toFont(ls *lua.LState, offset int) *gfx.Font {
	img := ls.CheckUserData(offset)
	if v, ok := img.Value.(*gfx.Font); ok {
		return v
	}
	ls.ArgError(offset, "font expected")
	return nil
}

func gfxNewFont(ls *lua.LState) int {
	newFont, err := gfx.NewFont(toString(ls, 1), toFloat(ls, 2))
	if err == nil {
		return returnUD(ls, "Font", newFont)
	}
	ls.Push(lua.LNil)
	return 1
}

func gfxFontGetWidth(ls *lua.LState) int {
	font := toFont(ls, 1)
	ls.Push(lua.LNumber(font.GetWidth(toString(ls, 2))))
	return 1
}

func gfxFontGetHeight(ls *lua.LState) int {
	font := toFont(ls, 1)
	ls.Push(lua.LNumber(font.GetHeight()))
	return 1
}

func gfxFontSetFallback(ls *lua.LState) int {
	font := toFont(ls, 1)
	start := 2
	fallbacks := []*gfx.Font{}
	for x := ls.Get(start); x != nil; start++ {
		img := ls.CheckUserData(start)
		if v, ok := img.Value.(*gfx.Font); ok {
			fallbacks = append(fallbacks, v)
		}
		ls.ArgError(start, "font expected")
	}
	font.SetFallbacks(fallbacks...)
	return 0
}

func gfxFontGetWrap(ls *lua.LState) int {
	font := toFont(ls, 1)
	wrap, strs := font.GetWrap(toString(ls, 1), toFloat(ls, 2))
	ls.Push(lua.LNumber(wrap))
	table := ls.NewTable()
	for _, str := range strs {
		table.Append(lua.LString(str))
	}
	ls.Push(table)
	return 2
}
