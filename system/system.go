package system

import (
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/skratchdot/open-golang/open"

	"github.com/tanema/amore/window"
)

func GetClipboardText() (string, error) {
	return window.GetCurrent().GetClipboardString()
}

func SetClipboardText(str string) {
	window.GetCurrent().GetClipboardString(str)
}

func GetOS() string {
	return runtime.GOOS
}

func GetPowerInfo() {

}

func GetProcessorCount() int {
	return runtime.NumCPU()
}

func OpenUrl(url string) {
	open.Run(url)
}
