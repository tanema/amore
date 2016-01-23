package gfx

type Volatile interface {
	LoadVolatile() bool
	UnloadVolatile()
}

var (
	all_volatile = []Volatile{}
)

func Register(new_volatile Volatile) {
	all_volatile = append(all_volatile, new_volatile)
	if gl_state.initialized {
		new_volatile.LoadVolatile()
	}
}

func Release(vol Volatile) {
	var pos int
	for i, v := range all_volatile {
		if v == vol {
			pos = i
			v.UnloadVolatile()
		}
	}

	all_volatile = all_volatile[:pos+copy(all_volatile[pos:], all_volatile[pos+1:])]
}

func ReleaseAllVolatile() {
	UnloadAllVolatile()
	all_volatile = []Volatile{}
}

func LoadAllVolatile() bool {
	success := true

	for _, v := range all_volatile {
		success = success && v.LoadVolatile()
	}

	return success
}

func UnloadAllVolatile() {
	for _, v := range all_volatile {
		v.UnloadVolatile()
	}
}
