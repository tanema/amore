package gfx

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl32/matstack"
)

const (
	ATTRIB_POS = iota
	ATTRIB_TEXCOORD
	ATTRIB_COLOR
	ATTRIB_MAX_ENUM
)

type BlendMode int

const (
	BLENDMODE_ALPHA BlendMode = iota
	BLENDMODE_MULTIPLICATIVE
	BLENDMODE_PREMULTIPLIED
	BLENDMODE_SUBTRACTIVE
	BLENDMODE_ADDITIVE
	BLENDMODE_SCREEN
	BLENDMODE_REPLACE
)

type Viewport [4]int32 //The Viewport Values (X, Y, Width, Height)

var (
	opengl_version         string
	opengl_vendor          string
	maxAnisotropy          float32
	maxTextureSize         int32
	maxRenderTargets       int32
	maxRenderbufferSamples int32
	maxTextureUnits        int32
	screen_width           = 0
	screen_height          = 0
	modelIdent             = mgl32.Ident4()
	defaultShader          *Shader

	gl_state = glState{}
	states   = displayStateStack{newDisplayState()}
)

func InitContext(w, h int) {
	if gl_state.initialized {
		return
	}

	// Okay, setup OpenGL.
	gl.Init()

	//Get system info
	opengl_version = gl.GoStr(gl.GetString(gl.VERSION))
	opengl_vendor = gl.GoStr(gl.GetString(gl.VENDOR))
	gl_state.framebufferSRGBEnabled = gl.IsEnabled(gl.FRAMEBUFFER_SRGB)
	gl.GetIntegerv(gl.VIEWPORT, &gl_state.viewport[0])
	// And the current scissor - but we need to compensate for GL scissors
	// starting at the bottom left instead of top left.
	gl.GetIntegerv(gl.SCISSOR_BOX, &states.back().scissorBox[0])
	states.back().scissorBox[1] = gl_state.viewport[3] - (states.back().scissorBox[1] + states.back().scissorBox[3])

	gl.GetFloatv(gl.POINT_SIZE, &states.back().pointSize)
	gl.GetFloatv(gl.MAX_TEXTURE_MAX_ANISOTROPY_EXT, &maxAnisotropy)
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &maxTextureSize)
	gl.GetIntegerv(gl.MAX_SAMPLES, &maxRenderbufferSamples)
	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &maxTextureUnits)
	gl.GetIntegerv(gl.MAX_DRAW_BUFFERS, &maxRenderTargets)
	var maxattachments int32
	gl.GetIntegerv(gl.MAX_COLOR_ATTACHMENTS, &maxattachments)
	if maxattachments < maxRenderTargets {
		maxRenderTargets = maxattachments
	}

	// Enable blending
	gl.Enable(gl.BLEND)
	// Auto-generated mipmaps should be the best quality possible
	gl.Hint(gl.GENERATE_MIPMAP_HINT, gl.NICEST)
	// Make sure antialiasing works when set elsewhere
	gl.Enable(gl.MULTISAMPLE)
	// Set pixel row alignment
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	//default matricies
	gl_state.projectionStack = matstack.NewMatStack()
	gl_state.viewStack = matstack.NewMatStack() //stacks are initialized with ident matricies on top

	SetViewportSize(w, h)
	SetBackgroundColor(0, 0, 0, 1)

	gl_state.boundTextures = make([]uint32, maxTextureUnits)
	var curgltextureunit int32
	gl.GetIntegerv(gl.ACTIVE_TEXTURE, &curgltextureunit)
	gl_state.curTextureUnit = uint32(curgltextureunit) - gl.TEXTURE0
	// Retrieve currently bound textures for each texture unit.
	for i := 0; i < len(gl_state.boundTextures); i++ {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		var boundTex int32
		gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &boundTex)
		gl_state.boundTextures[i] = uint32(boundTex)
	}
	gl.ActiveTexture(uint32(curgltextureunit))
	createDefaultTexture()
	setTextureUnit(0)

	// We always need a default shader.
	defaultShader = NewShader()
	SetShader(defaultShader)

	gl_state.initialized = true

	LoadAllVolatile()
}

// Set the 'default' texture (id 0) as a repeating white pixel. Otherwise,
// texture2D calls inside a shader would return black when drawing graphics
// primitives, which would create the need to use different "passthrough"
// shaders for untextured primitives vs images.
func createDefaultTexture() {
	gl.GenTextures(1, &gl_state.defaultTexture)
	bindTexture(gl_state.defaultTexture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	pix := []uint8{255, 255, 255, 255}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pix))
}

func prepareDraw(model *mgl32.Mat4) {
	if model == nil {
		model = &modelIdent
	}

	states.back().shader.SendMat4("ProjectionMat", gl_state.projectionStack.Peek())
	states.back().shader.SendMat4("ViewMat", gl_state.viewStack.Peek())
	states.back().shader.SendMat4("ModelMat", *model)
	states.back().shader.SendFloat("ScreenSize", float32(screen_width), float32(screen_height), 0, 0)
	states.back().shader.SendFloat("PointSize", states.back().pointSize)
}

