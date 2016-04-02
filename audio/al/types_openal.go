// +build !js

package al

import (
	"github.com/tanema/amore/audio/al/openal"
)

type (
	Source      struct{ openal.Source }
	Buffer      struct{ openal.Buffer }
	Cone        openal.Cone
	Orientation struct{ openal.Orientation }
)
