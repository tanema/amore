package audio

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/tanema/amore/audio/al"
	"github.com/tanema/amore/audio/decoding"
)

const (
	maxAttenuationDistance = 1000000.0 // upper limit of sound attentuation time.
	maxBuffers             = 8         //arbitrary limit of umber of buffers a source can use to stream
)

// Source manages decoding sound data, creates an openal sound and manages the
// data associated with the source.
type Source struct {
	decoder           *decoding.Decoder
	source            al.Source
	isStatic          bool
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
	offsetBytes       int32
}

// State indicates the current playing state of the source.
type State int

// Audio States
const (
	Unknown = State(0)
	Initial = State(al.Initial)
	Playing = State(al.Playing)
	Paused  = State(al.Paused)
	Stopped = State(al.Stopped)
)

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
		decoder:           decoder,
		isStatic:          static,
		pitch:             1,
		volume:            1,
		maxVolume:         1,
		referenceDistance: 1,
		rolloffFactor:     1,
		maxDistance:       maxAttenuationDistance,
		cone:              al.Cone{0, 0, 0},
		position:          al.Vector{},
		velocity:          al.Vector{},
		direction:         al.Vector{},
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
	if s.isStatic {
		s.source.SetLooping(s.looping)
	}
}

// GetAttenuationDistances returns the reference and maximum attenuation distances of the Source.
func (s *Source) GetAttenuationDistances() (float32, float32) {
	if s.isValid() {
		return s.source.ReferenceDistance(), s.source.MaxDistance()
	}
	return s.referenceDistance, s.maxDistance
}

// GetChannels returns the number of channels in the Source.
func (s *Source) GetChannels() int16 {
	return s.decoder.Channels
}

// GetDuration returns the total duration of the source.
func (s *Source) GetDuration() time.Duration {
	return s.decoder.Duration()
}

// GetCone returns the Source's directional volume cones by inner angle, outer angle,
// and outer volume.
func (s *Source) GetCone() (float32, float32, float32) {
	if s.isValid() {
		c := s.source.Cone()
		return mgl32.DegToRad(float32(c.InnerAngle)), mgl32.DegToRad(float32(c.OuterAngle)), c.OuterVolume
	}
	return mgl32.DegToRad(float32(s.cone.InnerAngle)), mgl32.DegToRad(float32(s.cone.OuterAngle)), s.cone.OuterVolume
}

// GetDirection returns the direction of the Source with a vector of x, y, z
func (s *Source) GetDirection() (float32, float32, float32) {
	if s.isValid() {
		d := s.source.Direction()
		return d[0], d[1], d[2]
	}
	return s.direction[0], s.direction[1], s.direction[2]
}

// GetPitch returns the current pitch of the Source in the range 0.0, 1.0
func (s *Source) GetPitch() float32 {
	if s.isValid() {
		return s.source.Pitch()
	}
	return s.pitch
}

// GetPosition returns the position of the Source in a point x, y, z
func (s *Source) GetPosition() (float32, float32, float32) {
	if s.isValid() {
		vec := s.source.Position()
		return vec[0], vec[1], vec[2]
	}
	return s.position[0], s.position[1], s.position[2]
}

// GetRolloff returns the rolloff factor of the source.
func (s *Source) GetRolloff() float32 {
	if s.isValid() {
		return s.source.Rolloff()
	}
	return s.rolloffFactor
}

// GetVelocity returns the velocity of the Source with a vector x, y, x
func (s *Source) GetVelocity() (float32, float32, float32) {
	if s.isValid() {
		vec := s.source.Velocity()
		return vec[0], vec[1], vec[2]
	}
	return s.velocity[0], s.velocity[1], s.velocity[2]
}

// GetVolume returns the current volume of the Source.
func (s *Source) GetVolume() float32 {
	if s.isValid() {
		return s.source.Gain()
	}
	return s.volume
}

// GetVolumeLimits returns the volume limits of the source, min, max.
func (s *Source) GetVolumeLimits() (float32, float32) {
	if s.isValid() {
		return s.source.MinGain(), s.source.MaxGain()
	}
	return s.minVolume, s.maxVolume
}

