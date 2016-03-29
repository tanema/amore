// +build js

package ui

import (
	"errors"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var (
	CurrentWindow      *Window
	document           = dom.GetWindow().Document().(dom.HTMLDocument)
	grabSupport        bool
	fullscreenSupport  bool
	animationFrameChan = make(chan struct{})
)

func InitWindowAndContext(config WindowConfig) (*Window, Context, error) {
	canvas := document.CreateElement("canvas").(*dom.HTMLCanvasElement)

	canvas.Style().SetProperty("position", "absolute", "")
	canvas.Style().SetProperty("width", fmt.Sprintf("%vpx", config.Width), "")
	canvas.Style().SetProperty("height", fmt.Sprintf("%vpx", config.Height), "")

	document.Body().AppendChild(canvas)

	// Create GL context.
	context, err := newContext(config, canvas.Underlying())
	if err != nil {
		return nil, Context{}, err
	}

	CurrentWindow := &Window{canvas}

	grabSupport = CurrentWindow.Underlying().Get("requestPointerLock") == js.Undefined || document.Underlying().Get("exitPointerLock") == js.Undefined
	fullscreenSupport = CurrentWindow.Underlying().Get("webkitRequestFullscreen") == js.Undefined || document.Underlying().Get("webkitExitFullscreen") == js.Undefined

	CurrentWindow.SetGrab(false)
	CurrentWindow.SetMinimumSize(config.Minwidth, config.Minheight)
	CurrentWindow.SetTitle(config.Title)

	if config.Centered {
		config.X = (dom.GetWindow().InnerWidth() - config.Width) / 2
		config.Y = (dom.GetWindow().InnerHeight() - config.Height) / 2
	}

	if !config.Fullscreen {
		CurrentWindow.SetPosition(config.X, config.Y)
	}

	CurrentWindow.Raise()

	// Request first animation frame.
	js.Global.Call("requestAnimationFrame", animationFrame)

	return CurrentWindow, Context{context}, nil
}

func newContext(config WindowConfig, canvas *js.Object) (*js.Object, error) {
	if js.Global.Get("WebGLRenderingContext") == js.Undefined {
		return nil, errors.New("Your browser doesn't appear to support WebGL.")
	}

	if gl := canvas.Call("getContext", "webgl"); gl != nil {
		return gl, nil
	} else if gl := canvas.Call("getContext", "experimental-webgl"); gl != nil {
		return gl, nil
	} else {
		return nil, errors.New("Creating a WebGL context has failed.")
	}
}

func (window *Window) SetTitle(title string) {
	document.SetTitle(title)
}

func (window *Window) Minimize() {
	window.Style().SetProperty("display", "none", "")
}

func (window *Window) Maximize() {
	window.Style().SetProperty("display", "block", "")
}

func (window *Window) WarpMouseInWindow(x, y int) {}

func (window *Window) SetGrab(grabbed bool) {
	if grabbed && grabSupport {
		window.Underlying().Call("requestPointerLock")
	} else if grabSupport {
		document.Underlying().Call("exitPointerLock")
	}
}

func (window *Window) SetMinimumSize(w, h int) {}

func (window *Window) SetPosition(x, y int) {
	window.Style().SetProperty("left", fmt.Sprintf("%vpx", x), "")
	window.Style().SetProperty("top", fmt.Sprintf("%vpx", y), "")
}

func (window *Window) GetPosition() (x, y int) {
	bounds := window.GetBoundingClientRect()
	return int(bounds.Left), int(bounds.Top)
}

func (window *Window) SwapBuffers() {
	<-animationFrameChan
	js.Global.Call("requestAnimationFrame", animationFrame)
}

func animationFrame() {
	go func() {
		animationFrameChan <- struct{}{}
	}()
}

func (window *Window) Raise() {
	window.Maximize()
}

func (window *Window) Destroy() {
	if window != nil {
		document.Body().RemoveChild(window.HTMLCanvasElement)
	}
}

func (window *Window) GetDrawableSize() (int, int) { return 0, 0 }

func (window *Window) SetIcon(path string) error {
	link := document.CreateElement("link").(*dom.HTMLLinkElement)
	link.Set("rel", "icon")
	link.Set("href", path)
	return nil
}

func (window *Window) IsMouseGrabbed() bool { return false }

func (window *Window) IsVisible() bool {
	return window.Style().GetPropertyValue("display") == "block"
}

func (window *Window) Alert(title, message string) error {
	dom.GetWindow().Alert(fmt.Sprintf("%v\n%v", title, message))
	return nil
}

func (window *Window) Confirm(title, message string) bool {
	return dom.GetWindow().Confirm(fmt.Sprintf("%v\n%v", title, message))
}

func (window *Window) ShowMessageBox(title, message string, buttons []string) string {
	println("ShowMessageBox not supported")
	return ""
}

func (window *Window) HasFocus() bool      { return false }
func (window *Window) HasMouseFocus() bool { return false }

func (window *Window) RequestAttention(continuous bool) {
	println("RequestAttention not supported")
}
