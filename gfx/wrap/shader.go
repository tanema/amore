package wrap

import (
	"github.com/yuin/gopher-lua"

	"github.com/tanema/amore/gfx"
)

func toShader(ls *lua.LState, offset int) *gfx.Shader {
	img := ls.CheckUserData(offset)
	if v, ok := img.Value.(*gfx.Shader); ok {
		return v
	}
	ls.ArgError(offset, "shader expected")
	return nil
}
