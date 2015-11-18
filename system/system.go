package system

import (
	"runtime"

	"github.com/skratchdot/open-golang/open"
	"github.com/veandco/go-sdl2/sdl"
)

func GetClipboardText() (string, error) {
	return sdl.GetClipboardText()
}

func SetClipboardText(str string) error {
	return sdl.SetClipboardText(str)
}

func GetOS() string {
	return runtime.GOOS
}

func GetPowerInfo() (PowerState, int, int) {
	state, seconds, percent := sdl.GetPowerInfo()
	return PowerState(state), seconds, percent
}

func GetProcessorCount() int {
	return runtime.NumCPU()
}

func OpenUrl(url string) {
	open.Run(url)
}
