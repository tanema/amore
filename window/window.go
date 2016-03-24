// The window Pacakge creates and manages the window and gl context.
package window

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"runtime"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/tanema/amore/file"
)

var (
	current_window *Window
	created        = false
)

type MessageBoxType uint32

const (
	MESSAGEBOX_ERROR   MessageBoxType = sdl.MESSAGEBOX_ERROR
	MESSAGEBOX_WARNING MessageBoxType = sdl.MESSAGEBOX_WARNING
	MESSAGEBOX_INFO    MessageBoxType = sdl.MESSAGEBOX_INFORMATION
)

type Window struct {
	sdl_window                *sdl.Window
	context                   sdl.GLContext
	pixel_width, pixel_height int
	should_close              bool
	config                    *WindowConfig
	refresh_rate              int32
	open                      bool
}

func newWindow() (*Window, error) {
	var config *WindowConfig
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
		current_window.Destroy()
	}

	created = false
	window, err := createWindowAndContext(config, sdlflags)
	if err != nil {
		return nil, err
	}
	created = true

	if window.config.Icon != "" {
		window.SetIcon(window.config.Icon)
	}

	window.SetMouseGrab(false)
	window.SetMinimumSize(config.Minwidth, config.Minheight)
	window.SetTitle(config.Title)

	if config.Centered && !config.Fullscreen {
		window.SetPosition(config.X, config.Y)
	}

	window.Raise()

	if config.Vsync {
		sdl.GL_SetSwapInterval(1)
	} else {
		sdl.GL_SetSwapInterval(0)
	}

	window.UpdateSettings()

	window.open = true

	return window, nil
}

func createWindowAndContext(config *WindowConfig, windowflags uint32) (*Window, error) {
	setGLFramebufferAttributes(config.Msaa, config.Srgb)
	_, debug := os.LookupEnv("AMORE_DEBUG")
	setGLContextAttributes(2, 1, debug)

	window, err := sdl.CreateWindow(config.Title, int(config.X), int(config.Y), int(config.Width), int(config.Height), windowflags)
	if err != nil {
		panic(err)
	}

	context, err := sdl.GL_CreateContext(window)
	if err != nil {
		panic(err)
	}

	current_window = &Window{
		sdl_window:   window,
		context:      context,
		should_close: false,
		config:       config,
	}
	return current_window, nil
}

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

func (window *Window) OnSizeChanged(width, height int32) {
	window.config.Width = width
	window.config.Height = height
	window.config.PixelWidth, window.config.PixelHeight = window.GetDrawableSize()
}

func (window *Window) GetDrawableSize() (int32, int32) {
	w, h := sdl.GL_GetDrawableSize(window.sdl_window)
	return int32(w), int32(h)
}

func (window *Window) UpdateSettings() {
	wflags := window.sdl_window.GetFlags()

	// Set the new display mode as the current display mode.
	w, h := window.sdl_window.GetSize()
	window.config.Width, window.config.Height = int32(w), int32(h)
	window.pixel_width, window.pixel_height = sdl.GL_GetDrawableSize(window.sdl_window)

	if (wflags & sdl.WINDOW_FULLSCREEN_DESKTOP) == sdl.WINDOW_FULLSCREEN_DESKTOP {
		window.config.Fullscreen = true
		window.config.Fstype = "desktop"
	} else if (wflags & sdl.WINDOW_FULLSCREEN) == sdl.WINDOW_FULLSCREEN {
		window.config.Fullscreen = true
		window.config.Fstype = "exclusive"
	} else {
		window.config.Fullscreen = false
		window.config.Fstype = "normal"
	}

	// The min width/height is set to 0 internally in SDL when in fullscreen.
	if window.config.Fullscreen {
		window.config.Minwidth = 1
		window.config.Minheight = 1
	} else {
		mw, mh := window.sdl_window.GetMinimumSize()
		window.config.Minwidth, window.config.Minheight = int32(mw), int32(mh)
	}

	window.config.Resizable = (wflags & sdl.WINDOW_RESIZABLE) != 0
	window.config.Borderless = (wflags & sdl.WINDOW_BORDERLESS) != 0
	window.config.Centered = true

	x, y := window.sdl_window.GetPosition()
	window.config.X, window.config.Y = int32(x), int32(y)

	window.config.Highdpi = (wflags & sdl.WINDOW_ALLOW_HIGHDPI) != 0

	// Only minimize on focus loss if the window is in exclusive-fullscreen
	// mode.
	if window.config.Fullscreen && window.config.Fstype == "exclusive" {
		sdl.SetHint(sdl.HINT_VIDEO_MINIMIZE_ON_FOCUS_LOSS, "1")
	} else {
		sdl.SetHint(sdl.HINT_VIDEO_MINIMIZE_ON_FOCUS_LOSS, "0")
	}

	window.config.Srgb = false
	interval, _ := sdl.GL_GetSwapInterval()
	window.config.Vsync = (interval != 0)

	var dmode sdl.DisplayMode
	sdl.GetCurrentDisplayMode(window.config.Display, &dmode)

	// May be 0 if the refresh rate can't be determined.
	window.refresh_rate = dmode.RefreshRate
}

