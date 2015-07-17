package window

import (
	"math"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/tanema/amore/gfx"
)

var (
	default_title  = "Amore Engine"
	current_window *Window
	created        = false
)

type Window struct {
	sdl_window                *sdl.Window
	context                   sdl.GLContext
	width, height             int
	pixel_width, pixel_height int
	title                     string
	should_close              bool
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

	x := config.X
	y := config.Y
	if !config.Fullscreen {
		// The position needs to be in the global coordinate space.
		var displaybounds sdl.Rect
		sdl.GetDisplayBounds(config.Display, &displaybounds)
		x += int(displaybounds.X)
		y += int(displaybounds.Y)
	} else {
		if config.Centered {
			x = sdl.WINDOWPOS_CENTERED
			y = sdl.WINDOWPOS_CENTERED
		} else {
			x = sdl.WINDOWPOS_UNDEFINED
			y = sdl.WINDOWPOS_UNDEFINED
		}
	}

	gfx.UnSetMode()

	if current_window != nil {
		current_window.Destroy()
	}

	created = false
	window, err := createWindowAndContext(x, y, config.Width, config.Height, sdlflags)
	if err != nil {
		return nil, err
	}
	created = true

	//SetIcon(curMode.icon.get());
	window.SetMouseGrab(false)
	window.SetMinimumSize(config.Minwidth, config.Minheight)
	window.SetTitle(config.Title)

	if config.Centered && !config.Fullscreen {
		window.SetPosition(x, y)
	}

	window.Raise()

	if config.Vsync {
		sdl.GL_SetSwapInterval(1)
	} else {
		sdl.GL_SetSwapInterval(0)
	}

	gfx.SetMode(config.Width, config.Height)

	return window, nil
}

func createWindowAndContext(x, y, w, h int, windowflags uint32) (*Window, error) {
	window, err := sdl.CreateWindow(default_title, x, y, w, h, windowflags)
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
		width:        w,
		height:       h,
		title:        default_title,
		should_close: false,
	}
	return current_window, nil
}

func setupGL() error {
	if err := gl.Init(); err != nil {
		return err
	}
	gl.Enable(gl.BLEND)
	// Auto-generated mipmaps should be the best quality possible
	gl.Hint(gl.GENERATE_MIPMAP_HINT, gl.NICEST)
	// Set pixel row alignment
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	return nil
}

func GetCurrent() *Window {
	return current_window
}

func GetDisplayCount() int {
	num, _ := sdl.GetNumVideoDisplays()
	return num
}

//func onResize(win *glfw.Window, w, h int) {
//Height = h
//Width = w
//gl.Viewport(0, 0, int32(Width), int32(Height))
//}

func (window *Window) SetTitle(title string) {
	window.sdl_window.SetTitle(title)
	window.title = title
}

func (window *Window) GetTitle() string {
	return window.title
}

func (window *Window) GetWidth() int {
	return window.width
}

func (window *Window) GetHeight() int {
	return window.height
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
	new_x := x * (float32(window.pixel_width) / float32(window.width))
	new_y := y * (float32(window.pixel_height) / float32(window.height))
	return new_x, new_y
}

func (window *Window) PixelToWindowCoords(x, y float32) (float32, float32) {
	new_x := x * (float32(window.width) / float32(window.pixel_width))
	new_y := y * (float32(window.height) / float32(window.pixel_height))
	return new_x, new_y
}

func (window *Window) GetMousePosition() (float32, float32) {
	mx, my, _ := sdl.GetMouseState()
	print(mx)
	println(my)
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
	sdl.GL_DeleteContext(window.context)
	window.sdl_window.Destroy()
	// The old window may have generated pending events which are no longer
	// relevant. Destroy them all!
	sdl.FlushEvent(sdl.WINDOWEVENT)
}
