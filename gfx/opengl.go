package gfx

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl32/matstack"
	"github.com/tanema/amore/font"

	"github.com/tanema/amore/gfx/gl"
	"github.com/tanema/amore/window"
)

var (
	openglVersion          string
	openglVendor           string
	maxAnisotropy          float32
	maxTextureSize         int32
	maxRenderTargets       int32
	maxRenderbufferSamples int32
	maxTextureUnits        int32
	screenWidth            = int32(0)
	screenHeight           = int32(0)
	modelIdent             = mgl32.Ident4()
	defaultShader          *Shader

	glState = openglState{
		viewport: make([]int32, 4),
	}
	states = displayStateStack{newDisplayState()}
)

// InitContext will initiate the opengl context with a viewport in the size of
// w x h. This is generally called from the game loop and wont generally need to
// be called unless you are rolling your own game loop.
func InitContext(w, h int32) {
	if glState.initialized {
		return
	}

	// Okay, setup OpenGL.
	gl.ContextWatcher.OnMakeCurrent(nil)

	//Get system info
	openglVersion = gl.GetString(gl.VERSION)
	openglVendor = gl.GetString(gl.VENDOR)
	glState.defaultFBO = gl.GetBoundFramebuffer()
	gl.GetIntegerv(gl.VIEWPORT, glState.viewport)
	// And the current scissor - but we need to compensate for GL scissors
	// starting at the bottom left instead of top left.
	gl.GetIntegerv(gl.SCISSOR_BOX, states.back().scissorBox)
	states.back().scissorBox[1] = glState.viewport[3] - (states.back().scissorBox[1] + states.back().scissorBox[3])

	initMaxValues() //check shim code

	glcolor := []float32{1.0, 1.0, 1.0, 1.0}
	gl.VertexAttrib4fv(attribColor, glcolor)
	gl.VertexAttrib4fv(attribConstantColor, glcolor)
	useVertexAttribArrays(0)

	// Enable blending
	gl.Enable(gl.BLEND)
	SetBlendMode(BlendModeAlpha)
	// Auto-generated mipmaps should be the best quality possible
	gl.Hint(gl.GENERATE_MIPMAP_HINT, gl.NICEST)
	// Make sure antialiasing works when set elsewhere
	enableMultisample() //check shim code
	// Set pixel row alignment
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	//default matricies
	glState.projectionStack = matstack.NewMatStack()
	glState.viewStack = matstack.NewMatStack() //stacks are initialized with ident matricies on top

	SetViewportSize(w, h)
	SetBackgroundColor(0, 0, 0, 255)

	glState.boundTextures = make([]gl.Texture, maxTextureUnits)
	curgltextureunit := gl.GetInteger(gl.ACTIVE_TEXTURE)
	glState.curTextureUnit = int(curgltextureunit - gl.TEXTURE0)
	// Retrieve currently bound textures for each texture unit.
	for i := 0; i < len(glState.boundTextures); i++ {
		gl.ActiveTexture(gl.Enum(gl.TEXTURE0 + uint32(i)))
		glState.boundTextures[i] = gl.Texture{Value: uint32(gl.GetInteger(gl.TEXTURE_BINDING_2D))}
	}
	gl.ActiveTexture(gl.Enum(curgltextureunit))
	createDefaultTexture()
	setTextureUnit(0)

	// We always need a default shader.
	defaultShader = NewShader()

	glState.initialized = true

	loadAllVolatile()

	//have to set this after loadallvolatile() so we are sure the  default shader is loaded
	SetShader(nil)
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

	pix := []byte{255, 255, 255, 255}
	gl.TexImage2D(gl.TEXTURE_2D, 0, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pix))
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
func useVertexAttribArrays(arraybits uint32) {
	diff := arraybits ^ glState.enabledAttribArrays

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

	glState.enabledAttribArrays = arraybits

	// glDisableVertexAttribArray will make the constant value for a vertex
	// attribute undefined. We rely on the per-vertex color attribute being
	// white when no per-vertex color is used, so we set it here.
	// FIXME: Is there a better place to do this?
	if (diff&attribFlagColor) > 0 && (arraybits&attribFlagColor) == 0 {
		gl.VertexAttrib4f(attribColor, 1.0, 1.0, 1.0, 1.0)
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
			glState.boundTextures[i] = gl.Texture{Value: 0}
		}
	}

	gl.DeleteTexture(texture)
}