func GetCurrent() *Window {
	if current_window == nil {
		new_window, _ := newWindow()
		return new_window
	} else {
		return current_window
	}
}

func GetDisplayCount() int {
	num, _ := sdl.GetNumVideoDisplays()
	return num
}

func GetDisplayName(displayindex int) string {
	return sdl.GetDisplayName(displayindex)
}

func (window *Window) GetFullscreenSizes(displayindex int) [][]int32 {
	var sizes [][]int32
	modes, _ := sdl.GetNumDisplayModes(displayindex)
	for i := 0; i < modes; i++ {
		var mode sdl.DisplayMode
		sdl.GetDisplayMode(displayindex, i, &mode)
		sizes = append(sizes, []int32{mode.W, mode.H})
	}
	return sizes
}

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

func (window *Window) SetTitle(title string) {
	window.sdl_window.SetTitle(title)
	window.config.Title = title
}

func (window *Window) GetTitle() string {
	return window.config.Title
}

func (window *Window) GetWidth() int32 {
	return window.config.Width
}

func (window *Window) GetHeight() int32 {
	return window.config.Height
}

func (window *Window) SetIcon(path string) error {
	window.config.Icon = path

	imgFile, new_err := file.NewFile(path)
	defer imgFile.Close()
	if new_err != nil {
		return new_err
	}

	decoded_img, _, img_err := image.Decode(imgFile)
	if img_err != nil {
		return img_err
	}

	bounds := decoded_img.Bounds()
	rgba := image.NewRGBA(decoded_img.Bounds())
	draw.Draw(rgba, bounds, decoded_img, image.Point{0, 0}, draw.Src)

	var rmask, gmask, bmask, amask uint32
	switch runtime.GOARCH {
	case "mips64", "ppc64":
		rmask = 0xFF000000
		gmask = 0x00FF0000
		bmask = 0x0000FF00
		amask = 0x000000FF
	default:
		rmask = 0x000000FF
		gmask = 0x0000FF00
		bmask = 0x00FF0000
		amask = 0xFF000000
	}

	surface, ic_err := sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&rgba.Pix[0]), bounds.Dx(), bounds.Dy(), 32, bounds.Dx()*4, rmask, gmask, bmask, amask)
	if ic_err != nil {
		return ic_err
	}

	window.sdl_window.SetIcon(surface)
	surface.Free()

	return nil
}

func (window *Window) GetIcon() string {
	return window.config.Icon
}

func (window *Window) Minimize() {
	window.sdl_window.Minimize()
}

func (window *Window) Maximize() {
	window.sdl_window.Maximize()
}

func (window *Window) ShouldClose() bool {
	return window.should_close
}

func (window *Window) SetShouldClose(should_close bool) {
	window.should_close = should_close
}

func (window *Window) SwapBuffers() {
	sdl.GL_SwapWindow(window.sdl_window)
}

func (window *Window) WindowToPixelCoords(x, y float32) (float32, float32) {
	new_x := x * (float32(window.pixel_width) / float32(window.config.Width))
	new_y := y * (float32(window.pixel_height) / float32(window.config.Height))
	return new_x, new_y
}

