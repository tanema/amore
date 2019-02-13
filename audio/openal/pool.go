package openal

import (
	"sync"
	"time"

	"github.com/tanema/amore/audio/openal/al"
)

const maxSources = 64

var pool *audioPool

// audioPool manages all openAL generated sources.
type audioPool struct {
	mutex        sync.Mutex
	totalSources int
	sources      [maxSources]al.Source
	available    []al.Source
	playing      map[al.Source]*Source
}

// init will open the audio interface.
func init() {
	if err := al.OpenDevice(); err != nil {
		panic(err)
	}
}

// GetVolume returns the master volume.
func GetVolume() float32 { return al.ListenerGain() }

// SetVolume sets the master volume
func SetVolume(gain float32) {
	al.SetListenerGain(gain)
}

// PauseAll will pause all sources
func PauseAll() {
	for _, source := range pool.playing {
		source.Pause()
	}
}

// PlayAll will play all sources
func PlayAll() {
	for _, source := range pool.playing {
		source.Play()
	}
}

// RewindAll will rewind all sources
func RewindAll() {
	for _, source := range pool.playing {
		source.Rewind()
	}
}

// StopAll stop all sources
func StopAll() {
	for _, source := range pool.playing {
		source.Stop()
	}
	pool.playing = make(map[al.Source]*Source)
}

// createPool generates a new pool and gets the max sources.
func createPool() {
	pool = &audioPool{
		sources:   [maxSources]al.Source{},
		available: []al.Source{},
		playing:   make(map[al.Source]*Source),
	}

	// Generate sources.
	for i := 0; i < maxSources; i++ {
		pool.sources[i] = al.GenSources(1)[0]

		// We might hit an implementation-dependent limit on the total number
		// of sources before reaching maxSources.
		if al.Error() != 0 {
			break
		}

		pool.totalSources++
	}

	if pool.totalSources < 4 {
		panic("Could not generate audio sources.")
	}

	// Make all sources available initially.
	for i := 0; i < pool.totalSources; i++ {
		pool.available = append(pool.available, pool.sources[i])
	}

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					pool.update()
				}
			}
		}()
	}()
}

// update will cycle through al the playing sources and updates the buffers
func (p *audioPool) update() {
	for _, source := range p.playing {
		if !source.update() {
			source.Stop()
		}
	}
}

// getSourceCount will get the count of playing sources
func (p *audioPool) getSourceCount() int {
	return len(p.playing)
}

// claim will take an openAL source for playing with an amore source
func (p *audioPool) claim(source *Source) bool {
	if len(p.available) > 0 {
		source.source, p.available = p.available[len(p.available)-1], p.available[:len(p.available)-1]
		p.playing[source.source] = source
		return true
	}
	return false
}

// release will put an openAL source back in to the available queue
func (p *audioPool) release(source *Source) {
	p.available = append(p.available, source.source)
	delete(p.playing, source.source)
	source.source = 0
}
