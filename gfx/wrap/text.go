package wrap

import (
	"github.com/yuin/gopher-lua"

	"github.com/tanema/amore/gfx"
)

func extractPrintable(ls *lua.LState, offset int) ([]string, [][]float32) {
	strs := []string{}
	colors := [][]float32{}

	val1 := ls.Get(offset)
	if val1.Type() == lua.LTTable {
		table := val1.(*lua.LTable)
		rawstrs, strok := (table.RawGetInt(1)).(*lua.LTable)
		rawclrs, clrok := table.RawGetInt(2).(*lua.LTable)
		if !strok || !clrok {
			ls.ArgError(offset, "unexpected argument type")
		}

		rawstrs.ForEach(func(index lua.LValue, str lua.LValue) {
			strs = append(strs, str.String())
		})

		rawclrs.ForEach(func(index lua.LValue, color lua.LValue) {
			setColor := []float32{}
			(color.(*lua.LTable)).ForEach(func(index lua.LValue, color lua.LValue) {
				setColor = append(setColor, float32(color.(lua.LNumber)))
			})
			if len(setColor) == 3 {
				setColor = append(setColor, 1)
			} else if len(setColor) < 4 {
				ls.ArgError(offset, "not enough values for a color")
			}
			colors = append(colors, setColor)
		})
	} else if val1.Type() == lua.LTString {
		return []string{val1.String()}, [][]float32{gfx.GetColor()}
	} else {
		ls.ArgError(offset, "unexpected argument type")
	}

	return strs, colors
}

func toText(ls *lua.LState, offset int) *gfx.Text {
	img := ls.CheckUserData(offset)
	if v, ok := img.Value.(*gfx.Text); ok {
		return v
	}
	ls.ArgError(offset, "text expected")
	return nil
}

func gfxPrint(ls *lua.LState) int {
	str, clrs := extractPrintable(ls, 1)
	args := extractFloatArray(ls, 2)
	gfx.Print(str, clrs, args...)
	return 0
}

func gfxPrintf(ls *lua.LState) int {
	str, clrs := extractPrintable(ls, 1)
	wrap := toFloat(ls, 2)
	align := toString(ls, 3)
	args := extractFloatArray(ls, 4)
	gfx.Printf(str, clrs, wrap, align, args...)
	return 0
}

func gfxNewText(ls *lua.LState) int {
	str, clrs := extractPrintable(ls, 2)
	text := gfx.NewText(toFont(ls, 1), str, clrs, toFloatD(ls, 3, -1), toStringD(ls, 4, "left"))
	ls.Push(toUD(ls, "Text", text))
	return 1
}

func gfxTextDraw(ls *lua.LState) int {
	txt := toText(ls, 1)
	txt.Draw(extractFloatArray(ls, 2)...)
	return 0
}

func gfxTextSet(ls *lua.LState) int {
	txt := toText(ls, 1)
	str, clrs := extractPrintable(ls, 2)
	txt.Set(str, clrs)
	return 0
}

func gfxTextGetWidth(ls *lua.LState) int {
	txt := toText(ls, 1)
	ls.Push(lua.LNumber(txt.GetWidth()))
	return 1
}

func gfxTextGetHeight(ls *lua.LState) int {
	txt := toText(ls, 1)
	ls.Push(lua.LNumber(txt.GetHeight()))
	return 1
}

func gfxTextGetDimensions(ls *lua.LState) int {
	txt := toText(ls, 1)
	w, h := txt.GetDimensions()
	ls.Push(lua.LNumber(w))
	ls.Push(lua.LNumber(h))
	return 2
}

func gfxTextGetFont(ls *lua.LState) int {
	txt := toText(ls, 1)
	ls.Push(toUD(ls, "Font", txt.GetFont()))
	return 1
}

func gfxTextSetFont(ls *lua.LState) int {
	txt := toText(ls, 1)
	txt.SetFont(toFont(ls, 2))
	return 0
}
