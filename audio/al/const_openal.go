// +build !js

package al

import (
	"github.com/tanema/amore/audio/al/openal"
)

const (
	NoError          = openal.NoError
	InvalidName      = openal.InvalidName
	InvalidEnum      = openal.InvalidEnum
	InvalidValue     = openal.InvalidValue
	InvalidOperation = openal.InvalidOperation
	OutOfMemory      = openal.OutOfMemory

	InverseDistance         = openal.InverseDistance
	InverseDistanceClamped  = openal.InverseDistanceClamped
	LinearDistance          = openal.LinearDistance
	LinearDistanceClamped   = openal.LinearDistanceClamped
	ExponentDistance        = openal.ExponentDistance
	ExponentDistanceClamped = openal.ExponentDistanceClamped

	Initial = openal.Initial
	Playing = openal.Playing
	Paused  = openal.Paused
	Stopped = openal.Stopped

	FormatMono8    = openal.FormatMono8
	FormatMono16   = openal.FormatMono16
	FormatStereo8  = openal.FormatStereo8
	FormatStereo16 = openal.FormatStereo16
)
