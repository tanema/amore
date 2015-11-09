package audio

type SourceType int

const (
	STATIC_SOURCE SourceType = iota
	STREAM_SOURCE
)

type Source struct {
	relative          bool
	looping           bool
	pitch             float32
	volume            float32
	coneInnerAngle    float32
	coneOuterAngle    float32
	coneOuterVolume   float32
	minVolume         float32
	maxVolume         float32
	referenceDistance float32
	rolloffFactor     float32
	maxDistance       float32
}

// Creates an identical copy of the Source in the stopped state.
func (s *Source) Clone() {}

// Gets the reference and maximum attenuation distances of the Source.
func (s *Source) GetAttenuationDistances() {}

// Gets the number of channels in the Source.
func (s *Source) GetChannels() {}

// Gets the Source's directional volume cones.
func (s *Source) GetCone() {}

//Gets the direction of the Source.
func (s *Source) GetDirection() {}

//Gets the current pitch of the Source.
func (s *Source) GetPitch() {}

// Gets the position of the Source.
func (s *Source) GetPosition() {}

//Returns the rolloff factor of the source.
func (s *Source) GetRolloff() {}

// Gets the velocity of the Source.
func (s *Source) GetVelocity() {}

// Gets the current volume of the Source.
func (s *Source) GetVolume() {}

// Returns the volume limits of the source.
func (s *Source) GetVolumeLimits() {}

// Returns whether the Source will loop.
func (s *Source) IsLooping() {}

//Returns whether the Source is paused.
func (s *Source) IsPaused() {}

// Returns whether the Source is playing.
func (s *Source) IsPlaying() {}

//Gets whether the Source's position and direction are relative to the listener.
func (s *Source) IsRelative() {}

//Returns whether the Source is static.
func (s *Source) IsStatic() {}

// Returns whether the Source is stopped.
func (s *Source) IsStopped() {}

// Pauses a source.
func (s *Source) Pause() {}

//Plays a source.
func (s *Source) Play() {}

//Resumes a paused source.
func (s *Source) Resume() {}

//Rewinds a source.
func (s *Source) Rewind() {}

//Sets the currently playing position of the Source.
func (s *Source) Seek() {}

// Sets the reference and maximum attenuation distances of the Source.
func (s *Source) SetAttenuationDistances() {}

// Sets the Source's directional volume cones.
func (s *Source) SetCone() {}

//Sets the direction of the Source.
func (s *Source) SetDirection() {}

//Sets whether the Source should loop.
func (s *Source) SetLooping() {}

//Sets the pitch of the Source.
func (s *Source) SetPitch() {}

// Sets the position of the Source.
func (s *Source) SetPosition() {}

// Sets whether the Source's position and direction are relative to the listener.
func (s *Source) SetRelative() {}

//Sets the rolloff factor.
func (s *Source) SetRolloff() {}

// Sets the velocity of the Source.
func (s *Source) SetVelocity() {}

// Sets the current volume of the Source.
func (s *Source) SetVolume() {}

// Sets the volume limits of the source.
func (s *Source) SetVolumeLimits() {}

//Stops a source.
func (s *Source) Stop() {}

//Gets the currently playing position of the Source.
func (s *Source) Tell() {}
