package audio

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/tanema/amore/audio/al"
)

type SourceType int

const (
	STATIC_SOURCE SourceType = iota
	STREAM_SOURCE
)

type Source struct {
	decoder           Decoder
	Source            al.Source
	src_type          SourceType
	pitch             float32
	volume            float32
	position          al.Vector
	velocity          al.Vector
	direction         al.Vector
	relative          bool
	looping           bool
	paused            bool
	valid             bool
	minVolume         float32
	maxVolume         float32
	referenceDistance float32
	rolloffFactor     float32
	maxDistance       float32
	cone              al.Cone
	toLoop            uint
}

//	Creates a new Source from a file, SoundData, or Decoder
func NewStaticSource(filepath string) (*Source, error) { return NewSource(filepath, STATIC_SOURCE) }
func NewStreamSource(filepath string) (*Source, error) { return NewSource(filepath, STREAM_SOURCE) }
func NewSource(filepath string, source_type SourceType) (*Source, error) {
	decoder, err := decode(filepath)
	if err != nil {
		return nil, err
	}

	new_source := &Source{
		decoder:  decoder,
		src_type: source_type,
	}

	return new_source, nil
}

func (s *Source) Update() bool {
	return false
}

func (s *Source) reset() {
	s.Source.SetPosition(s.position)
	s.Source.SetVelocity(s.velocity)
	s.Source.SetDirection(s.direction)
	s.SetPitch(s.pitch)
	s.SetVolume(s.volume)
	s.SetVolumeLimits(s.minVolume, s.maxVolume)
	s.SetAttenuationDistances(s.referenceDistance, s.maxDistance)
	s.SetRolloff(s.rolloffFactor)
	s.SetLooping(s.looping)
	s.SetRelative(s.relative)
	s.Source.SetCone(s.cone)
}

// Gets the reference and maximum attenuation distances of the Source.
func (s *Source) GetAttenuationDistances() (float32, float32) {
	return s.Source.ReferenceDistance(), s.Source.MaxDistance()
}

// Gets the number of channels in the Source.
func (s *Source) GetChannels() int16 {
	return s.decoder.GetChannels()
}

// Gets the Source's directional volume cones.
func (s *Source) GetCone() (float32, float32, float32) {
	c := s.Source.Cone()
	inner := mgl32.DegToRad(float32(c.InnerAngle))
	outter := mgl32.DegToRad(float32(c.OuterAngle))
	return inner, outter, c.OuterVolume
}

//Gets the direction of the Source.
func (s *Source) GetDirection() (float32, float32, float32) {
	d := s.Source.Direction()
	return d[0], d[1], d[2]
}

//Gets the current pitch of the Source.
func (s *Source) GetPitch() float32 {
	return s.Source.Pitch()
}

// Gets the position of the Source.
func (s *Source) GetPosition() (float32, float32, float32) {
	vec := s.Source.Position()
	return vec[0], vec[1], vec[2]
}

//Returns the rolloff factor of the source.
func (s *Source) GetRolloff() float32 {
	return s.Source.Rolloff()
}

// Gets the velocity of the Source.
func (s *Source) GetVelocity() (float32, float32, float32) {
	vec := s.Source.Velocity()
	return vec[0], vec[1], vec[2]
}

// Gets the current volume of the Source.
func (s *Source) GetVolume() float32 {
	return s.Source.Gain()
}

// Returns the volume limits of the source.
func (s *Source) GetVolumeLimits() (float32, float32) {
	return s.Source.MinGain(), s.Source.MaxGain()
}

func (s *Source) GetState() State {
	return State(s.Source.State())
}

// Returns whether the Source will loop.
func (s *Source) IsLooping() bool {
	return s.looping
}

//Returns whether the Source is paused.
func (s *Source) IsPaused() bool {
	if s.valid {
		return s.GetState() == Paused
	}
	return s.paused
}

// Returns whether the Source is playing.
func (s *Source) IsPlaying() bool {
	if s.valid {
		return s.GetState() == Playing
	}
	return false
}

//Gets whether the Source's position and direction are relative to the listener.
func (s *Source) IsRelative() bool {
	return s.Source.Relative()
}

//Returns whether the Source is static or stream.
func (s *Source) IsStatic() bool {
	return s.src_type == STATIC_SOURCE
}

// Returns whether the Source is stopped.
func (s *Source) IsStopped() bool {
	if s.valid {
		return s.GetState() == Stopped
	}
	return true
}

