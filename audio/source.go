package audio

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/tanema/amore/audio/al"
)

type SourceType int

const (
	STATIC_SOURCE SourceType = iota
	STREAM_SOURCE
)

const (
	MAX_ATTENUATION_DISTANCE = 1000000.0
	MAX_BUFFERS              = 8
)

type Source struct {
	decoder           Decoder
	source            al.Source
	src_type          SourceType
	pitch             float32
	volume            float32
	position          al.Vector
	velocity          al.Vector
	direction         al.Vector
	relative          bool
	looping           bool
	paused            bool
	minVolume         float32
	maxVolume         float32
	referenceDistance float32
	rolloffFactor     float32
	maxDistance       float32
	cone              al.Cone
	staticBuffer      al.Buffer
	streamBuffers     []al.Buffer
}

//	Creates a new Source from a file, SoundData, or Decoder
func NewStaticSource(filepath string) (*Source, error) { return NewSource(filepath, STATIC_SOURCE) }
func NewStreamSource(filepath string) (*Source, error) { return NewSource(filepath, STREAM_SOURCE) }
func NewSource(filepath string, source_type SourceType) (*Source, error) {
	if pool == nil {
		createPool()
	}

	decoder, err := decode(filepath)
	if err != nil {
		return nil, err
	}

	new_source := &Source{
		decoder:           decoder,
		src_type:          source_type,
		pitch:             1,
		volume:            1,
		maxVolume:         1,
		referenceDistance: 1,
		rolloffFactor:     1,
		maxDistance:       MAX_ATTENUATION_DISTANCE,
		cone:              al.Cone{0, 0, 0},
		position:          al.Vector{},
		velocity:          al.Vector{},
		direction:         al.Vector{},
	}

	if source_type == STATIC_SOURCE {
		new_source.staticBuffer = al.GenBuffers(1)[0]
		fmt.Println("")
		new_source.staticBuffer.BufferData(decoder.getFormat(), decoder.getData(), decoder.getSampleRate())
	} else if source_type == STREAM_SOURCE {
		new_source.streamBuffers = al.GenBuffers(MAX_BUFFERS)
	}

	return new_source, nil
}

func (s *Source) isValid() bool {
	return s.source != 0
}

func (s *Source) IsFinished() bool {
	if s.src_type == STATIC_SOURCE {
		return s.IsStopped()
	}
	return s.IsStopped() && !s.IsLooping() && s.decoder.isFinished()
}

func (s *Source) update() bool {
	if !s.isValid() {
		return false
	}

	if s.src_type == STATIC_SOURCE {
		return !s.IsStopped()
	} else if s.IsLooping() || !s.IsFinished() {
		pool.mutex.Lock()
		defer pool.mutex.Unlock()
		for i := s.source.BuffersProcessed(); i > 0; i-- {
			var buffer al.Buffer
			s.source.UnqueueBuffers(buffer)
			s.stream(buffer)
			s.source.QueueBuffers(buffer)
		}
		return true
	}

	return false
}

func (s *Source) reset() {
	if s.source == 0 {
		return
	}
	s.source.SetGain(s.volume)
	s.source.SetPosition(s.position)
	s.source.SetVelocity(s.velocity)
	s.source.SetDirection(s.direction)
	s.source.SetPitch(s.pitch)
	s.source.SetMinGain(s.minVolume)
	s.source.SetMaxGain(s.maxVolume)
	s.source.SetReferenceDistance(s.referenceDistance)
	s.source.SetMaxDistance(s.maxDistance)
	s.source.SetRolloff(s.rolloffFactor)
	s.source.SetRelative(s.relative)
	s.source.SetCone(s.cone)
	if s.src_type == STATIC_SOURCE {
		s.source.SetLooping(s.looping)
	}
}

// Gets the reference and maximum attenuation distances of the Source.
func (s *Source) GetAttenuationDistances() (float32, float32) {
	if s.isValid() {
		return s.source.ReferenceDistance(), s.source.MaxDistance()
	}
	return s.referenceDistance, s.maxDistance
}

// Gets the number of channels in the Source.
func (s *Source) GetChannels() int16 {
	return s.decoder.getChannels()
}

