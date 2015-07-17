package window

import (
	"github.com/BurntSushi/toml"
)

type WindowConfig struct {
	Identity   string `toml:"identity"`       // The name of the save directory (string)
	Title      string `toml:"title"`          // The window title (string)
	Icon       string `toml:"icon"`           // Filepath to an image to use as the window's icon (string)
	Width      int    `toml:"width"`          // The window width (number)
	Height     int    `toml:"height"`         // The window height (number)
	Borderless bool   `toml:"borderless"`     // Remove all border visuals from the window (boolean)
	Resizable  bool   `toml:"resizable"`      // Let the window be user-resizable (boolean)
	Minwidth   int    `toml:"minwidth"`       // Minimum window width if the window is resizable (number)
	Minheight  int    `toml:"minheight"`      // Minimum window height if the window is resizable (number)
	Fullscreen bool   `toml:"fullscreen"`     // Enable fullscreen (boolean)
	Fstype     string `toml:"fullscreentype"` // Standard fullscreen or desktop fullscreen mode (string)
	Vsync      bool   `toml:"vsync"`          // Enable vertical sync (boolean)
	Fsaa       int    `toml:"fsaa"`           // The number of samples to use with multi-sampled antialiasing (number)
	Display    int    `toml:"display"`        // Index of the monitor to show the window in (number)
	Highdpi    bool   `toml:"highdpi"`        // Enable high-dpi mode for the window on a Retina display (boolean)
	Srgb       bool   `toml:"srgb"`           // Enable sRGB gamma correction when drawing to the screen (boolean)
	Centered   bool   `toml:"centered"`       // Center the window in the display
	X          int    `toml:"x"`              // The x-coordinate of the window's position in the specified display (number)
	Y          int    `toml:"y"`              // The y-coordinate of the window's position in the specified display (number)
}

func loadConfig() (*WindowConfig, error) {
	var config WindowConfig
	if _, err := toml.DecodeFile("conf.toml", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
