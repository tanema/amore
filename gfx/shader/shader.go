package shader

import (
	"github.com/go-gl/gl/v2.1/gl"
)

const (
	ATTRIB_POS = iota
	ATTRIB_TEXCOORD
	ATTRIB_COLOR
	ATTRIB_MAX_ENUM
)

var (
	current *Shader
)

type Shader struct {
	vertext_code string
	pixel_code   string
	program      uint32
}

func New(code ...string) *Shader {
	new_shader := &Shader{}
	new_shader.vertext_code, new_shader.pixel_code = shaderCodeToGLSL(code...)

	if new_shader.vertext_code == "" && new_shader.pixel_code == "" {
		panic("Cannot create shader: no source code!")
	}

	//registerVolatile(new_shader)
	new_shader.LoadVolatile()

	return new_shader
}

func (shader *Shader) LoadVolatile() bool {
	vert := compileCode(gl.VERTEX_SHADER, shader.vertext_code)
	frag := compileCode(gl.FRAGMENT_SHADER, shader.pixel_code)
	shader.program = gl.CreateProgram()

	gl.AttachShader(shader.program, vert)
	gl.AttachShader(shader.program, frag)

	gl.LinkProgram(shader.program)

	gl.DeleteShader(vert)
	gl.DeleteShader(frag)

	var isLinked int32
	gl.GetProgramiv(shader.program, gl.LINK_STATUS, &isLinked)
	if isLinked == gl.FALSE {
		var maxLength int32
		gl.GetProgramiv(shader.program, gl.INFO_LOG_LENGTH, &maxLength)

		var log string
		gl.GetProgramInfoLog(shader.program, maxLength, &maxLength, gl.Str(log))

		gl.DeleteProgram(shader.program)
		panic(log)
	}

	// Retreive all active uniform variables in this shader from OpenGL.
	shader.mapActiveUniforms()

	//for int i = 0; i < int(ATTRIB_MAX_ENUM); i++ {
	//const char *name = nullptr;
	//if attribNames.find(VertexAttribID(i), name) {
	//builtinAttributes[i] = glGetAttribLocation(program, name);
	//} else {
	//builtinAttributes[i] = -1;
	//}
	//}

	if current == shader {
		// make sure glUseProgram gets called.
		current = nil
		shader.Attach()
		shader.checkSetScreenParams()
	}

	return true
}

func (shader *Shader) UnloadVolatile() {
	if current == shader {
		gl.UseProgram(0)
	}

	if shader.program != 0 {
		gl.DeleteProgram(shader.program)
		shader.program = 0
	}
}

func (shader *Shader) mapActiveUniforms() {
}

func (shader *Shader) Attach() {
}

func (shader *Shader) checkSetScreenParams() {
}
