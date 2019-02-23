package runtime

import (
	"runtime"

	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/yuin/gopher-lua"

	"github.com/tanema/amore/file"
	"github.com/tanema/amore/gfx"
)

// LuaFuncs is the declarations of functions within a module
type LuaFuncs map[string]lua.LGFunction

// LuaMetaTable is the declarations of metatables within a module
type LuaMetaTable map[string]LuaFuncs

// LuaLoadHook is a function that will be called before the gameloop starts with
// the state. This is good for fetching global callbacks to call later.
type LuaLoadHook func(*lua.LState, *glfw.Window)

type luaModule struct {
	name       string
	functions  LuaFuncs
	metatables map[string]LuaFuncs
}

var (
	registeredModules = []luaModule{}
	registeredHooks   = []LuaLoadHook{}
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread() //important OpenGl Demand it and stamp thier feet if you dont
}

// RegisterModule registers a lua module within the global namespace for easy access
func RegisterModule(name string, funcs LuaFuncs, metatables LuaMetaTable) {
	registeredModules = append(registeredModules, luaModule{name: name, functions: funcs, metatables: metatables})
}

// RegisterHook will add a lua state load hook
func RegisterHook(fn LuaLoadHook) {
	registeredHooks = append(registeredHooks, fn)
}

// Run starts the lua program
func Run(entrypoint string) error {
	if err := glfw.Init(gl.ContextWatcher); err != nil {
		return err
	}
	defer glfw.Terminate()
	win, err := createWindow(conf)
	if err != nil {
		return err
	}

	ls := lua.NewState()
	defer ls.Close()

	gfx.InitContext(win.Window)
	importGlobals(ls, win.Window)
	importModules(ls)
	runHooks(ls, win.Window)

	entryfile, err := file.Open(entrypoint)
	if err != nil {
		return err
	}

	if fn, err := ls.Load(entryfile, entrypoint); err != nil {
		return err
	} else {
		ls.Push(fn)
		if err := ls.PCall(0, lua.MultRet, nil); err != nil {
			return err
		}
	}

	if load := ls.GetGlobal("load"); load != nil {
		if err := ls.CallByParam(lua.P{Fn: load, Protect: true}); err != nil {
			return err
		}
	}

	return gameloop(ls, win)
}

func importGlobals(ls *lua.LState, win *glfw.Window) {
	ls.SetGlobal("getfps", ls.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(fps))
		return 1
	}))
	ls.SetGlobal("quit", ls.NewFunction(func(L *lua.LState) int {
		win.SetShouldClose(true)
		return 0
	}))
}

func importModules(ls *lua.LState) {
	for _, mod := range registeredModules {
		newMod := ls.NewTable()
		ls.SetFuncs(newMod, mod.functions)
		for tablename, metatable := range mod.metatables {
			mt := ls.NewTypeMetatable(tablename)
			newMod.RawSetString(tablename, mt)
			ls.SetField(mt, "__index", ls.SetFuncs(ls.NewTable(), metatable))
		}
		ls.SetGlobal(mod.name, newMod)
	}
}

func runHooks(ls *lua.LState, win *glfw.Window) {
	for _, fn := range registeredHooks {
		fn(ls, win)
	}
}

// Start creates a window and context for the game to run on and runs the game
// loop. As such this function should be put as the last call in your main function.
// update and draw will be called synchronously because calls to OpenGL that are
// not on the main thread will crash your program.
func gameloop(luaState *lua.LState, win window) error {
	for !win.ShouldClose() {
		if update := luaState.GetGlobal("update"); update != lua.LNil {
			dt := lua.LNumber(step())
			if err := luaState.CallByParam(lua.P{Fn: update, Protect: true}, dt); err != nil {
				return err
			}
		}
		if win.active {
			color := gfx.GetBackgroundColor()
			gfx.Clear(color[0], color[1], color[2], color[3])
			gfx.Origin()
			if draw := luaState.GetGlobal("draw"); draw != lua.LNil {
				if err := luaState.CallByParam(lua.P{Fn: draw, Protect: true}); err != nil {
					return err
				}
			}
			gfx.Present()
			win.SwapBuffers()
		}
		glfw.PollEvents()
	}
	return nil
}
