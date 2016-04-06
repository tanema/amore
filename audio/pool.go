package audio

import (
	"sync"
	"time"

	"github.com/tanema/amore/audio/al"
)

// An upper limit to try and reach while finding out how many sources the system
// can handle
const MAX_SOURCES = 64

var pool *audioPool

// audioPool manages all openAL generated sources.
type audioPool struct {
	mutex        sync.Mutex
	totalSources int
	sources      [MAX_SOURCES]al.Source
	available    []al.Source
	playing      map[al.Source]*Source
}

// createPool generates a new pool and gets the max sources.
func createPool() {
	pool = &audioPool{
		sources:   [MAX_SOURCES]al.Source{},
		available: []al.Source{},
		playing:   make(map[al.Source]*Source),
	}

	// Generate sources.
	for i := 0; i < MAX_SOURCES; i++ {
		pool.sources[i] = al.GenSources(1)[0]

		// We might hit an implementation-dependent limit on the total number
		// of sources before reaching MAX_SOURCES.
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
		ticker := time.NewTicker(5 * time.Millisecond)
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

// play will play a source in the pool
// if source is nil it will play all
func (p *audioPool) play(source *Source) bool {
	if source == nil {
		success := true
		for _, source := range p.playing {
			success = success && source.Play()
		}
		return success
	}
	return source.Play()
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

// stop will stop a playing source in the pool
// if source is nil it will stop all
func (p *audioPool) stop(source *Source) {
	if source == nil {
		for _, source := range p.playing {
			source.Stop()
		}
		p.playing = make(map[al.Source]*Source)
	} else {
		source.Stop()
	}
}

// pause will pause a playing source in the pool
// if source is nil it will pause all
func (p *audioPool) pause(source *Source) {
	if source == nil {
		for _, source := range p.playing {
			source.Pause()
		}
	} else {
		source.Pause()
	}
}

// resume will resume a stopped source in the pool
// if source is nil it will resume all
func (p *audioPool) resume(source *Source) {
	if source == nil {
		for _, source := range p.playing {
			source.Resume()
		}
	} else {
		source.Resume()
	}
}

// rewind will resume a stopped source in the pool
// if source is nil it will rewind all
func (p *audioPool) rewind(source *Source) {
	if source == nil {
		for _, source := range p.playing {
			source.Rewind()
		}
	} else {
		source.Rewind()
	}
}
