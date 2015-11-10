package audio

import (
	"sync"

	"golang.org/x/mobile/exp/audio/al"
)

const (
	MAX_SOURCES = 64
)

var (
	pool                audioPool
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
	pool = audioPool{
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
}

func (p *audioPool) IsAvailable() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	has := len(p.available) > 0
	return has
}

func (p *audioPool) IsPlaying(s *Source) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, is_playing := p.playing[s.Channel]
	return is_playing
}

func (p *audioPool) Update() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, source := range p.playing {
		if !source.Update() {
			p.removeSource(source)
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
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, alreadyPlaying := p.playing[source.Channel]; !alreadyPlaying {
		// Try to play.
		if len(p.available) > 0 {
			// Get the first available source and remove it
			var out al.Source
			out, p.available = p.available[len(p.available)-1], p.available[:len(p.available)-1]

			// Insert into map of playing sources.
			p.playing[out] = source
			source.Channel = out

			return source.PlayAtomic()
		} else {
			return false
		}
	}

	return true

}

func (p *audioPool) Stop(source *Source) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if source == nil {
		for _, source := range p.playing {
			p.removeSource(source)
		}

		p.playing = make(map[al.Source]*Source)
	} else {
		p.removeSource(source)
	}
}

func (p *audioPool) Pause(source *Source) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if source == nil {
		for _, source := range p.playing {
			source.PauseAtomic()
		}
	} else {
		source.PauseAtomic()
	}
}

func (p *audioPool) Resume(source *Source) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if source == nil {
		for _, source := range p.playing {
			source.ResumeAtomic()
		}
	} else {
		source.ResumeAtomic()
	}
}

func (p *audioPool) Rewind(source *Source) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if source == nil {
		for _, source := range p.playing {
			source.RewindAtomic()
		}
	} else {
		source.RewindAtomic()
	}
}

func (p *audioPool) Seek(source *Source, offset float32) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	source.SeekAtomic(offset)
}

func (p *audioPool) Tell(source *Source) float32 {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return source.TellAtomic()
}

func (p *audioPool) Release(source *Source) {
	p.available = append(p.available, source.Channel)
	delete(p.playing, source.Channel)
	source.Channel = 0
}

func (p *audioPool) removeSource(source *Source) {
	source.StopAtomic()
	source.RewindAtomic()
	source.Release()
	p.available = append(p.available, source.Channel)
	delete(p.playing, source.Channel)
	source.Channel = 0
}
