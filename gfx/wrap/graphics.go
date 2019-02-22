package wrap

import (
	"github.com/yuin/gopher-lua"

	"github.com/tanema/amore/gfx"
)

func gfxCirle(ls *lua.LState) int {
	gfx.Circle(extractMode(ls, 1), toFloat(ls, 2), toFloat(ls, 3), toFloat(ls, 4), toIntD(ls, 5, 30))
	return 0
}

func gfxArc(ls *lua.LState) int {
	gfx.Arc(extractMode(ls, 1), toFloat(ls, 2), toFloat(ls, 3), toFloat(ls, 4), toFloat(ls, 5), toFloat(ls, 6), toIntD(ls, 6, 30))
	return 0
}

func gfxEllipse(ls *lua.LState) int {
	gfx.Ellipse(extractMode(ls, 1), toFloat(ls, 2), toFloat(ls, 3), toFloat(ls, 4), toFloat(ls, 5), toIntD(ls, 6, 30))
	return 0
}

func gfxPoints(ls *lua.LState) int {
	gfx.Points(extractCoords(ls, 1))
	return 0
}

func gfxLine(ls *lua.LState) int {
	gfx.PolyLine(extractCoords(ls, 1))
	return 0
}

func gfxRectangle(ls *lua.LState) int {
	gfx.Rect(extractMode(ls, 1), toFloat(ls, 2), toFloat(ls, 3), toFloat(ls, 4), toFloat(ls, 5))
	return 0
}

func gfxPolygon(ls *lua.LState) int {
	gfx.Polygon(extractMode(ls, 1), extractCoords(ls, 2))
	return 0
}

func gfxScreenShot(ls *lua.LState) int {
	ls.Push(toUD(ls, "Image", gfx.NewScreenshot()))
	return 1
}

func gfxGetViewport(ls *lua.LState) int {
	for _, x := range gfx.GetViewport() {
		ls.Push(lua.LNumber(x))
	}
	return 4
}

func gfxSetViewport(ls *lua.LState) int {
	viewport := extractCoords(ls, 1)
	if len(viewport) == 2 {
		gfx.SetViewport(0, 0, int32(viewport[0]), int32(viewport[1]))
	} else if len(viewport) == 4 {
		gfx.SetViewport(int32(viewport[0]), int32(viewport[1]), int32(viewport[2]), int32(viewport[3]))
	} else {
		ls.ArgError(1, "either provide (x, y, w, h) or (w, h)")
	}
	return 0
}

func gfxGetWidth(ls *lua.LState) int {
	ls.Push(lua.LNumber(gfx.GetWidth()))
	return 1
}

func gfxGetHeight(ls *lua.LState) int {
	ls.Push(lua.LNumber(gfx.GetHeight()))
	return 1
}

func gfxGetDimensions(ls *lua.LState) int {
	w, h := gfx.GetDimensions()
	ls.Push(lua.LNumber(w))
	ls.Push(lua.LNumber(h))
	return 2
}

func gfxOrigin(ls *lua.LState) int {
	gfx.Origin()
	return 0
}

func gfxTranslate(ls *lua.LState) int {
	x, y := toFloat(ls, 1), toFloat(ls, 2)
	gfx.Translate(x, y)
	return 0
}

func gfxRotate(ls *lua.LState) int {
	gfx.Rotate(toFloat(ls, 1))
	return 0
}

func gfxScale(ls *lua.LState) int {
	sx := toFloat(ls, 1)
	sy := toFloatD(ls, 2, sx)
	gfx.Scale(sx, sy)
	return 0
}

func gfxShear(ls *lua.LState) int {
	kx := toFloat(ls, 1)
	ky := toFloatD(ls, 2, kx)
	gfx.Shear(kx, ky)
	return 0
}

func gfxPush(ls *lua.LState) int {
	gfx.Push()
	return 0
}

func gfxPop(ls *lua.LState) int {
	gfx.Pop()
	return 0
}

func gfxClear(ls *lua.LState) int {
	gfx.Clear(extractColor(ls, 1))
	return 0
}

