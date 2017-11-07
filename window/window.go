// Package window creates and manages the window and gl context.
package window

import (
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/tanema/amore/window/surface"
)

var (
	currentWindow *window
	created       = false
	initError     error
)

// MessageBoxType specifies the types of Message box to create
type MessageBoxType uint32

// MessageBoxTypes for opening an alert dialog
const (
	MessageBoxError   MessageBoxType = sdl.MESSAGEBOX_ERROR
	MessageBoxWarning MessageBoxType = sdl.MESSAGEBOX_WARNING
	MessageBoxInfo    MessageBoxType = sdl.MESSAGEBOX_INFORMATION
)

// Window contains the created window information
type window struct {
	SDLWindow   *sdl.Window
	context     sdl.GLContext
	shouldClose bool
	Config      *windowConfig
	open        bool
}

// NewWindow will create a new sdl window with a gl context. It will also destroy
// an old one if there is one already. Errors returned come straight from SDL so
// errors will be indicitive of SDL errors
func NewWindow() (*window, error) {
	if currentWindow != nil || initError != nil {
		return currentWindow, initError
	}

	var config *windowConfig
	var err error

	if err = sdl.InitSubSystem(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	if config, err = loadConfig(); err != nil {
		panic(err)
	}

	config.Minwidth = int32(math.Max(float64(config.Minwidth), 1.0))
	config.Minheight = int32(math.Max(float64(config.Minheight), 1.0))
	config.Display = int(math.Min(math.Max(float64(config.Display), 0.0), float64(GetDisplayCount()-1)))

	if config.Width == 0 || config.Height == 0 {
		var mode sdl.DisplayMode
		sdl.GetDesktopDisplayMode(config.Display, &mode)
		config.Width = mode.W
		config.Height = mode.H
	}

	sdlflags := uint32(sdl.WINDOW_OPENGL)

	if config.Fullscreen {
		if config.Fstype == "desktop" {
			sdlflags |= sdl.WINDOW_FULLSCREEN_DESKTOP
		} else {
			sdlflags |= sdl.WINDOW_FULLSCREEN

			mode := sdl.DisplayMode{W: int32(config.Width), H: int32(config.Height)}

			// Fullscreen window creation will bug out if no mode can be used.
			if _, err := sdl.GetClosestDisplayMode(config.Display, &mode, &mode); err != nil {
				// GetClosestDisplayMode will fail if we request a size larger
				// than the largest available display mode, so we'll try to use
				// the largest (first) mode in that case.
				if err := sdl.GetDisplayMode(config.Display, 0, &mode); err != nil {
					return nil, err
				}
			}

			config.Width = mode.W
			config.Height = mode.H
		}
	}

	if config.Resizable {
		sdlflags |= sdl.WINDOW_RESIZABLE
	}

	if config.Borderless {
		sdlflags |= sdl.WINDOW_BORDERLESS
	}

	if config.Highdpi {
		sdlflags |= sdl.WINDOW_ALLOW_HIGHDPI
	}

	if config.Fullscreen {
		// The position needs to be in the global coordinate space.
		var displaybounds sdl.Rect
		sdl.GetDisplayBounds(config.Display, &displaybounds)
		config.X += displaybounds.X
		config.Y += displaybounds.Y
	} else {
		if config.Centered {
			config.X = sdl.WINDOWPOS_CENTERED
			config.Y = sdl.WINDOWPOS_CENTERED
		} else {
			config.X = sdl.WINDOWPOS_UNDEFINED
			config.Y = sdl.WINDOWPOS_UNDEFINED
		}
	}

	if currentWindow != nil {
		Destroy()
	}

	created = false
	newWindow, err := createWindowAndContext(config, sdlflags)
	if err != nil {
		return nil, err
	}
	created = true

	if newWindow.Config.Icon != "" {
		SetIcon(config.Icon)
	}

	SetMouseGrab(false)
	SetMinimumSize(config.Minwidth, config.Minheight)
	SetTitle(config.Title)

	if config.Centered && !config.Fullscreen {
		SetPosition(config.X, config.Y)
	}

	getCurrent().SDLWindow.Raise()

	if config.Vsync {
		sdl.GL_SetSwapInterval(1)
	} else {
		sdl.GL_SetSwapInterval(0)
	}

	newWindow.open = true
	return newWindow, nil
}

// createWindowAndContext is the actual interface with SDL to create window and context
func createWindowAndContext(config *windowConfig, windowflags uint32) (*window, error) {
	setGLFramebufferAttributes(config.Msaa, config.Srgb)
	_, debug := os.LookupEnv("AMORE_DEBUG")
	setGLContextAttributes(2, 1, debug)

	newWindow, err := sdl.CreateWindow(config.Title, int(config.X), int(config.Y), int(config.Width), int(config.Height), windowflags)
	if err != nil {
		panic(err)
	}

	context, err := sdl.GL_CreateContext(newWindow)
	if err != nil {
		panic(err)
	}

	currentWindow = &window{
		SDLWindow:   newWindow,
		context:     context,
		shouldClose: false,
		Config:      config,
	}
	return currentWindow, nil
}

// Set hints on the sdl framebuffer
func setGLFramebufferAttributes(msaa int, sRGB bool) {
	// Set GL window / framebuffer attributes.
	sdl.GL_SetAttribute(sdl.GL_RED_SIZE, 8)
	sdl.GL_SetAttribute(sdl.GL_GREEN_SIZE, 8)
	sdl.GL_SetAttribute(sdl.GL_BLUE_SIZE, 8)
	sdl.GL_SetAttribute(sdl.GL_ALPHA_SIZE, 8)
	sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GL_SetAttribute(sdl.GL_STENCIL_SIZE, 1)
	sdl.GL_SetAttribute(sdl.GL_RETAINED_BACKING, 0)

	if msaa > 0 {
		sdl.GL_SetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
		sdl.GL_SetAttribute(sdl.GL_MULTISAMPLESAMPLES, msaa)
	} else {
		sdl.GL_SetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 0)
		sdl.GL_SetAttribute(sdl.GL_MULTISAMPLESAMPLES, 0)
	}
}

