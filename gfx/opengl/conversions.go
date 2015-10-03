package opengl

import (
	"github.com/go-gl/gl/v2.1/gl"
)

func Str(str string) *uint8 {
	return gl.Str(str)
}
