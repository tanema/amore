package mouse

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

func OnClick(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
}

func OnScroll(window *glfw.Window, xoff float64, yoff float64) {
}

//Checks whether a certain button is down.
func IsDown(button Button) bool {
	action := glfw.GetCurrentContext().GetMouseButton(glfw.MouseButton(button))
	return action == glfw.Press || action == glfw.Repeat
}

//Returns the current x-position of the mouse.
func GetX() float64 {
	x, _ := GetPosition()
	return x
}

//Returns the current y-position of the mouse.
func GetY() float64 {
	_, y := GetPosition()
	return y
}

//Sets the current X position of the mouse.
func SetX(x float64) {
	_, y := GetPosition()
	SetPosition(x, y)
}

//Sets the current Y position of the mouse.
func SetY(y float64) {
	x, _ := GetPosition()
	SetPosition(x, y)
}

//Returns the current position of the mouse.
func GetPosition() (float64, float64) {
	return glfw.GetCurrentContext().GetCursorPos()
}

//Sets the current position of the mouse.
func SetPosition(x, y float64) {
	glfw.GetCurrentContext().SetCursorPos(x, y)
}

//Gets whether relative mode is enabled for the mouse.
func GetRelativeMode() bool {
	current_mode := glfw.GetCurrentContext().GetInputMode(glfw.CursorMode)
	return current_mode == glfw.CursorDisabled
}

//	Sets whether relative mode is enabled for the mouse.
func SetRelativeMode(isvisible bool) {
	SetVisible(isvisible)
}

//Checks if the mouse is grabbed.
func IsGrabbed() bool {
	current_mode := glfw.GetCurrentContext().GetInputMode(glfw.CursorMode)
	return current_mode == glfw.CursorDisabled
}

//Grabs the mouse and confines it to the window.
func SetGrabbed(enabled bool) {
	if enabled {
		glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	} else {
		glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}
}

//Checks if the cursor is visible.
func IsVisible() bool {
	current_mode := glfw.GetCurrentContext().GetInputMode(glfw.CursorMode)
	return current_mode == glfw.CursorNormal
}

//Sets the current visibility of the cursor.
func SetVisible(isvisible bool) {
	if isvisible {
		glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	} else {
		glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}
}
