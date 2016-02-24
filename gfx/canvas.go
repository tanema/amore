package gfx

import (
	"fmt"
	"image"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Canvas struct {
	*Texture
	fbo              uint32
	depth_stencil    uint32
	status           uint32
	attachedCanvases []*Canvas
	width, height    int32
	systemViewport   Viewport
}

func NewCanvas(width, height int32) *Canvas {
	new_canvas := &Canvas{
		width:  width,
		height: height,
	}
	registerVolatile(new_canvas)
	return new_canvas
}

func (canvas *Canvas) loadVolatile() bool {
	canvas.status = gl.FRAMEBUFFER_COMPLETE

	// glTexImage2D is guaranteed to error in this case.
	if canvas.width > maxTextureSize || canvas.height > maxTextureSize {
		canvas.status = gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT
		return false
	}

	canvas.Texture = newTexture(canvas.width, canvas.height, false)
	//NULL means reserve texture memory, but texels are undefined
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, canvas.width, canvas.height, 0, gl.BGRA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
	if gl.GetError() != gl.NO_ERROR {
		canvas.Texture.Release()
		canvas.status = gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT
		return false
	}

	canvas.fbo, canvas.status = newFBO(canvas.GetHandle())

	if canvas.status != gl.FRAMEBUFFER_COMPLETE {
		if canvas.fbo != 0 {
			gl.DeleteFramebuffers(1, &canvas.fbo)
			canvas.fbo = 0
		}
		return false
	}

	return true
}

func (canvas *Canvas) unLoadVolatile() {
	if gl_state.currentCanvas == canvas {
		canvas.stopGrab(false)
	}
	gl.DeleteFramebuffers(1, &canvas.fbo)
	gl.DeleteRenderbuffers(1, &canvas.depth_stencil)

	canvas.fbo = 0
	canvas.depth_stencil = 0

	canvas.attachedCanvases = []*Canvas{}
	canvas.Texture.Release()
}

func (canvas *Canvas) Release() {
	releaseVolatile(canvas)
}

func (canvas *Canvas) isMultiCanvasSupported() bool {
	// system must support at least 4 simultaneous active canvases.
	return GetMaxRenderTargets() >= 4
}

func (canvas *Canvas) startGrab(canvases ...*Canvas) error {
	if gl_state.currentCanvas == canvas {
		return nil // already grabbing
	}

	if canvases != nil && len(canvases) > 0 {
		// Whether the new canvas list is different from the old one.
		// A more thorough check is done below.
		if !canvas.isMultiCanvasSupported() {
			return fmt.Errorf("Multi-canvas rendering is not supported on this system.")
		}

		if int32(len(canvases)+1) > GetMaxRenderTargets() {
			return fmt.Errorf("This system can't simultaneously render to %v canvases.", len(canvases)+1)
		}

		for i := 0; i < len(canvases); i++ {
			if canvases[i].width != canvas.width || canvases[i].height != canvas.height {
				return fmt.Errorf("All canvases must have the same dimensions.")
			}
		}
	}

	// cleanup after previous Canvas
	if gl_state.currentCanvas != nil {
		canvas.systemViewport = gl_state.currentCanvas.systemViewport
		gl_state.currentCanvas.stopGrab(true)
	} else {
		canvas.systemViewport = GetViewport()
	}

	// indicate we are using this Canvas.
	gl_state.currentCanvas = canvas
	// bind the framebuffer object.
	gl.BindFramebuffer(gl.FRAMEBUFFER, canvas.fbo)
	SetViewport(0, 0, canvas.width, canvas.height)
	// Set up the projection matrix
	gl_state.projectionStack.Push()
	gl_state.projectionStack.Load(mgl32.Ortho(0.0, float32(screen_width), 0.0, float32(screen_height), -1, 1))

	if canvases != nil && len(canvases) > 0 {
		// Attach the canvas textures to the active FBO and set up MRTs.
		drawbuffers := []uint32{gl.COLOR_ATTACHMENT0}
		// Attach the canvas textures to the currently bound framebuffer.
		for i := 0; i < len(canvases); i++ {
			buf := gl.COLOR_ATTACHMENT1 + uint32(i)
			gl.FramebufferTexture2D(gl.FRAMEBUFFER, buf, gl.TEXTURE_2D, canvases[i].GetHandle(), 0)
			drawbuffers = append(drawbuffers, buf)
		}
		// set up multiple render targets
		gl.DrawBuffers(int32(len(drawbuffers)), &drawbuffers[0])
	} else {
		// Make sure the FBO is only using a single draw buffer.
		gl.DrawBuffer(gl.COLOR_ATTACHMENT0)
	}

	canvas.attachedCanvases = canvases
	return nil
}

func (canvas *Canvas) stopGrab(switchingToOtherCanvas bool) error {
	// i am not grabbing. leave me alone
	if gl_state.currentCanvas != canvas {
		return nil
	}
	gl_state.projectionStack.Pop()
	if !switchingToOtherCanvas {
		// bind system framebuffer.
		gl_state.currentCanvas = nil
		gl.BindFramebuffer(gl.FRAMEBUFFER, getDefaultFBO())
		SetViewport(canvas.systemViewport[0], canvas.systemViewport[1], canvas.systemViewport[2], canvas.systemViewport[3])
	}
	return nil
}

func (canvas *Canvas) NewImageData(x, y, w, h int32) (image.Image, error) {
	if x < 0 || y < 0 || w <= 0 || h <= 0 || (x+w) > canvas.width || (y+h) > canvas.height {
		return nil, fmt.Errorf("Invalid ImageData rectangle dimensions.")
	}

	prev_canvas := GetCanvas()
	SetCanvas(canvas)

	screenshot := image.NewRGBA(image.Rect(int(x), int(y), int(w), int(h)))
	stride := int32(screenshot.Stride)
	pixels := make([]byte, len(screenshot.Pix))
	gl.ReadPixels(x, y, w, h, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&pixels[0]))

	for y := int32(0); y < h; y++ {
		i := (h - 1 - y) * stride
		copy(screenshot.Pix[y*stride:], pixels[i:i+w*4])
	}

	SetCanvas(prev_canvas...)

	// The new ImageData now owns the pixel data, so we don't delete it here.
	return screenshot, nil
}