// Set hints on the context
func setGLContextAttributes(versionMajor, versionMinor int, debug bool) {
	var profilemask, contextflags int

	if debug {
		profilemask = profilemask | sdl.GL_CONTEXT_PROFILE_COMPATIBILITY
		contextflags = contextflags | sdl.GL_CONTEXT_DEBUG_FLAG
	}

	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, versionMajor)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, versionMinor)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, profilemask)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_FLAGS, contextflags)
}

// getCurrent is a singleton method. It makes sure there is a window and if there
// is not it will create it. It also records if there were any errors for later
// return from NewWindow
func getCurrent() *window {
	if currentWindow == nil {
		currentWindow, initError = NewWindow()
	}

	return currentWindow
}

// GetDisplayCount gets the count of displays
func GetDisplayCount() int {
	num, _ := sdl.GetNumVideoDisplays()
	return num
}

// GetDisplayName gets the display name given an index
func GetDisplayName(displayindex int) string {
	return sdl.GetDisplayName(displayindex)
}

// GetDesktopDimensions gets the dimensions of a display for the given index
func GetDesktopDimensions(displayindex int) (int32, int32) {
	var width, height int32
	if displayindex >= 0 && displayindex < GetDisplayCount() {
		var mode sdl.DisplayMode
		sdl.GetDesktopDisplayMode(displayindex, &mode)
		width = mode.W
		height = mode.H
	}
	return width, height
}

// onSizeChanged is called from events and records the change in viewport
func onSizeChanged(width, height int32) {
	win := getCurrent()
	win.Config.Width = width
	win.Config.Height = height
	win.Config.PixelWidth, win.Config.PixelHeight = GetDrawableSize()
}

// GetDrawableSize gets the size of pixels in view
func GetDrawableSize() (int32, int32) {
	w, h := sdl.GL_GetDrawableSize(getCurrent().SDLWindow)
	return int32(w), int32(h)
}

// GetFullscreenSizes gets the size of the screen for the display index
func GetFullscreenSizes(displayindex int) [][]int32 {
	var sizes [][]int32
	modes, _ := sdl.GetNumDisplayModes(displayindex)
	for i := 0; i < modes; i++ {
		var mode sdl.DisplayMode
		sdl.GetDisplayMode(displayindex, i, &mode)
		sizes = append(sizes, []int32{mode.W, mode.H})
	}
	return sizes
}

