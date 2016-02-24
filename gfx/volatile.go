package gfx

type Volatile interface {
	loadVolatile() bool
	unloadVolatile()
}

var (
	all_volatile    = make(map[string][]Volatile)
	volatile_groups = []string{}
)

func init() {
	CreateGlObjectGroup("amore_default")
}

func CurrentGlObjectGroup() string {
	return volatile_groups[len(volatile_groups)-1]
}

func CreateGlObjectGroup(group string) {
	volatile_groups = append(volatile_groups, group)
	all_volatile[group] = []Volatile{}
}

func ReleaseGlObjectGroup(group string) {
	for _, vol := range all_volatile[group] {
		vol.unloadVolatile()
	}
	all_volatile[group] = []Volatile{}
}

func registerVolatile(new_volatile Volatile) {
	current_group := CurrentGlObjectGroup()
	all_volatile[current_group] = append(all_volatile[current_group], new_volatile)
	if gl_state.initialized {
		new_volatile.loadVolatile()
	}
}

func releaseVolatile(vol Volatile) {
	for group_name, group_items := range all_volatile {
		for i, v := range group_items {
			if v == vol {
				vol.unloadVolatile()
				all_volatile[group_name] = append(all_volatile[group_name][:i], all_volatile[group_name][i+1:]...)
			}
		}
	}
}

func releaseAllVolatile() {
	unloadAllVolatile()
	all_volatile = make(map[string][]Volatile)
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
