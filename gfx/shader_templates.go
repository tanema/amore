package gfx

import (
	"text/template"
)

const (
	VERSION    = "#version 120"
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
	DEFAULT_VERTEX_SHADER_CODE = `
vec4 position(mat4 transform_proj, vec4 vertpos) {
	return transform_proj * vertpos;
}`
	DEFAULT_PIXEL_SHADER_CODE = `
vec4 effect(lowp vec4 vcolor, Image tex, vec2 texcoord, vec2 pixcoord) {
	return Texel(tex, texcoord) * vcolor;
}`
)
