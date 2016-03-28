// +build js

package ui

func PollEvent() Event                                      { return nil }
func InitJoyStickAndGamePad() error                         { return nil }
func InitHaptic() bool                                      { return true }
func DisableScreenSaver()                                   {}
func EnableScreenSaver()                                    {}
func GetDisplayCount() int                                  { return 1 }
func GetDisplayName(displayindex int) string                { return "" }
func GetFullscreenSizes(displayindex int) [][]int32         { return [][]int32{} }
func GetDesktopDimensions(displayindex int) (int32, int32)  { return 0, 0 }
func GetMousePosition() (int, int)                          { return 0, 0 }
func SetMouseVisible(visible bool)                          {}
func GetMouseVisible() bool                                 { return false }
func GetRelativeMouseMode() bool                            { return false }
func SetRelativeMouseMode(isvisible bool)                   {}
func IsMouseDown(button MouseButton) bool                   { return false }
func GetClipboardText() (string, error)                     { return "", nil }
func SetClipboardText(str string) error                     { return nil }
func GetPowerInfo() (string, int, int)                      { return "", 0, 0 }
func NewCursor(filename string, hx, hy int) (Cursor, error) { return Cursor{}, nil }
func SetCursor(cursor Cursor)                               {}
func GetCursor() Cursor                                     { return Cursor{} }
func GetSystemCursor(name string) Cursor                    { return Cursor{} }
func SetTextInput(enabled bool)                             {}
func HasTextInput() bool                                    { return false }
func IsKeyDown(key Keycode) bool                            { return false }
func IsScancodeDown(scancode Scancode) bool                 { return false }
func GetKeyFromScancode(code Scancode) Keycode              { return "" }
func GetScancodeFromKey(key Keycode) Scancode               { return 0 }
func NumJoysticks() int                                     { return 0 }
func GetJoystickName(idx int) string                        { return "" }
