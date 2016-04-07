package gfx

// volatile is an interface for all items that are loaded into the gl context.
// The volatile wrapper allows you to instantiate any gl items before the context
// exist, and they will be completed after the context exists.
type volatile interface {
	loadVolatile() bool
	unloadVolatile()
}

var (
	// all registered items
	all_volatile = make(map[string][]volatile)
	// grouping of items.
	volatile_groups = []string{}
)

// init will crete a default object group for initial items to be grouped in.
func init() {
	CreateGlObjectGroup("amore_default")
}

// CurrentGlObjectGroup will get the name of the current GLObjectGroup. If you
// have not set one it will return the default group.
func CurrentGlObjectGroup() string {
	return volatile_groups[len(volatile_groups)-1]
}

// CreateGlObjectGroup will create a scope for all new gl object that will be
// instantiated. This is useful when for example when you are loading a stage/level
// that will be unloaded later. You would create the group before loading the items
// and call ReleaseGlObjectGroup when you unload the level. This allows for a bit
// of memory management.
func CreateGlObjectGroup(group string) {
	volatile_groups = append(volatile_groups, group)
	all_volatile[group] = []volatile{}
}

// ReleaseGlObjectGroup will release all the gl objects in the group with the
// group name.
func ReleaseGlObjectGroup(group string) {
	for _, vol := range all_volatile[group] {
		vol.unloadVolatile()
	}
	all_volatile[group] = []volatile{}
}

// registerVolatile will put the volatile in the current object group and call
// loadVolatile if the gl context is initialized.
func registerVolatile(new_volatile volatile) {
	current_group := CurrentGlObjectGroup()
	all_volatile[current_group] = append(all_volatile[current_group], new_volatile)
	if gl_state.initialized {
		new_volatile.loadVolatile()
	}
}

// releaseVolatile will call unloadVolatile on the volatile object
func releaseVolatile(vol volatile) {
	vol.unloadVolatile()
	for group_name, group_items := range all_volatile {
		for i, v := range group_items {
			if v == vol {
				all_volatile[group_name] = append(all_volatile[group_name][:i], all_volatile[group_name][i+1:]...)
				return
			}
		}
	}
}

func releaseAllVolatile() {
	unloadAllVolatile()
	all_volatile = make(map[string][]volatile)
	CreateGlObjectGroup("amore_default")
}

func loadAllVolatile() {
	for _, group := range all_volatile {
		for _, vol := range group {
			vol.loadVolatile()
		}
	}
}

func unloadAllVolatile() {
	for _, group := range all_volatile {
		for _, vol := range group {
			vol.unloadVolatile()
		}
	}
}
