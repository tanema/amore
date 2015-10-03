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
	current *Shader
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

	if current == shader {
		// make sure glUseProgram gets called.
		current = nil
		shader.Attach(false)
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
		fmt.Println(u.Name)

		u.Location = gl.GetUniformLocation(shader.program, gl.Str(u.Name+"\x00"))
		u.BaseType = getUniformBaseType(u.Type)

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
	if current != shader {
		gl.UseProgram(shader.program)
		current = shader
		// retain/release happens in gfx.SetShader.
	}

	if !temporary {
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

func (shader *Shader) checkSetScreenParams() {
	view := opengl.GetViewport()
	if view == shader.lastViewport && shader.lastCanvas == canvas.Current {
		return
	}

	// In the shader, we do pixcoord.y = gl_FragCoord.y * params.z + params.w.
	// This lets us flip pixcoord.y when needed, to be consistent (drawing with
	// no Canvas active makes the y-values for pixel coordinates flipped.)
	params := []float32{
		(float32)(view[2]),
		(float32)(view[3]),
		0.0,
		0.0,
	}

	if canvas.Current != nil {
		// No flipping: pixcoord.y = gl_FragCoord.y * 1.0 + 0.0.
		params[2] = 1.0
		params[3] = 0.0
	} else {
		// gl_FragCoord.y is flipped when drawing to the screen, so we un-flip:
		// pixcoord.y = gl_FragCoord.y * -1.0 + height.
		params[2] = -1.0
		params[3] = (float32)(view[3])
	}

	sendBuiltinFloat(BUILTIN_SCREEN_SIZE, 4, params, 1)

	shader.lastCanvas = canvas.Current
	shader.lastViewport = view
}

func sendBuiltinFloat(builtin BuiltinUniform, size int, m []float32, count int) bool {
	return true
}