func gfxSetScissor(ls *lua.LState) int {
	args := extractFloatArray(ls, 1)
	if len(args) == 2 {
		gfx.SetScissor(0, 0, int32(args[0]), int32(args[1]))
	} else if len(args) == 4 {
		gfx.SetScissor(int32(args[0]), int32(args[1]), int32(args[2]), int32(args[3]))
	} else if len(args) == 0 {
		gfx.ClearScissor()
	} else {
		ls.ArgError(1, "either pass 2, 4 or no arguments")
	}
	return 0
}

func gfxGetScissor(ls *lua.LState) int {
	x, y, w, h := gfx.GetScissor()
	ls.Push(lua.LNumber(x))
	ls.Push(lua.LNumber(y))
	ls.Push(lua.LNumber(w))
	ls.Push(lua.LNumber(h))
	return 4
}

func gfxSetLineWidth(ls *lua.LState) int {
	gfx.SetLineWidth(toFloat(ls, 1))
	return 0
}

func gfxSetLineJoin(ls *lua.LState) int {
	gfx.SetLineJoin(extractLineJoin(ls, 1))
	return 0
}

func gfxGetLineWidth(ls *lua.LState) int {
	ls.Push(lua.LNumber(gfx.GetLineWidth()))
	return 1
}

func gfxGetLineJoin(ls *lua.LState) int {
	ls.Push(lua.LString(gfx.GetLineJoin()))
	return 1
}

func gfxSetPointSize(ls *lua.LState) int {
	gfx.SetPointSize(toFloat(ls, 1))
	return 0
}

func gfxGetPointSize(ls *lua.LState) int {
	ls.Push(lua.LNumber(gfx.GetPointSize()))
	return 1
}

func gfxSetColor(ls *lua.LState) int {
	gfx.SetColor(extractColor(ls, 1))
	return 0
}

func gfxSetBackgroundColor(ls *lua.LState) int {
	gfx.SetBackgroundColor(extractColor(ls, 1))
	return 0
}

func gfxGetColor(ls *lua.LState) int {
	for _, x := range gfx.GetColor() {
		ls.Push(lua.LNumber(x))
	}
	return 4
}

func gfxGetBackgroundColor(ls *lua.LState) int {
	for _, x := range gfx.GetBackgroundColor() {
		ls.Push(lua.LNumber(x))
	}
	return 4
}

func gfxGetColorMask(ls *lua.LState) int {
	r, g, b, a := gfx.GetColorMask()
	ls.Push(lua.LBool(r))
	ls.Push(lua.LBool(g))
	ls.Push(lua.LBool(b))
	ls.Push(lua.LBool(a))
	return 4
}

func gfxSetColorMask(ls *lua.LState) int {
	args := []bool{}
	offset := 1
	for x := ls.Get(offset); x != nil; offset++ {
		val := ls.Get(offset)
		if lv, ok := val.(lua.LBool); ok {
			args = append(args, bool(lv))
		} else if val.Type() == lua.LTNil {
			break
		} else {
			ls.ArgError(offset, "argument wrong type, should be boolean")
		}
	}

	if len(args) == 4 {
		gfx.SetColorMask(args[0], args[1], args[2], args[3])
	} else if len(args) == 3 {
		gfx.SetColorMask(args[0], args[1], args[2], true)
	} else if len(args) == 1 {
		gfx.SetColorMask(args[0], args[0], args[0], args[0])
	} else if len(args) == 0 {
		gfx.ClearStencilTest()
	} else {
		ls.ArgError(offset, "invalid argument count")
	}

	return 0
}

func gfxSetFont(ls *lua.LState) int {
	gfx.SetFont(toFont(ls, 1))
	return 0
}

func gfxGetFont(ls *lua.LState) int {
	ls.Push(toUD(ls, "Font", gfx.GetFont()))
	return 1
}

func gfxSetBlendMode(ls *lua.LState) int {
	gfx.SetBlendMode(extractBlendmode(ls, 1))
	return 0
}

//SetDefaultFilter
//func gfxGetStencilTest(ls *lua.LState) int {
//	return 0
//}

//func gfxSetStencilTest(ls *lua.LState) int {
//	//ClearStencilTest
//	return 0
//}
// Stencil, SetShader, GetCanvas, SetCanvas,
