package wrap

import "github.com/tanema/amore/runtime"

var graphicsFunctions = runtime.LuaFuncs{
	"circle":             gfxCirle,
	"arc":                gfxArc,
	"ellipse":            gfxEllipse,
	"points":             gfxPoints,
	"line":               gfxLine,
	"rectangle":          gfxRectangle,
	"polygon":            gfxPolygon,
	"getviewport":        gfxGetViewport,
	"setviewport":        gfxSetViewport,
	"getwidth":           gfxGetWidth,
	"getheight":          gfxGetHeight,
	"getdimensions":      gfxGetDimensions,
	"origin":             gfxOrigin,
	"translate":          gfxTranslate,
	"rotate":             gfxRotate,
	"scale":              gfxScale,
	"shear":              gfxShear,
	"push":               gfxPush,
	"pop":                gfxPop,
	"clear":              gfxClear,
	"setscissor":         gfxSetScissor,
	"getscissor":         gfxGetScissor,
	"setlinewidth":       gfxSetLineWidth,
	"setlinejoin":        gfxSetLineJoin,
	"getlinewidth":       gfxGetLineWidth,
	"getlinejoin":        gfxGetLineJoin,
	"setpointsize":       gfxSetPointSize,
	"getpointsize":       gfxGetPointSize,
	"setcolor":           gfxSetColor,
	"setbackgroundcolor": gfxSetBackgroundColor,
	"getcolor":           gfxGetColor,
	"getbackgroundcolor": gfxGetBackgroundColor,
	"getcolormask":       gfxGetColorMask,
	"setcolormask":       gfxSetColorMask,
	"print":              gfxPrint,
	"printf":             gfxPrintf,
	"getfont":            gfxGetFont,
	"setfont":            gfxSetFont,
	"setblendmode":       gfxSetBlendMode,

	// metatable entries
	"newimage":  gfxNewImage,
	"newtext":   gfxNewText,
	"newfont":   gfxNewFont,
	"newquad":   gfxNewQuad,
	"newcanvas": gfxNewCanvas,
}

var graphicsMetaTables = runtime.LuaMetaTable{
	"Image": {
		"draw":          gfxTextureDraw,
		"drawq":         gfxTextureDrawq,
		"getwidth":      gfxTextureGetWidth,
		"getheigth":     gfxTextureGetHeight,
		"getDimensions": gfxTextureGetDimensions,
		"setwrap":       gfxTextureSetWrap,
		"setfilter":     gfxTextureSetFilter,
	},
	"Text": {
		"set":           gfxTextSet,
		"draw":          gfxTextDraw,
		"getwidth":      gfxTextGetWidth,
		"getheight":     gfxTextGetHeight,
		"getdimensions": gfxTextGetDimensions,
		"getfont":       gfxTextGetFont,
		"setfont":       gfxTextSetFont,
	},
	"Font": {
		"getwidth":    gfxFontGetWidth,
		"getheight":   gfxFontGetHeight,
		"setfallback": gfxFontSetFallback,
		"getwrap":     gfxFontGetWrap,
	},
	"Quad": {
		"getwidth":    gfxQuadGetWidth,
		"geteheight":  gfxQuadGetHeight,
		"getviewport": gfxQuadGetViewport,
		"setviewport": gfxQuadSetViewport,
	},
	"Canvas": {
		"newimage":      gfxCanvasNewImage,
		"draw":          gfxTextureDraw,
		"drawq":         gfxTextureDrawq,
		"getwidth":      gfxTextureGetWidth,
		"getheigth":     gfxTextureGetHeight,
		"getDimensions": gfxTextureGetDimensions,
		"setwrap":       gfxTextureSetWrap,
		"setfilter":     gfxTextureSetFilter,
	},
	"SpriteBatch": {},
	"Shader":      {},
}

func init() {
	runtime.RegisterModule("gfx", graphicsFunctions, graphicsMetaTables)
}
