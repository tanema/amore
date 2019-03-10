package gfx

import (
	"text/template"
)

var (
	shaderTemplate, _ = template.New("shader").Parse(`
#ifdef GL_ES
	precision highp float;
#endif
uniform mat4 TransformMat;
uniform vec4 ScreenSize;
{{.Header}}
#line 1
{{.Code}}
{{.Footer}}
`)
)

const (
	vertexHeader = `
attribute vec4 VertexPosition;
attribute vec4 VertexTexCoord;
attribute vec4 VertexColor;
attribute vec4 ConstantColor;
varying vec4 VaryingTexCoord;
varying vec4 VaryingColor;
uniform float PointSize;
`

	defaultVertexShaderCode = `
vec4 position(mat4 transformMatrix, vec4 vertexPosition) {
	return transformMatrix * vertexPosition;
}`

	vertexFooter = `
void main() {
	VaryingTexCoord = VertexTexCoord;
	VaryingColor = VertexColor * ConstantColor;
	gl_PointSize = PointSize;
	gl_Position = position(TransformMat, VertexPosition);
}`

	fragmentHeader = `
varying vec4 VaryingTexCoord;
varying vec4 VaryingColor;
uniform sampler2D Texture0;
`

	defaultFragmentShaderCode = `
vec4 effect(vec4 color, sampler2D texture, vec2 textureCoordinate, vec2 pixcoord) {
	return texture2D(texture, textureCoordinate) * color;
}`

	fragmentFooter = `
void main() {
	vec2 pixelcoord = vec2(gl_FragCoord.x, (gl_FragCoord.y * ScreenSize.z) + ScreenSize.w);
	gl_FragColor = effect(VaryingColor, Texture0, VaryingTexCoord.st, pixelcoord);
}`
)
