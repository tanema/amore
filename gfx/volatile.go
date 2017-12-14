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
	allVolatile = []volatile{}
)

// registerVolatile will put the volatile in the current object group and call
// loadVolatile if the gl context is initialized.
func registerVolatile(newVolatile volatile) {
	allVolatile = append(allVolatile, newVolatile)
	if glState.initialized {
		newVolatile.loadVolatile()
		runtime.SetFinalizer(newVolatile, func(vol volatile) {
			vol.unloadVolatile()
		})
	}
}

func loadAllVolatile() {
	for _, vol := range allVolatile {
		vol.loadVolatile()
		runtime.SetFinalizer(vol, func(vol volatile) {
			vol.unloadVolatile()
		})
	}
}
