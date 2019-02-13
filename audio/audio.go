// Package audio is use for creating audio sources, managing/pooling resources,
// and playback of those audio sources.
package audio

import (
	"github.com/tanema/amore/audio/openal"
)

// Source is an playable audio source
type Source = openal.Source

// NewSource creates a new Source from a file at the path provided. If you
// specify a static source it will all be buffered into a single buffer. If
// false then it will create many buffers a cycle through them with data chunks.
// This allows a smaller memory footprint while playing bigger music files. You
// may want a static file if the sound is less than 2 seconds. It allows for faster
// cleaning playing of shorter sounds like footsteps.
func NewSource(filepath string, static bool) (*Source, error) {
	return openal.NewSource(filepath, static)
}

// GetVolume returns the master volume.
func GetVolume() float32 { return openal.GetVolume() }

// SetVolume sets the master volume
func SetVolume(gain float32) { openal.SetVolume(gain) }

// PauseAll will pause all sources
func PauseAll() { openal.PauseAll() }

// PlayAll will play all sources
func PlayAll() { openal.PlayAll() }

// RewindAll will rewind all sources
func RewindAll() { openal.RewindAll() }

// StopAll stop all sources
func StopAll() { openal.StopAll() }
