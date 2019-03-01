package wrap

import (
	"github.com/yuin/gopher-lua"

	"github.com/tanema/amore/gfx"
)

func extractCoords(ls *lua.LState, offset int) []float32 {
	coords := extractFloatArray(ls, offset)
	if len(coords)%2 != 0 {
		ls.ArgError(0, "coordinates need to be given in pairs, arguments are not even")
	}
	return coords
}

func toFloat(ls *lua.LState, offset int) float32 {
	val := ls.Get(offset)
	lv, ok := val.(lua.LNumber)
	if !ok {
		ls.ArgError(offset, "invalid required argument")
		return 0
	}
	return float32(lv)
}

func toFloatD(ls *lua.LState, offset int, fallback float32) float32 {
	val := ls.Get(offset)
	if lv, ok := val.(lua.LNumber); ok {
		return float32(lv)
	}
	return fallback
}

func toInt(ls *lua.LState, offset int) int {
	val := ls.Get(offset)
	lv, ok := val.(lua.LNumber)
	if !ok {
		ls.ArgError(offset, "invalid required argument")
		return 0
	}
	return int(lv)
}

func toIntD(ls *lua.LState, offset int, fallback int) int {
	val := ls.Get(offset)
	if lv, ok := val.(lua.LNumber); ok {
		return int(lv)
	}
	return fallback
}

func toString(ls *lua.LState, offset int) string {
	val := ls.Get(offset)
	lv, ok := val.(lua.LString)
	if !ok {
		ls.ArgError(offset, "invalid required argument")
		return ""
	}
	return string(lv)
}

func toStringD(ls *lua.LState, offset int, fallback string) string {
	val := ls.Get(offset)
	if lv, ok := val.(lua.LString); ok {
		return string(lv)
	}
	return fallback
}

func extractMode(ls *lua.LState, offset int) string {
	mode := ls.ToString(offset)
	if mode == "" || (mode != "fill" && mode != "line") {
		ls.ArgError(offset, "invalid drawmode")
	}
	return mode
}

func extractBlendmode(ls *lua.LState, offset int) string {
	mode := ls.ToString(offset)
	if mode == "" || (mode != "multiplicative" && mode != "premultiplied" &&
		mode != "subtractive" && mode != "additive" && mode != "screen" && mode != "replace" && mode != "alpha") {
		ls.ArgError(offset, "invalid blendmode")
	}
	return mode
}

func extractLineJoin(ls *lua.LState, offset int) string {
	join := ls.ToString(offset)
	if join == "" || (join != "bevel" && join != "miter") {
		ls.ArgError(offset, "invalid drawmode")
	}
	return join
}

func extractFloatArray(ls *lua.LState, offset int) []float32 {
	args := []float32{}
	for x := ls.Get(offset); x != nil; offset++ {
		val := ls.Get(offset)
		if lv, ok := val.(lua.LNumber); ok {
			args = append(args, float32(lv))
		} else if val.Type() == lua.LTNil {
			break
		} else {
			ls.ArgError(offset, "argument wrong type, should be number")
		}
	}

	return args
}

func extractIntArray(ls *lua.LState, offset int) []int32 {
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

	return args
}

func returnUD(ls *lua.LState, metatable string, item interface{}) int {
	f := ls.NewUserData()
	f.Value = item
	ls.SetMetatable(f, ls.GetTypeMetatable(metatable))
	ls.Push(f)
	return 1
}

func extractColor(ls *lua.LState, offset int) (r, g, b, a float32) {
	args := extractFloatArray(ls, offset)
	if len(args) == 0 {
		return 1, 1, 1, 1
	}
	argsLength := len(args)
	switch argsLength {
	case 4:
		return args[0], args[1], args[2], args[3]
	case 3:
		return args[0], args[1], args[2], 1
	case 2:
		ls.ArgError(offset, "argument wrong type, should be number")
	case 1:
		return args[0], args[0], args[0], 1
	}
	return 1, 1, 1, 1
}

func toWrap(ls *lua.LState, offset int) gfx.WrapMode {
	wrapStr := toStringD(ls, offset, "clamp")
	switch wrapStr {
	case "clamp":
		return gfx.WrapClamp
	case "repeat":
		return gfx.WrapRepeat
	case "mirror":
		return gfx.WrapMirroredRepeat
	default:
		ls.ArgError(offset, "invalid wrap mode")
	}
	return gfx.WrapClamp
}

func fromWrap(mode gfx.WrapMode) string {
	switch mode {
	case gfx.WrapRepeat:
		return "repeat"
	case gfx.WrapMirroredRepeat:
		return "mirror"
	case gfx.WrapClamp:
		fallthrough
	default:
		return "clamp"
	}
}

func toFilter(ls *lua.LState, offset int) gfx.FilterMode {
	wrapStr := toStringD(ls, offset, "near")
	switch wrapStr {
	case "none":
		return gfx.FilterNone
	case "near":
		return gfx.FilterNearest
	case "linear":
		return gfx.FilterLinear
	default:
		ls.ArgError(offset, "invalid filter mode")
	}
	return gfx.FilterNearest
}

func fromFilter(mode gfx.FilterMode) string {
	switch mode {
	case gfx.FilterLinear:
		return "linear"
	case gfx.FilterNone:
		return "none"
	case gfx.FilterNearest:
		fallthrough
	default:
		return "near"
	}
}

func toUsage(ls *lua.LState, offset int) gfx.Usage {
	wrapStr := toStringD(ls, offset, "dynamic")
	switch wrapStr {
	case "dynamic":
		return gfx.UsageDynamic
	case "static":
		return gfx.UsageStatic
	case "stream":
		return gfx.UsageStream
	default:
		ls.ArgError(offset, "invalid usage mode")
	}
	return gfx.UsageDynamic
}

func fromUsage(mode gfx.Usage) string {
	switch mode {
	case gfx.UsageStatic:
		return "static"
	case gfx.UsageStream:
		return "stream"
	case gfx.UsageDynamic:
		fallthrough
	default:
		return "dynamic"
	}
}

func toCompareMode(wrapStr string, offset int) gfx.CompareMode {
	switch wrapStr {
	case "always":
		return gfx.CompareAlways
	case "greater":
		return gfx.CompareGreater
	case "equal":
		return gfx.CompareEqual
	case "gequal":
		return gfx.CompareGequal
	case "less":
		return gfx.CompareLess
	case "nequal":
		return gfx.CompareNotequal
	case "lequal":
		return gfx.CompareLequal
	}
	return gfx.CompareAlways
}

func fromCompareMode(mode gfx.CompareMode) string {
	switch mode {
	case gfx.CompareLequal:
		return "lequal"
	case gfx.CompareNotequal:
		return "nequal"
	case gfx.CompareLess:
		return "less"
	case gfx.CompareGequal:
		return "gequal"
	case gfx.CompareEqual:
		return "equal"
	case gfx.CompareGreater:
		return "greater"
	case gfx.CompareAlways:
		fallthrough
	default:
		return "always"
	}
}

func toStencilAction(wrapStr string, offset int) gfx.StencilAction {
	switch wrapStr {
	case "replace":
		return gfx.StencilReplace
	case "increment":
		return gfx.StencilIncrement
	case "decrement":
		return gfx.StencilDecrement
	case "incrementwrap":
		return gfx.StencilIncrementWrap
	case "decrementwrap":
		return gfx.StencilDecrementWrap
	case "invert":
		return gfx.StencilInvert
	}
	return gfx.StencilReplace
}
