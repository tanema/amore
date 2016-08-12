package gfx

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/tanema/amore/file"
	"github.com/tanema/amore/gfx/gl"
)

type Shader struct {
	vertex_code    string
	pixel_code     string
	program        gl.Program
	uniforms       map[string]uniform // uniform location buffer map
	texUnitPool    map[string]int
	activeTexUnits []gl.Texture
}

func NewShader(paths ...string) *Shader {
	new_shader := &Shader{}
	code := pathsToCode(paths...)
	new_shader.vertex_code, new_shader.pixel_code = shaderCodeToGLSL(code...)
	registerVolatile(new_shader)
	return new_shader
}

func (shader *Shader) loadVolatile() bool {
	vert := compileCode(gl.VERTEX_SHADER, shader.vertex_code)
	frag := compileCode(gl.FRAGMENT_SHADER, shader.pixel_code)
	shader.program = gl.CreateProgram()
	shader.texUnitPool = make(map[string]int)
	shader.activeTexUnits = make([]gl.Texture, maxTextureUnits)

	gl.AttachShader(shader.program, vert)
	gl.AttachShader(shader.program, frag)

	gl.BindAttribLocation(shader.program, attrib_pos, "VertexPosition")
	gl.BindAttribLocation(shader.program, attrib_texcoord, "VertexTexCoord")
	gl.BindAttribLocation(shader.program, attrib_color, "VertexColor")
	gl.BindAttribLocation(shader.program, attrib_constantcolor, "ConstantColor")

	gl.LinkProgram(shader.program)
	gl.DeleteShader(vert)
	gl.DeleteShader(frag)

	if gl.GetProgrami(shader.program, gl.LINK_STATUS) == 0 {
		gl.DeleteProgram(shader.program)
		panic(fmt.Errorf("shader link error: %s", gl.GetProgramInfoLog(shader.program)))
	}

	shader.mapUniforms()

	return true
}

func (shader *Shader) unloadVolatile() {
	// active texture list is probably invalid, clear it
	gl.DeleteProgram(shader.program)

	// decrement global texture id counters for texture units which had textures bound from this shader
	for i := 0; i < len(shader.activeTexUnits); i++ {
		if shader.activeTexUnits[i].Valid() {
			gl_state.textureCounters[i] = gl_state.textureCounters[i] - 1
		}
	}
}

func (shader *Shader) Release() {
	releaseVolatile(shader)
}

func (shader *Shader) mapUniforms() {
	// Built-in uniform locations default to -1 (nonexistent.)
	shader.uniforms = map[string]uniform{}

	for i := 0; i < gl.GetProgrami(shader.program, gl.ACTIVE_UNIFORMS); i++ {
		u := uniform{}

		u.Name, u.Count, u.Type = gl.GetActiveUniform(shader.program, uint32(i))
		u.Location = gl.GetUniformLocation(shader.program, u.Name)
		u.CalculateTypeInfo()

		// glGetActiveUniform appends "[0]" to the end of array uniform names...
		if len(u.Name) > 3 {
			if strings.Contains(u.Name, "[0]") {
				u.Name = u.Name[:len(u.Name)-3]
			}
		}

		if u.Location.Value != -1 {
			shader.uniforms[u.Name] = u
		}
	}
}

func (shader *Shader) attach(temporary bool) {
	if gl_state.currentShader != shader {
		gl.UseProgram(shader.program)
		gl_state.currentShader = shader
	}
	if !temporary {
		// make sure all sent textures are properly bound to their respective texture units
		// note: list potentially contains texture ids of deleted/invalid textures!
		for i := 0; i < len(shader.activeTexUnits); i++ {
			if shader.activeTexUnits[i].Valid() {
				bindTextureToUnit(shader.activeTexUnits[i], i+1, false)
			}
		}

		// We always want to use texture unit 0 for everyhing else.
		setTextureUnit(0)
	}
}

func (shader *Shader) getUniformAndCheck(name string, expected_type UniformType, value_count int) (uniform, error) {
	u, ok := shader.uniforms[name]
	if !ok {
		return u, errors.New(fmt.Sprintf("No uniform with the name %v", name))
	}
	if u.BaseType != expected_type {
		return u, errors.New("Invalid type for uniform " + name + ". expected " + translateUniformBaseType(u.BaseType) + " and got " + translateUniformBaseType(expected_type))
	}
	if value_count != u.Count*u.TypeSize {
		return u, errors.New(fmt.Sprintf("Invalid number of arguments for uniform  %v expected %v and got %v", name, (u.Count * u.TypeSize), value_count))
	}
	return u, nil
}

func (shader *Shader) SendInt(name string, values ...int32) error {
	shader.attach(true)
	defer states.back().shader.attach(false)

	u, err := shader.getUniformAndCheck(name, UNIFORM_INT, len(values))
	if err != nil {
		return err
	}

	switch u.TypeSize {
	case 4:
		gl.Uniform4iv(u.Location, values)
		return nil
	case 3:
		gl.Uniform3iv(u.Location, values)
		return nil
	case 2:
		gl.Uniform2iv(u.Location, values)
		return nil
	case 1:
		gl.Uniform1iv(u.Location, values)
		return nil
	}
	return errors.New("Invalid type size for uniform: " + name)
}

func (shader *Shader) SendFloat(name string, values ...float32) error {
	shader.attach(true)
	defer states.back().shader.attach(false)

	u, err := shader.getUniformAndCheck(name, UNIFORM_FLOAT, len(values))
	if err != nil {
		return err
	}

	switch u.TypeSize {
	case 4:
		gl.Uniform4fv(u.Location, values)
		return nil
	case 3:
		gl.Uniform3fv(u.Location, values)
		return nil
	case 2:
		gl.Uniform2fv(u.Location, values)
		return nil
	case 1:
		gl.Uniform1fv(u.Location, values)
		return nil
	}
	return errors.New("Invalid type size for uniform: " + name)
}

