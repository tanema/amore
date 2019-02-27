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
