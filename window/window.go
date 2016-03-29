// The window Pacakge creates and manages the window and gl context.
package window

import (
	"math"

	"github.com/BurntSushi/toml"

	"github.com/tanema/amore/file"
	"github.com/tanema/amore/gfx"
	"github.com/tanema/amore/joystick"
	"github.com/tanema/amore/window/ui"
)

const (
	config_file_name = "conf.toml"
)

var (
	currentWindow  *ui.Window
	currentContext ui.Context
	initError      error
	currentConfig  ui.WindowConfig
	shouldClose    bool
)

func Init() error {
	if initError != nil {
		return initError
	}
	if currentWindow == nil {
		config_file, file_err := file.NewFile(config_file_name)
		if file_err != nil {
			return file_err
		}

		if _, err := toml.DecodeReader(config_file, &currentConfig); err != nil {
			return err
		}

		//normlize values
		currentConfig.Minwidth = int(math.Max(float64(currentConfig.Minwidth), 1.0))
		currentConfig.Minheight = int(math.Max(float64(currentConfig.Minheight), 1.0))
		currentConfig.Display = int(math.Min(math.Max(float64(currentConfig.Display), 0.0), float64(GetDisplayCount()-1)))

		currentWindow, currentContext, initError = ui.InitWindowAndContext(currentConfig)

		if initError != nil {
			return initError
		}

		gfx.InitContext(int32(currentConfig.Width), int32(currentConfig.Height), currentContext)
		currentConfig.PixelWidth, currentConfig.PixelHeight = curWin().GetDrawableSize()
		joystick.Init()
	}
	return initError
}

func curWin() *ui.Window {
	if currentWindow == nil {
		Init()
	}
	return currentWindow
}

func OnSizeChanged(width, height int32) {
	currentConfig.Width = int(width)
	currentConfig.Height = int(height)
	currentConfig.PixelWidth, currentConfig.PixelHeight = curWin().GetDrawableSize()
}

func GetDrawableSize() (int, int) {
	return curWin().GetDrawableSize()
}

func GetDisplayCount() int {
	return ui.GetDisplayCount()
}

func GetDisplayName(displayindex int) string {
	return ui.GetDisplayName(displayindex)
}

func GetFullscreenSizes(displayindex int) [][]int32 {
	return ui.GetFullscreenSizes(displayindex)
}

func GetDesktopDimensions(displayindex int) (int32, int32) {
	return ui.GetDesktopDimensions(displayindex)
}

func SetTitle(title string) {
	curWin().SetTitle(title)
	currentConfig.Title = title
}

func GetTitle() string {
	return currentConfig.Title
}

func SetIcon(path string) error {
	return curWin().SetIcon(path)
}

func GetIcon() string {
	return currentConfig.Icon
}

func Minimize() {
	curWin().Minimize()
}

func Maximize() {
	curWin().Maximize()
}

func ShouldClose() bool {
	return shouldClose
}

func SetShouldClose(should_close bool) {
	shouldClose = should_close
}

func SwapBuffers() {
	curWin().SwapBuffers()
}

func WindowToPixelCoords(x, y float32) (float32, float32) {
	new_x := x * (float32(currentConfig.PixelWidth) / float32(currentConfig.Width))
	new_y := y * (float32(currentConfig.PixelHeight) / float32(currentConfig.Height))
	return new_x, new_y
}

func PixelToWindowCoords(x, y float32) (float32, float32) {
	new_x := x * (float32(currentConfig.Width) / float32(currentConfig.PixelWidth))
	new_y := y * (float32(currentConfig.Height) / float32(currentConfig.PixelHeight))
	return new_x, new_y
}

func GetMousePosition() (float32, float32) {
	mx, my := ui.GetMousePosition()
	return WindowToPixelCoords(float32(mx), float32(my))
}

func SetMousePosition(x, y float32) {
	wx, wy := PixelToWindowCoords(x, y)
	curWin().WarpMouseInWindow(int(wx), int(wy))
}

func IsMouseGrabbed() bool {
	return curWin().IsMouseGrabbed()
}

func IsVisible() bool {
	return curWin().IsVisible()
}

func SetMouseVisible(visible bool) {
	ui.SetMouseVisible(visible)
}

func GetMouseVisible() bool {
	return ui.GetMouseVisible()
}

func GetPixelDimensions() (int, int) {
	return currentConfig.PixelWidth, currentConfig.PixelHeight
}

func GetPixelScale() float32 {
	return float32(currentConfig.PixelHeight) / float32(currentConfig.Height)
}

func ToPixels(x float32) float32 {
	return x * GetPixelScale()
}

func ToPixelsPoint(x, y float32) (float32, float32) {
	scale := GetPixelScale()
	return x * scale, y * scale
}

func FromPixels(x float32) float32 {
	return x / GetPixelScale()
}

func FromPixelsPoint(x, y float32) (float32, float32) {
	scale := GetPixelScale()
	return x / scale, y / scale
}

func SetMouseGrab(grabbed bool) {
	curWin().SetGrab(grabbed)
}

func SetMinimumSize(w, h int) {
	currentConfig.Minwidth = w
	currentConfig.Minheight = h
	curWin().SetMinimumSize(w, h)
}

func SetPosition(x, y int) {
	currentConfig.X = x
	currentConfig.Y = y
	curWin().SetPosition(x, y)
}

func GetPosition() (int, int) {
	return curWin().GetPosition()
}

func Alert(title, message string) error {
	return curWin().Alert(title, message)
}

func Confirm(title, message string) bool {
	return curWin().Confirm(title, message)
}

func ShowMessageBox(title, message string, buttons []string) string {
	return curWin().ShowMessageBox(title, message, buttons)
}

func HasFocus() bool {
	return curWin().HasFocus()
}

func RequestAttention(continuous bool) {
	curWin().RequestAttention(continuous)
}

func HasMouseFocus() bool {
	return curWin().HasMouseFocus()
}

func IsOpen() bool {
	return currentWindow != nil
}

func Raise() {
	curWin().Raise()
}

func Destroy() {
	gfx.DeInit()
	curWin().Destroy()
	currentWindow = nil
}
