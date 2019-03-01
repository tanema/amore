package wrap

import (
	"fmt"

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

func gfxNewShader(ls *lua.LState) int {
	return returnUD(ls, "Shader", gfx.NewShader(toString(ls, 1), toStringD(ls, 2, "")))
}

func gfxShaderSend(ls *lua.LState) int {
	program := toShader(ls, 1)
	name := toString(ls, 2)
	uniformType, found := program.GetUniformType(name)
	if !found {
		ls.ArgError(2, fmt.Sprintf("unknown uniform with name [%s]", name))
	}
	switch uniformType {
	case gfx.UniformFloat:
		program.SendFloat(name, extractFloatArray(ls, 3)...)
	case gfx.UniformInt:
		program.SendInt(name, extractIntArray(ls, 3)...)
	case gfx.UniformSampler:
		program.SendTexture(name, toTexture(ls, 3))
	}
	return 0
}
