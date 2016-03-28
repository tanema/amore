// +build js

package al

import (
	"github.com/tanema/amore/audio/al/webal"
)

type (
	Source      struct{ webal.Source }
	Buffer      struct{ webal.Buffer }
	Cone        webal.Cone
	Orientation struct{ webal.Orientation }
)
