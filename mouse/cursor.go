package mouse

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/tanema/amore/window/surface"
)

// NewCursor creates a new hardware Cursor object from an image.
func NewCursor(filename string, hx, hy int32) (*sdl.Cursor, error) {
	newSurface, err := surface.Load(filename)
	if err != nil {
		return nil, err
	}
	cursor := sdl.CreateColorCursor(newSurface, hx, hy)
	newSurface.Free()
	return cursor, nil
}

// SetCursor sets the current mouse cursor.
func SetCursor(cursor *sdl.Cursor) {
	sdl.SetCursor(cursor)
}

// GetCursor gets the current Cursor the program is using.
func GetCursor() *sdl.Cursor {
	return sdl.GetCursor()
}

// GetSystemCursor Gets a Cursor object representing a system-native hardware cursor.
func GetSystemCursor(name string) *sdl.Cursor {
	var cursorType sdl.SystemCursor
	switch name {
	case "hand":
		cursorType = sdl.SYSTEM_CURSOR_HAND
	case "ibeam":
		cursorType = sdl.SYSTEM_CURSOR_IBEAM
	case "crosshair":
		cursorType = sdl.SYSTEM_CURSOR_CROSSHAIR
	case "wait":
		cursorType = sdl.SYSTEM_CURSOR_WAIT
	case "waitarrow":
		cursorType = sdl.SYSTEM_CURSOR_WAITARROW
	case "sizenwse":
		cursorType = sdl.SYSTEM_CURSOR_SIZENWSE
	case "sizenesw":
		cursorType = sdl.SYSTEM_CURSOR_SIZENESW
	case "sizewe":
		cursorType = sdl.SYSTEM_CURSOR_SIZEWE
	case "sizens":
		cursorType = sdl.SYSTEM_CURSOR_SIZENS
	case "sizeall":
		cursorType = sdl.SYSTEM_CURSOR_SIZEALL
	case "no":
		cursorType = sdl.SYSTEM_CURSOR_NO
	case "arrow":
		fallthrough
	default:
		cursorType = sdl.SYSTEM_CURSOR_ARROW
	}
	return sdl.CreateSystemCursor(cursorType)
}