// GetState returns the playing state of the source.
//     source.GetState() == audio.Playing
func (s *Source) GetState() State {
	return State(s.source.State())
}

// IsLooping returns whether the Source will loop.
func (s *Source) IsLooping() bool {
	return s.looping
}

// IsPaused returns whether the Source is paused.
func (s *Source) IsPaused() bool {
	if s.isValid() {
		return s.GetState() == Paused
	}
	return false
}

// IsPlaying returns whether the Source is playing.
func (s *Source) IsPlaying() bool {
	if s.isValid() {
		return s.GetState() == Playing
	}
	return false
}

// IsRelative returns whether the Source's position and direction are relative
// to the listener.
func (s *Source) IsRelative() bool {
	if s.isValid() {
		return s.source.Relative()
	}
	return s.relative
}

// IsStatic returns whether the Source is static or stream.
func (s *Source) IsStatic() bool {
	return s.isStatic
}

// IsStopped returns whether the Source is stopped.
func (s *Source) IsStopped() bool {
	if s.isValid() {
		return s.GetState() == Stopped
	}
	return true
}

// SetAttenuationDistances sets the reference and maximum attenuation distances of the Source.
func (s *Source) SetAttenuationDistances(ref, max float32) {
	s.referenceDistance = ref
	s.maxDistance = max
	s.reset()
}

// SetCone sets the Source's directional volume cones with the inner angle, outer
// angle, and outer volume
func (s *Source) SetCone(innerAngle, outerAngle, outerVolume float32) {
	s.cone = al.Cone{
		InnerAngle:  int32(mgl32.RadToDeg(innerAngle)),
		OuterAngle:  int32(mgl32.RadToDeg(outerAngle)),
		OuterVolume: outerVolume,
	}
	s.reset()
}

// SetDirection sets the direction of the Source with the vector x, y, z
func (s *Source) SetDirection(x, y, z float32) {
	if s.GetChannels() > 1 {
		panic("This spatial audio functionality is only available for mono Sources. Ensure the Source is not multi-channel before calling this function.")
	}

	s.direction = al.Vector{x, y, z}
	s.reset()
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

// SetPosition sets the position of the Source at the point x, y, z
func (s *Source) SetPosition(x, y, z float32) {
	s.position = al.Vector{x, y, z}
	s.reset()
}

// SetRelative sets whether the Source's position and direction are relative to
// the listener.
func (s *Source) SetRelative(isRelative bool) {
	s.relative = isRelative
	s.reset()
}

// SetRolloff sets the rolloff factor.
func (s *Source) SetRolloff(rolloff float32) {
	s.rolloffFactor = rolloff
	s.reset()
}

// SetVelocity sets the velocity of the Source with the vector x, y, z
func (s *Source) SetVelocity(x, y, z float32) {
	s.velocity = al.Vector{x, y, z}
	s.reset()
}

// SetVolume sets the current volume of the Source.
func (s *Source) SetVolume(v float32) {
	s.volume = v
	s.reset()
}

// SetVolumeLimits sets the volume limits of the source both min and max
func (s *Source) SetVolumeLimits(min, max float32) {
	s.minVolume = min
	s.maxVolume = max
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
	if s.isValid() {
		s.offsetBytes = s.decoder.DurToByteOffset(offset)
		if !s.isStatic {
			waspaused := s.paused
			s.decoder.Seek(int64(s.offsetBytes))
			// Because we still have old data from before the seek in the buffers let's empty them.
			s.Stop()
			s.Play()
			if waspaused {
				s.Pause()
			}
		} else {
			pool.mutex.Lock()
			defer pool.mutex.Unlock()
			s.source.SetOffsetBytes(s.offsetBytes)
		}
	}
}

// Stop stops a playing source.
func (s *Source) Stop() {
	if s.isValid() {
		pool.mutex.Lock()
		al.StopSources(s.source)
		if !s.isStatic {
			queued := s.source.BuffersQueued()
			for i := queued; i > 0; i-- {
				buffer := s.source.UnqueueBuffer()
				al.DeleteBuffers(buffer)
			}
		}
		s.source.ClearBuffers()
		pool.release(s)
		pool.mutex.Unlock()
	}
	s.Rewind()
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
