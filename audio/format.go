package audio

import (
	"github.com/tanema/amore/audio/al"
)

// Format represents a PCM data format.
type Format int

const (
	Mono8 Format = iota + 1
	Mono16
	Stereo8
	Stereo16
)

func (f Format) String() string { return formatStrings[f] }

// formatBytes is the product of bytes per sample and number of channels.
var formatBytes = [...]int64{
	Mono8:    1,
	Mono16:   2,
	Stereo8:  2,
	Stereo16: 4,
}

var formatCodes = [...]uint32{
	Mono8:    al.FormatMono8,
	Mono16:   al.FormatMono16,
	Stereo8:  al.FormatStereo8,
	Stereo16: al.FormatStereo16,
}

var formatStrings = [...]string{
	0:        "unknown",
	Mono8:    "mono8",
	Mono16:   "mono16",
	Stereo8:  "stereo8",
	Stereo16: "stereo16",
}
