package gfx

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl32/matstack"
	"github.com/goxjs/gl"

	"github.com/tanema/amore/window"
)

var (
	opengl_version         string
	opengl_vendor          string
	maxAnisotropy          float32
	maxTextureSize         int32
	maxRenderTargets       int32
	maxRenderbufferSamples int32
	maxTextureUnits        int32
	screen_width           = int32(0)
	screen_height          = int32(0)
	modelIdent             = mgl32.Ident4()
	defaultShader          *Shader

	gl_state = glState{
		viewport: make([]int32, 4),
	}
	states = displayStateStack{newDisplayState()}
)

func InitContext(w, h int32) {
	if gl_state.initialized {
		return
	}

	// Okay, setup OpenGL.
	gl.ContextWatcher.OnMakeCurrent(nil)

	//Get system info
	opengl_version = gl.GetString(gl.VERSION)
	opengl_vendor = gl.GetString(gl.VENDOR)
	gl_state.defaultFBO = getCurrentFBO()
	gl.GetIntegerv(gl.VIEWPORT, gl_state.viewport)
	// And the current scissor - but we need to compensate for GL scissors
	// starting at the bottom left instead of top left.
	gl.GetIntegerv(gl.SCISSOR_BOX, states.back().scissorBox)
	states.back().scissorBox[1] = gl_state.viewport[3] - (states.back().scissorBox[1] + states.back().scissorBox[3])

	initMaxValues() //check shim code

	glcolor := []float32{1.0, 1.0, 1.0, 1.0}
	gl.VertexAttrib4fv(ATTRIB_COLOR, glcolor)
	gl.VertexAttrib4fv(ATTRIB_CONSTANTCOLOR, glcolor)
	useVertexAttribArrays(0)

	// Enable blending
	gl.Enable(gl.BLEND)
	SetBlendMode(BLENDMODE_ALPHA)
	// Auto-generated mipmaps should be the best quality possible
	gl.Hint(gl.GENERATE_MIPMAP_HINT, gl.NICEST)
	// Make sure antialiasing works when set elsewhere
	enableMultisample() //check shim code
	// Set pixel row alignment
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	//default matricies
	gl_state.projectionStack = matstack.NewMatStack()
	gl_state.viewStack = matstack.NewMatStack() //stacks are initialized with ident matricies on top

	SetViewportSize(w, h)
	SetBackgroundColor(0, 0, 0, 1)

	gl_state.boundTextures = make([]gl.Texture, maxTextureUnits)
	curgltextureunit := gl.GetInteger(gl.ACTIVE_TEXTURE)
	gl_state.curTextureUnit = int(curgltextureunit - gl.TEXTURE0)
	// Retrieve currently bound textures for each texture unit.
	for i := 0; i < len(gl_state.boundTextures); i++ {
		gl.ActiveTexture(gl.Enum(gl.TEXTURE0 + uint32(i)))
		gl_state.boundTextures[i] = gl.Texture{Value: uint32(gl.GetInteger(gl.TEXTURE_BINDING_2D))}
	}
	gl.ActiveTexture(gl.Enum(curgltextureunit))
	createDefaultTexture()
	setTextureUnit(0)

	// We always need a default shader.
	defaultShader = NewShader()

	gl_state.initialized = true

	loadAllVolatile()

	//have to set this after loadallvolatile() so we are sure the  default shader is loaded
	SetShader(nil)
}

