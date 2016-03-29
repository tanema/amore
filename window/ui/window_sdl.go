// +build !js

package ui

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var CurrentWindow *Window

func InitWindowAndContext(config WindowConfig) (*Window, Context, error) {
	var err error

	if err = sdl.InitSubSystem(sdl.INIT_VIDEO); err != nil {
		return nil, nil, err
	}

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
					return nil, nil, err
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

	if config.Highdpi {
		sdlflags |= sdl.WINDOW_ALLOW_HIGHDPI
	}

	if config.Fullscreen {
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

	setGLFramebufferAttributes(config.Msaa, config.Srgb)
	_, debug := os.LookupEnv("AMORE_DEBUG")
	setGLContextAttributes(2, 1, debug)

	newWindow, err := sdl.CreateWindow(config.Title, int(config.X), int(config.Y), int(config.Width), int(config.Height), sdlflags)
	if err != nil {
		return nil, nil, err
	}
	CurrentWindow = &Window{newWindow}

	newContext, err := sdl.GL_CreateContext(newWindow)
	if err != nil {
		return nil, nil, err
	}

	if config.Icon != "" {
		CurrentWindow.SetIcon(config.Icon)
	}

	CurrentWindow.SetGrab(false)
	CurrentWindow.SetMinimumSize(config.Minwidth, config.Minheight)
	CurrentWindow.SetTitle(config.Title)
	if !config.Fullscreen {
		CurrentWindow.SetPosition(config.X, config.Y)
	}
	CurrentWindow.Raise()

	if config.Vsync {
		sdl.GL_SetSwapInterval(1)
	} else {
		sdl.GL_SetSwapInterval(0)
	}

	return CurrentWindow, Context(newContext), nil
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

func (window *Window) SwapBuffers() {
	sdl.GL_SwapWindow(window.Window)
}

func (window *Window) GetDrawableSize() (int, int) {
	return sdl.GL_GetDrawableSize(window.Window)
}

func (w *Window) SetIcon(path string) error {
	surface, err := loadSurface(path)
	if err != nil {
		return err
	}
	w.Window.SetIcon(surface)
	surface.Free()
	return nil
}

func (window *Window) IsMouseGrabbed() bool {
	return window.Window.GetGrab() != false
}

func (window *Window) IsVisible() bool {
	return (window.Window.GetFlags() & sdl.WINDOW_SHOWN) != 0
}

func (window *Window) Alert(title, message string) error {
	return sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_INFORMATION, title, message, window.Window)
}

func (window *Window) Confirm(title, message string) bool {
	return window.ShowMessageBox(title, message, []string{"ok", "cancel"}) == "ok"
}

func (window *Window) ShowMessageBox(title, message string, buttons []string) string {
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
		Flags:      sdl.MESSAGEBOX_INFORMATION,
		Window:     window.Window,
		Title:      title,
		Message:    message,
		NumButtons: int32(len(sdl_buttons)),
		Buttons:    sdl_buttons,
	}

	var _, buttonid = sdl.ShowMessageBox(&messageboxdata)
	return buttons[buttonid]
}

func (window *Window) HasFocus() bool {
	return sdl.GetKeyboardFocus() == window.Window
}

func (window *Window) HasMouseFocus() bool {
	return sdl.GetMouseFocus() == window.Window
}

func (window *Window) RequestAttention(continuous bool) {
	if window.HasFocus() {
		return
	}
	requestAttention(continuous)
}
