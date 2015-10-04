package shader

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"

	"github.com/tanema/amore/gfx/canvas"
	"github.com/tanema/amore/gfx/opengl"
	"github.com/tanema/amore/gfx/volatile"
)

type BuiltinUniform int

const (
	BUILTIN_TRANSFORM_MATRIX BuiltinUniform = iota
	BUILTIN_PROJECTION_MATRIX
	BUILTIN_TRANSFORM_PROJECTION_MATRIX
	BUILTIN_POINT_SIZE
	BUILTIN_SCREEN_SIZE
	BUILTIN_MAX_ENUM
)

var (
	attribNames = map[string]int{
		"VertexPosition": opengl.ATTRIB_POS,
		"VertexTexCoord": opengl.ATTRIB_TEXCOORD,
		"VertexColor":    opengl.ATTRIB_COLOR,
	}
	builtinNames = map[string]BuiltinUniform{
		"TransformMatrix":           BUILTIN_TRANSFORM_MATRIX,
		"ProjectionMatrix":          BUILTIN_PROJECTION_MATRIX,
		"TransformProjectionMatrix": BUILTIN_TRANSFORM_PROJECTION_MATRIX,
		"amore_PointSize":           BUILTIN_POINT_SIZE,
		"amore_ScreenSize":          BUILTIN_SCREEN_SIZE,
	}
	Current       *Shader
	DefaultShader *Shader
)

type Shader struct {
	vertext_code      string
	pixel_code        string
	program           uint32
	builtinUniforms   [BUILTIN_MAX_ENUM]int32       // Location values for any built-in uniform variables.
	builtinAttributes [opengl.ATTRIB_MAX_ENUM]int32 // Location values for any generic vertex attribute variables.
	attributes        map[string]uint32
	uniforms          map[string]Uniform // Uniform location buffer map
	lastViewport      opengl.Viewport
	lastCanvas        *canvas.Canvas
	lastPointSize     float32

	// Texture unit pool for setting images
	texUnitPool      map[string]uint32      // texUnitPool[name] = textureunit
	activeTexUnits   []uint32               // activeTexUnits[textureunit-1] = textureid
	boundRetainables map[string]interface{} // Uniform name to retainable objects
}

func New(code ...string) *Shader {
	new_shader := &Shader{}
	new_shader.vertext_code, new_shader.pixel_code = shaderCodeToGLSL(code...)

	if new_shader.vertext_code == "" && new_shader.pixel_code == "" {
		panic("Cannot create shader: no source code!")
	}

	volatile.Register(new_shader)

	return new_shader
}

func (shader *Shader) LoadVolatile() bool {
	shader.lastPointSize = 0.0

	vert := compileCode(gl.VERTEX_SHADER, shader.vertext_code)
	frag := compileCode(gl.FRAGMENT_SHADER, shader.pixel_code)
	shader.program = gl.CreateProgram()

	gl.AttachShader(shader.program, vert)
	gl.AttachShader(shader.program, frag)

	// Bind generic vertex attribute indices to names in the shader.
	for name, i := range attribNames {
		gl.BindAttribLocation(shader.program, (uint32)(i), gl.Str(name+"\x00"))
	}

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

	// Retreive all active uniform variables in this shader from OpenGL.
	shader.mapActiveUniforms()

	for name, i := range attribNames {
		shader.builtinAttributes[i] = gl.GetAttribLocation(shader.program, gl.Str(name+"\x00"))
	}

	if Current == shader {
		// make sure glUseProgram gets called.
		Current = nil
		shader.Attach(false)
		shader.checkSetScreenParams()
	}

	return true
}

func (shader *Shader) UnloadVolatile() {
	if Current == shader {
		gl.UseProgram(0)
	}

	if shader.program != 0 {
		gl.DeleteProgram(shader.program)
		shader.program = 0
	}
}

