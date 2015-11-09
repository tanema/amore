package audio

import (
	"sync"

	//"golang.org/x/mobile/exp/audio/al"
)

const (
	MAX_SOURCES = 64
)

type audioPool struct {
	mu           sync.Mutex
	totalSources int
	sources      [MAX_SOURCES]uint32
	available    []uint32
	playing      map[uint32]*Source
}

func createPool() {}

func (p *audioPool) IsAvailable()                        {}
func (p *audioPool) IsPlaying(s *Source)                 {}
func (p *audioPool) Update()                             {}
func (p *audioPool) GetSourceCount()                     {}
func (p *audioPool) GetMaxSources()                      {}
func (p *audioPool) Play(source *Source)                 {}
func (p *audioPool) Stop(source *Source)                 {}
func (p *audioPool) Pause(source *Source)                {}
func (p *audioPool) Resume(source *Source)               {}
func (p *audioPool) Rewind(source *Source)               {}
func (p *audioPool) SoftRewind(source *Source)           {}
func (p *audioPool) Seek(source *Source, offset float32) {}
func (p *audioPool) Tell(source *Source)                 {}
