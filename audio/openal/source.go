package openal

import (
	"time"

	"github.com/tanema/amore/audio/openal/al"
	"github.com/tanema/amore/audio/openal/decoding"
)

const (
	maxAttenuationDistance = 1000000.0 // upper limit of sound attentuation time.
	maxBuffers             = 8         //arbitrary limit of umber of buffers a source can use to stream
)

// Source manages decoding sound data, creates an openal sound and manages the
// data associated with the source.
type Source struct {
	decoder      *decoding.Decoder
	source       al.Source
	isStatic     bool
	pitch        float32
	volume       float32
	looping      bool
	paused       bool
	staticBuffer al.Buffer
	offsetBytes  int32
}

// NewSource creates a new Source from a file at the path provided. If you
// specify a static source it will all be buffered into a single buffer. If
// false then it will create many buffers a cycle through them with data chunks.
// This allows a smaller memory footprint while playing bigger music files. You
// may want a static file if the sound is less than 2 seconds. It allows for faster
// cleaning playing of shorter sounds like footsteps.
func NewSource(filepath string, static bool) (*Source, error) {
	if pool == nil {
		createPool()
	}

	decoder, err := decoding.Decode(filepath)
	if err != nil {
		return nil, err
	}

	newSource := &Source{
		decoder:  decoder,
		isStatic: static,
		pitch:    1,
		volume:   1,
	}

	if static {
		newSource.staticBuffer = al.GenBuffers(1)[0]
		newSource.staticBuffer.BufferData(decoder.Format, decoder.GetData(), decoder.SampleRate)
	}

	return newSource, nil
}

// isValid will return true if the source is associated with an openal source anymore.
// if not it will return false and will disable most funtionality
func (s *Source) isValid() bool {
	return s.source != 0
}

// IsFinished will return true if the source is at the end of its duration and
// it is not a looping source.
func (s *Source) IsFinished() bool {
	if s.isStatic {
		return s.IsStopped()
	}
	return s.IsStopped() && !s.IsLooping() && s.decoder.IsFinished()
}

// update will return true if successfully updated the source. If the source is
// static it will return if the item is still playing. If the item is a streamed
// source it will return true if it is still playing but after updating it's buffers.
func (s *Source) update() bool {
	if !s.isValid() {
		return false
	}

	if s.isStatic {
		return !s.IsStopped()
	} else if s.IsLooping() || !s.IsFinished() {
		pool.mutex.Lock()
		defer pool.mutex.Unlock()
		for i := s.source.BuffersProcessed(); i > 0; i-- {
			curOffsetBytes := s.source.OffsetByte()
			buffer := s.source.UnqueueBuffer()
			newOffsetBytes := s.source.OffsetByte()
			s.offsetBytes += (curOffsetBytes - newOffsetBytes)
			if s.stream(buffer) > 0 {
				s.source.QueueBuffers(buffer)
			}
		}
		return true
	}

	return false
}

// reset sets all the source's values in openal to the preset values.
func (s *Source) reset() {
	if !s.isValid() {
		return
	}
	s.source.SetGain(s.volume)
	s.source.SetPitch(s.pitch)
	if s.isStatic {
		s.source.SetLooping(s.looping)
	}
}

// GetDuration returns the total duration of the source.
func (s *Source) GetDuration() time.Duration {
	return s.decoder.Duration()
}

// GetPitch returns the current pitch of the Source in the range 0.0, 1.0
func (s *Source) GetPitch() float32 {
	if s.isValid() {
		return s.source.Pitch()
	}
	return s.pitch
}

// GetVolume returns the current volume of the Source.
func (s *Source) GetVolume() float32 {
	if s.isValid() {
		return s.source.Gain()
	}
	return s.volume
}

// GetState returns the playing state of the source.
func (s *Source) GetState() string {
	switch s.source.State() {
	case al.Initial:
		return "initial"
	case al.Playing:
		return "playing"
	case al.Paused:
		return "paused"
	case al.Stopped:
		return "stopped"
	default:
		return "unknown"
	}
}

// IsLooping returns whether the Source will loop.
func (s *Source) IsLooping() bool {
	return s.looping
}

// IsPaused returns whether the Source is paused.
func (s *Source) IsPaused() bool {
	if s.isValid() {
		return s.GetState() == "paused"
	}
	return false
}

// IsPlaying returns whether the Source is playing.
func (s *Source) IsPlaying() bool {
	if s.isValid() {
		return s.GetState() == "playing"
	}
	return false
}

// IsStatic returns whether the Source is static or stream.
func (s *Source) IsStatic() bool {
	return s.isStatic
}

