package gfx

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
)

type ShaderSource struct {
	vertex, pixel string
}

const (
	DEFAULT_VERTEX_SHADER = `vec4 position(mat4 transform_proj, vec4 vertpos) {
						return transform_proj * vertpos;
					}`
	DEFAULT_PIXEL_SHADER = `vec4 effect(lowp vec4 vcolor, Image tex, vec2 texcoord, vec2 pixcoord) {
						return Texel(tex, texcoord) * vcolor;
					}`
)

type Shader struct {
	program uint32
}

var (
	defaultShader         *Shader
	currentShader         *Shader
	default_shader_source = &ShaderSource{
		vertex: DEFAULT_VERTEX_SHADER,
		pixel:  DEFAULT_PIXEL_SHADER,
	}
)

func NewShader(vert_string, frag_string string) (*Shader, error) {
	vert, vert_err := generateGenericShader(vert_string, gl.VERTEX_SHADER)
	frag, frag_err := generateGenericShader(frag_string, gl.FRAGMENT_SHADER)
	if vert_err != nil || frag_err != nil {
		return nil, fmt.Errorf("error compiling shader(s) (%v) \n (%v)", vert_err, frag_err)
	}
	program, err := generateProgram(vert, frag)
	if err != nil {
		return nil, err
	}

	return &Shader{
		program: program,
	}, nil
}

func generateProgram(vert, frag uint32) (uint32, error) {
	program := gl.CreateProgram()

	gl.AttachShader(program, vert)
	gl.AttachShader(program, frag)

	gl.LinkProgram(program)

	var isLinked int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &isLinked)
	if isLinked == gl.FALSE {
		var maxLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &maxLength)

		var log string
		gl.GetProgramInfoLog(program, maxLength, &maxLength, gl.Str(log))

		gl.DeleteProgram(program)
		gl.DeleteShader(vert)
		gl.DeleteShader(frag)

		return 0, errors.New(log)
	}

	gl.DetachShader(program, vert)
	gl.DetachShader(program, frag)

	return program, nil
}

func generateGenericShader(source string, shader_type uint32) (uint32, error) {
	id := gl.CreateShader(shader_type)
	csource := gl.Str(source)
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
		return 0, errors.New(log)
	}

	return id, nil
}

func (shader *Shader) Use() {
	gl.UseProgram(shader.program)
}

func (shader *Shader) Sendi(var_name string, data int32) {
	gl.ProgramUniform1i(shader.program, shader.getVarNameLocation(var_name), data)
}

func (shader *Shader) Sendiv(var_name string, data []int32) {
	location := shader.getVarNameLocation(var_name)
	switch len(data) {
	case 4:
		gl.ProgramUniform4i(shader.program, location, data[0], data[1], data[2], data[3])
	case 3:
		gl.ProgramUniform3i(shader.program, location, data[0], data[1], data[2])
	case 2:
		gl.ProgramUniform2i(shader.program, location, data[0], data[1])
	case 1:
		shader.Sendi(var_name, data[0])
	}
}

func (shader *Shader) Sendf(var_name string, data float32) {
	gl.ProgramUniform1f(shader.program, shader.getVarNameLocation(var_name), data)
}

func (shader *Shader) Sendfv(var_name string, data []float32) {
	location := shader.getVarNameLocation(var_name)
	switch len(data) {
	case 4:
		gl.ProgramUniform4f(shader.program, location, data[0], data[1], data[2], data[3])
	case 3:
		gl.ProgramUniform3f(shader.program, location, data[0], data[1], data[2])
	case 2:
		gl.ProgramUniform2f(shader.program, location, data[0], data[1])
	case 1:
		shader.Sendf(var_name, data[0])
	}
}

func (shader *Shader) Sendb(var_name string, data bool) {
	if data {
		shader.Sendf(var_name, 1)
	} else {
		shader.Sendf(var_name, 0)
	}
}

func (shader *Shader) getVarNameLocation(var_name string) int32 {
	location := gl.GetUniformLocation(shader.program, gl.Str(var_name))
	if location == -1 {
		panic(fmt.Errorf("There is no uniform %v in this shader", var_name))
	}
	return location
}
