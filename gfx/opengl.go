package gfx

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl32/matstack"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"

	"github.com/tanema/amore/gfx/font"
)

var (
	maxAnisotropy          float32
	maxTextureSize         int32
	maxRenderTargets       int32
	maxRenderbufferSamples int32
	maxTextureUnits        int32
	screenWidth            = int32(0)
	screenHeight           = int32(0)
	modelIdent             = mgl32.Ident4()
	defaultShader          *Shader
	defaultFace, _         = font.Bold(20)
	defaultFont            = newFont(defaultFace)

	glState = openglState{
		viewport: make([]int32, 4),
	}
	states = displayStateStack{stack: []displayState{newDisplayState()}}
)

// InitContext will initiate the opengl context with a viewport in the size of
// w x h. This is generally called from the game loop and wont generally need to
// be called unless you are rolling your own game loop.
func InitContext(window *glfw.Window) {
	if glState.initialized {
		return
	}

	//Get system info
	glState.defaultFBO = gl.GetBoundFramebuffer()
	gl.GetIntegerv(gl.VIEWPORT, glState.viewport)
	// And the current scissor - but we need to compensate for GL scissors
	// starting at the bottom left instead of top left.
	fmt.Println(states)
	fmt.Println(states.back())
	fmt.Println(states.back().scissorBox)
	gl.GetIntegerv(gl.SCISSOR_BOX, states.back().scissorBox)
	states.back().scissorBox[1] = glState.viewport[3] - (states.back().scissorBox[1] + states.back().scissorBox[3])

	maxTextureSize = int32(gl.GetInteger(gl.MAX_TEXTURE_SIZE))
	maxTextureUnits = int32(gl.GetInteger(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS))
	glState.textureCounters = make([]int, maxTextureUnits)

	glcolor := []float32{1.0, 1.0, 1.0, 1.0}
	gl.VertexAttrib4fv(shaderColor, glcolor)
	gl.VertexAttrib4fv(shaderConstantColor, glcolor)
	useVertexAttribArrays()

	// Enable blending
	gl.Enable(gl.BLEND)
	SetBlendMode("alpha")
	// Auto-generated mipmaps should be the best quality possible
	gl.Hint(gl.GENERATE_MIPMAP_HINT, gl.NICEST)
	// Set pixel row alignment
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	//default matricies
	glState.projectionStack = matstack.NewMatStack()
	glState.viewStack = matstack.NewMatStack() //stacks are initialized with ident matricies on top

	w, h := window.GetFramebufferSize()
	SetViewportSize(int32(w), int32(h))
	SetBackgroundColor(0, 0, 0, 1)

	glState.boundTextures = make([]gl.Texture, maxTextureUnits)
	curgltextureunit := gl.GetInteger(gl.ACTIVE_TEXTURE)
	glState.curTextureUnit = int(curgltextureunit - gl.TEXTURE0)
	// Retrieve currently bound textures for each texture unit.
	// for i := 0; i < len(glState.boundTextures); i++ {
	//	gl.ActiveTexture(gl.Enum(gl.TEXTURE0 + uint32(i)))
	//	glState.boundTextures[i] = gl.Texture{Value: uint32(gl.GetInteger(gl.TEXTURE_BINDING_2D))}
	// }
	gl.ActiveTexture(gl.Enum(curgltextureunit))
	createDefaultTexture()
	setTextureUnit(0)

	// We always need a default shader.
	defaultShader = NewShader()

	glState.initialized = true

	callbackHandlers(window)
	loadAllVolatile()

	//have to set this after loadallvolatile() so we are sure the  default shader is loaded
	SetShader(nil)
}

func deInit(w *glfw.Window) {
	gl.DeleteTexture(glState.defaultTexture)
	glState.defaultTexture = gl.Texture{}
	glState.initialized = false
}

func callbackHandlers(window *glfw.Window) {
	window.SetCloseCallback(deInit)
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		SetViewport(0, 0, int32(width), int32(height))
	})
}