// Gets the Source's directional volume cones.
func (s *Source) GetCone() (float32, float32, float32) {
	if s.isValid() {
		c := s.source.Cone()
		return mgl32.DegToRad(float32(c.InnerAngle)), mgl32.DegToRad(float32(c.OuterAngle)), c.OuterVolume
	}
	return mgl32.DegToRad(float32(s.cone.InnerAngle)), mgl32.DegToRad(float32(s.cone.OuterAngle)), s.cone.OuterVolume
}

//Gets the direction of the Source.
func (s *Source) GetDirection() (float32, float32, float32) {
	if s.isValid() {
		d := s.source.Direction()
		return d[0], d[1], d[2]
	}
	return s.direction[0], s.direction[1], s.direction[2]
}

//Gets the current pitch of the Source.
func (s *Source) GetPitch() float32 {
	if s.isValid() {
		return s.source.Pitch()
	}
	return s.pitch
}

// Gets the position of the Source.
func (s *Source) GetPosition() (float32, float32, float32) {
	if s.isValid() {
		vec := s.source.Position()
		return vec[0], vec[1], vec[2]
	}
	return s.position[0], s.position[1], s.position[2]
}

//Returns the rolloff factor of the source.
func (s *Source) GetRolloff() float32 {
	if s.isValid() {
		return s.source.Rolloff()
	}
	return s.rolloffFactor
}

// Gets the velocity of the Source.
func (s *Source) GetVelocity() (float32, float32, float32) {
	if s.isValid() {
		vec := s.source.Velocity()
		return vec[0], vec[1], vec[2]
	}
	return s.velocity[0], s.velocity[1], s.velocity[2]
}

// Gets the current volume of the Source.
func (s *Source) GetVolume() float32 {
	if s.isValid() {
		return s.source.Gain()
	}
	return s.volume
}

// Returns the volume limits of the source.
func (s *Source) GetVolumeLimits() (float32, float32) {
	if s.isValid() {
		return s.source.MinGain(), s.source.MaxGain()
	}
	return s.minVolume, s.maxVolume
}

func (s *Source) GetState() State {
	return State(s.source.State())
}

// Returns whether the Source will loop.
func (s *Source) IsLooping() bool {
	return s.looping
}

//Returns whether the Source is paused.
func (s *Source) IsPaused() bool {
	if s.isValid() {
		return s.GetState() == Paused
	}
	return false
}

// Returns whether the Source is playing.
func (s *Source) IsPlaying() bool {
	if s.isValid() {
		return s.GetState() == Playing
	}
	return false
}

//Gets whether the Source's position and direction are relative to the listener.
func (s *Source) IsRelative() bool {
	if s.isValid() {
		return s.source.Relative()
	}
	return s.relative
}

//Returns whether the Source is static or stream.
func (s *Source) IsStatic() bool {
	return s.src_type == STATIC_SOURCE
}

// Returns whether the Source is stopped.
func (s *Source) IsStopped() bool {
	if s.isValid() {
		return s.GetState() == Stopped
	}
	return true
}

// Sets the reference and maximum attenuation distances of the Source.
func (s *Source) SetAttenuationDistances(ref, max float32) {
	s.referenceDistance = ref
	s.maxDistance = max
	s.reset()
}

// Sets the Source's directional volume cones.
func (s *Source) SetCone(innerAngle, outerAngle, outerVolume float32) {
	s.cone = al.Cone{
		InnerAngle:  int32(mgl32.RadToDeg(innerAngle)),
		OuterAngle:  int32(mgl32.RadToDeg(outerAngle)),
		OuterVolume: outerVolume,
	}
	s.reset()
}

//Sets the direction of the Source.
func (s *Source) SetDirection(x, y, z float32) {
	if s.GetChannels() > 1 {
		panic("This spatial audio functionality is only available for mono Sources. Ensure the Source is not multi-channel before calling this function.")
	}

	s.direction = al.Vector{x, y, z}
	s.reset()
}

//Sets whether the Source should loop.
func (s *Source) SetLooping(do_loop bool) {
	s.looping = do_loop
	s.reset()
}

//Sets the pitch of the Source.
func (s *Source) SetPitch(p float32) {
	s.pitch = p
	s.reset()
}

