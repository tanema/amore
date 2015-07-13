package mouse

//import (
//"image"
//"os"

//"github.com/go-gl/glfw/v3.1/glfw"
//)

////Creates a new hardware Cursor object from an image.
//func NewCursor(filename string, hx, hy int) (*glfw.Cursor, error) {
//imgFile, err := os.Open(filename)
//if err != nil {
//return nil, err
//}
//defer imgFile.Close()

//img, _, err := image.Decode(imgFile)
//if err != nil {
//return nil, err
//}

//return NewImageCursor(img, hx, hy), nil
//}

////Creates a new hardware Cursor object from an image.
//func NewImageCursor(img image.Image, hx, hy int) *glfw.Cursor {
//return glfw.CreateCursor(img, hx, hy)
//}

////Sets the current mouse cursor.
//func SetCursor(cursor *glfw.Cursor) {
//glfw.GetCurrentContext().SetCursor(cursor)
//}

////Gets the current Cursor.
//func GetCursor() (*glfw.Cursor, error) {

//return nil, nil
//}

////Gets a Cursor object representing a system-native hardware cursor.
//func GetSystemCursor(name string) *glfw.Cursor {
//var cursor_type glfw.StandardCursor
//switch name {
//case "hand":
//cursor_type = glfw.HandCursor
//case "ibeam":
//cursor_type = glfw.IBeamCursor
//case "crosshair":
//cursor_type = glfw.CrosshairCursor
//case "sizeh":
//cursor_type = glfw.HResizeCursor
//case "sizev":
//cursor_type = glfw.VResizeCursor
//case "arrow":
//fallthrough
//default:
//cursor_type = glfw.ArrowCursor
//}
//return glfw.CreateStandardCursor(int(cursor_type))
//}
