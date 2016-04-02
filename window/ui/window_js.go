// +build js

package ui

import (
	"errors"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var (
	CurrentWindow     *Window
	document          = dom.GetWindow().Document().(dom.HTMLDocument)
	grabSupport       bool
	fullscreenSupport bool
	contextConfig     = map[string]bool{
		"alpha":                 false,
		"depth":                 true,
		"stencil":               true,
		"premultipliedAlpha":    false,
		"preserveDrawingBuffer": true,
	}
)

func InitWindowAndContext(config WindowConfig) (*Window, Context, error) {
	if js.Global.Get("WebGLRenderingContext") == js.Undefined {
		return nil, Context{}, errors.New("Your browser doesn't appear to support WebGL.")
	}

	canvas := document.CreateElement("canvas").(*dom.HTMLCanvasElement)
	canvas.Width = config.Width
	canvas.Height = config.Height
	canvas.Style().SetProperty("position", "absolute", "")
	document.Body().AppendChild(canvas)

	// Create GL context.
	contextConfig["antialias"] = (config.Msaa > 0)

	var context *js.Object
	if context = canvas.Call("getContext", "webgl", contextConfig); context == nil {
		if context = canvas.Call("getContext", "experimental-webgl", contextConfig); context == nil {
			return nil, Context{}, errors.New("Creating a WebGL context has failed.")
		}
	}

	CurrentWindow = &Window{
		HTMLCanvasElement: canvas,
		focused:           true,
		windowListeners:   make(map[string]func(*js.Object)),
		documentListeners: make(map[string]func(*js.Object)),
	}

	grabSupport = CurrentWindow.Underlying().Get("requestPointerLock") == js.Undefined || document.Underlying().Get("exitPointerLock") == js.Undefined
	fullscreenSupport = CurrentWindow.Underlying().Get("webkitRequestFullscreen") == js.Undefined || document.Underlying().Get("webkitExitFullscreen") == js.Undefined

	CurrentWindow.SetGrab(false)
	CurrentWindow.SetMinimumSize(config.Minwidth, config.Minheight)
	CurrentWindow.SetTitle(config.Title)

	if config.Centered {
		if dom.GetWindow().InnerWidth() > config.Width {
			config.X = (dom.GetWindow().InnerWidth() - config.Width) / 2
		} else {
			config.Y = 0
		}
		if dom.GetWindow().InnerHeight() > config.Height {
			config.Y = (dom.GetWindow().InnerHeight() - config.Height) / 2
		} else {
			config.Y = 0
		}
	}

	if !config.Fullscreen {
		CurrentWindow.SetPosition(config.X, config.Y)
	}

	CurrentWindow.Raise()

	CurrentWindow.bindEvents()

	return CurrentWindow, Context{context}, nil
}

var windowEventWrappers = map[string]func(event dom.Event) Event{
	"resize":              func(event dom.Event) Event { return nil },
	"gamepadconnected":    func(event dom.Event) Event { return nil },
	"gamepaddisconnected": func(event dom.Event) Event { return nil },
}

var documentEventWrappers = map[string]func(event dom.Event) Event{
	"keydown":     func(event dom.Event) Event { return nil },
	"keyup":       func(event dom.Event) Event { return nil },
	"mousedown":   func(event dom.Event) Event { return nil },
	"contextmenu": func(event dom.Event) Event { return nil },
	"mouseup":     func(event dom.Event) Event { return nil },
	"mousemove":   func(event dom.Event) Event { return nil },
	"wheel":       func(event dom.Event) Event { return nil },
	"touchstart":  func(event dom.Event) Event { return nil },
	"touchmove":   func(event dom.Event) Event { return nil },
	"touchend":    func(event dom.Event) Event { return nil },
	"focus":       func(event dom.Event) Event { return nil },
	"blur":        func(event dom.Event) Event { return nil },
}

func (window *Window) bindEvents() {
	for eventName, wrapper := range windowEventWrappers {
		window.windowListeners[eventName] = dom.GetWindow().AddEventListener(eventName, false, func(event dom.Event) {
			event_buffer = append(event_buffer, wrapper(event))
			event.PreventDefault()
		})
	}
	for eventName, wrapper := range documentEventWrappers {
		window.documentListeners[eventName] = document.AddEventListener(eventName, false, func(event dom.Event) {
			event_buffer = append(event_buffer, wrapper(event))
			event.PreventDefault()
		})
	}
}

func (window *Window) SetTitle(title string) {
	document.SetTitle(title)
}

func (window *Window) Minimize() {
	event_buffer = append(event_buffer, &WindowEvent{Event: WINDOWEVENT_HIDDEN})
	window.Style().SetProperty("display", "none", "")
}

func (window *Window) Maximize() {
	event_buffer = append(event_buffer, &WindowEvent{Event: WINDOWEVENT_SHOWN})
	window.Style().SetProperty("display", "block", "")
}

func (window *Window) IsMouseGrabbed() bool {
	return window.grabbed
}

func (window *Window) SetGrab(grabbed bool) {
	if grabbed && grabSupport {
		window.Underlying().Call("requestPointerLock")
	} else if grabSupport {
		document.Underlying().Call("exitPointerLock")
	}
	window.grabbed = grabbed
}

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

var animationFrameChan = make(chan int)

func animationFrame() {
	go func() {
		animationFrameChan <- 0
	}()
}

func (window *Window) Raise() {
	window.Maximize()
	js.Global.Call("requestAnimationFrame", animationFrame)
}

func (window *Window) Destroy() {
	if window != nil {
		document.Body().RemoveChild(window.HTMLCanvasElement)
		for eventName, listener := range window.windowListeners {
			dom.GetWindow().RemoveEventListener(eventName, false, listener)
		}
		for eventName, listener := range window.documentListeners {
			document.RemoveEventListener(eventName, false, listener)
		}
	}
}

func (window *Window) GetDrawableSize() (int, int) {
	return window.Width, window.Height
}

func (window *Window) SetIcon(path string) error {
	link := document.CreateElement("link").(*dom.HTMLLinkElement)
	link.Set("rel", "icon")
	link.Set("href", path)
	return nil
}

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

func (window *Window) HasFocus() bool {
	return window.focused
}

func (window *Window) HasMouseFocus() bool {
	return window.focused
}

// NOT SUPPORTED

func (window *Window) WarpMouseInWindow(x, y int) {}
func (window *Window) SetMinimumSize(w, h int)    {}

func (window *Window) ShowMessageBox(title, message string, buttons []string) string {
	println("ShowMessageBox not supported")
	return ""
}

func (window *Window) RequestAttention(continuous bool) {
	println("RequestAttention not supported")
}