// Set the 'default' texture (id 0) as a repeating white pixel. Otherwise,
// texture2D calls inside a shader would return black when drawing graphics
// primitives, which would create the need to use different "passthrough"
// shaders for untextured primitives vs images.
func createDefaultTexture() {
	gl_state.defaultTexture = gl.CreateTexture()
	bindTexture(gl_state.defaultTexture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexImage2D(gl.TEXTURE_2D, 0, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, []byte{255, 255, 255, 255})
}

func prepareDraw(model *mgl32.Mat4) {
	if model == nil {
		model = &modelIdent
	}

	gl_state.currentShader.SendMat4("ProjectionMat", gl_state.projectionStack.Peek())
	gl_state.currentShader.SendMat4("ViewMat", gl_state.viewStack.Peek())
	gl_state.currentShader.SendMat4("ModelMat", *model)
	gl_state.currentShader.SendFloat("ScreenSize", float32(screen_width), float32(screen_height), 0, 0)
	gl_state.currentShader.SendFloat("PointSize", states.back().pointSize)
}

func useVertexAttribArrays(arraybits uint32) {
	diff := arraybits ^ gl_state.enabledAttribArrays

	if diff == 0 {
		return
	}

	// Max 32 attributes. As of when this was written, no GL driver exposes more
	// than 32. Lets hope that doesn't change...
	for i := uint32(0); i < 32; i++ {
		bit := uint32(1 << i)
		if (diff & bit) > 0 {
			if (arraybits & bit) > 0 {
				gl.EnableVertexAttribArray(gl.Attrib{Value: uint(i)})
			} else {
				gl.DisableVertexAttribArray(gl.Attrib{Value: uint(i)})
			}
		}
	}

	gl_state.enabledAttribArrays = arraybits

	// glDisableVertexAttribArray will make the constant value for a vertex
	// attribute undefined. We rely on the per-vertex color attribute being
	// white when no per-vertex color is used, so we set it here.
	// FIXME: Is there a better place to do this?
	if (diff&ATTRIBFLAG_COLOR) > 0 && (arraybits&ATTRIBFLAG_COLOR) == 0 {
		gl.VertexAttrib4f(ATTRIB_COLOR, 1.0, 1.0, 1.0, 1.0)
	}
}

func setTextureUnit(textureunit int) error {
	if textureunit < 0 || textureunit >= len(gl_state.boundTextures) {
		return fmt.Errorf("Invalid texture unit index (%v).", textureunit)
	}

	if textureunit != gl_state.curTextureUnit {
		gl.ActiveTexture(gl.Enum(gl.TEXTURE0 + uint32(textureunit)))
	}

	gl_state.curTextureUnit = textureunit
	return nil
}

func bindTexture(texture gl.Texture) {
	if texture != gl_state.boundTextures[gl_state.curTextureUnit] {
		gl_state.boundTextures[gl_state.curTextureUnit] = texture
		gl.BindTexture(gl.TEXTURE_2D, texture)
	}
}

func bindTextureToUnit(texture gl.Texture, textureunit int, restoreprev bool) error {
	if texture != gl_state.boundTextures[textureunit] {
		oldtextureunit := gl_state.curTextureUnit
		if err := setTextureUnit(textureunit); err != nil {
			return err
		}
		gl_state.boundTextures[textureunit] = texture
		gl.BindTexture(gl.TEXTURE_2D, gl_state.boundTextures[textureunit])
		if restoreprev {
			return setTextureUnit(oldtextureunit)
		}
	}
	return nil
}

func HasFramebufferSRGB() bool {
	return gl_state.framebufferSRGBEnabled
}

func getDefaultFBO() gl.Framebuffer {
	return gl_state.defaultFBO
}

func getCurrentFBO() gl.Framebuffer {
	current_fbo := gl.GetInteger(gl.FRAMEBUFFER_BINDING)
	return gl.Framebuffer{Value: uint32(current_fbo)}
}

func GetMaxTextureSize() int32 {
	return maxTextureSize
}

func GetMaxRenderTargets() int32 {
	return maxRenderTargets
}

func GetMaxRenderbufferSamples() int32 {
	return maxRenderbufferSamples
}

func GetMaxTextureUnits() int32 {
	return maxTextureUnits
}

func GetVendor() string {
	return opengl_vendor
}

func hasFramebufferSRGB() bool {
	return gl_state.framebufferSRGBEnabled
}

func deleteTexture(texture gl.Texture) {
	// glDeleteTextures binds texture 0 to all texture units the deleted texture
	// was bound to before deletion.
	for i, texid := range gl_state.boundTextures {
		if texid == texture {
			gl_state.boundTextures[i] = gl.Texture{Value: 0}
		}
	}

	gl.DeleteTexture(texture)
}

func DeInit() {
	unloadAllVolatile()
	gl.DeleteTexture(gl_state.defaultTexture)
	gl_state.defaultTexture = gl.Texture{}
	gl_state.initialized = false
	gl.ContextWatcher.OnDetach()
}

func GetViewport() []int32 {
	return gl_state.viewport
}

func SetViewportSize(w, h int32) {
	SetViewport(0, 0, w, h)
}

func SetViewport(x, y, w, h int32) {
	screen_width = w
	screen_height = h
	// Set the viewport to top-left corner.
	gl.Viewport(int(y), int(x), int(screen_width), int(screen_height))
	gl_state.viewport = []int32{y, x, screen_width, screen_height}
	gl_state.projectionStack.Load(mgl32.Ortho(float32(x), float32(screen_width), float32(screen_height), float32(y), -1, 1))
	SetScissor(states.back().scissorBox[0], states.back().scissorBox[1], states.back().scissorBox[2], states.back().scissorBox[3])
}

func GetWidth() float32 {
	return float32(screen_width)
}

func GetHeight() float32 {
	return float32(screen_height)
}

func Clear(r, g, b, a float32) {
	gl.ClearColor(r/255.0, g/255.0, b/255.0, a/255.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func ClearC(c Color) {
	Clear(c[0], c[1], c[2], c[3])
}

func Present() {
	if !IsActive() {
		return
	}

	// Make sure we don't have a canvas active.
	canvases := states.back().canvases
	SetCanvas()
	window.GetCurrent().SwapBuffers()
	// Restore the currently active canvas, if there is one.
	SetCanvas(canvases...)
}

func IsActive() bool {
	// The graphics module is only completely 'active' if there's a window, a
	// context, and the active variable is set.
	return gl_state.active && isCreated() && window.GetCurrent().IsOpen()
}

func SetActive(enable bool) {
	// Make sure all pending OpenGL commands have fully executed before
	// returning, when going from active to inactive. This is required on iOS.
	if isCreated() && gl_state.active && !enable {
		gl.Finish()
	}

	gl_state.active = enable
}

func isCreated() bool {
	return gl_state.initialized
}

func Origin() {
	gl_state.viewStack.LoadIdent()
	states.back().pixelSize = 1.0
}

func Translate(x, y float32) {
	gl_state.viewStack.LeftMul(mgl32.Translate3D(x, y, 0))
}

func Rotate(angle float32) {
	gl_state.viewStack.LeftMul(mgl32.HomogRotate3D(angle, mgl32.Vec3{0, 0, 1}))
}

func Scale(args ...float32) {
	if args == nil || len(args) == 0 {
		panic("not enough params passed to scale call")
	}
	var sx, sy float32
	sx = args[0]
	if len(args) > 1 {
		sy = args[1]
	} else {
		sy = sx
	}

	gl_state.viewStack.LeftMul(mgl32.Scale3D(sx, sy, 1))

	states.back().pixelSize *= (2.0 / (mgl32.Abs(sx) + mgl32.Abs(sy)))
}

func Shear(args ...float32) {
	if args == nil || len(args) == 0 {
		panic("not enough params passed to scale call")
	}
	var kx, ky float32
	kx = args[0]
	if len(args) > 1 {
		ky = args[1]
	} else {
		ky = kx
	}

	gl_state.viewStack.LeftMul(mgl32.ShearX3D(kx, ky))
}

func Push() {
	gl_state.viewStack.Push()
	states.push(*states.back())
}

func Pop() {
	gl_state.viewStack.Pop()
	states.pop()
}

func SetScissor(args ...int32) {
	if args == nil {
		gl.Disable(gl.SCISSOR_TEST)
		states.back().scissor = false
	} else if len(args) == 4 {
		x, y, width, height := args[0], args[1], args[2], args[3]
		gl.Enable(gl.SCISSOR_TEST)
		if gl_state.currentCanvas != nil {
			gl.Scissor(x, y, width, height)
		} else {
			// With no Canvas active, we need to compensate for glScissor starting
			// from the lower left of the viewport instead of the top left.
			gl.Scissor(x, gl_state.viewport[3]-(y+height), width, height)
		}
		states.back().scissorBox = []int32{x, y, width, height}
		states.back().scissor = true
	} else {
		panic("incorrect number of arguments to setscissor")
	}
}

func IntersectScissor(x, y, width, height int32) {
	rect := states.back().scissorBox

	if !states.back().scissor {
		rect[0] = 0
		rect[1] = 0
		rect[2] = math.MaxInt32
		rect[3] = math.MaxInt32
	}

	x1 := int32(math.Max(float64(rect[0]), float64(x)))
	y1 := int32(math.Max(float64(rect[1]), float64(y)))
	x2 := int32(math.Min(float64(rect[0]+rect[2]), float64(x+width)))
	y2 := int32(math.Min(float64(rect[1]+rect[3]), float64(y+height)))

	SetScissor(x1, y1, int32(math.Max(0, float64(x2-x1))), int32(math.Max(0, float64(y2-y1))))
}

func ClearScissor() {
	gl.Disable(gl.SCISSOR_TEST)
	states.back().scissor = false
}

func Stencil(stencil_func func()) {
	StencilExt(stencil_func, STENCIL_REPLACE, 1, false)
}

func StencilExt(stencil_func func(), action StencilAction, value int32, keepvalues bool) {
	gl_state.writingToStencil = true
	if !keepvalues {
		clearStencil()
	}
	if gl_state.currentCanvas != nil {
		gl_state.currentCanvas.checkCreateStencil()
	}
	gl.Enable(gl.STENCIL_TEST)
	gl.ColorMask(false, false, false, false)
	gl.StencilFunc(gl.ALWAYS, int(value), 0xFF)
	gl.StencilOp(gl.KEEP, gl.KEEP, gl.Enum(action))

	stencil_func()

	gl_state.writingToStencil = false
	SetColorMask(states.back().colorMask)
	SetStencilTest(states.back().stencilCompare, states.back().stencilTestValue)
}

func SetStencilTest(compare CompareMode, value int32) {
	if gl_state.writingToStencil {
		return
	}

	states.back().stencilCompare = compare
	states.back().stencilTestValue = value
	if compare == COMPARE_ALWAYS {
		gl.Disable(gl.STENCIL_TEST)
		return
	}

	if gl_state.currentCanvas != nil {
		gl_state.currentCanvas.checkCreateStencil()
	}

	gl.Enable(gl.STENCIL_TEST)
	gl.StencilFunc(gl.Enum(compare), int(value), 0xFF)
	gl.StencilOp(gl.KEEP, gl.KEEP, gl.REPLACE)
}

func ClearStencilTest() {
	SetStencilTest(COMPARE_ALWAYS, 0)
}

func GetStencilTest() (CompareMode, int32) {
	return states.back().stencilCompare, states.back().stencilTestValue
}

func clearStencil() {
	gl.Clear(gl.STENCIL_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func GetScissor() (int32, int32, int32, int32) {
	return states.back().scissorBox[0], states.back().scissorBox[1], states.back().scissorBox[2], states.back().scissorBox[3]
}

func SetColorMask(mask ColorMask) {
	gl.ColorMask(mask.r, mask.g, mask.b, mask.a)
	states.back().colorMask = mask
}

func GetColorMask() ColorMask {
	return states.back().colorMask
}

func SetLineWidth(width float32) {
	states.back().line_width = width
}

func SetLineStyle(style LineStyle) {
	states.back().line_style = style
}

func SetLineJoin(join LineJoin) {
	states.back().line_join = join
}

func GetLineWidth() float32 {
	return states.back().line_width
}

func GetLineStyle() LineStyle {
	return states.back().line_style
}

func GetLineJoin() LineJoin {
	return states.back().line_join
}

func SetPointSize(size float32) {
	states.back().pointSize = size
}

func GetPointSize() float32 {
	return states.back().pointSize
}

func IsWireframe() bool {
	return states.back().wireframe
}

func SetShader(shader *Shader) {
	if shader == nil {
		states.back().shader = defaultShader
	} else {
		states.back().shader = shader
	}
	states.back().shader.attach(false)
}

func SetBackgroundColor(r, g, b, a float32) {
	states.back().background_color = Color{r / 255.0, g / 255.0, b / 255.0, a / 255.0}
	gl.ClearColor(r/255.0, g/255.0, b/255.0, a/255.0)
}

func SetBackgroundColorC(c Color) {
	states.back().background_color = c
	gl.ClearColor(c[0], c[1], c[2], c[3])
}

func GetBackgroundColor() (r, g, b, a float32) {
	bc := states.back().background_color
	return bc[0], bc[1], bc[2], bc[3]
}

func GetBackgroundColorC() Color {
	return states.back().background_color
}

func setColor(r, g, b, a float32) {
	SetColorC(Color{r, g, b, a})
}

func SetColor(r, g, b, a float32) {
	setColor(r/255.0, g/255.0, b/255.0, a/255.0)
}

func SetColorC(c Color) {
	states.back().color = c

	gl.VertexAttrib4f(ATTRIB_CONSTANTCOLOR, c[0], c[1], c[2], c[3])
}

func GetColor() Color {
	return states.back().color
}

func SetFont(font *Font) {
	states.back().font = font
}

func GetFont() *Font {
	//if no font set, use default font
	if states.back().font == nil {
		SetFont(NewFont("arialbd.ttf", 15))
	}
	return states.back().font
}

func SetCanvas(canvases ...*Canvas) error {
	if canvases == nil || len(canvases) < 0 {
		states.back().canvases = nil
		if gl_state.currentCanvas != nil {
			gl_state.currentCanvas.stopGrab(false)
			gl_state.currentCanvas = nil
		}
	} else {
		states.back().canvases = canvases
		return states.back().canvases[0].startGrab(canvases[1:]...)
	}
	return nil
}

func GetCanvas() []*Canvas {
	return states.back().canvases
}

func SetBlendMode(mode BlendMode) {
	fn := gl.FUNC_ADD
	srcRGB := gl.ONE
	srcA := gl.ONE
	dstRGB := gl.ZERO
	dstA := gl.ZERO

	switch mode {
	case BLENDMODE_ALPHA:
		srcRGB = gl.SRC_ALPHA
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_ALPHA
		dstA = gl.ONE_MINUS_SRC_ALPHA
	case BLENDMODE_MULTIPLICATIVE:
		srcRGB = gl.DST_COLOR
		srcA = gl.DST_COLOR
		dstRGB = gl.ZERO
		dstA = gl.ZERO
	case BLENDMODE_PREMULTIPLIED:
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_ALPHA
		dstA = gl.ONE_MINUS_SRC_ALPHA
	case BLENDMODE_SUBTRACTIVE:
		fn = gl.FUNC_REVERSE_SUBTRACT
	case BLENDMODE_ADDITIVE:
		srcRGB = gl.SRC_ALPHA
		srcA = gl.SRC_ALPHA
		dstRGB = gl.ONE
		dstA = gl.ONE
	case BLENDMODE_SCREEN:
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_COLOR
		dstA = gl.ONE_MINUS_SRC_COLOR
		break
	case BLENDMODE_REPLACE:
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ZERO
		dstA = gl.ZERO
	}

	gl.BlendEquation(gl.Enum(fn))
	gl.BlendFuncSeparate(gl.Enum(srcRGB), gl.Enum(dstRGB), gl.Enum(srcA), gl.Enum(dstA))
	states.back().blend_mode = mode
}

func SetDefaultFilter(min, mag FilterMode, anisotropy float32) {
	states.back().defaultFilter = Filter{
		min:        min,
		mag:        mag,
		anisotropy: float32(math.Min(math.Max(float64(anisotropy), 1.0), float64(maxAnisotropy))),
	}
}

func SetDefaultFilterF(f Filter) {
	states.back().defaultFilter = f
}

func GetDefaultFilter() Filter {
	return states.back().defaultFilter
}

func SetDefaultMipmapFilter(filter FilterMode, sharpness float32) {
	states.back().defaultMipmapFilter = filter
	states.back().defaultMipmapSharpness = sharpness
}

func GetDefaultMipmapFilter() (filter FilterMode, sharpness float32) {
	return states.back().defaultMipmapFilter, states.back().defaultMipmapSharpness
}
