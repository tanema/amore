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
