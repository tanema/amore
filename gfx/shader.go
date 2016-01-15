package gfx

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type BuiltinUniform int

type Shader struct {
	vertex_code string
	pixel_code  string
	program     uint32
	uniforms    map[string]Uniform // Uniform location buffer map
}

func NewShader(code ...string) *Shader {
	new_shader := &Shader{}
	new_shader.vertex_code, new_shader.pixel_code = shaderCodeToGLSL(code...)

	Register(new_shader)
	return new_shader
}

func (shader *Shader) LoadVolatile() bool {
	vert := compileCode(gl.VERTEX_SHADER, shader.vertex_code)
	frag := compileCode(gl.FRAGMENT_SHADER, shader.pixel_code)
	shader.program = gl.CreateProgram()

	gl.AttachShader(shader.program, vert)
	gl.AttachShader(shader.program, frag)

	gl.BindAttribLocation(shader.program, ATTRIB_POS, gl.Str("VertexPosition\x00"))
	gl.BindAttribLocation(shader.program, ATTRIB_TEXCOORD, gl.Str("VertexTexCoord\x00"))
	gl.BindAttribLocation(shader.program, ATTRIB_COLOR, gl.Str("VertexColor\x00"))

	gl.LinkProgram(shader.program)
	gl.DeleteShader(vert)
	gl.DeleteShader(frag)

	var isLinked int32
	gl.GetProgramiv(shader.program, gl.LINK_STATUS, &isLinked)
	if isLinked == gl.FALSE {
		var maxLength int32
		gl.GetProgramiv(shader.program, gl.INFO_LOG_LENGTH, &maxLength)

		log := make([]uint8, maxLength)
		gl.GetProgramInfoLog(shader.program, maxLength, nil, &log[0])

		gl.DeleteProgram(shader.program)
		panic(fmt.Sprintf("Cannot link program: %v", gl.GoStr(&log[0])))
	}

	shader.mapUniforms()

	return true
}

func (shader *Shader) UnloadVolatile() {
	gl.DeleteProgram(shader.program)
}

func (shader *Shader) mapUniforms() {
	// Built-in uniform locations default to -1 (nonexistant.)
	shader.uniforms = map[string]Uniform{}

	var numuniforms int32
	gl.GetProgramiv(shader.program, gl.ACTIVE_UNIFORMS, &numuniforms)

	nameBuf := make([]uint8, 256)

	for i := 0; i < int(numuniforms); i++ {
		u := Uniform{}

		gl.GetActiveUniform(shader.program, (uint32)(i), 256, nil, &u.Count, &u.Type, &nameBuf[0])
		u.Name = gl.GoStr(&nameBuf[0])
		u.Location = gl.GetUniformLocation(shader.program, gl.Str(u.Name+"\x00"))
		u.CalculateTypeInfo()

		// glGetActiveUniform appends "[0]" to the end of array uniform names...
		if len(u.Name) > 3 {
			if strings.Contains(u.Name, "[0]") {
				u.Name = u.Name[:len(u.Name)-3]
			}
		}

		if u.Location != -1 {
			shader.uniforms[u.Name] = u
		}
	}
}

func (shader *Shader) Attach() {
	gl.UseProgram(shader.program)
}

func (shader *Shader) getUniformAndCheck(name string, expected_type UniformType, value_count int) (Uniform, error) {
	uniform, ok := shader.uniforms[name]
	if !ok {
		return uniform, errors.New(fmt.Sprintf("No uniform with the name %v", name))
	}
	if uniform.BaseType != expected_type {
		return uniform, errors.New("Invalid type for uniform " + name + ". expected " + translateUniformBaseType(uniform.BaseType) + " and got " + translateUniformBaseType(expected_type))
	}
	if value_count != (int)(uniform.Count*uniform.TypeSize) {
		return uniform, errors.New(fmt.Sprintf("Invalid number of arguments for uniform  %v expected %v and got %v", name, (uniform.Count * uniform.TypeSize), value_count))
	}
	return uniform, nil
}

func (shader *Shader) SendInt(name string, values ...int32) error {
	uniform, err := shader.getUniformAndCheck(name, UNIFORM_INT, len(values))
	if err != nil {
		return err
	}

	switch uniform.TypeSize {
	case 4:
		gl.Uniform4iv(uniform.Location, uniform.Count, &values[0])
		return nil
	case 3:
		gl.Uniform3iv(uniform.Location, uniform.Count, &values[0])
		return nil
	case 2:
		gl.Uniform2iv(uniform.Location, uniform.Count, &values[0])
		return nil
	case 1:
		gl.Uniform1iv(uniform.Location, uniform.Count, &values[0])
		return nil
	}
	return errors.New("Invalid type size for uniform: " + name)
}

func (shader *Shader) SendFloat(name string, values ...float32) error {
	uniform, err := shader.getUniformAndCheck(name, UNIFORM_FLOAT, len(values))
	if err != nil {
		return err
	}

	switch uniform.TypeSize {
	case 4:
		gl.Uniform4fv(uniform.Location, uniform.Count, &values[0])
		return nil
	case 3:
		gl.Uniform3fv(uniform.Location, uniform.Count, &values[0])
		return nil
	case 2:
		gl.Uniform2fv(uniform.Location, uniform.Count, &values[0])
		return nil
	case 1:
		gl.Uniform1fv(uniform.Location, uniform.Count, &values[0])
		return nil
	}
	return errors.New("Invalid type size for uniform: " + name)
}

func (shader *Shader) SendMat4(name string, mat mgl32.Mat4) error {
	uniform, err := shader.getUniformAndCheck(name, UNIFORM_FLOAT, 4)
	if err != nil {
		return err
	}
	gl.UniformMatrix4fv(uniform.Location, uniform.Count, false, &mat[0])
	return nil
}

func (shader *Shader) SendMat3(name string, mat mgl32.Mat3) error {
	uniform, err := shader.getUniformAndCheck(name, UNIFORM_FLOAT, 3)
	if err != nil {
		return err
	}
	gl.UniformMatrix3fv(uniform.Location, uniform.Count, false, &mat[0])
	return nil
}

func (shader *Shader) SendMat2(name string, mat mgl32.Mat2) error {
	uniform, err := shader.getUniformAndCheck(name, UNIFORM_FLOAT, 3)
	if err != nil {
		return err
	}
	gl.UniformMatrix2fv(uniform.Location, uniform.Count, false, &mat[0])
	return nil
}

func (shader *Shader) sendTexture(name string, tex *Texture) error {
	uniform, err := shader.getUniformAndCheck(name, UNIFORM_SAMPLER, 1)
	if err != nil {
		return err
	}

	//REDO with pool, this is just a quick hack to make thing work right away
	gl.Uniform1i(uniform.Location, 0)
	return nil
}

func createVertexCode(code string) string {
	codes := struct {
		Syntax, Header, Uniforms, Code, Footer string
	}{
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
		Syntax, Header, Uniforms, Line, Footer, Code string
	}{
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