func (window *Window) PixelToWindowCoords(x, y float32) (float32, float32) {
	new_x := x * (float32(window.config.Width) / float32(window.pixel_width))
	new_y := y * (float32(window.config.Height) / float32(window.pixel_height))
	return new_x, new_y
}

func (window *Window) GetMousePosition() (float32, float32) {
	mx, my, _ := sdl.GetMouseState()
	return window.WindowToPixelCoords(float32(mx), float32(my))
}

func (window *Window) SetMousePosition(x, y float32) {
	wx, wy := window.PixelToWindowCoords(x, y)
	window.sdl_window.WarpMouseInWindow(int(wx), int(wy))
}

func (window *Window) IsMouseGrabbed() bool {
	return window.sdl_window.GetGrab() != false
}

func (window *Window) IsVisible() bool {
	return (window.sdl_window.GetFlags() & sdl.WINDOW_SHOWN) != 0
}

func (window *Window) SetMouseVisible(visible bool) {
	if visible {
		sdl.ShowCursor(sdl.ENABLE)
	} else {
		sdl.ShowCursor(sdl.DISABLE)
	}
}

func (window *Window) GetMouseVisible() bool {
	return sdl.ShowCursor(sdl.QUERY) == sdl.ENABLE
}

func (window *Window) GetPixelDimensions() (int32, int32) {
	return window.config.PixelWidth, window.config.PixelHeight
}

func (window *Window) GetPixelScale() float32 {
	return float32(window.config.PixelHeight) / float32(window.config.Height)
}

func (window *Window) ToPixels(x float32) float32 {
	return x * window.GetPixelScale()
}

func (window *Window) ToPixelsPoint(x, y float32) (float32, float32) {
	scale := window.GetPixelScale()
	return x * scale, y * scale
}

func (window *Window) FromPixels(x float32) float32 {
	return x / window.GetPixelScale()
}

func (window *Window) FromPixelsPoint(x, y float32) (float32, float32) {
	scale := window.GetPixelScale()
	return x / scale, y / scale
}

func (window *Window) SetMouseGrab(grabbed bool) {
	window.sdl_window.SetGrab(grabbed)
}

func (window *Window) SetMinimumSize(w, h int32) {
	window.config.Minwidth = w
	window.config.Minheight = h
	window.sdl_window.SetMinimumSize(int(w), int(h))
}

func (window *Window) SetPosition(x, y int32) {
	window.config.X = x
	window.config.Y = y
	window.sdl_window.SetPosition(int(x), int(y))
}

func (window *Window) GetPosition() (int, int) {
	return window.sdl_window.GetPosition()
}

func (window *Window) ShowSimpleMessageBox(title, message string, box_type MessageBoxType, attachtowindow bool) error {
	var sdlwindow *sdl.Window
	if attachtowindow {
		sdlwindow = window.sdl_window
	}
	return sdl.ShowSimpleMessageBox(uint32(box_type), title, message, sdlwindow)
}

func (window *Window) ShowMessageBox(title, message string, buttons []string, box_type MessageBoxType, attachtowindow bool) string {
	var sdlwindow *sdl.Window
	if attachtowindow {
		sdlwindow = window.sdl_window
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

func (window *Window) HasFocus() bool {
	return sdl.GetKeyboardFocus() == window.sdl_window
}

func (window *Window) RequestAttention(continuous bool) {
	if window.HasFocus() {
		return
	}
	requestAttention(continuous)
}

func (window *Window) HasMouseFocus() bool {
	return sdl.GetMouseFocus() == window.sdl_window
}

func (window *Window) IsOpen() bool {
	return window.open
}

func (window *Window) Raise() {
	window.sdl_window.Raise()
}

func (window *Window) Destroy() {
	window.open = false
	sdl.GL_DeleteContext(window.context)
	window.sdl_window.Destroy()
	// The old window may have generated pending events which are no longer
	// relevant. Destroy them all!
	sdl.FlushEvent(sdl.WINDOWEVENT)
}
