package mouse

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/tanema/amore/window/surface"
)

// NewCursor creates a new hardware Cursor object from an image.
func NewCursor(filename string, hx, hy int) (*sdl.Cursor, error) {
	new_surface, err := surface.Load(filename)
	if err != nil {
		return nil, err
	}
	cursor := sdl.CreateColorCursor(new_surface, hx, hy)
	new_surface.Free()
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
	return sdl.CreateSystemCursor(cursor_type)
}