// Sets the position of the Source.
func (s *Source) SetPosition(x, y, z float32) {
	s.position = al.Vector{x, y, z}
	s.reset()
}

// Sets whether the Source's position and direction are relative to the listener.
func (s *Source) SetRelative(is_relative bool) {
	s.relative = is_relative
	s.reset()
}

//Sets the rolloff factor.
func (s *Source) SetRolloff(roll_off float32) {
	s.rolloffFactor = roll_off
	s.reset()
}

// Sets the velocity of the Source.
func (s *Source) SetVelocity(x, y, z float32) {
	s.velocity = al.Vector{x, y, z}
	s.reset()
}

// Sets the current volume of the Source.
func (s *Source) SetVolume(v float32) {
	s.volume = v
	s.reset()
}

// Sets the volume limits of the source.
func (s *Source) SetVolumeLimits(min, max float32) {
	s.minVolume = min
	s.maxVolume = max
	s.reset()
}

// Pauses a source.
func (s *Source) Pause() {
	if s.isValid() {
		pool.mutex.Lock()
		defer pool.mutex.Unlock()
		al.PauseSources(s.source)
		s.paused = true
	}
}

//Plays a source.
func (s *Source) Play() bool {
	if s.IsPlaying() {
		return true
	}

	if s.IsPaused() {
		pool.Resume(s)
		return true
	}

	//claim a source for ourselves and make sure it worked
	if !pool.claim(s) || !s.isValid() {
		return false
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if s.src_type == STATIC_SOURCE {
		s.source.SetBuffer(s.staticBuffer)
	} else if s.src_type == STREAM_SOURCE {
		for i := 0; i <= MAX_BUFFERS; i++ {
			buffer := al.GenBuffers(1)[0]
			s.stream(buffer)
			s.streamBuffers = append(s.streamBuffers, buffer)
			if s.decoder.isFinished() {
				break
			}
		}
		if len(s.streamBuffers) > 0 {
			s.source.QueueBuffers(s.streamBuffers...)
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

func (s *Source) stream(buffer al.Buffer) int {
	decoded := s.decoder.decode() //get more data
	buffer.BufferData(s.decoder.getFormat(), s.decoder.getBuffer(), s.decoder.getSampleRate())
	if s.decoder.isFinished() && s.IsLooping() {
		s.decoder.rewind()
	}
	return decoded
}

//Resumes a paused source.
func (s *Source) Resume() {
	if s.isValid() && s.paused {
		pool.mutex.Lock()
		defer pool.mutex.Unlock()
		al.PlaySources(s.source)
		s.paused = false
	}
}

//Rewinds a source.
func (s *Source) Rewind() { s.Seek(0) }

//Sets the currently playing position of the Source.
func (s *Source) Seek(offset time.Duration) {
	if s.isValid() {
		size := s.decoder.durToByteOffset(offset)
		if s.src_type == STREAM_SOURCE {
			s.decoder.seek(int64(size))
			waspaused := s.paused
			// Because we still have old data from before the seek in the buffers let's empty them.
			s.Stop()
			s.Play()
			if waspaused {
				s.Pause()
			}
		} else {
			pool.mutex.Lock()
			defer pool.mutex.Unlock()
			s.source.SetOffsetBytes(size)
		}
	}
}

//Stops a source.
func (s *Source) Stop() {
	if !s.IsStopped() && s.isValid() {
		pool.mutex.Lock()
		if s.src_type == STATIC_SOURCE {
			al.StopSources(s.source)
		} else if s.src_type == STREAM_SOURCE {
			al.StopSources(s.source)
			s.source.UnqueueBuffers(s.streamBuffers...)
		}
		s.source.ClearBuffers()
		pool.mutex.Unlock()
	}
	s.Rewind()
	if s.src_type == STREAM_SOURCE {
		al.DeleteBuffers(s.streamBuffers...)
	}
	pool.release(s)
}

//Gets the currently playing position of the Source.
func (s *Source) Tell() time.Duration {
	if s.isValid() {
		pool.mutex.Lock()
		defer pool.mutex.Unlock()
		return s.decoder.byteOffsetToDur(s.source.OffsetByte())
	}
	return time.Duration(0.0)
}
