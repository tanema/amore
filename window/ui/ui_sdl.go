// +build !js

package ui

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"runtime"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/tanema/amore/file"
)

func PollEvent() Event {
	event := sdl.PollEvent()
	switch e := event.(type) {
	case *sdl.WindowEvent:
		return (*WindowEvent)(e)
	case *sdl.KeyDownEvent:
		return (*KeyDownEvent)(e)
	case *sdl.KeyUpEvent:
		return (*KeyUpEvent)(e)
	case *sdl.TextEditingEvent:
		return (*TextEditingEvent)(e)
	case *sdl.TextInputEvent:
		return (*TextInputEvent)(e)
	case *sdl.MouseMotionEvent:
		return (*MouseMotionEvent)(e)
	case *sdl.MouseButtonEvent:
		return (*MouseButtonEvent)(e)
	case *sdl.MouseWheelEvent:
		return (*MouseWheelEvent)(e)
	case *sdl.JoyAxisEvent:
		return (*JoyAxisEvent)(e)
	case *sdl.JoyBallEvent:
		return (*JoyBallEvent)(e)
	case *sdl.JoyHatEvent:
		return (*JoyHatEvent)(e)
	case *sdl.JoyButtonEvent:
		return (*JoyButtonEvent)(e)
	case *sdl.JoyDeviceEvent:
		return (*JoyDeviceEvent)(e)
	case *sdl.ControllerAxisEvent:
		return (*ControllerAxisEvent)(e)
	case *sdl.ControllerButtonEvent:
		return (*ControllerButtonEvent)(e)
	case *sdl.ControllerDeviceEvent:
		return (*ControllerDeviceEvent)(e)
	case *sdl.TouchFingerEvent:
		return (*TouchFingerEvent)(e)
	case *sdl.QuitEvent:
		return (*QuitEvent)(e)
	case *sdl.DropEvent:
		return (*DropEvent)(e)
	case *sdl.RenderEvent:
		return (*RenderEvent)(e)
	case *sdl.UserEvent:
		return (*UserEvent)(e)
	case *sdl.ClipboardEvent:
		return (*ClipboardEvent)(e)
	case *sdl.OSEvent:
		return (*OSEvent)(e)
	case *sdl.CommonEvent:
		return (*CommonEvent)(e)
	}
	return nil
}

func InitJoyStickAndGamePad() error {
	if err := sdl.InitSubSystem(sdl.INIT_JOYSTICK | sdl.INIT_GAMECONTROLLER); err != nil {
		return err
	}
	sdl.JoystickEventState(sdl.ENABLE)
	sdl.GameControllerEventState(sdl.ENABLE)
	return nil
}

func InitHaptic() bool {
	if sdl.WasInit(sdl.INIT_HAPTIC) == 0 && sdl.InitSubSystem(sdl.INIT_HAPTIC) != nil {
		return false
	}
	return true
}

func loadSurface(path string) (*sdl.Surface, error) {
	imgFile, new_err := file.NewFile(path)
	defer imgFile.Close()
	if new_err != nil {
		return nil, new_err
	}

	decoded_img, _, img_err := image.Decode(imgFile)
	if img_err != nil {
		return nil, img_err
	}

	bounds := decoded_img.Bounds()
	rgba := image.NewRGBA(decoded_img.Bounds())
	draw.Draw(rgba, bounds, decoded_img, image.Point{0, 0}, draw.Src)

	var rmask, gmask, bmask, amask uint32
	switch runtime.GOARCH {
	case "mips64", "ppc64":
		rmask = 0xFF000000
		gmask = 0x00FF0000
		bmask = 0x0000FF00
		amask = 0x000000FF
	default:
		rmask = 0x000000FF
		gmask = 0x0000FF00
		bmask = 0x00FF0000
		amask = 0xFF000000
	}

	return sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&rgba.Pix[0]), bounds.Dx(), bounds.Dy(), 32, bounds.Dx()*4, rmask, gmask, bmask, amask)
}

func DisableScreenSaver() {
	sdl.DisableScreenSaver()
}

func EnableScreenSaver() {
	sdl.EnableScreenSaver()
}

func GetDisplayCount() int {
	num, _ := sdl.GetNumVideoDisplays()
	return num
}

func GetDisplayName(displayindex int) string {
	return sdl.GetDisplayName(displayindex)
}

func GetFullscreenSizes(displayindex int) [][]int32 {
	var sizes [][]int32
	modes, _ := sdl.GetNumDisplayModes(displayindex)
	for i := 0; i < modes; i++ {
		var mode sdl.DisplayMode
		sdl.GetDisplayMode(displayindex, i, &mode)
		sizes = append(sizes, []int32{mode.W, mode.H})
	}
	return sizes
}

func GetDesktopDimensions(displayindex int) (int32, int32) {
	var width, height int32
	if displayindex >= 0 && displayindex < GetDisplayCount() {
		var mode sdl.DisplayMode
		sdl.GetDesktopDisplayMode(displayindex, &mode)
		width = mode.W
		height = mode.H
	}
	return width, height
}

