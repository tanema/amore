package wrap

import "github.com/yuin/gopher-lua"

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

func toUD(ls *lua.LState, metatable string, item interface{}) *lua.LUserData {
	f := ls.NewUserData()
	f.Value = item
	ls.SetMetatable(f, ls.GetTypeMetatable(metatable))
	return f
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