// Pauses a source.
func (s *Source) Pause() {
	pool.Pause(s)
}

//Plays a source.
func (s *Source) Play() bool {
	if s.valid && s.IsPaused() {
		pool.Resume(s)
		return true
	}

	s.valid = pool.Play(s)
	return s.valid
}

//Resumes a paused source.
func (s *Source) Resume() {
	pool.Resume(s)
}

//Rewinds a source.
func (s *Source) Rewind() {
	pool.Rewind(s)
}

//Sets the currently playing position of the Source.
func (s *Source) Seek() {}

// Sets the reference and maximum attenuation distances of the Source.
func (s *Source) SetAttenuationDistances(ref, max float32) {
	s.referenceDistance = ref
	s.maxDistance = max
	if s.valid {
		s.Source.SetReferenceDistance(s.referenceDistance)
		s.Source.SetMaxDistance(s.maxDistance)
	}
}

// Sets the Source's directional volume cones.
func (s *Source) SetCone(innerAngle, outerAngle, outerVolume float32) {
	s.cone = al.Cone{
		InnerAngle:  int32(mgl32.RadToDeg(innerAngle)),
		OuterAngle:  int32(mgl32.RadToDeg(outerAngle)),
		OuterVolume: outerVolume,
	}

	if s.valid {
		s.Source.SetCone(s.cone)
	}
}

//Sets the direction of the Source.
func (s *Source) SetDirection(x, y, z float32) {
	if s.GetChannels() > 1 {
		panic("This spatial audio functionality is only available for mono Sources. Ensure the Source is not multi-channel before calling this function.")
	}

	s.direction = al.Vector{x, y, z}
	if s.valid {
		s.Source.SetDirection(s.direction)
	}
}

//Sets whether the Source should loop.
func (s *Source) SetLooping(do_loop bool) {
	s.looping = do_loop
	if s.valid && s.src_type == STATIC_SOURCE {
		s.Source.SetLooping(s.looping)
	}
}

//Sets the pitch of the Source.
func (s *Source) SetPitch(p float32) {
	s.pitch = p
	if s.valid {
		s.Source.SetPitch(s.pitch)
	}
}

// Sets the position of the Source.
func (s *Source) SetPosition(x, y, z float32) {
	s.position = al.Vector{x, y, z}
	if s.valid {
		s.Source.SetPosition(s.position)
	}
}

// Sets whether the Source's position and direction are relative to the listener.
func (s *Source) SetRelative(is_relative bool) {
	s.relative = is_relative
	if s.valid {
		s.Source.SetRelative(s.relative)
	}
}

//Sets the rolloff factor.
func (s *Source) SetRolloff(roll_off float32) {
	s.rolloffFactor = roll_off
	if s.valid {
		s.Source.SetRolloff(s.rolloffFactor)
	}
}

// Sets the velocity of the Source.
func (s *Source) SetVelocity(x, y, z float32) {
	s.velocity = al.Vector{x, y, z}
	if s.valid {
		s.Source.SetVelocity(s.velocity)
	}
}

// Sets the current volume of the Source.
func (s *Source) SetVolume(v float32) {
	s.volume = v
	if s.valid {
		s.Source.SetGain(s.volume)
	}
}

// Sets the volume limits of the source.
func (s *Source) SetVolumeLimits(min, max float32) {
	s.minVolume = min
	s.maxVolume = max
	if s.valid {
		s.Source.SetMinGain(s.minVolume)
		s.Source.SetMaxGain(s.maxVolume)
	}
}

//Stops a source.
func (s *Source) Stop() {
	if !s.IsStopped() {
		pool.Stop(s)
		pool.Rewind(s)
	}
}

//Gets the currently playing position of the Source.
func (s *Source) Tell() float32 {
	return pool.Tell(s)
}

func (s *Source) Release() {
	pool.Release(s)
}

// Pauses a source.
func (s *Source) PauseAtomic() {}

//Plays a source.
func (s *Source) PlayAtomic() bool {
	return true
}

//Resumes a paused source.
func (s *Source) ResumeAtomic() {}

//Rewinds a source.
func (s *Source) RewindAtomic() {}

//Sets the currently playing position of the Source.
func (s *Source) SeekAtomic(offset float32) {}

//Stops a source.
func (s *Source) StopAtomic() {}

//Gets the currently playing position of the Source.
func (s *Source) TellAtomic() float32 {
	return 0.0
}
