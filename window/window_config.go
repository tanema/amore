package window

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"

	"github.com/tanema/amore/file"
)

const (
	config_file_name = "conf.toml" // default config file name.
)

// windowConfig is the struct that the config file will unmarshalled into
type windowConfig struct {
	Identity    string `toml:"identity"` // The name of the save directory (string)
	Title       string `toml:"title"`    // The window title (string)
	Icon        string `toml:"icon"`     // Filepath to an image to use as the window's icon (string)
	Width       int32  `toml:"width"`    // The window width (number)
	Height      int32  `toml:"height"`   // The window height (number)
	PixelWidth  int32  `toml:"-"`
	PixelHeight int32  `toml:"-"`
	Borderless  bool   `toml:"borderless"`     // Remove all border visuals from the window (boolean)
	Resizable   bool   `toml:"resizable"`      // Let the window be user-resizable (boolean)
	Minwidth    int32  `toml:"minwidth"`       // Minimum window width if the window is resizable (number)
	Minheight   int32  `toml:"minheight"`      // Minimum window height if the window is resizable (number)
	Fullscreen  bool   `toml:"fullscreen"`     // Enable fullscreen (boolean)
	Fstype      string `toml:"fullscreentype"` // Standard fullscreen or desktop fullscreen mode (string)
	Vsync       bool   `toml:"vsync"`          // Enable vertical sync (boolean)
	Msaa        int    `toml:"msaa"`           // The number of samples to use with multi-sampled antialiasing (number)
	Display     int    `toml:"display"`        // Index of the monitor to show the window in (number)
	Highdpi     bool   `toml:"highdpi"`        // Enable high-dpi mode for the window on a Retina display (boolean)
	Srgb        bool   `toml:"srgb"`           // Enable sRGB gamma correction when drawing to the screen (boolean)
	Centered    bool   `toml:"centered"`       // Center the window in the display
	X           int32  `toml:"x"`              // The x-coordinate of the window's position in the specified display (number)
	Y           int32  `toml:"y"`              // The y-coordinate of the window's position in the specified display (number)
}

// loadConfig will find the config file (works with bundle) and load it into the
// struct and the return it. If the config does not exist it will return a default
// config. If there was an issue reading the file it will return the error, possibly
// permission errors.
func loadConfig() (*windowConfig, error) {
	var config windowConfig
	path := path.Join(".", config_file_name)

	if _, err := file.NewFileData(path); os.IsNotExist(err) {
		return &config, nil //no added config bail out early
	}

	config_file, file_err := file.NewFile(config_file_name)
	if file_err != nil {
		return nil, file_err
	}

	if _, err := toml.DecodeReader(config_file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
