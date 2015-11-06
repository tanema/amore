package gfx

import (
	"text/template"
)

const (
	SYNTAX = `
#version 120
#define number float
#define Image sampler2D
#define extern uniform
#define Texel texture2D
#pragma optionNV(strict on)
`
	//Uniforms shared by the vertex and pixel shader stages.
	UNIFORMS = `
uniform mat4 ProjectionMat;
uniform mat4 ViewMat;
uniform mat4 ModelMat;
uniform vec4 ScreenSize;
`

	VERTEX_HEADER = `
attribute vec3 VertexPosition;
attribute vec2 VertexTexCoord;
attribute vec4 VertexColor;
varying vec2 VaryingTexCoord;
varying vec4 VaryingColor;
uniform float PointSize;
`

	VERTEX_FOOTER = `
void main() {
	VaryingTexCoord = VertexTexCoord;
	VaryingColor = VertexColor;
	gl_PointSize = PointSize;
	gl_Position = position(ProjectionMat, ViewMat, ModelMat, VertexPosition);
}`

	PIXEL_HEADER = `
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
	vec2 pixelcoord = vec2(gl_FragCoord.x, (gl_FragCoord.y * amore_ScreenSize.z) + amore_ScreenSize.w);
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
	SHADER_TEMPLATE, _ = template.New("shader").Parse(
		`{{.Syntax}}
{{.Header}}
{{.Uniforms}}
#line 1
{{.Code}}
{{.Footer}}
`)
)
