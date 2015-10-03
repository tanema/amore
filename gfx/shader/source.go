package shader

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/go-gl/gl/v2.1/gl"
)

func createVertexCode(code string) string {
	codes := struct {
		Version, Syntax, Header, Uniforms, Code, Footer string
	}{
		Version:  VERSION,
		Syntax:   SYNTAX,
		Header:   VERTEX_HEADER,
		Uniforms: UNIFORMS,
		Code:     code,
		Footer:   VERTEX_FOOTER,
	}

	var template_writer bytes.Buffer
	err := SHADER_TEMPLATE.Execute(&template_writer, codes)
	if err != nil {
		panic(err)
	}

	return template_writer.String()
}

func createPixelCode(code string, is_multicanvas bool) string {
	codes := struct {
		Version, Syntax, Header, Uniforms, Line, Footer, Code string
	}{
		Version:  VERSION,
		Syntax:   SYNTAX,
		Header:   PIXEL_HEADER,
		Uniforms: UNIFORMS,
		Code:     code,
	}

	if is_multicanvas {
		codes.Footer = FOOTER_MULTI_CANVAS
	} else {
		codes.Footer = PIXEL_FOOTER
	}

	var template_writer bytes.Buffer
	err := SHADER_TEMPLATE.Execute(&template_writer, codes)
	if err != nil {
		panic(err)
	}

	return template_writer.String()
}

func isVertexCode(code string) bool {
	match, _ := regexp.MatchString(`vec4\s+position\s*\(`, code)
	return match
}

func isPixelCode(code string) (bool, bool) {
	if match, _ := regexp.MatchString(`vec4\s+effect\s*\(`, code); match {
		return true, false
	} else if match, _ := regexp.MatchString(`vec4\s+effects\s*\(`, code); match {
		// function for rendering to multiple canvases simultaneously
		return true, true
	}
	return false, false
}

func shaderCodeToGLSL(code ...string) (string, string) {
	vertexcode := DEFAULT_VERTEX_SHADER_CODE
	pixelcode := DEFAULT_PIXEL_SHADER_CODE
	is_multicanvas := false // whether pixel code has "effects" function instead of "effect"

	if code != nil {
		for _, shader_code := range code {
			if isVertexCode(shader_code) {
				vertexcode = shader_code
			}

			ispixel, isMultiCanvas := isPixelCode(shader_code)
			if ispixel {
				pixelcode = shader_code
				is_multicanvas = isMultiCanvas
			}
		}
	}

	return createVertexCode(vertexcode), createPixelCode(pixelcode, is_multicanvas)
}

func compileCode(shader_type uint32, code string) uint32 {
	id := gl.CreateShader(shader_type)
	csource := gl.Str(code + "\x00")
	gl.ShaderSource(id, 1, &csource, nil)
	gl.CompileShader(id)

	var isCompiled int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &isCompiled)
	if isCompiled == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &logLength)

		logBuffer := make([]uint8, logLength)
		gl.GetShaderInfoLog(id, logLength, nil, &logBuffer[0])
		panic(fmt.Sprintf("Cannot compile shader code: %v", gl.GoStr(&logBuffer[0])))
	}

	return id
}
