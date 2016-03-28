package openal

// Source represents an individual sound source in 3D-space.
// They take PCM data, apply modifications and then submit them to
// be mixed according to their spatial location.
type Source uint32

// Gain returns the source gain.
func (s Source) IsValid() bool {
	return s != 0
}

// Gain returns the source gain.
func (s Source) Gain() float32 {
	return alGetSourcef(s, paramGain)
}

// SetGain sets the source gain.
func (s Source) SetGain(v float32) {
	alSourcef(s, paramGain, v)
}

func (s Source) Pitch() float32 {
	return alGetSourcef(s, paramPitch)
}

func (s Source) SetPitch(p float32) {
	alSourcef(s, paramPitch, p)
}

func (s Source) Rolloff() float32 {
	return alGetSourcef(s, paramRolloffFactor)
}

func (s Source) SetRolloff(roll_off float32) {
	alSourcef(s, paramRolloffFactor, roll_off)
}

func (s Source) ReferenceDistance() float32 {
	return alGetSourcef(s, paramReferenceDistance)
}

func (s Source) SetReferenceDistance(dis float32) {
	alSourcef(s, paramReferenceDistance, dis)
}

func (s Source) MaxDistance() float32 {
	return alGetSourcef(s, paramMaxDistance)
}

func (s Source) SetMaxDistance(dis float32) {
	alSourcef(s, paramMaxDistance, dis)
}

// Looping returns the source looping.
func (s Source) Looping() bool {
	return alGetSourcei(s, paramLooping) == 1
}

// SetLooping sets the source looping
func (s Source) SetLooping(should_loop bool) {
	if should_loop {
		alSourcei(s, paramLooping, 1)
	} else {
		alSourcei(s, paramLooping, 0)
	}
}

func (s Source) Relative() bool {
	return alGetSourcei(s, paramSourceRelative) == 1
}

func (s Source) SetRelative(is_relative bool) {
	if is_relative {
		alSourcei(s, paramSourceRelative, 1)
	} else {
		alSourcei(s, paramSourceRelative, 0)
	}
}

// MinGain returns the source's minimum gain setting.
func (s Source) MinGain() float32 {
	return alGetSourcef(s, paramMinGain)
}

// SetMinGain sets the source's minimum gain setting.
func (s Source) SetMinGain(v float32) {
	alSourcef(s, paramMinGain, v)
}

// MaxGain returns the source's maximum gain setting.
func (s Source) MaxGain() float32 {
	return alGetSourcef(s, paramMaxGain)
}

// SetMaxGain sets the source's maximum gain setting.
func (s Source) SetMaxGain(v float32) {
	alSourcef(s, paramMaxGain, v)
}

// Position returns the position of the source.
func (s Source) Position() [3]float32 {
	v := [3]float32{}
	alGetSourcefv(s, paramPosition, v[:])
	return v
}

// SetPosition sets the position of the source.
func (s Source) SetPosition(v [3]float32) {
	alSourcefv(s, paramPosition, v[:])
}

// Position returns the position of the source.
func (s Source) Direction() [3]float32 {
	v := [3]float32{}
	alGetSourcefv(s, paramDirection, v[:])
	return v
}

// SetDirection sets the direction of the source.
func (s Source) SetDirection(v [3]float32) {
	alSourcefv(s, paramDirection, v[:])
}

func (s Source) Cone() Cone {
	return Cone{
		InnerAngle:  alGetSourcei(s, paramConeInnerAngle),
		OuterAngle:  alGetSourcei(s, paramConeOuterAngle),
		OuterVolume: alGetSourcef(s, paramConeOuterGain),
	}
}

func (s Source) SetCone(c Cone) {
	alSourcei(s, paramConeInnerAngle, c.InnerAngle)
	alSourcei(s, paramConeOuterAngle, c.OuterAngle)
	alSourcef(s, paramConeOuterGain, c.OuterVolume)
}

// Velocity returns the source's velocity.
func (s Source) Velocity() [3]float32 {
	v := [3]float32{}
	alGetSourcefv(s, paramVelocity, v[:])
	return v
}

// SetVelocity sets the source's velocity.
func (s Source) SetVelocity(v [3]float32) {
	alSourcefv(s, paramVelocity, v[:])
}

// Orientation returns the orientation of the source.
func (s Source) Orientation() Orientation {
	v := make([]float32, 6)
	alGetSourcefv(s, paramOrientation, v)
	return orientationFromSlice(v)
}

// SetOrientation sets the orientation of the source.
func (s Source) SetOrientation(o Orientation) {
	alSourcefv(s, paramOrientation, o.slice())
}

// State returns the playing state of the source.
func (s Source) State() int32 {
	return alGetSourcei(s, paramSourceState)
}

// BuffersQueued returns the number of the queued buffers.
func (s Source) SetBuffer(b Buffer) {
	alSourcei(s, paramBuffer, int32(b))
}

// BuffersQueued returns the number of the queued buffers.
func (s Source) Buffer() Buffer {
	return Buffer(alGetSourcei(s, paramBuffer))
}

// BuffersQueued returns the number of the queued buffers.
func (s Source) ClearBuffers() {
	alSourcei(s, paramBuffer, 0)
}

// BuffersQueued returns the number of the queued buffers.
func (s Source) BuffersQueued() int32 {
	return alGetSourcei(s, paramBuffersQueued)
}

// BuffersProcessed returns the number of the processed buffers.
func (s Source) BuffersProcessed() int32 {
	return alGetSourcei(s, paramBuffersProcessed)
}

// OffsetSeconds returns the current playback position of the source in seconds.
func (s Source) OffsetSeconds() float32 {
	return alGetSourcef(s, paramSecOffset)
}

// OffsetSeconds returns the current playback position of the source in seconds.
func (s Source) SetOffsetSeconds(seconds float32) {
	alSourcef(s, paramSecOffset, seconds)
}

// OffsetSample returns the sample offset of the current playback position.
func (s Source) OffsetSample() float32 {
	return alGetSourcef(s, paramSampleOffset)
}

// OffsetSample returns the sample offset of the current playback position.
func (s Source) SetOffsetSample(samples float32) {
	alSourcef(s, paramSampleOffset, samples)
}

// OffsetByte returns the byte offset of the current playback position.
func (s Source) OffsetByte() int32 {
	return alGetSourcei(s, paramByteOffset)
}

// OffsetSample returns the sample offset of the current playback position.
func (s Source) SetOffsetBytes(bytes int32) {
	alSourcei(s, paramByteOffset, bytes)
}

// QueueBuffers adds the buffers to the buffer queue.
func (s Source) QueueBuffers(buffer ...Buffer) {
	alSourceQueueBuffers(s, buffer)
}

// UnqueueBuffers removes the specified buffers from the buffer queue.
func (s Source) UnqueueBuffer() Buffer {
	buffers := make([]Buffer, 1)
	alSourceUnqueueBuffers(s, buffers)
	return buffers[0]
}