func (shader *Shader) mapActiveUniforms() {
	// Built-in uniform locations default to -1 (nonexistant.)
	shader.builtinUniforms = [BUILTIN_MAX_ENUM]int32{}
	shader.uniforms = map[string]Uniform{}

	var numuniforms int32
	gl.GetProgramiv(shader.program, gl.ACTIVE_UNIFORMS, &numuniforms)

	nameBuf := make([]uint8, 256)

	for i := 0; i < int(numuniforms); i++ {
		u := Uniform{}

		gl.GetActiveUniform(shader.program, (uint32)(i), 256, nil, &u.Count, &u.Type, &nameBuf[0])
		u.Name = gl.GoStr(&nameBuf[0])

		u.Location = gl.GetUniformLocation(shader.program, gl.Str(u.Name+"\x00"))
		u.SetBaseType()

		// glGetActiveUniform appends "[0]" to the end of array uniform names...
		if len(u.Name) > 3 {
			if strings.Contains(u.Name, "[0]") {
				u.Name = u.Name[:len(u.Name)-3]
			}

			// If this is a built-in (amore-created) uniform, store the location.
			if builtin, ok := builtinNames[u.Name]; ok {
				shader.builtinUniforms[builtin] = u.Location
			}

			if u.Location != -1 {
				shader.uniforms[u.Name] = u
			}
		}
	}
}

func (shader *Shader) Attach(temporary bool) {
	if Current != shader {
		gl.UseProgram(shader.program)
	}

	if !temporary {
		Current = shader
		// make sure all sent textures are properly bound to their respective texture units
		// note: list potentially contains texture ids of deleted/invalid textures!
		for i := 0; i < len(shader.activeTexUnits); i++ {
			if shader.activeTexUnits[i] > 0 {
				opengl.BindTextureToUnit(shader.activeTexUnits[i], i+1, false)
			}
		}

		// We always want to use texture unit 0 for everyhing else.
		opengl.SetTextureUnit(0)
	}
}

func (shader *Shader) Detach() {
	if Current != nil {
		gl.UseProgram(0)
	}
	Current = nil
}

func (shader *Shader) getUniform(name string) Uniform {
	return shader.uniforms[name]
}

func (shader *Shader) checkSetUniformError(u Uniform, size, count int32, sendtype UniformType) {
	if shader.program == 0 {
		panic("No active shader program.")
	}

	if realsize := u.GetTypeSize(); size != realsize {
		panic(fmt.Sprintf("Value size of %v does not match variable size of %v.", size, realsize))
	}

	if (u.Count == 1 && count > 1) || count < 0 {
		panic(fmt.Sprintf("Invalid number of values (expected %v, got %v).", u.Count, count))
	}

	if u.BaseType == UNIFORM_SAMPLER && sendtype != u.BaseType {
		panic("Cannot send a value of this type to an Image variable.")
	}

	if (sendtype == UNIFORM_FLOAT && u.BaseType == UNIFORM_INT) || (sendtype == UNIFORM_INT && u.BaseType == UNIFORM_FLOAT) {
		panic("Cannot convert between float and int.")
	}
}

func (shader *Shader) checkSetScreenParams() {
	view := opengl.GetViewport()
	if view == shader.lastViewport && shader.lastCanvas == canvas.Current {
		return
	}

	// In the shader, we do pixcoord.y = gl_FragCoord.y * params.z + params.w.
	// This lets us flip pixcoord.y when needed, to be consistent (drawing with
	// no Canvas active makes the y-values for pixel coordinates flipped.)
	if canvas.Current != nil {
		// No flipping: pixcoord.y = gl_FragCoord.y * 1.0 + 0.0.
		shader.sendBuiltinFloat(BUILTIN_SCREEN_SIZE, 4, 1, (float32)(view[2]), (float32)(view[3]), 1.0, 0.0)
	} else {
		// gl_FragCoord.y is flipped when drawing to the screen, so we un-flip:
		// pixcoord.y = gl_FragCoord.y * -1.0 + height.
		shader.sendBuiltinFloat(BUILTIN_SCREEN_SIZE, 4, 1, (float32)(view[2]), (float32)(view[3]), -1.0, (float32)(view[3]))
	}

	shader.lastCanvas = canvas.Current
	shader.lastViewport = view
}

