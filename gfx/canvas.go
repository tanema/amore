package gfx

import (
	"fmt"
	"image"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/goxjs/gl"
)

// Canvas is an off-screen render target.
type Canvas struct {
	*Texture
	fbo            gl.Framebuffer
	depthStencil   gl.Renderbuffer
	status         uint32
	width, height  int32
	systemViewport []int32
}

// NewCanvas creates a pointer to a new canvas with the privided width and height
func NewCanvas(width, height int32) *Canvas {
	newCanvas := &Canvas{
		width:  width,
		height: height,
	}
	registerVolatile(newCanvas)
	return newCanvas
}

// loadVolatile will create the framebuffer and return true if successful
func (canvas *Canvas) loadVolatile() bool {
	canvas.status = gl.FRAMEBUFFER_COMPLETE

	// glTexImage2D is guaranteed to error in this case.
	if canvas.width > maxTextureSize || canvas.height > maxTextureSize {
		canvas.status = gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT
		return false
	}

	canvas.Texture = newTexture(canvas.width, canvas.height, false)
	//NULL means reserve texture memory, but texels are undefined
	gl.TexImage2D(gl.TEXTURE_2D, 0, int(canvas.width), int(canvas.height), gl.RGBA, gl.UNSIGNED_BYTE, nil)
	if gl.GetError() != gl.NO_ERROR {
		canvas.status = gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT
		return false
	}

	canvas.fbo, canvas.status = newFBO(canvas.getHandle())

	if canvas.status != gl.FRAMEBUFFER_COMPLETE {
		if canvas.fbo.Valid() {
			gl.DeleteFramebuffer(canvas.fbo)
			canvas.fbo = gl.Framebuffer{}
		}
		return false
	}

	return true
}

// unLoadVolatile will release the texture, framebuffer and depth buffer
func (canvas *Canvas) unLoadVolatile() {
	if glState.currentCanvas == canvas {
		canvas.stopGrab(false)
	}
	gl.DeleteFramebuffer(canvas.fbo)
	gl.DeleteRenderbuffer(canvas.depthStencil)

	canvas.fbo = gl.Framebuffer{}
	canvas.depthStencil = gl.Renderbuffer{}
}

// startGrab will bind this canvas to grab all drawing operations
func (canvas *Canvas) startGrab() error {
	if glState.currentCanvas == canvas {
		return nil // already grabbing
	}

	// cleanup after previous Canvas
	if glState.currentCanvas != nil {
		canvas.systemViewport = glState.currentCanvas.systemViewport
		glState.currentCanvas.stopGrab(true)
	} else {
		canvas.systemViewport = GetViewport()
	}

	// indicate we are using this Canvas.
	glState.currentCanvas = canvas
	// bind the framebuffer object.
	gl.BindFramebuffer(gl.FRAMEBUFFER, canvas.fbo)
	SetViewport(0, 0, canvas.width, canvas.height)
	// Set up the projection matrix
	glState.projectionStack.Push()
	glState.projectionStack.Load(mgl32.Ortho(0.0, float32(screenWidth), 0.0, float32(screenHeight), -1, 1))

	return nil
}

// stopGrab will bind the context back to the default framebuffer and set back
// all the settings
func (canvas *Canvas) stopGrab(switchingToOtherCanvas bool) error {
	// i am not grabbing. leave me alone
	if glState.currentCanvas != canvas {
		return nil
	}
	glState.projectionStack.Pop()
	if !switchingToOtherCanvas {
		// bind system framebuffer.
		glState.currentCanvas = nil
		gl.BindFramebuffer(gl.FRAMEBUFFER, getDefaultFBO())
		SetViewport(canvas.systemViewport[0], canvas.systemViewport[1], canvas.systemViewport[2], canvas.systemViewport[3])
	}
	return nil
}

// NewImageData will create an image from the canvas data. It will return an error
// only if the dimensions given are invalid
func (canvas *Canvas) NewImageData(x, y, w, h int32) (image.Image, error) {
	if x < 0 || y < 0 || w <= 0 || h <= 0 || (x+w) > canvas.width || (y+h) > canvas.height {
		return nil, fmt.Errorf("invalid ImageData rectangle dimensions")
	}

	prevCanvas := GetCanvas()
	SetCanvas(canvas)

	screenshot := image.NewRGBA(image.Rect(int(x), int(y), int(w), int(h)))
	stride := int32(screenshot.Stride)
	pixels := make([]byte, len(screenshot.Pix))
	gl.ReadPixels(pixels, int(x), int(y), int(w), int(h), gl.RGBA, gl.UNSIGNED_BYTE)

	for y := int32(0); y < h; y++ {
		i := (h - 1 - y) * stride
		copy(screenshot.Pix[y*stride:], pixels[i:i+w*4])
	}

	SetCanvas(prevCanvas)

	// The new ImageData now owns the pixel data, so we don't delete it here.
	return screenshot, nil
}

// checkCreateStencil if a stencil is set on a canvas then we need to create
// some buffers to handle this.
func (canvas *Canvas) checkCreateStencil() bool {
	// Do nothing if we've already created the stencil buffer.
	if canvas.depthStencil.Valid() {
		return true
	}

	if glState.currentCanvas != canvas {
		gl.BindFramebuffer(gl.FRAMEBUFFER, canvas.fbo)
	}

	format := gl.STENCIL_INDEX8
	attachment := gl.STENCIL_ATTACHMENT

	canvas.depthStencil = gl.CreateRenderbuffer()
	gl.BindRenderbuffer(gl.RENDERBUFFER, canvas.depthStencil)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.Enum(format), int(canvas.width), int(canvas.height))

	// Attach the stencil buffer to the framebuffer object.
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.Enum(attachment), gl.RENDERBUFFER, canvas.depthStencil)
	gl.BindRenderbuffer(gl.RENDERBUFFER, gl.Renderbuffer{})

	success := (gl.CheckFramebufferStatus(gl.FRAMEBUFFER) == gl.FRAMEBUFFER_COMPLETE)

	// We don't want the stencil buffer filled with garbage.
	if success {
		gl.Clear(gl.STENCIL_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	} else {
		gl.DeleteRenderbuffer(canvas.depthStencil)
		canvas.depthStencil = gl.Renderbuffer{}
	}

	if glState.currentCanvas != nil && glState.currentCanvas != canvas {
		gl.BindFramebuffer(gl.FRAMEBUFFER, glState.currentCanvas.fbo)
	} else if glState.currentCanvas == nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, getDefaultFBO())
	}

	return success
}

// newFBO will generate a new Frame Buffer Object for use with the canvas
func newFBO(texture gl.Texture) (gl.Framebuffer, uint32) {
	// get currently bound fbo to reset to it later
	currentFBO := gl.GetBoundFramebuffer()

	framebuffer := gl.CreateFramebuffer()
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	if texture.Valid() {
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture, 0)
		// Initialize the texture to transparent black.
		gl.ClearColor(0.0, 0.0, 0.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	}
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)

	// unbind framebuffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, currentFBO)

	return framebuffer, uint32(status)
}
