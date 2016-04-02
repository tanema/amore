// +build js

package ui

import (
	"fmt"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var (
	event_buffer   []Event
	currentCursor  Cursor
	mousePos       = [2]float32{}
	mouseRelative  = false
	mouseButtonMap = make(map[MouseButton]bool)
	textInput      = true
	keyMap         = make(map[Scancode]bool)
	keyMeaningMap  = make(map[string]Scancode)
	joysticks      = []Joystick{}
)

func PollEvent() Event {
	var current_event Event
	if len(event_buffer) > 0 {
		current_event, event_buffer = event_buffer[0], event_buffer[1:]
	}
	return current_event
}

func InitJoyStickAndGamePad() error {
	return nil
}

func InitHaptic() bool {
	return true
}

func GetDisplayCount() int {
	return 1
}

func GetDisplayName(displayindex int) string {
	return "browser"
}

func GetFullscreenSizes(displayindex int) [][]int32 {
	return [][]int32{
		[]int32{int32(dom.GetWindow().InnerWidth()), int32(dom.GetWindow().InnerHeight())},
	}
}

func GetDesktopDimensions(displayindex int) (int32, int32) {
	return int32(dom.GetWindow().InnerWidth()), int32(dom.GetWindow().InnerHeight())
}

func GetMousePosition() (int, int) {
	return int(mousePos[0]), int(mousePos[1])
}

func SetMouseVisible(visible bool) {
	if visible {
		SetCursor(GetCursor())
	} else {
		//make sure we have a back up cursor
		if currentCursor.icon == "" {
			currentCursor = GetCursor()
		}
		document.Body().Style().SetProperty("cursor", "none", "")
	}
}

func GetMouseVisible() bool {
	return document.Body().Style().GetPropertyValue("cursor") != "none"
}

func SetRelativeMouseMode(is_relative bool) {
	mouseRelative = is_relative
}

func IsMouseDown(button MouseButton) bool {
	down, ok := mouseButtonMap[button]
	return down && ok
}

func NewCursor(filename string, hx, hy int) (Cursor, error) {
	return Cursor{icon: fmt.Sprintf("url(%v)", filename), hx: hx, hy: hy}, nil
}

func SetCursor(cursor Cursor) {
	if strings.HasPrefix(cursor.icon, "url(") {
		document.Body().Style().SetProperty("cursor", fmt.Sprintf("%v %v %v", cursor.icon, cursor.hx, cursor.hy), "")
	} else {
		document.Body().Style().SetProperty("cursor", cursor.icon, "")
	}
	currentCursor = cursor
}

func GetCursor() Cursor {
	if currentCursor.icon == "" {
		return Cursor{icon: document.Body().Style().GetPropertyValue("cursor")}
	}
	return currentCursor
}

func GetSystemCursor(name string) Cursor {
	var cursor_type string
	switch name {
	case "hand":
		cursor_type = "pointer"
	case "ibeam":
		cursor_type = "text"
	case "wait":
		cursor_type = "progress"
	case "waitarrow":
		cursor_type = "wait"
	case "sizenwse":
		cursor_type = "nwse-resize"
	case "sizenesw":
		cursor_type = "nesw-resize"
	case "sizewe":
		cursor_type = "ew-resize"
	case "sizens":
		cursor_type = "ns-resize"
	case "sizeall":
		cursor_type = "move"
	case "no":
		cursor_type = "not-allowed"
	case "arrow":
		cursor_type = "default"
	case "crosshair":
		cursor_type = name
	}
	return Cursor{icon: cursor_type}
}

func SetTextInput(enabled bool) {
	textInput = enabled
}

func HasTextInput() bool {
	return textInput
}

func IsKeyDown(key Keycode) bool {
	return IsScancodeDown(GetScancodeFromKey(key))
}

func IsScancodeDown(scancode Scancode) bool {
	down, ok := keyMap[scancode]
	return down && ok
}

func GetKeyFromScancode(code Scancode) Keycode {
	return Keycode(js.Global.Get("String").Call("fromCharCode", int(code)).String())
}

func GetScancodeFromKey(key Keycode) Scancode {
	return Scancode(js.MakeWrapper(string(key)).Call("charCodeAt", 0).Int())
}

func NumJoysticks() int {
	return len(joysticks)
}

func GetJoystickName(idx int) string {
	if idx < 0 || idx > len(joysticks) {
		return ""
	}
	return joysticks[idx].name
}

// NOT SUPPORTED

func DisableScreenSaver()               {}
func EnableScreenSaver()                {}
func GetClipboardText() (string, error) { return "", nil }
func SetClipboardText(str string) error { return nil }
func GetPowerInfo() (string, int, int)  { return "", 0, 0 }
