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
	currentConfig     WindowConfig
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
		canvasListeners:   make(map[string]func(*js.Object)),
	}

	grabSupport = CurrentWindow.Underlying().Get("requestPointerLock") == js.Undefined || document.Underlying().Get("exitPointerLock") == js.Undefined
	fullscreenSupport = CurrentWindow.Underlying().Get("webkitRequestFullscreen") == js.Undefined || document.Underlying().Get("webkitExitFullscreen") == js.Undefined

	CurrentWindow.SetGrab(false)
	CurrentWindow.SetMinimumSize(config.Minwidth, config.Minheight)
	CurrentWindow.SetTitle(config.Title)

	currentConfig = config
	onResize(nil)
	CurrentWindow.Raise()
	CurrentWindow.bindEvents()

	return CurrentWindow, Context{context}, nil
}

var windowEventWrappers = map[string]func(event dom.Event) Event{
	"resize":              onResize,
	"gamepadconnected":    onGamepadConnected,
	"gamepaddisconnected": onGamepadDisconnected,
}

var documentEventWrappers = map[string]func(event dom.Event) Event{
	"keydown": onKeyDown, "keyup": onKeyUp,
	"visibilitychange": onVisibilityChange,
}

var canvasEventWrappers = map[string]func(event dom.Event) Event{
	"mousedown": onMouseDown, "mouseup": onMouseUp,
	"mousemove": onMouseMove, "contextmenu": onContextMenu,
	"mouseleave": onMouseOut, "mouseenter": onMouseEnter,
	"wheel": onMouseWheel, "touchstart": onTouchStart,
	"touchmove": onTouchMove, "touchend": onTouchEnd,
}

func (window *Window) bindEvents() {
	for eventName, _ := range windowEventWrappers {
		window.windowListeners[eventName] = dom.GetWindow().AddEventListener(eventName, false, func(event dom.Event) {
			event_buffer = append(event_buffer, windowEventWrappers[event.Type()](event))
			event.PreventDefault()
		})
	}
	for eventName, _ := range documentEventWrappers {
		window.documentListeners[eventName] = document.AddEventListener(eventName, false, func(event dom.Event) {
			event_buffer = append(event_buffer, documentEventWrappers[event.Type()](event))
			event.PreventDefault()
		})
	}
	for eventName, _ := range canvasEventWrappers {
		window.canvasListeners[eventName] = window.AddEventListener(eventName, false, func(event dom.Event) {
			event_buffer = append(event_buffer, canvasEventWrappers[event.Type()](event))
			event.PreventDefault()
		})
	}
}

func onResize(event dom.Event) Event {
	if currentConfig.Fullscreen {
		CurrentWindow.Style().SetProperty("width", fmt.Sprintf("%vpx", dom.GetWindow().InnerWidth()), "")
		CurrentWindow.Style().SetProperty("height", fmt.Sprintf("%vpx", dom.GetWindow().InnerHeight()), "")
	}

	var x, y int
	if currentConfig.Centered && !currentConfig.Fullscreen {
		if dom.GetWindow().InnerWidth() > CurrentWindow.Width {
			x = (dom.GetWindow().InnerWidth() - CurrentWindow.Width) / 2
		}
		if dom.GetWindow().InnerHeight() > CurrentWindow.Height {
			y = (dom.GetWindow().InnerHeight() - CurrentWindow.Height) / 2
		}
	} else if !currentConfig.Centered && !currentConfig.Fullscreen {
		x, y = currentConfig.X, currentConfig.Y
	}
	CurrentWindow.SetPosition(x, y)

	return nil
}

func onGamepadConnected(event dom.Event) Event    { return nil }
func onGamepadDisconnected(event dom.Event) Event { return nil }

func onKeyDown(event dom.Event) Event {
	ke := event.(*dom.KeyboardEvent)
	keycode := Keycode(ke.Get("code").String())
	scancode := Scancode(ke.KeyCode)

	keyMap[scancode] = true
	keyMeaningMap[keycode] = scancode
	scancodeMeaningMap[scancode] = keycode

	repeat := 0
	if ke.Repeat {
		repeat = 1
	}

	return &KeyDownEvent{
		Timestamp: int(ke.Get("timeStamp").Float() * 100),
		Repeat:    repeat,
		Keysym: Keysym{
			Scancode: ke.KeyCode,
		},
	}
}

