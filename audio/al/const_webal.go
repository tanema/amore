// +build js

package al

const (
	NoError = iota
	InvalidName
	InvalidEnum
	InvalidValue
	InvalidOperation
	OutOfMemory

	Initial = iota
	Playing
	Paused
	Stopped

	FormatMono8 = iota
	FormatMono16
	FormatStereo8
	FormatStereo16
)
