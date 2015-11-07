package window

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/tanema/amore/gfx"
)

var (
	current_window *Window
	created        = false
)

type Window struct {
	sdl_window                *sdl.Window
	context                   sdl.GLContext
	pixel_width, pixel_height int
	should_close              bool
	config                    *WindowConfig
	refresh_rate              int32
}

func New() (*Window, error) {
	var config *WindowConfig
	var err error

	if err = sdl.InitSubSystem(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	if config, err = loadConfig(); err != nil {
		panic(err)
	}

	config.Minwidth = int(math.Max(float64(config.Minwidth), 1.0))
	config.Minheight = int(math.Max(float64(config.Minheight), 1.0))
	config.Display = int(math.Min(math.Max(float64(config.Display), 0.0), float64(GetDisplayCount()-1)))

	if config.Width == 0 || config.Height == 0 {
		var mode sdl.DisplayMode
		sdl.GetDesktopDisplayMode(config.Display, &mode)
		config.Width = int(mode.W)
		config.Height = int(mode.H)
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

			config.Width = int(mode.W)
			config.Height = int(mode.H)
		}
	}

	if config.Resizable {
		sdlflags |= sdl.WINDOW_RESIZABLE
	}

	if config.Borderless {
		sdlflags |= sdl.WINDOW_BORDERLESS
	}

	//if config.Highdpi {
	//sdlflags |= sdl.WINDOW_ALLOW_HIGHDPI
	//}

	if !config.Fullscreen {
		// The position needs to be in the global coordinate space.
		var displaybounds sdl.Rect
		sdl.GetDisplayBounds(config.Display, &displaybounds)
		config.X += int(displaybounds.X)
		config.Y += int(displaybounds.Y)
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

	//SetIcon(curMode.icon.get());
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

	gfx.InitContext(config.Width, config.Height)

	return window, nil
}

func createWindowAndContext(config *WindowConfig, windowflags uint32) (*Window, error) {
	window, err := sdl.CreateWindow(config.Title, config.X, config.Y, config.Width, config.Height, windowflags)
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

func (window *Window) UpdateSettings() {
	wflags := window.sdl_window.GetFlags()

	// Set the new display mode as the current display mode.
	window.config.Width, window.config.Height = window.sdl_window.GetSize()
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
		window.config.Minwidth, window.config.Minheight = window.sdl_window.GetMinimumSize()
	}

	window.config.Resizable = (wflags & sdl.WINDOW_RESIZABLE) != 0
	window.config.Borderless = (wflags & sdl.WINDOW_BORDERLESS) != 0
	window.config.Centered = true

	window.config.X, window.config.Y = window.sdl_window.GetPosition()

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
	return current_window
}

func GetDisplayCount() int {
	num, _ := sdl.GetNumVideoDisplays()
	return num
}

func (window *Window) SetTitle(title string) {
	window.sdl_window.SetTitle(title)
	window.config.Title = title
}

func (window *Window) GetTitle() string {
	return window.config.Title
}

func (window *Window) GetWidth() int {
	return window.config.Width
}

func (window *Window) GetHeight() int {
	return window.config.Height
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

func (window *Window) SetMouseGrab(grabbed bool) {
	window.sdl_window.SetGrab(grabbed)
}

func (window *Window) SetMinimumSize(w, h int) {
	window.sdl_window.SetMinimumSize(w, h)
}

func (window *Window) SetPosition(x, y int) {
	window.sdl_window.SetPosition(x, y)
}

func (window *Window) Raise() {
	window.sdl_window.Raise()
}

func (window *Window) Destroy() {
	gfx.DeInit()
	sdl.GL_DeleteContext(window.context)
	window.sdl_window.Destroy()
	// The old window may have generated pending events which are no longer
	// relevant. Destroy them all!
	sdl.FlushEvent(sdl.WINDOWEVENT)
}
