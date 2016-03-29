package gfx

import (
	"runtime"
	"text/template"
)

type shaderTemplateData struct {
	Header string
	Line   string
	Footer string
	Code   string
}

func (std shaderTemplateData) ES() bool {
	if runtime.GOARCH == "js" {
		return true
	} else {
		return false
	}
}

func (std shaderTemplateData) NES() bool {
	return !std.ES()
}

func (std shaderTemplateData) Version() string {
	if std.ES() {
		return "100"
	} else {
		return "120"
	}
}

const (
	VERTEX_HEADER = `
attribute vec3 VertexPosition;
attribute vec2 VertexTexCoord;
attribute vec4 VertexColor;
attribute vec4 ConstantColor;
varying vec2 VaryingTexCoord;
varying vec4 VaryingColor;
uniform float PointSize;
`

	VERTEX_FOOTER = `
void main() {
	VaryingTexCoord = VertexTexCoord;
	VaryingColor = VertexColor * ConstantColor;
	gl_PointSize = PointSize;
	gl_Position = position(ProjectionMat, ViewMat, ModelMat, VertexPosition);
}`

	PIXEL_HEADER = `
#define PIXEL

#ifdef GL_ES
precision mediump float;
#endif

#define Canvases gl_FragData
varying vec2 VaryingTexCoord;
varying vec4 VaryingColor;
uniform sampler2D _tex0_;
`

	PIXEL_FOOTER = `
void main() {
	// fix crashing issue in OSX when _tex0_ is unused within effect()
	float dummy = Texel(_tex0_, vec2(.5)).r;
	vec2 pixelcoord = vec2(gl_FragCoord.x, (gl_FragCoord.y * ScreenSize.z) + ScreenSize.w);
	gl_FragColor = effect(VaryingColor, _tex0_, VaryingTexCoord, pixelcoord);
}`

	FOOTER_MULTI_CANVAS = `
void main() {
	// fix crashing issue in OSX when _tex0_ is unused within effect()
	float dummy = Texel(_tex0_, vec2(.5)).r;
	vec2 pixelcoord = vec2(gl_FragCoord.x, (gl_FragCoord.y * ScreenSize.z) + ScreenSize.w);
	effects(VaryingColor, _tex0_, VaryingTexCoord.st, pixelcoord);
}`

	DEFAULT_VERTEX_SHADER_CODE = `
vec4 position(mat4 projection, mat4 view, mat4 model, vec3 vertpos) {
	return projection * view * model * vec4(vertpos, 1); 
}`

	DEFAULT_PIXEL_SHADER_CODE = `
vec4 effect(vec4 vcolor, Image tex, vec2 texcoord, vec2 pixcoord) {
	return Texel(tex, texcoord) * vcolor;
}`
)

var (
	SHADER_TEMPLATE, temp_err = template.New("shader").Parse(
		`#version {{.Version}}
#ifndef GL_ES
#define lowp
#define mediump
#define highp
#pragma optionNV(strict on)
#endif
#define number float
#define Image sampler2D
#define extern uniform
#define Texel texture2D

{{.Header}}

// According to the GLSL ES 1.0 spec, uniform precision must match between stages,
// but we can't guarantee that highp is always supported in fragment shaders...
// We *really* don't want to use mediump for these in vertex shaders though.
#if defined(VERTEX) || defined(GL_FRAGMENT_PRECISION_HIGH)
#define AMORE_UNIFORM_PRECISION highp

#else
#define AMORE_UNIFORM_PRECISION mediump
#endif

uniform AMORE_UNIFORM_PRECISION mat4 ProjectionMat;
uniform AMORE_UNIFORM_PRECISION mat4 ViewMat;
uniform AMORE_UNIFORM_PRECISION mat4 ModelMat;
uniform mediump vec4 ScreenSize;

#line 1
{{.Code}}
{{.Footer}}
`)
)
