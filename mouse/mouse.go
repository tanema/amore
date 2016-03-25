// The mouse Package handles the mouse events from the gl context
package mouse

import (
	"github.com/tanema/amore/window"
	"github.com/tanema/amore/window/ui"
)

//Checks whether a certain button is down.
func IsDown(button MouseButton) bool {
	return ui.IsMouseDown(ui.MouseButton(button))
}

//Returns the current x-position of the mouse.
func GetX() float32 {
	x, _ := GetPosition()
	return x
}

//Returns the current y-position of the mouse.
func GetY() float32 {
	_, y := GetPosition()
	return y
}

//Sets the current X position of the mouse.
func SetX(x float32) {
	_, y := GetPosition()
	SetPosition(x, y)
}

//Sets the current Y position of the mouse.
func SetY(y float32) {
	x, _ := GetPosition()
	SetPosition(x, y)
}

//Returns the current position of the mouse.
func GetPosition() (float32, float32) {
	return window.GetMousePosition()
}

//Sets the current position of the mouse.
func SetPosition(x, y float32) {
	window.SetMousePosition(x, y)
}

//Gets whether relative mode is enabled for the mouse.
func GetRelativeMode() bool {
	return ui.GetRelativeMouseMode()
}

//	Sets whether relative mode is enabled for the mouse.
func SetRelativeMode(isvisible bool) {
	ui.SetRelativeMouseMode(isvisible)
}

//Checks if the mouse is grabbed.
func IsGrabbed() bool {
	return window.IsMouseGrabbed()
}

//Grabs the mouse and confines it to the window.
func SetGrabbed(enabled bool) {
	window.SetMouseGrab(enabled)
}

//Checks if the cursor is visible.
func IsVisible() bool {
	return ui.GetMouseVisible()
}

//Sets the current visibility of the cursor.
func SetVisible(isvisible bool) {
	ui.SetMouseVisible(isvisible)
}

//Creates a new hardware Cursor object from an image.
func NewCursor(filename string, hx, hy int) (ui.Cursor, error) {
	return ui.NewCursor(filename, hx, hy)
}

//Sets the current mouse cursor.
func SetCursor(cursor ui.Cursor) {
	ui.SetCursor(cursor)
}

//Gets the current Cursor.
func GetCursor() ui.Cursor {
	return ui.GetCursor()
}

//Gets a Cursor object representing a system-native hardware cursor.
func GetSystemCursor(name string) ui.Cursor {
	return ui.GetSystemCursor(name)
}
