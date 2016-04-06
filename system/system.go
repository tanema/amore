// The system Pacakge is a utility pacakge meant to interface with system operations
// like opening an application or setting the clipboard text
package system

import (
	"runtime"

	"github.com/skratchdot/open-golang/open"
	"github.com/veandco/go-sdl2/sdl"
)

// GetClipboardText returns a string of the clipboard text. It may return an error
// if unable to get the contents.
func GetClipboardText() (string, error) {
	return sdl.GetClipboardText()
}

// SetClipboardText will attempt to set the clipboard text with the provided string.
// It will return an error if it was unable to.
func SetClipboardText(str string) error {
	return sdl.SetClipboardText(str)
}

// GetOS will return runtime.GOOS this is the running program's operating system
// target: one of darwin, freebsd, linux, and so on.
func GetOS() string {
	return runtime.GOOS
}

// GetPowerInfo will return the PowerSate, seconds of battery left, and percentage
// of battery power.
func GetPowerInfo() (PowerState, int, int) {
	state, seconds, percent := sdl.GetPowerInfo()
	return PowerState(state), seconds, percent
}

// GetProcessorCount returns the number of logical CPUs usable by the current process.
func GetProcessorCount() int {
	return runtime.NumCPU()
}

// OpenUrl will open a file, directory, or URI using the OS's default application
// for that object type.
//
// See github.com/skratchdot/open-golang
func OpenUrl(url string) {
	open.Start(url)
}
