package shader

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"

	"github.com/go-gl/gl/v2.1/gl"
)

const (
	VERSION_ES = "#version 100"
	SYNTAX     = `
#ifndef GL_ES
#define lowp
#define mediump
#define highp
#endif
#define number float
#define Image sampler2D
#define extern uniform
#define Texel texture2D
#pragma optionNV(strict on)
`
	//Uniforms shared by the vertex and pixel shader stages.
	UNIFORMS = `
#ifdef GL_ES
// According to the GLSL ES 1.0 spec, uniform precision must match between stages,
// but we can't guarantee that highp is always supported in fragment shaders...
// We *really* don't want to use mediump for these in vertex shaders though.
#if defined(VERTEX) || defined(GL_FRAGMENT_PRECISION_HIGH)
#define AMORE_UNIFORM_PRECISION highp
#else
#define AMORE_UNIFORM_PRECISION mediump
#endif
uniform AMORE_UNIFORM_PRECISION mat4 TransformMatrix;
uniform AMORE_UNIFORM_PRECISION mat4 ProjectionMatrix;
uniform AMORE_UNIFORM_PRECISION mat4 TransformProjectionMatrix;
#else
#define TransformMatrix gl_ModelViewMatrix
#define ProjectionMatrix gl_ProjectionMatrix
#define TransformProjectionMatrix gl_ModelViewProjectionMatrix
#endif
uniform mediump vec4 amore_ScreenSize; `

	VERTEX_HEADER = `
#define VERTEX

attribute vec4 VertexPosition;
attribute vec4 VertexTexCoord;
attribute vec4 VertexColor;

varying vec4 VaryingTexCoord;
varying vec4 VaryingColor;

#ifdef GL_ES
uniform mediump float amore_PointSize;
#endif`

	VERTEX_FOOTER = `
void main() {
	VaryingTexCoord = VertexTexCoord;
	VaryingColor = VertexColor;
#ifdef GL_ES
	gl_PointSize = amore_PointSize;
#endif
	gl_Position = position(TransformProjectionMatrix, VertexPosition);
} `

	PIXEL_HEADER = `
#define PIXEL

#ifdef GL_ES
precision mediump float;
#endif

varying mediump vec4 VaryingTexCoord;
varying lowp vec4 VaryingColor;

#define amore_Canvases gl_FragData

uniform sampler2D _tex0_;`

	PIXEL_FOOTER = `
void main() {
	// fix crashing issue in OSX when _tex0_ is unused within effect()
	float dummy = Texel(_tex0_, vec2(.5)).r;

	// See Shader::checkSetScreenParams in Shader.cpp.
	vec2 pixelcoord = vec2(gl_FragCoord.x, (gl_FragCoord.y * amore_ScreenSize.z) + amore_ScreenSize.w);

	gl_FragColor = effect(VaryingColor, _tex0_, VaryingTexCoord.st, pixelcoord);
}`

	FOOTER_MULTI_CANVAS = `void main() {
	// fix crashing issue in OSX when _tex0_ is unused within effect()
	float dummy = Texel(_tex0_, vec2(.5)).r;

	// See Shader::checkSetScreenParams in Shader.cpp.
	vec2 pixelcoord = vec2(gl_FragCoord.x, (gl_FragCoord.y * amore_ScreenSize.z) + amore_ScreenSize.w);

	effects(VaryingColor, _tex0_, VaryingTexCoord.st, pixelcoord);
}`
)

var (
	SHADER_TEMPLATE, _ = template.New("shader").Parse(
		`{{.Version}}
{{.Syntax}}
{{.Header}}
{{.Uniforms}}
#line 1
{{.Code}}
{{.Footer}}
`)
	DEFAULT_VERTEX_SHADER_CODE = createVertexCode(`
vec4 position(mat4 transform_proj, vec4 vertpos) {
	return transform_proj * vertpos;
}`)
	DEFAULT_PIXEL_SHADER_CODE = createPixelCode(`
vec4 effect(lowp vec4 vcolor, Image tex, vec2 texcoord, vec2 pixcoord) {
	return Texel(tex, texcoord) * vcolor;
}`, false)
)

func init() {
	println(DEFAULT_VERTEX_SHADER_CODE)
	println(DEFAULT_PIXEL_SHADER_CODE)
}

func createVertexCode(code string) string {
	codes := struct {
		Version, Syntax, Header, Uniforms, Code, Footer string
	}{
		Version:  VERSION_ES,
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
		Version:  VERSION_ES,
		Syntax:   SYNTAX,
		Header:   VERTEX_HEADER,
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
	csource := gl.Str(code)
	gl.ShaderSource(id, 1, &csource, nil)
	gl.CompileShader(id)

	var isCompiled int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &isCompiled)
	if isCompiled == gl.FALSE {
		var maxLength int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &maxLength)

		var log string
		gl.GetShaderInfoLog(id, maxLength, &maxLength, gl.Str(log))
		gl.DeleteShader(id)
		fmt.Println("Cannot compile shader code")
		panic(log)
	}

	return id
}