func (shader *Shader) checkSetPointSize(size float32) {
	if size == shader.lastPointSize {
		return
	}

	shader.sendBuiltinFloat(BUILTIN_POINT_SIZE, 1, 1, size)

	shader.lastPointSize = size
}

func (shader *Shader) sendInt(name string, size, count int32, vec ...int32) {
	shader.Attach(true)

	u := shader.getUniform(name)
	shader.checkSetUniformError(u, size, count, UNIFORM_INT)

	switch size {
	case 4:
		gl.Uniform4iv(u.Location, count, &vec[0])
	case 3:
		gl.Uniform3iv(u.Location, count, &vec[0])
	case 2:
		gl.Uniform2iv(u.Location, count, &vec[0])
	case 1:
	default:
		gl.Uniform1iv(u.Location, count, &vec[0])
	}

	Current.Attach(false)
}

func (shader *Shader) sendFloat(name string, size, count int32, vec ...float32) {
	shader.Attach(true)

	u := shader.getUniform(name)
	shader.checkSetUniformError(u, size, count, UNIFORM_FLOAT)

	switch size {
	case 4:
		gl.Uniform4fv(u.Location, count, &vec[0])
	case 3:
		gl.Uniform3fv(u.Location, count, &vec[0])
	case 2:
		gl.Uniform2fv(u.Location, count, &vec[0])
	case 1:
	default:
		gl.Uniform1fv(u.Location, count, &vec[0])
	}

	Current.Attach(false)
}

func (shader *Shader) sendMatrix(name string, size, count int32, m ...float32) {
	shader.Attach(true)

	if size < 2 || size > 4 {
		panic(fmt.Sprintf("Invalid matrix size: %dx%d (can only set 2x2, 3x3 or 4x4 matrices.)", size, size))
	}

	u := shader.getUniform(name)
	shader.checkSetUniformError(u, size, count, UNIFORM_FLOAT)

	switch size {
	case 4:
		gl.UniformMatrix4fv(u.Location, count, false, &m[0])
	case 3:
		gl.UniformMatrix3fv(u.Location, count, false, &m[0])
	case 2:
	default:
		gl.UniformMatrix2fv(u.Location, count, false, &m[0])
	}
	Current.Attach(false)
}

func (shader *Shader) sendBuiltinMatrix(builtin BuiltinUniform, size, count int32, m ...float32) bool {
	location := shader.builtinUniforms[builtin]

	if shader.builtinUniforms[builtin] == -1 {
		return false
	}

	shader.Attach(true)

	switch size {
	case 2:
		gl.UniformMatrix2fv(location, count, false, &m[0])
	case 3:
		gl.UniformMatrix3fv(location, count, false, &m[0])
	case 4:
		gl.UniformMatrix4fv(location, count, false, &m[0])
	default:
		Current.Attach(false)
		return false
	}

	Current.Attach(false)
	return true
}

func (shader *Shader) sendBuiltinFloat(builtin BuiltinUniform, size, count int32, vec ...float32) bool {
	location := shader.builtinUniforms[builtin]

	if shader.builtinUniforms[builtin] == -1 {
		return false
	}

	shader.Attach(true)

	switch size {
	case 1:
		gl.Uniform1fv(location, count, &vec[0])
	case 2:
		gl.Uniform2fv(location, count, &vec[0])
	case 3:
		gl.Uniform3fv(location, count, &vec[0])
	case 4:
		gl.Uniform4fv(location, count, &vec[0])
	default:
		Current.Attach(false)
		return false
	}

	Current.Attach(false)
	return true
}