func (canvas *Canvas) checkCreateStencil() bool {
	// Do nothing if we've already created the stencil buffer.
	if canvas.depth_stencil != 0 {
		return true
	}

	if gl_state.currentCanvas != canvas {
		gl.BindFramebuffer(gl.FRAMEBUFFER, canvas.fbo)
	}

	format := gl.STENCIL_INDEX8
	attachment := gl.STENCIL_ATTACHMENT

	gl.GenRenderbuffers(1, &canvas.depth_stencil)
	gl.BindRenderbuffer(gl.RENDERBUFFER, canvas.depth_stencil)
	gl.RenderbufferStorage(gl.RENDERBUFFER, uint32(format), canvas.width, canvas.height)

	// Attach the stencil buffer to the framebuffer object.
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, uint32(attachment), gl.RENDERBUFFER, canvas.depth_stencil)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	success := (gl.CheckFramebufferStatus(gl.FRAMEBUFFER) == gl.FRAMEBUFFER_COMPLETE)

	// We don't want the stencil buffer filled with garbage.
	if success {
		gl.Clear(gl.STENCIL_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	} else {
		gl.DeleteRenderbuffers(1, &canvas.depth_stencil)
		canvas.depth_stencil = 0
	}

	if gl_state.currentCanvas != nil && gl_state.currentCanvas != canvas {
		gl.BindFramebuffer(gl.FRAMEBUFFER, gl_state.currentCanvas.fbo)
	} else if gl_state.currentCanvas == nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, getDefaultFBO())
	}

	return success
}

func (canvas *Canvas) GetStatus() uint32 {
	return canvas.status
}

func newFBO(texture uint32) (uint32, uint32) {
	// get currently bound fbo to reset to it later
	current_fbo := getCurrentFBO()

	var framebuffer uint32
	gl.GenFramebuffers(1, &framebuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	if texture != 0 {
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture, 0)
		// Initialize the texture to transparent black.
		gl.ClearColor(0.0, 0.0, 0.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	}
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)

	// unbind framebuffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(current_fbo))

	return framebuffer, status
}