// SetTitle sets the title of the window displayed on the display bar and task bar
func SetTitle(title string) {
	win := getCurrent()
	win.SDLWindow.SetTitle(title)
	win.Config.Title = title
}

// GetTitle gets the current title of the window
func GetTitle() string {
	return getCurrent().Config.Title
}

// SetIcon sets the small icon on the display bar as well as the large icon on the
// task bar and icon on the dock in osx
func SetIcon(path string) error {
	win := getCurrent()
	win.Config.Icon = path
	newSurface, err := surface.Load(path)
	if err != nil {
		return err
	}
	win.SDLWindow.SetIcon(newSurface)
	newSurface.Free()
	return nil
}

// GetIcon returns the path to the currently set icon
func GetIcon() string {
	return getCurrent().Config.Icon
}

// Minimize hides the program in the task bar
func Minimize() {
	getCurrent().SDLWindow.Minimize()
}

// Maximize unhides the program and sets it to max size
func Maximize() {
	getCurrent().SDLWindow.Maximize()
}

// ShouldClose returns if the window has beed told to close and is in the process
// of shutting down
func ShouldClose() bool {
	return getCurrent().shouldClose
}

// Close prepares the window to shut down gracefully
func Close(shouldClose bool) {
	getCurrent().shouldClose = shouldClose
}

// SwapBuffers swaps the current frame buffer in the window. This is generally
// only used by the engine but you could use it to roll your own game loop.
func SwapBuffers() {
	sdl.GL_SwapWindow(getCurrent().SDLWindow)
}

// ToPixelCoords translates window coords to pixel coords for pixel perfect
// operations
func ToPixelCoords(x, y float32) (float32, float32) {
	config := getCurrent().Config
	newX := x * (float32(config.PixelWidth) / float32(config.Width))
	newY := y * (float32(config.PixelHeight) / float32(config.Height))
	return newX, newY
}

// PixelToWindowCoords translates pixel coords to window coords
func PixelToWindowCoords(x, y float32) (float32, float32) {
	config := getCurrent().Config
	newX := x * (float32(config.Width) / float32(config.PixelWidth))
	newY := y * (float32(config.Height) / float32(config.PixelHeight))
	return newX, newY
}

// GetMousePosition return the current position of the mouse
func GetMousePosition() (float32, float32) {
	getCurrent() // must call getCurrent to ensure there is a window
	mx, my, _ := sdl.GetMouseState()
	return ToPixelCoords(float32(mx), float32(my))
}

// SetMousePosition warps the mouse position to a new position x, y in the window
func SetMousePosition(x, y float32) {
	wx, wy := PixelToWindowCoords(x, y)
	getCurrent().SDLWindow.WarpMouseInWindow(int(wx), int(wy))
}

// IsMouseGrabbed returns true if pointer lock is enabled and false otherwise
func IsMouseGrabbed() bool {
	return getCurrent().SDLWindow.GetGrab() != false
}

// SetMouseGrab enables or disables pointer lock on the window
func SetMouseGrab(grabbed bool) {
	getCurrent().SDLWindow.SetGrab(grabbed)
}

// IsVisible returns true if the window is not minimized
func IsVisible() bool {
	return (getCurrent().SDLWindow.GetFlags() & sdl.WINDOW_SHOWN) != 0
}

// SetMouseVisible will hide the mouse if passed false, and show the mouse cursor
// if passed true
func SetMouseVisible(visible bool) {
	getCurrent() // must call getCurrent to ensure there is a window
	if visible {
		sdl.ShowCursor(sdl.ENABLE)
	} else {
		sdl.ShowCursor(sdl.DISABLE)
	}
}

// GetMouseVisible will return true if the mouse is not hidden and will return
// false otherwise
func GetMouseVisible() bool {
	getCurrent() // must call getCurrent to ensure there is a window
	return sdl.ShowCursor(sdl.QUERY) == sdl.ENABLE
}

// GetPixelDimensions will return the size of the window in pixels
func GetPixelDimensions() (int32, int32) {
	win := getCurrent()
	return win.Config.PixelWidth, win.Config.PixelHeight
}

// GetPixelScale will return ratio of pixels in the window
func GetPixelScale() float32 {
	win := getCurrent()
	return float32(win.Config.PixelHeight) / float32(win.Config.Height)
}

