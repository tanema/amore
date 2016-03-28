package audio

import (
	"sync"
	"time"

	"github.com/tanema/amore/audio/al"
)

const (
	MAX_SOURCES = 64
)

var (
	pool                *audioPool
	supportedExtentions []string
)

type audioPool struct {
	mutex        sync.Mutex
	totalSources int
	sources      [MAX_SOURCES]al.Source
	available    []al.Source
	playing      map[al.Source]*Source
}

func createPool() {
	pool = &audioPool{
		sources:   [MAX_SOURCES]al.Source{},
		available: []al.Source{},
		playing:   make(map[al.Source]*Source),
	}

	// Generate sources.
	for i := 0; i < MAX_SOURCES; i++ {
		pool.sources[i] = al.CreateSource()

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

func (p *audioPool) update() {
	for _, source := range p.playing {
		if !source.update() {
			source.Stop()
		}
	}
}

func (p *audioPool) GetSourceCount() int {
	return len(p.playing)
}

func (p *audioPool) GetMaxSources() int {
	return p.totalSources
}

func (p *audioPool) Play(source *Source) bool {
	if source == nil {
		success := true
		for _, source := range p.playing {
			success = success && source.Play()
		}
		return success
	}
	return source.Play()
}

func (p *audioPool) claim(source *Source) bool {
	if len(p.available) > 0 {
		source.source, p.available = p.available[len(p.available)-1], p.available[:len(p.available)-1]
		p.playing[source.source] = source
		return true
	}
	return false
}

func (p *audioPool) release(source *Source) {
	p.available = append(p.available, source.source)
	delete(p.playing, source.source)
	source.source = al.Source{}
}

func (p *audioPool) Stop(source *Source) {
	if source == nil {
		for _, source := range p.playing {
			source.Stop()
		}
		p.playing = make(map[al.Source]*Source)
	} else {
		source.Stop()
	}
}

func (p *audioPool) Pause(source *Source) {
	if source == nil {
		for _, source := range p.playing {
			source.Pause()
		}
	} else {
		source.Pause()
	}
}

func (p *audioPool) Resume(source *Source) {
	if source == nil {
		for _, source := range p.playing {
			source.Resume()
		}
	} else {
		source.Resume()
	}
}

func (p *audioPool) Rewind(source *Source) {
	if source == nil {
		for _, source := range p.playing {
			source.Rewind()
		}
	} else {
		source.Rewind()
	}
}
