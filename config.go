package amore

import (
	"math"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/goxjs/glfw"

	"github.com/tanema/amore/file"
)

const (
	configFileName = "conf.toml" // default config file name.
)

// Config is the struct that the config file will unmarshalled into
type config struct {
	Title      string `toml:"title"`       // The window title (string)
	Width      int    `toml:"width"`       // The window width (number)
	Height     int    `toml:"height"`      // The window height (number)
	Resizable  bool   `toml:"resizable"`   // Let the window be user-resizable (boolean)
	Fullscreen bool   `toml:"fullscreen"`  // Enable fullscreen (boolean)
	Vsync      bool   `toml:"vsync"`       // Enable vertical sync (boolean)
	Msaa       int    `toml:"msaa"`        // The number of samples to use with multi-sampled antialiasing (number)
	MouseShown bool   `toml:"mouse_shown"` // show the mouse
}

// loadConfig will find the config file (works with bundle) and load it into the
// struct and the return it. If the config does not exist it will return a default
// config. If there was an issue reading the file it will return the error, possibly
// permission errors.
func loadConfig() (*config, error) {
	var c config
	if configFile, err := file.NewFile(configFileName); err == nil {
		if _, err := toml.DecodeReader(configFile, &c); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	if c.Width == 0 || c.Height == 0 {
		c.Width = 800
		c.Height = 600
	}

	glfw.WindowHint(glfw.Samples, int(math.Max(float64(c.Msaa), 4.0)))
	if c.Resizable {
		glfw.WindowHint(glfw.Resizable, 1)
	} else {
		glfw.WindowHint(glfw.Resizable, 0)
	}

	return &c, nil
}