func onKeyUp(event dom.Event) Event {
	ke := event.(*dom.KeyboardEvent)
	keyMap[Scancode(ke.KeyCode)] = false

	repeat := 0
	if ke.Repeat {
		repeat = 1
	}

	return &KeyUpEvent{
		Timestamp: int(ke.Get("timeStamp").Float() * 100),
		Repeat:    repeat,
		Keysym: Keysym{
			Scancode: ke.KeyCode,
		},
	}
}

func onMouseDown(event dom.Event) Event {
	me := event.(*dom.MouseEvent)
	mouseButtonMap[MouseButton(me.Button)] = false
	return &MouseButtonEvent{
		Timestamp: int(me.Get("timeStamp").Float() * 100),
		Type:      MOUSEBUTTONDOWN,
		Button:    me.Button,
		X:         me.Get("offsetX").Int(),
		Y:         me.Get("offsetY").Int(),
	}
}

func onMouseUp(event dom.Event) Event {
	me := event.(*dom.MouseEvent)
	mouseButtonMap[MouseButton(me.Button)] = false
	return &MouseButtonEvent{
		Timestamp: int(me.Get("timeStamp").Float() * 100),
		Type:      MOUSEBUTTONUP,
		Button:    me.Button,
		X:         me.Get("offsetX").Int(),
		Y:         me.Get("offsetY").Int(),
	}
}

func onMouseMove(event dom.Event) Event {
	me := event.(*dom.MouseEvent)
	x, y := me.Get("offsetX").Int(), me.Get("offsetY").Int()
	mdx, mdy := x-mousePos[0], y-mousePos[1]
	mousePos[0], mousePos[1] = x, y
	return &MouseMotionEvent{
		Timestamp: int(me.Get("timeStamp").Float() * 100),
		X:         x,
		Y:         y,
		XRel:      mdx,
		YRel:      mdy,
	}
}

//not used but context menu prevented
func onContextMenu(event dom.Event) Event { return nil }

func onMouseWheel(event dom.Event) Event {
	me := event.(*dom.WheelEvent)
	return &MouseWheelEvent{
		Timestamp: int(me.Get("timeStamp").Float() * 100),
		X:         int(me.DeltaX),
		Y:         int(me.DeltaY),
	}
}

func onMouseOut(event dom.Event) Event {
	return &WindowEvent{Event: WINDOWEVENT_LEAVE}
}

func onMouseEnter(event dom.Event) Event {
	return &WindowEvent{Event: WINDOWEVENT_ENTER}
}

func onTouchStart(event dom.Event) Event { return nil }
func onTouchEnd(event dom.Event) Event   { return nil }
func onTouchMove(event dom.Event) Event  { return nil }

func onVisibilityChange(event dom.Event) Event {
	if document.Underlying().Get("hidden").Bool() {
		return &WindowEvent{Event: WINDOWEVENT_FOCUS_LOST}
	} else {
		return &WindowEvent{Event: WINDOWEVENT_FOCUS_GAINED}
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
		for eventName, listener := range window.windowListeners {
			dom.GetWindow().RemoveEventListener(eventName, false, listener)
		}
		for eventName, listener := range window.documentListeners {
			document.RemoveEventListener(eventName, false, listener)
		}
		for eventName, listener := range window.canvasListeners {
			window.RemoveEventListener(eventName, false, listener)
		}
		document.Body().RemoveChild(window.HTMLCanvasElement)
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

func (window *Window) WarpMouseInWindow(x, y int) {
	mousePos[0], mousePos[1] = x, y
}

func (window *Window) SetMinimumSize(w, h int) {
	window.minWidth = w
	window.minHeight = h
}

func (window *Window) ShowMessageBox(title, message string, buttons []string) string {
	println("ShowMessageBox not supported")
	return ""
}

func (window *Window) RequestAttention(continuous bool) {
	println("RequestAttention not supported")
}
