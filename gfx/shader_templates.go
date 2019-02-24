package gfx

import (
	"text/template"
)

const (
	shaderSyntax = `
//#version 120
#define number float
#define Image sampler2D
#define extern uniform
#define Texel texture2D
//#pragma optionNV(strict on)

#ifdef GL_ES
	precision highp float;
#endif

`
	//Uniforms shared by the vertex and pixel shader stages.
	shaderUniforms = `
uniform mat4 ProjectionMat;
uniform mat4 ViewMat;
uniform mat4 ModelMat;
uniform mat4 PreMulMat;
uniform vec4 ScreenSize;
`

	vertexHeader = `
attribute vec4 VertexPosition;
attribute vec4 VertexTexCoord;
attribute vec4 VertexColor;
attribute vec4 ConstantColor;
varying vec4 VaryingTexCoord;
varying vec4 VaryingColor;
uniform float PointSize;
`

	vertexFooter = `
void main() {
	VaryingTexCoord = VertexTexCoord;
	VaryingColor = VertexColor * ConstantColor;
	gl_PointSize = PointSize;
	gl_Position = position(PreMulMat, VertexPosition);
}`

	pixelHeader = `
#define Canvases gl_FragData
varying vec4 VaryingTexCoord;
varying vec4 VaryingColor;
uniform sampler2D _tex0_;
`

	pixelFooter = `
void main() {
	// fix crashing issue in OSX when _tex0_ is unused within effect()
	float dummy = Texel(_tex0_, vec2(.5)).r;
	vec2 pixelcoord = vec2(gl_FragCoord.x, (gl_FragCoord.y * ScreenSize.z) + ScreenSize.w);
	gl_FragColor = effect(VaryingColor, _tex0_, VaryingTexCoord.st, pixelcoord);
}`

	defaultVertexShaderCode = `
vec4 position(mat4 transform, vec4 vertpos) {
	return transform * vertpos;
}`

	defaultPixelShaderCode = `
vec4 effect(vec4 vcolor, Image tex, vec2 texcoord, vec2 pixcoord) {
	return Texel(tex, texcoord) * vcolor;
}`
)

var (
	shaderTemplate, _ = template.New("shader").Parse(`{{.Syntax}}
{{.Header}}
{{.Uniforms}}
#line 1
{{.Code}}
{{.Footer}}
`)
)