func (shader *Shader) SendMat4(name string, mat ...float32) error {
	shader.attach(true)
	defer states.back().shader.attach(false)

	u, err := shader.getUniformAndCheck(name, UNIFORM_FLOAT, 4)
	if err != nil {
		return err
	}
	gl.UniformMatrix4fv(u.Location, []float32{
		mat[0], mat[1], mat[2], mat[3],
		mat[4], mat[5], mat[6], mat[7],
		mat[8], mat[9], mat[10], mat[11],
		mat[12], mat[13], mat[14], mat[15],
	})
	return nil
}

func (shader *Shader) SendMat3(name string, mat ...float32) error {
	shader.attach(true)
	defer states.back().shader.attach(false)

	u, err := shader.getUniformAndCheck(name, UNIFORM_FLOAT, 3)
	if err != nil {
		return err
	}
	gl.UniformMatrix3fv(u.Location, []float32{
		mat[0], mat[1], mat[2],
		mat[3], mat[4], mat[5],
		mat[6], mat[7], mat[8],
	})
	return nil
}

func (shader *Shader) SendMat2(name string, mat ...float32) error {
	shader.attach(true)
	defer states.back().shader.attach(false)

	u, err := shader.getUniformAndCheck(name, UNIFORM_FLOAT, 3)
	if err != nil {
		return err
	}
	gl.UniformMatrix2fv(u.Location, []float32{
		mat[0], mat[1],
		mat[2], mat[3],
	})
	return nil
}

func (shader *Shader) SendTexture(name string, texture iTexture) error {
	shader.attach(true)
	defer states.back().shader.attach(false)

	gltex := texture.getHandle()
	texunit := shader.getTextureUnit(name)

	u, err := shader.getUniformAndCheck(name, UNIFORM_SAMPLER, 1)
	if err != nil {
		return err
	}

	bindTextureToUnit(gltex, texunit, true)

	gl.Uniform1i(u.Location, int(texunit))

	// increment global shader texture id counter for this texture unit, if we haven't already
	if !shader.activeTexUnits[texunit-1].Valid() {
		gl_state.textureCounters[texunit-1]++
	}

	// store texture id so it can be re-bound to the proper texture unit later
	shader.activeTexUnits[texunit-1] = gltex

	return nil
}

func (shader *Shader) getTextureUnit(name string) int {
	unit, found := shader.texUnitPool[name]
	if found {
		return unit
	}

	texunit := -1
	// prefer texture units which are unused by all other shaders
	for i := 0; i < len(gl_state.textureCounters); i++ {
		if gl_state.textureCounters[i] == 0 {
			texunit = i + 1
			break
		}
	}

	if texunit == -1 {
		// no completely unused texture units exist, try to use next free slot in our own list
		for i := 0; i < len(shader.activeTexUnits); i++ {
			if !shader.activeTexUnits[i].Valid() {
				texunit = i + 1
				break
			}
		}

		if texunit == -1 {
			panic("No more texture units available for shader.")
		}
	}

	shader.texUnitPool[name] = texunit
	return shader.texUnitPool[name]
}

func createVertexCode(code string) string {
	codes := struct {
		Syntax, Header, Uniforms, Code, Footer string
	}{
		Syntax:   shader_syntax,
		Header:   vertex_header,
		Uniforms: shader_uniforms,
		Code:     code,
		Footer:   vertex_footer,
	}

	var template_writer bytes.Buffer
	err := shader_template.Execute(&template_writer, codes)
	if err != nil {
		panic(err)
	}

	return template_writer.String()
}

func createPixelCode(code string, is_multicanvas bool) string {
	codes := struct {
		Syntax, Header, Uniforms, Line, Footer, Code string
	}{
		Syntax:   shader_syntax,
		Header:   pixel_header,
		Uniforms: shader_uniforms,
		Code:     code,
	}

	if is_multicanvas {
		codes.Footer = footer_multi_canvas
	} else {
		codes.Footer = pixel_footer
	}

	var template_writer bytes.Buffer
	err := shader_template.Execute(&template_writer, codes)
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

//convert paths to strings of code
//if string is already code just pass it along
func pathsToCode(paths ...string) []string {
	code := []string{}
	if paths != nil {
		for _, path := range paths {
			//if this is not code it must be a path
			is_pixel, _ := isPixelCode(path)
			if !isVertexCode(path) && !is_pixel {
				code = append(code, file.ReadString(path))
			} else { //it is code!
				code = append(code, path)
			}
		}
	}
	return code
}

func shaderCodeToGLSL(code ...string) (string, string) {
	vertexcode := default_vertex_shader_code
	pixelcode := default_pixel_shader_code
	is_multicanvas := false // whether pixel code has "effects" function instead of "effect"

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

	return createVertexCode(vertexcode), createPixelCode(pixelcode, is_multicanvas)
}

func compileCode(shaderType gl.Enum, src string) gl.Shader {
	shader := gl.CreateShader(shaderType)
	if !shader.Valid() {
		panic(fmt.Errorf("could not create shader (type %v)", shaderType))
	}
	gl.ShaderSource(shader, src)
	gl.CompileShader(shader)
	if gl.GetShaderi(shader, gl.COMPILE_STATUS) == 0 {
		defer gl.DeleteShader(shader)
		panic(fmt.Errorf("shader compile: %s", gl.GetShaderInfoLog(shader)))
	}
	return shader
}
