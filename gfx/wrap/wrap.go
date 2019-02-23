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
	"newimage": gfxNewImage,
	"newtext":  gfxNewText,
}

var graphicsMetaTables = runtime.LuaMetaTable{
	"Image": {
		"draw":          gfxImageDraw,
		"getwidth":      gfxImageGetWidth,
		"getheigth":     gfxImageGetHeight,
		"getDimensions": gfxImageGetDimensions,
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
		// getWidth, getHeight, setFallback
	},
}

func init() {
	runtime.RegisterModule("gfx", graphicsFunctions, graphicsMetaTables)
}
