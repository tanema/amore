package amore

import (
	"strings"

	"github.com/goxjs/glfw"
)

type (
	keyboardInstance struct {
		keys map[string]bool
	}
	mouseInstance struct {
		isInside bool
		x, y     float64
		scrollx  float64
		scrolly  float64
		buttons  map[string]bool
	}
	InputListener func(device, button, action string, modifiers []string)
)

var (
	mouse     mouseInstance
	keyboard  keyboardInstance
	listeners = map[string][]InputListener{}
)

func On(device, button, action string, listener InputListener) {
	tag := strings.Join([]string{device, button, action}, ":")
	if _, ok := listeners[tag]; !ok {
		listeners[tag] = []InputListener{}
	}
	listeners[tag] = append(listeners[tag], listener)
}

func captureInput(window *glfw.Window) {
	mouse = mouseInstance{buttons: map[string]bool{}}
	keyboard = keyboardInstance{keys: map[string]bool{}}
	window.SetCursorEnterCallback(mouse.enter)
	window.SetMouseMovementCallback(mouse.move)
	window.SetScrollCallback(mouse.scroll)
	window.SetMouseButtonCallback(mouse.button)
	window.SetKeyCallback(keyboard.key)
}

func dispatch(device, button, action string, modifiers []string) {
	tag := strings.Join([]string{device, button, action}, ":")
	if _, ok := listeners[tag]; !ok {
		return
	}
	for _, fn := range listeners[tag] {
		fn(device, button, action, modifiers)
	}
}

func (mouse *mouseInstance) enter(w *glfw.Window, entered bool) { mouse.isInside = entered }

func (mouse *mouseInstance) scroll(w *glfw.Window, xoff, yoff float64) {
	mouse.scrollx, mouse.scrolly = xoff, yoff
}

func (mouse *mouseInstance) move(w *glfw.Window, xpos, ypos, xdelta, ydelta float64) {
	mouse.x, mouse.y = xpos, ypos
}

func (mouse *mouseInstance) button(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	buttonName := mouseButtons[button]
	if action == glfw.Press {
		mouse.buttons[buttonName] = true
	} else if action == glfw.Release {
		mouse.buttons[buttonName] = false
	}
	dispatch("mouse", buttonName, actions[action], expandModifiers(mods))
}

func (keyboard *keyboardInstance) key(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	buttonName := keyboardMap[key]
	if action == glfw.Press {
		keyboard.keys[buttonName] = true
	} else if action == glfw.Release {
		keyboard.keys[buttonName] = false
	}
	dispatch("keyboard", buttonName, actions[action], expandModifiers(mods))
}

func expandModifiers(keys glfw.ModifierKey) []string {
	mods := []string{}
	if keys&glfw.ModShift == glfw.ModShift {
		mods = append(mods, "shift")
	}
	if keys&glfw.ModControl == glfw.ModControl {
		mods = append(mods, "control")
	}
	if keys&glfw.ModAlt == glfw.ModAlt {
		mods = append(mods, "alt")
	}
	if keys&glfw.ModSuper == glfw.ModSuper {
		mods = append(mods, "super")
	}
	return mods
}