// IsStopped returns whether the Source is stopped.
func (s *Source) IsStopped() bool {
	if s.isValid() {
		return s.GetState() == "stopped"
	}
	return true
}

// SetLooping sets whether the Source should loop when the source is complete.
func (s *Source) SetLooping(loop bool) {
	s.looping = loop
	s.reset()
}

// SetPitch sets the pitch of the Source, the value should be between 0.0, 1.0
func (s *Source) SetPitch(p float32) {
	s.pitch = p
	s.reset()
}

// SetVolume sets the current volume of the Source.
func (s *Source) SetVolume(v float32) {
	s.volume = v
	s.reset()
}

// Play starts playing the source.
func (s *Source) Play() bool {
	if s.IsPlaying() {
		return true
	}

	if s.IsPaused() {
		s.Resume()
		return true
	}

	//claim a source for ourselves and make sure it worked
	if !pool.claim(s) || !s.isValid() {
		return false
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if s.isStatic {
		s.source.SetBuffer(s.staticBuffer)
	} else {
		buffers := []al.Buffer{}
		for i := 0; i < maxBuffers; i++ {
			buffer := al.GenBuffers(1)[0]
			if s.stream(buffer) > 0 {
				buffers = append(buffers, buffer)
			}
			if s.decoder.IsFinished() {
				break
			}
		}
		if len(buffers) > 0 {
			s.source.QueueBuffers(buffers...)
		}
	}

	// This Source may now be associated with an OpenAL source that still has
	// the properties of another Source. Let's reset it to the settings
	// of the new one.
	s.reset()

	// Clear errors.
	al.Error()

	al.PlaySources(s.source)

	// alSourcePlay may fail if the system has reached its limit of simultaneous
	// playing sources.
	return al.Error() == al.NoError
}

// stream fills a buffer with the next chunk of data
func (s *Source) stream(buffer al.Buffer) int {
	decoded := s.decoder.Decode() //get more data
	if decoded > 0 {
		buffer.BufferData(s.decoder.Format, s.decoder.Buffer, s.decoder.SampleRate)
	}
	if s.decoder.IsFinished() && s.IsLooping() {
		s.Rewind()
	}
	return decoded
}

// Pause pauses the source.
func (s *Source) Pause() {
	if s.isValid() {
		pool.mutex.Lock()
		defer pool.mutex.Unlock()
		al.PauseSources(s.source)
		s.paused = true
	}
}

// Resume resumes a paused source.
func (s *Source) Resume() {
	if s.isValid() && s.paused {
		pool.mutex.Lock()
		defer pool.mutex.Unlock()
		al.PlaySources(s.source)
		s.paused = false
	}
}

// Rewind rewinds the source source to its start time.
func (s *Source) Rewind() { s.Seek(0) }

// Seek sets the currently playing position of the Source.
func (s *Source) Seek(offset time.Duration) {
	if !s.isValid() {
		return
	}
	s.offsetBytes = s.decoder.DurToByteOffset(offset)
	if !s.isStatic {
		al.StopSources(s.source)
		s.decoder.Seek(int64(s.offsetBytes))
		for i := s.source.BuffersQueued(); i > 0; i-- {
			buffer := s.source.UnqueueBuffer()
			if s.stream(buffer) > 0 {
				s.source.QueueBuffers(buffer)
			} else {
				al.DeleteBuffers(buffer)
			}
		}
		if !s.paused {
			al.PlaySources(s.source)
		}
	} else {
		pool.mutex.Lock()
		defer pool.mutex.Unlock()
		s.source.SetOffsetBytes(s.offsetBytes)
	}
}

// Stop stops a playing source.
func (s *Source) Stop() {
	if s.isValid() {
		pool.mutex.Lock()
		defer pool.mutex.Unlock()
		al.StopSources(s.source)
		s.offsetBytes = 0
		if !s.isStatic {
			queued := s.source.BuffersQueued()
			for i := queued; i > 0; i-- {
				buffer := s.source.UnqueueBuffer()
				al.DeleteBuffers(buffer)
			}
			s.decoder.Seek(0)
		} else {
			s.source.SetOffsetBytes(0)
		}
		s.source.ClearBuffers()
		pool.release(s)
	}
}

// Tell returns the currently playing position of the Source.
func (s *Source) Tell() time.Duration {
	if s.isValid() {
		pool.mutex.Lock()
		defer pool.mutex.Unlock()
		if s.isStatic {
			return s.decoder.ByteOffsetToDur(s.source.OffsetByte())
		}
		return s.decoder.ByteOffsetToDur(s.offsetBytes + s.source.OffsetByte())
	}
	return time.Duration(0.0)
}