func GetMousePosition() (int, int) {
	mx, my, _ := sdl.GetMouseState()
	return mx, my
}

func SetMouseVisible(visible bool) {
	if visible {
		sdl.ShowCursor(sdl.ENABLE)
	} else {
		sdl.ShowCursor(sdl.DISABLE)
	}
}

func GetMouseVisible() bool {
	return sdl.ShowCursor(sdl.QUERY) == sdl.ENABLE
}

//Gets whether relative mode is enabled for the mouse.
func GetRelativeMouseMode() bool {
	return sdl.GetRelativeMouseMode() != false
}

//	Sets whether relative mode is enabled for the mouse.
func SetRelativeMouseMode(isvisible bool) {
	sdl.SetRelativeMouseMode(isvisible)
}

func IsMouseDown(button MouseButton) bool {
	_, _, state := sdl.GetMouseState()

	if (uint32(button) & state) == 1 {
		return true
	}

	return false
}

func GetClipboardText() (string, error) {
	return sdl.GetClipboardText()
}

func SetClipboardText(str string) error {
	return sdl.SetClipboardText(str)
}

func GetPowerInfo() (string, int, int) {
	state, seconds, percent := sdl.GetPowerInfo()
	state_str := ""
	switch state {
	case sdl.POWERSTATE_UNKNOWN:
		state_str = "unknown"
	case sdl.POWERSTATE_ON_BATTERY:
		state_str = "battery"
	case sdl.POWERSTATE_NO_BATTERY:
		state_str = "no battery"
	case sdl.POWERSTATE_CHARGING:
		state_str = "charging"
	case sdl.POWERSTATE_CHARGED:
		state_str = "charged"
	}
	return state_str, seconds, percent
}

//Creates a new hardware Cursor object from an image.
func NewCursor(filename string, hx, hy int) (*sdl.Cursor, error) {
	surface, err := loadSurface(filename)
	if err != nil {
		return nil, err
	}
	return Cursor(sdl.CreateColorCursor(surface, hx, hy)), nil
}

//Sets the current mouse cursor.
func SetCursor(cursor Cursor) {
	sdl.SetCursor(cursor)
}

//Gets the current Cursor.
func GetCursor() Cursor {
	return Cursor(sdl.GetCursor())
}

//Gets a Cursor object representing a system-native hardware cursor.
func GetSystemCursor(name string) Cursor {
	var cursor_type sdl.SystemCursor
	switch name {
	case "hand":
		cursor_type = sdl.SYSTEM_CURSOR_HAND
	case "ibeam":
		cursor_type = sdl.SYSTEM_CURSOR_IBEAM
	case "crosshair":
		cursor_type = sdl.SYSTEM_CURSOR_CROSSHAIR
	case "wait":
		cursor_type = sdl.SYSTEM_CURSOR_WAIT
	case "waitarrow":
		cursor_type = sdl.SYSTEM_CURSOR_WAITARROW
	case "sizenwse":
		cursor_type = sdl.SYSTEM_CURSOR_SIZENWSE
	case "sizenesw":
		cursor_type = sdl.SYSTEM_CURSOR_SIZENESW
	case "sizewe":
		cursor_type = sdl.SYSTEM_CURSOR_SIZEWE
	case "sizens":
		cursor_type = sdl.SYSTEM_CURSOR_SIZENS
	case "sizeall":
		cursor_type = sdl.SYSTEM_CURSOR_SIZEALL
	case "no":
		cursor_type = sdl.SYSTEM_CURSOR_NO
	case "arrow":
		fallthrough
	default:
		cursor_type = sdl.SYSTEM_CURSOR_ARROW
	}
	return Cursor(sdl.CreateSystemCursor(cursor_type))
}

//Enables or disables text input events.
func SetTextInput(enabled bool) {
	if enabled {
		sdl.StartTextInput()
	} else {
		sdl.StopTextInput()
	}
}

//Gets whether text input events are enabled.
func HasTextInput() bool {
	return sdl.IsTextInputActive() != false
}

//Gets whether text input events are enabled.
func IsKeyDown(key Keycode) bool {
	return IsScancodeDown(GetScancodeFromKey(key))
}

//Checks whether a certain scancode is down.
func IsScancodeDown(scancode Scancode) bool {
	state := sdl.GetKeyboardState()
	for _, code := range state {
		if code == uint8(scancode) {
			return true
		}
	}
	return false
}

func GetKeyFromScancode(code Scancode) Keycode {
	return Keycode(sdl.GetKeyFromScancode(sdl.Scancode(code)))
}

func GetScancodeFromKey(key Keycode) Scancode {
	return Scancode(sdl.GetScancodeFromKey(sdl.Keycode(key)))
}

func NumJoysticks() int {
	return sdl.NumJoysticks()
}

func GetJoystickName(idx int) string {
	return sdl.JoystickGetGUIDString(sdl.JoystickGetDeviceGUID(idx))
}
