// The system Pacakge is a utility pacakge meant to interface with system operations
// like opening an application or setting the clipboard text
package system

import (
	"runtime"

	"github.com/skratchdot/open-golang/open"

	"github.com/tanema/amore/window/ui"
)

func GetClipboardText() (string, error) {
	return ui.GetClipboardText()
}

func SetClipboardText(str string) error {
	return ui.SetClipboardText(str)
}

func GetOS() string {
	return runtime.GOOS
}

func GetPowerInfo() (string, int, int) {
	return ui.GetPowerInfo()
}

func GetProcessorCount() int {
	return runtime.NumCPU()
}

func OpenUrl(url string) {
	open.Run(url)
}
