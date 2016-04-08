package gfx

import (
	"text/template"
)

const (
	shader_syntax = `
#version 120
#define number float
#define Image sampler2D
#define extern uniform
#define Texel texture2D
#pragma optionNV(strict on)
`
	//Uniforms shared by the vertex and pixel shader stages.
	shader_uniforms = `
uniform mat4 ProjectionMat;
uniform mat4 ViewMat;
uniform mat4 ModelMat;
uniform vec4 ScreenSize;
`

	vertex_header = `
attribute vec3 VertexPosition;
attribute vec2 VertexTexCoord;
attribute vec4 VertexColor;
attribute vec4 ConstantColor;
varying vec2 VaryingTexCoord;
varying vec4 VaryingColor;
uniform float PointSize;
`

	vertex_footer = `
void main() {
	VaryingTexCoord = VertexTexCoord;
	VaryingColor = VertexColor * ConstantColor;
	gl_PointSize = PointSize;
	gl_Position = position(ProjectionMat, ViewMat, ModelMat, VertexPosition);
}`

	pixel_header = `
#define Canvases gl_FragData
varying vec2 VaryingTexCoord;
varying vec4 VaryingColor;
uniform sampler2D _tex0_;
`

	pixel_footer = `
void main() {
	// fix crashing issue in OSX when _tex0_ is unused within effect()
	float dummy = Texel(_tex0_, vec2(.5)).r;
	vec2 pixelcoord = vec2(gl_FragCoord.x, (gl_FragCoord.y * ScreenSize.z) + ScreenSize.w);
	gl_FragColor = effect(VaryingColor, _tex0_, VaryingTexCoord, pixelcoord);
}`

	footer_multi_canvas = `
void main() {
	// fix crashing issue in OSX when _tex0_ is unused within effect()
	float dummy = Texel(_tex0_, vec2(.5)).r;
	vec2 pixelcoord = vec2(gl_FragCoord.x, (gl_FragCoord.y * amore_ScreenSize.z) + amore_ScreenSize.w);
	effects(VaryingColor, _tex0_, VaryingTexCoord.st, pixelcoord);
}`

	default_vertex_shader_code = `
vec4 position(mat4 projection, mat4 view, mat4 model, vec3 vertpos) {
	return projection * view * model * vec4(vertpos, 1); 
}`

	default_pixel_shader_code = `
vec4 effect(vec4 vcolor, Image tex, vec2 texcoord, vec2 pixcoord) {
	return Texel(tex, texcoord) * vcolor;
}`
)

var (
	shader_template, _ = template.New("shader").Parse(`{{.Syntax}}
{{.Header}}
{{.Uniforms}}
#line 1
{{.Code}}
{{.Footer}}
`)
)
