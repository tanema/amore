package audio

import (
	"os"

	"golang.org/x/mobile/exp/audio/al"
)

type SourceType int

const (
	STATIC_SOURCE SourceType = iota
	STREAM_SOURCE
)

type Cone struct {
	innerAngle  int
	outerAngle  int
	outerVolume float32
}

type Source struct {
	Channel           al.Source
	_type             SourceType
	src               os.File
	pitch             float32
	volume            float32
	position          [3]float32
	velocity          [3]float32
	direction         [3]float32
	relative          bool
	looping           bool
	paused            bool
	valid             bool
	minVolume         float32
	maxVolume         float32
	referenceDistance float32
	rolloffFactor     float32
	maxDistance       float32
	cone              Cone
	offsetSamples     float32
	offsetSeconds     float32
	sampleRate        int
	channels          int
	toLoop            uint
}

//	Creates a new Source from a file, SoundData, or Decoder
func NewSource(filepath string) (*Source, error) {
	_, err := decode(filepath)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Creates an identical copy of the Source in the stopped state.
func (s *Source) Clone() {}

func (s *Source) Update() bool {
	return false
}

// Gets the reference and maximum attenuation distances of the Source.
func (s *Source) GetAttenuationDistances() {}

// Gets the number of channels in the Source.
func (s *Source) GetChannels() {}

// Gets the Source's directional volume cones.
func (s *Source) GetCone() {}

//Gets the direction of the Source.
func (s *Source) GetDirection() {
}

//Gets the current pitch of the Source.
func (s *Source) GetPitch() {}

// Gets the position of the Source.
func (s *Source) GetPosition() (float32, float32, float32) {
	vec := s.Channel.Position()
	return vec[0], vec[1], vec[2]
}

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

// Pauses a source.
func (s *Source) PauseAtomic() {}

//Plays a source.
func (s *Source) Play() {}

//Plays a source.
func (s *Source) PlayAtomic() bool {
	return true
}

//Resumes a paused source.
func (s *Source) Resume() {}

//Resumes a paused source.
func (s *Source) ResumeAtomic() {}

//Rewinds a source.
func (s *Source) Rewind() {}

//Rewinds a source.
func (s *Source) RewindAtomic() {}

//Sets the currently playing position of the Source.
func (s *Source) Seek() {}

//Sets the currently playing position of the Source.
func (s *Source) SeekAtomic(offset float32) {}

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

//Stops a source.
func (s *Source) StopAtomic() {}

//Gets the currently playing position of the Source.
func (s *Source) Tell() {}

//Gets the currently playing position of the Source.
func (s *Source) TellAtomic() float32 {
	return 0.0
}

func (s *Source) Release() {}
