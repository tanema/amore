package gfx

import "runtime"

// volatile is an interface for all items that are loaded into the gl context.
// The volatile wrapper allows you to instantiate any gl items before the context
// exist, and they will be completed after the context exists.
type volatile interface {
	loadVolatile() bool
	unloadVolatile()
}

var (
	allVolatile     = []volatile{}
	unloadcallQueue = make(chan func(), 20)
)

// registerVolatile will put the volatile in the current object group and call
// loadVolatile if the gl context is initialized.
func registerVolatile(newVolatile volatile) {
	if !glState.initialized {
		allVolatile = append(allVolatile, newVolatile)
		return
	}
	loadVolatile(newVolatile)
}

func loadAllVolatile() {
	for _, vol := range allVolatile {
		loadVolatile(vol)
	}
	allVolatile = []volatile{}
}

func loadVolatile(newVolatile volatile) {
	newVolatile.loadVolatile()
	runtime.SetFinalizer(newVolatile, func(vol volatile) {
		unloadcallQueue <- vol.unloadVolatile
	})
}

// used to make sure that the unload functions are called on the main thread
// and are called after each game loop. This is a bit more overhead but it
// make sure that the users don't need to release resources explicitly like
// some old ass technology
func cleanupVolatile() {
	for {
		select {
		case f := <-unloadcallQueue:
			f()
		default:
			return
		}
	}
}