func setTextureUnit(textureunit uint32) error {
	if textureunit < 0 || int(textureunit) >= len(gl_state.boundTextures) {
		return fmt.Errorf("Invalid texture unit index (%v).", textureunit)
	}

	if textureunit != gl_state.curTextureUnit {
		gl.ActiveTexture(gl.TEXTURE0 + textureunit)
	}

	gl_state.curTextureUnit = textureunit
	return nil
}

func bindTexture(texture uint32) {
	if texture != gl_state.boundTextures[gl_state.curTextureUnit] {
		gl_state.boundTextures[gl_state.curTextureUnit] = texture
		gl.BindTexture(gl.TEXTURE_2D, texture)
	}
}

func bindTextureToUnit(texture, textureunit uint32, restoreprev bool) error {
	if texture != gl_state.boundTextures[textureunit] {
		oldtextureunit := gl_state.curTextureUnit
		if err := setTextureUnit(textureunit); err != nil {
			return err
		}
		gl_state.boundTextures[textureunit] = texture
		gl.BindTexture(gl.TEXTURE_2D, texture)
		if restoreprev {
			return setTextureUnit(oldtextureunit)
		}
	}
	return nil
}

func deleteTexture(texture uint32) {
	// glDeleteTextures binds texture 0 to all texture units the deleted texture
	// was bound to before deletion.
	for i, texid := range gl_state.boundTextures {
		if texid == texture {
			gl_state.boundTextures[i] = 0
		}
	}

	gl.DeleteTextures(1, &texture)
}

func DeInit() {
	UnloadAllVolatile()
	gl.DeleteTextures(1, &gl_state.defaultTexture)
	gl_state.defaultTexture = 0
}

func GetViewport() Viewport {
	return gl_state.viewport
}

func SetViewportSize(w, h int) {
	screen_width = w
	screen_height = h
	// Set the viewport to top-left corner.
	gl.Viewport(0, 0, int32(screen_width), int32(screen_height))
	gl_state.viewport = Viewport{0, 0, int32(screen_width), int32(screen_height)}
	gl_state.projectionStack.Load(mgl32.Ortho(0, float32(screen_width), float32(screen_height), 0, -1, 1))
	setScissor(states.back().scissorBox[0], states.back().scissorBox[1], states.back().scissorBox[2], states.back().scissorBox[3])
}

func Reset() {
	Origin()
	SetBlendMode(BLENDMODE_ALPHA)
	Clear()
}

func Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
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

func setScissor(x, y, width, height int32) {
	if gl_state.currentCanvas != nil {
		gl.Scissor(x, y, width, height)
	} else {
		// With no Canvas active, we need to compensate for glScissor starting
		// from the lower left of the viewport instead of the top left.
		gl.Scissor(x, gl_state.viewport[3]-(y+height), width, height)
	}
	states.back().scissorBox = Viewport{x, y, width, height}
	states.back().scissor = true
}

func SetScissor(x, y, width, height int32) {
	gl.Enable(gl.SCISSOR_TEST)
	// OpenGL's reversed y-coordinate is compensated for in OpenGL::setScissor.
	setScissor(x, y, width, height)
}

func ClearScissor() {
	gl.Disable(gl.SCISSOR_TEST)
	states.back().scissor = false
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

func SetWireframe(enable bool) {
	if enable {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
	states.back().wireframe = enable
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
	states.back().shader.Attach()
}

func SetBackgroundColor(r, g, b, a float32) {
	states.back().background_color = Color{r / 255.0, g / 255.0, b / 255.0, a / 255.0}
	gl.ClearColor(r/255.0, g/255.0, b/255.0, a/255.0)
}

func SetColor(r, g, b, a float32) {
	states.back().color = Color{r / 255.0, g / 255.0, b / 255.0, a / 255.0}
	gl.VertexAttrib4f(ATTRIB_COLOR, r/255.0, g/255.0, b/255.0, a/255.0)
}

func SetColorC(c Color) {
	states.back().color = c
	gl.VertexAttrib4f(ATTRIB_COLOR, c[0], c[1], c[2], c[3])
}

func GetColor() Color {
	return states.back().color
}

func SetFont(font *Font) {
	states.back().font = font
}

func GetFont() *Font {
	return states.back().font
}

func SetCanvas(canvases ...Canvas) {
	if canvases == nil || len(canvases) < 0 {
		states.back().canvases = nil
		if gl_state.currentCanvas != nil {
			gl_state.currentCanvas.stopGrab()
		}
	} else {
		states.back().canvases = canvases
		states.back().canvases[0].startGrab(canvases[1:]...)
	}
}

func GetCanvas() []Canvas {
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

	gl.BlendEquation(uint32(fn))
	gl.BlendFuncSeparate(uint32(srcRGB), uint32(dstRGB), uint32(srcA), uint32(dstA))
	states.back().blend_mode = mode
}

func SetDefaultFilter(min, mag FilterMode, anisotropy float32) {
	states.back().defaultFilter = Filter{
		min:        min,
		mag:        mag,
		anisotropy: float32(math.Min(math.Max(float64(anisotropy), 1.0), float64(maxAnisotropy))),
	}
}

func GetDefaultFilter() Filter {
	return states.back().defaultFilter
}