// Set the 'default' texture (id 0) as a repeating white pixel. Otherwise,
// texture2D calls inside a shader would return black when drawing graphics
// primitives, which would create the need to use different "passthrough"
// shaders for untextured primitives vs images.
func createDefaultTexture() {
	glState.defaultTexture = gl.CreateTexture()
	bindTexture(glState.defaultTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(gl.TEXTURE_2D, 0, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, []byte{255, 255, 255, 255})
}

// prepareDraw will upload all the transformations to the current shader
func prepareDraw(model *mgl32.Mat4) {
	if model == nil {
		model = &modelIdent
	}

	pmMat := glState.projectionStack.Peek().Mul4(glState.viewStack.Peek().Mul4(*model))

	glState.currentShader.SendMat4("ProjectionMat", glState.projectionStack.Peek())
	glState.currentShader.SendMat4("ViewMat", glState.viewStack.Peek())
	glState.currentShader.SendMat4("ModelMat", *model)
	glState.currentShader.SendMat4("PreMulMat", pmMat)
	if glState.currentCanvas != nil {
		glState.currentShader.SendFloat("ScreenSize", float32(screenWidth), float32(screenHeight), 1, 0)
	} else {
		glState.currentShader.SendFloat("ScreenSize", float32(screenWidth), float32(screenHeight), -1, float32(screenHeight))
	}
	glState.currentShader.SendFloat("PointSize", states.back().pointSize)
}

// useVertexAttribArrays will enable the vertex attrib array for the flags passed
// and if the flags were not passed it will disabled them. This make sure that only
// those attributes are enabled.
func useVertexAttribArrays(enabledAttribs ...gl.Attrib) {
	attribs := map[gl.Attrib]bool{
		shaderPos:           false,
		shaderTexCoord:      false,
		shaderColor:         false,
		shaderConstantColor: false,
	}

	for _, enabledAttrib := range enabledAttribs {
		attribs[enabledAttrib] = true
	}

	for attrib, enabled := range attribs {
		if enabled {
			gl.EnableVertexAttribArray(attrib)
		} else {
			gl.DisableVertexAttribArray(attrib)
			if attrib == shaderColor {
				gl.VertexAttrib4f(shaderColor, 1.0, 1.0, 1.0, 1.0)
			}
		}
	}
}

// setTextureUnit activates a texture unit
func setTextureUnit(textureunit int) error {
	if textureunit < 0 || textureunit >= len(glState.boundTextures) {
		return fmt.Errorf("invalid texture unit index (%v)", textureunit)
	}

	if textureunit != glState.curTextureUnit {
		gl.ActiveTexture(gl.Enum(gl.TEXTURE0 + uint32(textureunit)))
	}

	glState.curTextureUnit = textureunit
	return nil
}

// bindTexture will bind a texture to the current context if it isnt already bound
func bindTexture(texture gl.Texture) {
	if texture != glState.boundTextures[glState.curTextureUnit] {
		glState.boundTextures[glState.curTextureUnit] = texture
		gl.BindTexture(gl.TEXTURE_2D, texture)
	}
}

// bindTextureToUnit will bind a texture to a texture unit. If restorprev is true
// it will enable the current texture unit after completing
func bindTextureToUnit(texture gl.Texture, textureunit int, restoreprev bool) error {
	if texture != glState.boundTextures[textureunit] {
		oldtextureunit := glState.curTextureUnit
		if err := setTextureUnit(textureunit); err != nil {
			return err
		}
		glState.boundTextures[textureunit] = texture
		gl.BindTexture(gl.TEXTURE_2D, glState.boundTextures[textureunit])
		if restoreprev {
			return setTextureUnit(oldtextureunit)
		}
	}
	return nil
}

// HasFramebufferSRGB will return true if standard RGB color space is suporrted on
// this system. Only supported on non ES environments.
func HasFramebufferSRGB() bool {
	return glState.framebufferSRGBEnabled
}

// getDefaultFBO will return the framebuffer that was bound at startup.
func getDefaultFBO() gl.Framebuffer {
	return glState.defaultFBO
}

// deleteTexture will clean up the texture if it was bound before and also clean
// up the open gl data.
func deleteTexture(texture gl.Texture) {
	// glDeleteTextures binds texture 0 to all texture units the deleted texture
	// was bound to before deletion.
	for i, texid := range glState.boundTextures {
		if texid == texture {
			glState.boundTextures[i] = gl.Texture{}
		}
	}

	gl.DeleteTexture(texture)
}

// GetViewport will return the x, y, w, h of the opengl viewport. This is interfaced
// directly with opengl and used by the framework. Only use this if you know what
// you are doing
func GetViewport() []int32 {
	return glState.viewport
}

// SetViewportSize will set the viewport to 0, 0, w, h. This is interfaced
// directly with opengl and used by the framework. Only use this if you know what
// you are doing
func SetViewportSize(w, h int32) {
	SetViewport(0, 0, w, h)
}

// SetViewport will set the viewport to x, y, w, h. This is interfaced
// directly with opengl and used by the framework. Only use this if you know what
// you are doing
func SetViewport(x, y, w, h int32) {
	screenWidth = w
	screenHeight = h
	// Set the viewport to top-left corner.
	if glState.currentCanvas == nil {
		gl.Viewport(int(y), int(x), int(screenWidth), int(screenHeight))
		glState.viewport = []int32{y, x, screenWidth, screenHeight}
		glState.projectionStack.Load(mgl32.Ortho(float32(x), float32(screenWidth), float32(screenHeight), float32(y), -1, 1))
		if states.back().scissor {
			SetScissor(states.back().scissorBox[0], states.back().scissorBox[1], states.back().scissorBox[2], states.back().scissorBox[3])
		}
	}
}

// GetWidth will return the width of the rendering context.
func GetWidth() float32 {
	return float32(screenWidth)
}

// GetHeight will return the height of the rendering context.
func GetHeight() float32 {
	return float32(screenHeight)
}

// GetDimensions will return the width and height of screen
func GetDimensions() (float32, float32) {
	return float32(screenWidth), float32(screenHeight)
}

// Clear will clear everything already rendered to the screen and set is all to
// the r, g, b, a provided.
func Clear(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// Present is used at the end of the game loop to swap the frame buffers and display
// the next rendered frame. This is normally used by the game loop and should not
// be used unless rolling your own game loop.
func Present() {
	if !glState.initialized {
		return
	}

	// Make sure we don't have a canvas active.
	canvas := states.back().canvas
	SetCanvas(nil)
	// Restore the currently active canvas, if there is one.
	SetCanvas(canvas)

	// Cleanup after each loop
	cleanupVolatile()
}

// Origin will reset all translations and transformations back to defaults.
// This function is always used to reverse any previous calls to Rotate, Scale,
// Shear or Translate.
func Origin() {
	glState.viewStack.LoadIdent()
}

// Translate will translate the rendering origin to the point x, y.
// When this function is called with two numbers, dx, and dy, all the following
// drawing operations take effect as if their x and y coordinates were x+dx and y+dy.
// Scale and translate are not commutative operations, therefore, calling them in
// different orders will change the outcome. This change lasts until drawing completes
// or else a Pop reverts to a previous graphics state. Translating using whole
// numbers will prevent tearing/blurring of images and fonts draw after translating.
func Translate(x, y float32) {
	glState.viewStack.LeftMul(mgl32.Translate3D(x, y, 0))
}

// Rotate rotates the coordinate system in two dimensions. Calling this function
// affects all future drawing operations by rotating the coordinate system around
// the origin by the given amount of radians. This change lasts until drawing completes
func Rotate(angle float32) {
	glState.viewStack.LeftMul(mgl32.HomogRotate3DZ(angle))
}

// Scale scales the coordinate system in two dimensions. By default the coordinate system
/// in amore corresponds to the display pixels in horizontal and vertical directions
// one-to-one, and the x-axis increases towards the right while the y-axis increases
// downwards. Scaling the coordinate system changes this relation. After scaling by
// sx and sy, all coordinates are treated as if they were multiplied by sx and sy.
// Every result of a drawing operation is also correspondingly scaled, so scaling by
// (2, 2) for example would mean making everything twice as large in both x- and y-directions. Scaling by a negative value flips the coordinate system in the corresponding direction, which also means everything will be drawn flipped or upside down, or both. Scaling by zero is not a useful operation.
// Scale and translate are not commutative operations, therefore, calling them
// in different orders will change the outcome. Scaling lasts until drawing completes
func Scale(sx, sy float32) {
	glState.viewStack.LeftMul(mgl32.Scale3D(sx, sy, 1))
}

// Shear shears the coordinate system.
func Shear(kx, ky float32) {
	glState.viewStack.LeftMul(mgl32.ShearX3D(kx, ky))
}

// Push copies and pushes the current coordinate transformation to the transformation
// stack. This function is always used to prepare for a corresponding pop operation
// later. It stores the current coordinate transformation state into the transformation
// stack and keeps it active. Later changes to the transformation can be undone by
// using the pop operation, which returns the coordinate transform to the state
// it was in before calling push.
func Push() {
	glState.viewStack.Push()
	states.push(*states.back())
}

// Pop pops the current coordinate transformation from the transformation stack.
// This function is always used to reverse a previous push operation. It returns
// the current transformation state to what it was before the last preceding push.
func Pop() {
	glState.viewStack.Pop()
	states.pop()
}

// SetScissor Sets or disables scissor. The scissor limits the drawing area to a
// specified rectangle. This affects all graphics calls, including Clear.  The
// dimensions of the scissor is unaffected by graphical transformations
// (translate, scale, ...). if no arguments are given it will disable the scissor.
// if x, y, w, h are given it will enable the scissor
func SetScissor(x, y, width, height int32) {
	gl.Enable(gl.SCISSOR_TEST)
	if glState.currentCanvas != nil {
		gl.Scissor(x, y, width, height)
	} else {
		// With no Canvas active, we need to compensate for glScissor starting
		// from the lower left of the viewport instead of the top left.
		gl.Scissor(x, glState.viewport[3]-(y+height), width, height)
	}
	states.back().scissorBox = []int32{x, y, width, height}
	states.back().scissor = true
}

// ClearScissor will disable all set scissors.
func ClearScissor() {
	gl.Disable(gl.SCISSOR_TEST)
	states.back().scissor = false
}

// Stencil operates like stencil but with access to change the stencil action,
// value, and keepvalues.
//
// action: How to modify any stencil values of pixels that are touched by what's
// drawn in the stencil function.
//
// value: The new stencil value to use for pixels if the "replace" stencil action
// is used. Has no effect with other stencil actions. Must be between 0 and 1.
//
// keepvalues: True to preserve old stencil values of pixels, false to re-set
// every pixel's stencil value to 0 before executing the stencil function. Clear
// will also re-set all stencil values.
func Stencil(stencilFunc func(), action StencilAction, value int32, keepvalues bool) {
	// StencilReplace, 1, false
	glState.writingToStencil = true
	if !keepvalues {
		gl.Clear(gl.STENCIL_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	}
	if glState.currentCanvas != nil {
		glState.currentCanvas.checkCreateStencil()
	}
	gl.Enable(gl.STENCIL_TEST)
	gl.ColorMask(false, false, false, false)
	gl.StencilFunc(gl.ALWAYS, int(value), 0xFF)
	gl.StencilOp(gl.KEEP, gl.KEEP, gl.Enum(action))

	stencilFunc()

	glState.writingToStencil = false
	mask := states.back().colorMask
	SetColorMask(mask.r, mask.g, mask.b, mask.a)
	SetStencilTest(states.back().stencilCompare, states.back().stencilTestValue)
}

// SetStencilTest configures or disables stencil testing. When stencil testing is
// enabled, the geometry of everything that is drawn afterward will be clipped/stencilled
// out based on a comparison between the arguments of this function and the stencil
// value of each pixel that the geometry touches. The stencil values of pixels are
// affected via Stencil/StencilEXT.
func SetStencilTest(compare CompareMode, value int32) {
	if glState.writingToStencil {
		return
	}

	states.back().stencilCompare = compare
	states.back().stencilTestValue = value
	if compare == CompareAlways {
		gl.Disable(gl.STENCIL_TEST)
		return
	}

	if glState.currentCanvas != nil {
		glState.currentCanvas.checkCreateStencil()
	}

	gl.Enable(gl.STENCIL_TEST)
	gl.StencilFunc(gl.Enum(compare), int(value), 0xFF)
	gl.StencilOp(gl.KEEP, gl.KEEP, gl.REPLACE)
}

// ClearStencilTest stops the stencil test from operating
func ClearStencilTest() {
	SetStencilTest(CompareAlways, 0)
}

// GetStencilTest will return the current compare mode and the stencil test value.
func GetStencilTest() (CompareMode, int32) {
	return states.back().stencilCompare, states.back().stencilTestValue
}

// GetScissor will return the current scissor rectangle
func GetScissor() (x, y, w, h int32) {
	return states.back().scissorBox[0], states.back().scissorBox[1], states.back().scissorBox[2], states.back().scissorBox[3]
}

// SetColorMask will set a mask for each r, g, b, and alpha component.
func SetColorMask(r, g, b, a bool) {
	gl.ColorMask(r, g, b, a)
	states.back().colorMask = ColorMask{r, g, b, a}
}

// GetColorMask will return the current color mask
func GetColorMask() (bool, bool, bool, bool) {
	mask := states.back().colorMask
	return mask.r, mask.g, mask.b, mask.a
}

// SetLineWidth changes the width in pixels that the lines will render when using
// Line or PolyLine
func SetLineWidth(width float32) {
	states.back().lineWidth = width
}

// SetLineJoin will change how each line joins. options are None, Bevel or Miter.
func SetLineJoin(join string) {
	states.back().lineJoin = join
}

// GetLineWidth will return the current line width. Default line width is 1
func GetLineWidth() float32 {
	return states.back().lineWidth
}

// GetLineJoin will return the current line join. Default line join is miter.
func GetLineJoin() string {
	return states.back().lineJoin
}

// SetPointSize will set the size of points drawn by Point
func SetPointSize(size float32) {
	states.back().pointSize = size
}

// GetPointSize will return the current point size
func GetPointSize() float32 {
	return states.back().pointSize
}

// SetShader sets or resets a Shader as the current pixel effect or vertex shaders.
// All drawing operations until the next SetShader will be drawn using the Shader
// object specified.
func SetShader(shader *Shader) {
	if shader == nil {
		states.back().shader = defaultShader
	} else {
		states.back().shader = shader
	}
	states.back().shader.attach(false)
}

// SetBackgroundColor sets the background color.
func SetBackgroundColor(vals ...float32) {
	states.back().backgroundColor = vals
}

// GetBackgroundColor gets the background color.
func GetBackgroundColor() []float32 {
	return states.back().backgroundColor
}

// SetColor will sets the color used for drawing.
func SetColor(r, g, b, a float32) {
	states.back().color = []float32{r, g, b, a}
	gl.VertexAttrib4f(shaderConstantColor, r, g, b, a)
}

// GetColor returns the current drawing color.
func GetColor() []float32 {
	return states.back().color
}

// SetFont will set a font for the next print call
func SetFont(font *Font) {
	if font == nil {
		font = defaultFont
	}
	states.back().font = font
}

// GetFont will return the currenly bound font or the frameworks font if none has
// be bound.
func GetFont() *Font {
	return states.back().font
}

// SetCanvas will set the render target to a specified Canvas. All drawing operations
// until the next SetCanvas call will be redirected to the Canvas and not shown
// on the screen. Call with a no params to enable drawing to screen again.
func SetCanvas(canvas *Canvas) error {
	states.back().canvas = canvas

	if canvas != nil {
		return canvas.startGrab()
	}

	if glState.currentCanvas != nil {
		glState.currentCanvas.stopGrab(false)
		glState.currentCanvas = nil
	}
	return nil
}

// GetCanvas returns the currently bound canvases
func GetCanvas() *Canvas {
	return states.back().canvas
}

// SetBlendMode sets the blending mode. Blending modes are different ways to do
// color blending. See BlendMode constants to see how they operate.
func SetBlendMode(mode string) {
	fn := gl.FUNC_ADD
	srcRGB := gl.ONE
	srcA := gl.ONE
	dstRGB := gl.ZERO
	dstA := gl.ZERO

	switch mode {
	case "multiplicative":
		srcRGB = gl.DST_COLOR
		srcA = gl.DST_COLOR
		dstRGB = gl.ZERO
		dstA = gl.ZERO
	case "premultiplied":
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_ALPHA
		dstA = gl.ONE_MINUS_SRC_ALPHA
	case "subtractive":
		fn = gl.FUNC_REVERSE_SUBTRACT
	case "additive":
		srcRGB = gl.SRC_ALPHA
		srcA = gl.SRC_ALPHA
		dstRGB = gl.ONE
		dstA = gl.ONE
	case "screen":
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_COLOR
		dstA = gl.ONE_MINUS_SRC_COLOR
		break
	case "replace":
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ZERO
		dstA = gl.ZERO
	case "alpha":
		fallthrough
	default:
		srcRGB = gl.SRC_ALPHA
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_ALPHA
		dstA = gl.ONE_MINUS_SRC_ALPHA
	}

	gl.BlendEquation(gl.Enum(fn))
	gl.BlendFuncSeparate(gl.Enum(srcRGB), gl.Enum(dstRGB), gl.Enum(srcA), gl.Enum(dstA))
	states.back().blendMode = mode
}