// ToPixels will convert window unit to pixel unit
func ToPixels(x float32) float32 {
	return x * GetPixelScale()
}

// ToPixelsPoint will convert window point to pixel point
func ToPixelsPoint(x, y float32) (float32, float32) {
	scale := GetPixelScale()
	return x * scale, y * scale
}

// FromPixels will convert from pixel units to window units
func FromPixels(x float32) float32 {
	return x / GetPixelScale()
}

// FromPixelsPoint will convert a point from pixels to window units
func FromPixelsPoint(x, y float32) (float32, float32) {
	scale := GetPixelScale()
	return x / scale, y / scale
}

// SetMinimumSize sets the minimim size on the window when resizable
func SetMinimumSize(w, h int32) {
	win := getCurrent()
	win.Config.Minwidth = w
	win.Config.Minheight = h
	win.SDLWindow.SetMinimumSize(int(w), int(h))
}

// SetPosition will set the window position on screen
func SetPosition(x, y int32) {
	win := getCurrent()
	win.Config.X = x
	win.Config.Y = y
	win.SDLWindow.SetPosition(int(x), int(y))
}

// GetPosition will return the position of the window on the screen
func GetPosition() (int, int) {
	return getCurrent().SDLWindow.GetPosition()
}

// HasFocus will return true if the user has the program focused and return false
// if the user is focused on another program
func HasFocus() bool {
	return sdl.GetKeyboardFocus() == getCurrent().SDLWindow
}

// RequestAttention makes the application flash for attention to the user. If
// continuous is true, it will continue to do so until the user focuses the program.
func RequestAttention(continuous bool) {
	getCurrent() // must call getCurrent to ensure there is a window
	if HasFocus() {
		return
	}
	requestAttention(continuous)
}

// HasMouseFocus will return true if the user has clicked on the application and
// has it focues. It will return false otherwise.
func HasMouseFocus() bool {
	return sdl.GetMouseFocus() == getCurrent().SDLWindow
}

// IsOpen will return true if the window and been initialized. It will return false
// if it has not.
func IsOpen() bool {
	return getCurrent().open
}

// Destroy cleans up after the window has been closed. This is normally used by
// the game loop after the game loop has been stopped.
func Destroy() {
	win := getCurrent()
	win.open = false
	sdl.GL_DeleteContext(win.context)
	win.SDLWindow.Destroy()
	// The old window may have generated pending events which are no longer
	// relevant. Destroy them all!
	sdl.FlushEvent(sdl.WINDOWEVENT)
}

// ShowSimpleMessageBox will present a confirm style message box at the top of the screen.
// if attached to the window it will be attached to the top bar. the box type will style
// the message box and icon
func ShowSimpleMessageBox(title, message string, boxType MessageBoxType, attachtowindow bool) error {
	var sdlwindow *sdl.Window
	if attachtowindow {
		sdlwindow = getCurrent().SDLWindow
	}
	return sdl.ShowSimpleMessageBox(uint32(boxType), title, message, sdlwindow)
}

// ShowMessageBox will present a message box at the top of the screen.
// if attached to the window it will be attached to the top bar. the box type will style
// the message box and icon. The first button will be the primary enter action, and
// the last button will be the escape/cancel action
func ShowMessageBox(title, message string, buttons []string, boxType MessageBoxType, attachtowindow bool) string {
	var sdlwindow *sdl.Window
	if attachtowindow {
		sdlwindow = getCurrent().SDLWindow
	}

	SDLButtons := []sdl.MessageBoxButtonData{}
	for i, buttonText := range buttons {
		newButton := sdl.MessageBoxButtonData{
			ButtonId: int32(i),
			Text:     buttonText,
		}
		if i == 0 {
			newButton.Flags |= sdl.MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT
		}
		if i == len(buttons)-1 {
			newButton.Flags |= sdl.MESSAGEBOX_BUTTON_ESCAPEKEY_DEFAULT
		}
		SDLButtons = append(SDLButtons, newButton)
	}

	messageboxdata := sdl.MessageBoxData{
		Flags:      uint32(boxType),
		Window:     sdlwindow,
		Title:      title,
		Message:    message,
		NumButtons: int32(len(SDLButtons)),
		Buttons:    SDLButtons,
	}

	var _, buttonid = sdl.ShowMessageBox(&messageboxdata)
	return buttons[buttonid]
}
