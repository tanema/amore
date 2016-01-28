package gfx

type Volatile interface {
	loadVolatile() bool
	unloadVolatile()
}

var (
	all_volatile = []Volatile{}
)

func registerVolatile(new_volatile Volatile) {
	all_volatile = append(all_volatile, new_volatile)
	if gl_state.initialized {
		new_volatile.loadVolatile()
	}
}

func releaseVolatile(vol Volatile) {
	var pos int
	for i, v := range all_volatile {
		if v == vol {
			pos = i
			v.unloadVolatile()
		}
	}

	all_volatile = all_volatile[:pos+copy(all_volatile[pos:], all_volatile[pos+1:])]
}

func releaseAllVolatile() {
	unloadAllVolatile()
	all_volatile = []Volatile{}
}

func loadAllVolatile() {
	for _, v := range all_volatile {
		v.loadVolatile()
	}
}

func unloadAllVolatile() {
	for _, v := range all_volatile {
		v.unloadVolatile()
	}
}