// DeInit will do the clean up for the context.
func DeInit() {
	gl.DeleteTexture(glState.defaultTexture)
	glState.defaultTexture = gl.Texture{}
	glState.initialized = false
	gl.ContextWatcher.OnDetach()
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
	gl.Viewport(int(y), int(x), int(screenWidth), int(screenHeight))
	glState.viewport = []int32{y, x, screenWidth, screenHeight}
	glState.projectionStack.Load(mgl32.Ortho(float32(x), float32(screenWidth), float32(screenHeight), float32(y), -1, 1))
	SetScissor(states.back().scissorBox[0], states.back().scissorBox[1], states.back().scissorBox[2], states.back().scissorBox[3])
}

// GetWidth will return the width of the rendering context.
func GetWidth() float32 {
	return float32(screenWidth)
}

// GetHeight will return the height of the rendering context.
func GetHeight() float32 {
	return float32(screenHeight)
}

// Clear will clear everything already rendered to the screen and set is all to
// the r, g, b, a provided.
func Clear(r, g, b, a float32) {
	ClearC(NewColor(r, g, b, a))
}

// ClearC will clear everything already rendered to the screen and set is all to
// the *Color provided.
func ClearC(c *Color) {
	gl.ClearColor(c[0], c[1], c[2], c[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// Present is used at the end of the game loop to swap the frame buffers and display
// the next rendered frame. This is normally used by the game loop and should not
// be used unless rolling your own game loop.
func Present() {
	if !IsActive() {
		return
	}

	// Make sure we don't have a canvas active.
	canvases := states.back().canvases
	SetCanvas()
	window.SwapBuffers()
	// Restore the currently active canvas, if there is one.
	SetCanvas(canvases...)
}

// IsActive will return true of the context has been initialized and the window
// is open.
func IsActive() bool {
	// The graphics module is only completely 'active' if there's a window, a
	// context, and the active variable is set.
	return glState.active && IsCreated() && window.IsOpen()
}

// SetActive will enable or disable the rendering of the the game. Mainly this is
// used by the framework to disable rendering when not in view.
func SetActive(enable bool) {
	// Make sure all pending OpenGL commands have fully executed before
	// returning, when going from active to inactive. This is required on iOS.
	if IsCreated() && glState.active && !enable {
		gl.Finish()
	}

	glState.active = enable
}

// IsCreated checks if the opengl context has be initialized.
func IsCreated() bool {
	return glState.initialized
}

// Origin will reset all translations and transformations back to defaults.
// This function is always used to reverse any previous calls to Rotate, Scale,
// Shear or Translate.
func Origin() {
	glState.viewStack.LoadIdent()
	states.back().pixelSize = 1.0
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

	glState.viewStack.LeftMul(mgl32.Scale3D(sx, sy, 1))

	states.back().pixelSize *= (2.0 / (mgl32.Abs(sx) + mgl32.Abs(sy)))
}

// Shear shears the coordinate system.
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
func SetScissor(args ...int32) {
	if args == nil {
		gl.Disable(gl.SCISSOR_TEST)
		states.back().scissor = false
	} else if len(args) == 4 {
		x, y, width, height := args[0], args[1], args[2], args[3]
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
	} else {
		panic("incorrect number of arguments to setscissor")
	}
}

// IntersectScissor sets the scissor to the rectangle created by the intersection
// of the specified rectangle with the existing scissor. If no scissor is active
// yet, it behaves like SetScissor. The scissor limits the drawing area to a
// specified rectangle. This affects all graphics calls, including Clear. The
// dimensions of the scissor is unaffected by graphical transformations (translate, scale, ...).
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

// ClearScissor will disable all set scissors.
func ClearScissor() {
	gl.Disable(gl.SCISSOR_TEST)
	states.back().scissor = false
}

// Stencil draws geometry as a stencil. The geometry drawn by the supplied function
// sets invisible stencil values of pixels, instead of setting pixel colors.
// The stencil values of pixels can act like a mask / stencil - SetStencilTest can
// be used afterward to determine how further rendering is affected by the stencil
// values in each pixel. Each Canvas has its own per-pixel stencil values. Stencil
// values are within the range of [0, 255]. This stencil has the defaults of
// StencilAction: STENCIL_REPLACE, value: 1, keepvalues: true
func Stencil(stencilFunc func()) {
	StencilExt(stencilFunc, StencilReplace, 1, false)
}

// StencilExt operates like stencil but with access to change the stencil action,
// value, and keepvalues.
//
// action: How to modify any stencil values of pixels that are touched by what's
// drawn in the stencil function.
//
// value: The new stencil value to use for pixels if the "replace" stencil action
// is used. Has no effect with other stencil actions. Must be between 0 and 255.
//
// keepvalues: True to preserve old stencil values of pixels, false to re-set
// every pixel's stencil value to 0 before executing the stencil function. Clear
// will also re-set all stencil values.
func StencilExt(stencilFunc func(), action StencilAction, value int32, keepvalues bool) {
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
	SetColorMask(states.back().colorMask)
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
func SetColorMask(mask ColorMask) {
	gl.ColorMask(mask.r, mask.g, mask.b, mask.a)
	states.back().colorMask = mask
}

// GetColorMask will return the current color mask
func GetColorMask() ColorMask {
	return states.back().colorMask
}

// SetLineWidth changes the width in pixels that the lines will render when using
// Line or PolyLine
func SetLineWidth(width float32) {
	states.back().lineWidth = width
}

// SetLineStyle will set the line style either smooth (overdraw) or rough.
func SetLineStyle(style LineStyle) {
	states.back().lineStyle = style
}

// SetLineJoin will change how each line joins. options are None, Bevel or Miter.
func SetLineJoin(join LineJoin) {
	states.back().lineJoin = join
}

// GetLineWidth will return the current line width. Default line width is 1
func GetLineWidth() float32 {
	return states.back().lineWidth
}

// GetLineStyle will return the current line style. Default line style is smooth
func GetLineStyle() LineStyle {
	return states.back().lineStyle
}

// GetLineJoin will return the current line join. Default line join is miter.
func GetLineJoin() LineJoin {
	return states.back().lineJoin
}

// SetPointSize will set the size of points drawn by Point
func SetPointSize(size float32) {
	states.back().pointSize = size
	gl.PointSize(size)
}

// GetPointSize will return the current point size
func GetPointSize() float32 {
	return states.back().pointSize
}

// IsWireframe will return true if wirefame is current enabled. Wireframe is only
// available in non ES enviroments.
func IsWireframe() bool {
	return states.back().wireframe
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
func SetBackgroundColor(r, g, b, a float32) {
	states.back().backgroundColor = NewColor(r, g, b, a)
}

// SetBackgroundColorC sets the background color.
func SetBackgroundColorC(c *Color) {
	states.back().backgroundColor = c
}

// GetBackgroundColor gets the background color.
func GetBackgroundColor() (r, g, b, a float32) {
	bc := states.back().backgroundColor
	return bc[0], bc[1], bc[2], bc[3]
}

// GetBackgroundColorC gets the background color.
func GetBackgroundColorC() *Color {
	return states.back().backgroundColor
}

// setColor translates r,g,b,a to Color
func setColor(r, g, b, a float32) {
	SetColorC(&Color{r, g, b, a})
}

// SetColor will sets the color used for drawing.
func SetColor(r, g, b, a float32) {
	setColor(r/255.0, g/255.0, b/255.0, a/255.0)
}

// SetColorC will sets the color used for drawing.
func SetColorC(c *Color) {
	states.back().color = c

	gl.VertexAttrib4f(attribConstantColor, c[0], c[1], c[2], c[3])
}

// GetColor returns the current drawing color.
func GetColor() *Color {
	return states.back().color
}

// SetFont will set a font for the next print call
func SetFont(font *Font) {
	states.back().font = font
}

// GetFont will return the currenly bound font or the frameworks font if none has
// be bound.
func GetFont() *Font {
	//if no font set, use default font
	if states.back().font == nil {
		f, err := font.Bold(20)
		if err == nil {
			SetFont(NewFont(f))
		}
	}
	return states.back().font
}

// SetCanvas will set the render target to a specified Canvas. All drawing operations
// until the next SetCanvas call will be redirected to the Canvas and not shown
// on the screen. Call with a no params to enable drawing to screen again.
func SetCanvas(canvases ...*Canvas) error {
	if canvases == nil || len(canvases) <= 0 {
		states.back().canvases = nil
		if glState.currentCanvas != nil {
			glState.currentCanvas.stopGrab(false)
			glState.currentCanvas = nil
		}
	} else {
		states.back().canvases = canvases
		return states.back().canvases[0].startGrab(canvases[1:]...)
	}
	return nil
}

// GetCanvas returns the currently bound canvases
func GetCanvas() []*Canvas {
	return states.back().canvases
}

// SetBlendMode sets the blending mode. Blending modes are different ways to do
// color blending. See BlendMode constants to see how they operate.
func SetBlendMode(mode BlendMode) {
	fn := gl.FUNC_ADD
	srcRGB := gl.ONE
	srcA := gl.ONE
	dstRGB := gl.ZERO
	dstA := gl.ZERO

	switch mode {
	case BlendModeAlpha:
		srcRGB = gl.SRC_ALPHA
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_ALPHA
		dstA = gl.ONE_MINUS_SRC_ALPHA
	case BlendModeMultiplicative:
		srcRGB = gl.DST_COLOR
		srcA = gl.DST_COLOR
		dstRGB = gl.ZERO
		dstA = gl.ZERO
	case BlendModePremultiplied:
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_ALPHA
		dstA = gl.ONE_MINUS_SRC_ALPHA
	case BlendModeSubtractive:
		fn = gl.FUNC_REVERSE_SUBTRACT
	case BlendModeAdditive:
		srcRGB = gl.SRC_ALPHA
		srcA = gl.SRC_ALPHA
		dstRGB = gl.ONE
		dstA = gl.ONE
	case BlendModeScreen:
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ONE_MINUS_SRC_COLOR
		dstA = gl.ONE_MINUS_SRC_COLOR
		break
	case BlendModeReplace:
		srcRGB = gl.ONE
		srcA = gl.ONE
		dstRGB = gl.ZERO
		dstA = gl.ZERO
	}

	gl.BlendEquation(gl.Enum(fn))
	gl.BlendFuncSeparate(gl.Enum(srcRGB), gl.Enum(dstRGB), gl.Enum(srcA), gl.Enum(dstA))
	states.back().blendMode = mode
}

// SetDefaultFilter sets the default scaling filters used with Images, Canvases, and Fonts.
func SetDefaultFilter(min, mag FilterMode, anisotropy float32) {
	states.back().defaultFilter = Filter{
		min:        min,
		mag:        mag,
		anisotropy: float32(math.Min(math.Max(float64(anisotropy), 1.0), float64(maxAnisotropy))),
	}
}

// SetDefaultFilterF sets the default scaling filters used with Images, Canvases, and Fonts.
func SetDefaultFilterF(f Filter) {
	states.back().defaultFilter = f
}

// GetDefaultFilter returns the scaling filters used with Images, Canvases, and Fonts.
func GetDefaultFilter() Filter {
	return states.back().defaultFilter
}

// SetDefaultMipmapFilter will set the default mipmap filters for when new images
// are created
func SetDefaultMipmapFilter(filter FilterMode, sharpness float32) {
	states.back().defaultMipmapFilter = filter
	states.back().defaultMipmapSharpness = sharpness
}

// GetDefaultMipmapFilter will return the default filtermode and sharpness for
// all new images
func GetDefaultMipmapFilter() (filter FilterMode, sharpness float32) {
	return states.back().defaultMipmapFilter, states.back().defaultMipmapSharpness
}

// GetDebugInfo gets information on what version of opengl we are using.
func GetDebugInfo() string {
	return "OpenGL " + openglVersion + "\nVendor: " + openglVendor
}
