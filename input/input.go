package input

import (
	"github.com/goxjs/glfw"
	"github.com/yuin/gopher-lua"

	"github.com/tanema/amore/runtime"
)

type inputCapture struct {
	ls             *lua.LState
	isInside       bool
	mousex, mousey float64
	scrollx        float64
	scrolly        float64
	mouseButtons   map[string]bool
	keys           map[string]bool
}

var currentCapture inputCapture

func init() {
	runtime.RegisterHook(func(ls *lua.LState, window *glfw.Window) {
		currentCapture = inputCapture{
			ls:           ls,
			mouseButtons: map[string]bool{},
			keys:         map[string]bool{},
		}
		window.SetCursorEnterCallback(currentCapture.mouseEnter)
		window.SetMouseMovementCallback(currentCapture.mouseMove)
		window.SetScrollCallback(currentCapture.mouseScroll)
		window.SetMouseButtonCallback(currentCapture.mouseButton)
		window.SetKeyCallback(currentCapture.key)
	})
}

func (input *inputCapture) dispatch(device, button, action string, modifiers []string) {
	callback := input.ls.GetGlobal("oninput")
	if callback == nil {
		return
	}

	luaModifiers := input.ls.NewTable()
	for _, mod := range modifiers {
		luaModifiers.Append(lua.LString(mod))
	}

	input.ls.CallByParam(
		lua.P{Fn: callback, Protect: true},
		lua.LString(device),
		lua.LString(button),
		lua.LString(action),
		luaModifiers,
	)
}

func (input *inputCapture) mouseEnter(w *glfw.Window, entered bool) { input.isInside = entered }

func (input *inputCapture) mouseScroll(w *glfw.Window, xoff, yoff float64) {
	input.scrollx, input.scrolly = xoff, yoff
}

func (input *inputCapture) mouseMove(w *glfw.Window, xpos, ypos, xdelta, ydelta float64) {
	input.mousex, input.mousey = xpos, ypos
}

func (input *inputCapture) mouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	buttonName := mouseButtons[button]
	if action == glfw.Press {
		input.mouseButtons[buttonName] = true
	} else if action == glfw.Release {
		input.mouseButtons[buttonName] = false
	}
	input.dispatch("mouse", buttonName, actions[action], expandModifiers(mods))
}

func (input *inputCapture) key(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	buttonName := keyboardMap[key]
	if action == glfw.Press {
		input.keys[buttonName] = true
	} else if action == glfw.Release {
		input.keys[buttonName] = false
	}
	input.dispatch("keyboard", buttonName, actions[action], expandModifiers(mods))
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
