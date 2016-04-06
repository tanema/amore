// The window Pacakge creates and manages the window and gl context.
package window

import (
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/tanema/amore/window/surface"
)

var (
	current_window *window
	created        = false
	initError      error
)

// MessageBoxType specifies the types of Message box to create
type MessageBoxType uint32

const (
	MESSAGEBOX_ERROR   MessageBoxType = sdl.MESSAGEBOX_ERROR
	MESSAGEBOX_WARNING MessageBoxType = sdl.MESSAGEBOX_WARNING
	MESSAGEBOX_INFO    MessageBoxType = sdl.MESSAGEBOX_INFORMATION
)

// Window contains the created window information
type window struct {
	sdl_window                *sdl.Window
	context                   sdl.GLContext
	pixel_width, pixel_height int
	should_close              bool
	Config                    *windowConfig
	refresh_rate              int32
	open                      bool
}

// NewWindow will create a new sdl window with a gl context. It will also destroy
// an old one if there is one already. Errors returned come straight from SDL so
// errors will be indicitive of SDL errors
func NewWindow() (*window, error) {
	if current_window != nil || initError != nil {
		return current_window, initError
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

	if current_window != nil {
		Destroy()
	}

	created = false
	new_window, err := createWindowAndContext(config, sdlflags)
	if err != nil {
		return nil, err
	}
	created = true

	if new_window.Config.Icon != "" {
		SetIcon(config.Icon)
	}

	SetMouseGrab(false)
	SetMinimumSize(config.Minwidth, config.Minheight)
	SetTitle(config.Title)

	if config.Centered && !config.Fullscreen {
		SetPosition(config.X, config.Y)
	}

	getCurrent().sdl_window.Raise()

	if config.Vsync {
		sdl.GL_SetSwapInterval(1)
	} else {
		sdl.GL_SetSwapInterval(0)
	}

	new_window.open = true
	return new_window, nil
}

// createWindowAndContext is the actual interface with SDL to create window and context
func createWindowAndContext(config *windowConfig, windowflags uint32) (*window, error) {
	setGLFramebufferAttributes(config.Msaa, config.Srgb)
	_, debug := os.LookupEnv("AMORE_DEBUG")
	setGLContextAttributes(2, 1, debug)

	new_window, err := sdl.CreateWindow(config.Title, int(config.X), int(config.Y), int(config.Width), int(config.Height), windowflags)
	if err != nil {
		panic(err)
	}

	context, err := sdl.GL_CreateContext(new_window)
	if err != nil {
		panic(err)
	}

	current_window = &window{
		sdl_window:   new_window,
		context:      context,
		should_close: false,
		Config:       config,
	}
	return current_window, nil
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
	if current_window == nil {
		current_window, initError = NewWindow()
	}

	return current_window
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
	w, h := sdl.GL_GetDrawableSize(getCurrent().sdl_window)
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
	win.sdl_window.SetTitle(title)
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
	new_surface, err := surface.Load(path)
	if err != nil {
		return err
	}
	win.sdl_window.SetIcon(new_surface)
	new_surface.Free()
	return nil
}

// GetIcon returns the path to the currently set icon
func GetIcon() string {
	return getCurrent().Config.Icon
}

// Minimize hides the program in the task bar
func Minimize() {
	getCurrent().sdl_window.Minimize()
}

// Maximize unhides the program and sets it to max size
func Maximize() {
	getCurrent().sdl_window.Maximize()
}

// ShouldClose returns if the window has beed told to close and is in the process
// of shutting down
func ShouldClose() bool {
	return getCurrent().should_close
}

// Close prepares the window to shut down gracefully
func Close(should_close bool) {
	getCurrent().should_close = should_close
}

// SwapBuffers swaps the current frame buffer in the window. This is generally
// only used by the engine but you could use it to roll your own game loop.
func SwapBuffers() {
	sdl.GL_SwapWindow(getCurrent().sdl_window)
}

// WindowToPixelCoords translates window coords to pixel coords for pixel perfect
// operations
func WindowToPixelCoords(x, y float32) (float32, float32) {
	config := getCurrent().Config
	new_x := x * (float32(config.PixelWidth) / float32(config.Width))
	new_y := y * (float32(config.PixelHeight) / float32(config.Height))
	return new_x, new_y
}

// PixelToWindowCoords translates pixel coords to window coords
func PixelToWindowCoords(x, y float32) (float32, float32) {
	config := getCurrent().Config
	new_x := x * (float32(config.Width) / float32(config.PixelWidth))
	new_y := y * (float32(config.Height) / float32(config.PixelHeight))
	return new_x, new_y
}

// GetMousePosition return the current position of the mouse
func GetMousePosition() (float32, float32) {
	mx, my, _ := sdl.GetMouseState()
	return WindowToPixelCoords(float32(mx), float32(my))
}

// SetMousePosition warps the mouse position to a new position x, y in the window
func SetMousePosition(x, y float32) {
	wx, wy := PixelToWindowCoords(x, y)
	getCurrent().sdl_window.WarpMouseInWindow(int(wx), int(wy))
}

// IsMouseGrabbed returns true if pointer lock is enabled and false otherwise
func IsMouseGrabbed() bool {
	return getCurrent().sdl_window.GetGrab() != false
}

// SetMouseGrab enables or disables pointer lock on the window
func SetMouseGrab(grabbed bool) {
	getCurrent().sdl_window.SetGrab(grabbed)
}

// IsVisible returns true if the window is not minimized
func IsVisible() bool {
	return (getCurrent().sdl_window.GetFlags() & sdl.WINDOW_SHOWN) != 0
}

// SetMouseVisible will hide the mouse if passed false, and show the mouse cursor
// if passed true
func SetMouseVisible(visible bool) {
	if visible {
		sdl.ShowCursor(sdl.ENABLE)
	} else {
		sdl.ShowCursor(sdl.DISABLE)
	}
}

// GetMouseVisible will return true if the mouse is not hidden and will return
// false otherwise
func GetMouseVisible() bool {
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
	win.sdl_window.SetMinimumSize(int(w), int(h))
}

// SetPosition will set the window position on screen
func SetPosition(x, y int32) {
	win := getCurrent()
	win.Config.X = x
	win.Config.Y = y
	win.sdl_window.SetPosition(int(x), int(y))
}

// GetPosition will return the position of the window on the screen
func GetPosition() (int, int) {
	return getCurrent().sdl_window.GetPosition()
}

// HasFocus will return true if the user has the program focused and return false
// if the user is focused on another program
func HasFocus() bool {
	return sdl.GetKeyboardFocus() == getCurrent().sdl_window
}

// RequestAttention makes the application flash for attention to the user. If
// continuous is true, it will continue to do so until the user focuses the program.
func RequestAttention(continuous bool) {
	if HasFocus() {
		return
	}
	requestAttention(continuous)
}

// HasMouseFocus will return true if the user has clicked on the application and
// has it focues. It will return false otherwise.
func HasMouseFocus() bool {
	return sdl.GetMouseFocus() == getCurrent().sdl_window
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
	win.sdl_window.Destroy()
	// The old window may have generated pending events which are no longer
	// relevant. Destroy them all!
	sdl.FlushEvent(sdl.WINDOWEVENT)
}

// ShowSimpleMessageBox will present a confirm style message box at the top of the screen.
// if attached to the window it will be attached to the top bar. the box type will style
// the message box and icon
func ShowSimpleMessageBox(title, message string, box_type MessageBoxType, attachtowindow bool) error {
	var sdlwindow *sdl.Window
	if attachtowindow {
		sdlwindow = getCurrent().sdl_window
	}
	return sdl.ShowSimpleMessageBox(uint32(box_type), title, message, sdlwindow)
}

// ShowMessageBox will present a message box at the top of the screen.
// if attached to the window it will be attached to the top bar. the box type will style
// the message box and icon. The first button will be the primary enter action, and
// the last button will be the escape/cancel action
func ShowMessageBox(title, message string, buttons []string, box_type MessageBoxType, attachtowindow bool) string {
	var sdlwindow *sdl.Window
	if attachtowindow {
		sdlwindow = getCurrent().sdl_window
	}

	sdl_buttons := []sdl.MessageBoxButtonData{}
	for i, button_text := range buttons {
		new_button := sdl.MessageBoxButtonData{
			ButtonId: int32(i),
			Text:     button_text,
		}
		if i == 0 {
			new_button.Flags |= sdl.MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT
		}
		if i == len(buttons)-1 {
			new_button.Flags |= sdl.MESSAGEBOX_BUTTON_ESCAPEKEY_DEFAULT
		}
		sdl_buttons = append(sdl_buttons, new_button)
	}

	messageboxdata := sdl.MessageBoxData{
		Flags:      uint32(box_type),
		Window:     sdlwindow,
		Title:      title,
		Message:    message,
		NumButtons: int32(len(sdl_buttons)),
		Buttons:    sdl_buttons,
	}

	var _, buttonid = sdl.ShowMessageBox(&messageboxdata)
	return buttons[buttonid]
}
