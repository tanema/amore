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
	all_volatile = []volatile{}
)

// registerVolatile will put the volatile in the current object group and call
// loadVolatile if the gl context is initialized.
func registerVolatile(new_volatile volatile) {
	all_volatile = append(all_volatile, new_volatile)
	if gl_state.initialized {
		new_volatile.loadVolatile()
		runtime.SetFinalizer(new_volatile, new_volatile.unloadVolatile)
	}
}

func loadAllVolatile() {
	for _, vol := range all_volatile {
		vol.loadVolatile()
		runtime.SetFinalizer(vol, func(vol volatile) {
			vol.unloadVolatile()
		})
	}
}
