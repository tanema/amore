// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux

// Package al provides OpenAL Soft bindings for Go.
//
// Calls are not safe for concurrent use.
//
// More information about OpenAL Soft is available at
// http://www.openal.org/documentation/openal-1.1-specification.pdf.
//
// In order to use this package on Linux desktop distros,
// you will need OpenAL library as an external dependency.
// On Ubuntu 14.04 'Trusty', you may have to install this library
// by running the command below.
//
// 		sudo apt-get install libopenal-dev
//
// When compiled for Android, this package uses OpenAL Soft. Please add its
// license file to the open source notices of your application.
// OpenAL Soft's license file could be found at
// http://repo.or.cz/w/openal-soft.git/blob/HEAD:/COPYING.
//
// Vendored and extended from github.com/golang/mobile for extending functionality

package al

// Capability represents OpenAL extension capabilities.
type Capability int32

// Enable enables a capability.
func Enable(c Capability) {
	alEnable(int32(c))
}

// Disable disables a capability.
func Disable(c Capability) {
	alDisable(int32(c))
}

// Enabled returns true if the specified capability is enabled.
func Enabled(c Capability) bool {
	return alIsEnabled(int32(c))
}

// Vector represents an vector in a Cartesian coordinate system.
type Vector [3]float32

type Cone struct {
	InnerAngle  int32
	OuterAngle  int32
	OuterVolume float32
}

// Orientation represents the angular position of an object in a
// right-handed Cartesian coordinate system.
// A cross product between the forward and up vector returns a vector
// that points to the right.
type Orientation struct {
	// Forward vector is the direction that the object is looking at.
	Forward Vector
	// Up vector represents the rotation of the object.
	Up Vector
}

func orientationFromSlice(v []float32) Orientation {
	return Orientation{
		Forward: Vector{v[0], v[1], v[2]},
		Up:      Vector{v[3], v[4], v[5]},
	}
}

func (v Orientation) slice() []float32 {
	return []float32{v.Forward[0], v.Forward[1], v.Forward[2], v.Up[0], v.Up[1], v.Up[2]}
}

func geti(param int) int32 {
	return alGetInteger(param)
}

func getf(param int) float32 {
	return alGetFloat(param)
}

func getString(param int) string {
	return alGetString(param)
}

// DistanceModel returns the distance model.
func DistanceModel() int32 {
	return geti(paramDistanceModel)
}

// SetDistanceModel sets the distance model.
func SetDistanceModel(v int32) {
	alDistanceModel(v)
}

// DopplerFactor returns the doppler factor.
func DopplerFactor() float32 {
	return getf(paramDopplerFactor)
}

// SetDopplerFactor sets the doppler factor.
func SetDopplerFactor(v float32) {
	alDopplerFactor(v)
}

// DopplerVelocity returns the doppler velocity.
func DopplerVelocity() float32 {
	return getf(paramDopplerVelocity)
}

// SetDopplerVelocity sets the doppler velocity.
func SetDopplerVelocity(v float32) {
	alDopplerVelocity(v)
}

// SpeedOfSound is the speed of sound in meters per second (m/s).
func SpeedOfSound() float32 {
	return getf(paramSpeedOfSound)
}

// SetSpeedOfSound sets the speed of sound, its unit should be meters per second (m/s).
func SetSpeedOfSound(v float32) {
	alSpeedOfSound(v)
}

// Vendor returns the vendor.
func Vendor() string {
	return getString(paramVendor)
}

// Version returns the version string.
func Version() string {
	return getString(paramVersion)
}

// Renderer returns the renderer information.
func Renderer() string {
	return getString(paramRenderer)
}

// Extensions returns the enabled extensions.
func Extensions() string {
	return getString(paramExtensions)
}

// Error returns the most recently generated error.
func Error() int32 {
	return alGetError()
}

// Source represents an individual sound source in 3D-space.
// They take PCM data, apply modifications and then submit them to
// be mixed according to their spatial location.
type Source uint32

// GenSources generates n new sources. These sources should be deleted
// once they are not in use.
func GenSources(n int) []Source {
	return alGenSources(n)
}

// PlaySources plays the sources.
func PlaySources(source ...Source) {
	alSourcePlayv(source)
}

// PauseSources pauses the sources.
func PauseSources(source ...Source) {
	alSourcePausev(source)
}

// StopSources stops the sources.
func StopSources(source ...Source) {
	alSourceStopv(source)
}

// RewindSources rewinds the sources to their beginning positions.
func RewindSources(source ...Source) {
	alSourceRewindv(source)
}

// DeleteSources deletes the sources.
func DeleteSources(source ...Source) {
	alDeleteSources(source)
}

// Gain returns the source gain.
func (s Source) Gain() float32 {
	return getSourcef(s, paramGain)
}

// SetGain sets the source gain.
func (s Source) SetGain(v float32) {
	setSourcef(s, paramGain, v)
}

func (s Source) Pitch() float32 {
	return getSourcef(s, paramPitch)
}

func (s Source) SetPitch(p float32) {
	setSourcef(s, paramPitch, p)
}

func (s Source) Rolloff() float32 {
	return getSourcef(s, paramRolloffFactor)
}

func (s Source) SetRolloff(roll_off float32) {
	setSourcef(s, paramRolloffFactor, roll_off)
}

func (s Source) ReferenceDistance() float32 {
	return getSourcef(s, paramReferenceDistance)
}

func (s Source) SetReferenceDistance(dis float32) {
	setSourcef(s, paramReferenceDistance, dis)
}

func (s Source) MaxDistance() float32 {
	return getSourcef(s, paramMaxDistance)
}

func (s Source) SetMaxDistance(dis float32) {
	setSourcef(s, paramMaxDistance, dis)
}

// Looping returns the source looping.
func (s Source) Looping() bool {
	return getSourcei(s, paramLooping) == 1
}

// SetLooping sets the source looping
func (s Source) SetLooping(should_loop bool) {
	if should_loop {
		setSourcei(s, paramLooping, 1)
	} else {
		setSourcei(s, paramLooping, 0)
	}
}

func (s Source) Relative() bool {
	return getSourcei(s, paramSourceRelative) == 1
}

func (s Source) SetRelative(is_relative bool) {
	if is_relative {
		setSourcei(s, paramSourceRelative, 1)
	} else {
		setSourcei(s, paramSourceRelative, 0)
	}
}

// MinGain returns the source's minimum gain setting.
func (s Source) MinGain() float32 {
	return getSourcef(s, paramMinGain)
}

// SetMinGain sets the source's minimum gain setting.
func (s Source) SetMinGain(v float32) {
	setSourcef(s, paramMinGain, v)
}

// MaxGain returns the source's maximum gain setting.
func (s Source) MaxGain() float32 {
	return getSourcef(s, paramMaxGain)
}

// SetMaxGain sets the source's maximum gain setting.
func (s Source) SetMaxGain(v float32) {
	setSourcef(s, paramMaxGain, v)
}

// Position returns the position of the source.
func (s Source) Position() Vector {
	v := Vector{}
	getSourcefv(s, paramPosition, v[:])
	return v
}

// SetPosition sets the position of the source.
func (s Source) SetPosition(v Vector) {
	setSourcefv(s, paramPosition, v[:])
}

// Position returns the position of the source.
func (s Source) Direction() Vector {
	v := Vector{}
	getSourcefv(s, paramDirection, v[:])
	return v
}

// SetDirection sets the direction of the source.
func (s Source) SetDirection(v Vector) {
	setSourcefv(s, paramDirection, v[:])
}

func (s Source) Cone() Cone {
	return Cone{
		InnerAngle:  getSourcei(s, paramConeInnerAngle),
		OuterAngle:  getSourcei(s, paramConeOuterAngle),
		OuterVolume: getSourcef(s, paramConeOuterGain),
	}
}

func (s Source) SetCone(c Cone) {
	setSourcei(s, paramConeInnerAngle, c.InnerAngle)
	setSourcei(s, paramConeOuterAngle, c.OuterAngle)
	setSourcef(s, paramConeOuterGain, c.OuterVolume)
}

// Velocity returns the source's velocity.
func (s Source) Velocity() Vector {
	v := Vector{}
	getSourcefv(s, paramVelocity, v[:])
	return v
}

// SetVelocity sets the source's velocity.
func (s Source) SetVelocity(v Vector) {
	setSourcefv(s, paramVelocity, v[:])
}

// Orientation returns the orientation of the source.
func (s Source) Orientation() Orientation {
	v := make([]float32, 6)
	getSourcefv(s, paramOrientation, v)
	return orientationFromSlice(v)
}

// SetOrientation sets the orientation of the source.
func (s Source) SetOrientation(o Orientation) {
	setSourcefv(s, paramOrientation, o.slice())
}

// State returns the playing state of the source.
func (s Source) State() int32 {
	return getSourcei(s, paramSourceState)
}

// BuffersQueued returns the number of the queued buffers.
func (s Source) SetBuffer(b Buffer) {
	setSourcei(s, paramBuffer, int32(b))
}

// BuffersQueued returns the number of the queued buffers.
func (s Source) Buffer() Buffer {
	return Buffer(getSourcei(s, paramBuffer))
}

// BuffersQueued returns the number of the queued buffers.
func (s Source) ClearBuffers() {
	setSourcei(s, paramBuffer, 0)
}

// BuffersQueued returns the number of the queued buffers.
func (s Source) BuffersQueued() int32 {
	return getSourcei(s, paramBuffersQueued)
}

// BuffersProcessed returns the number of the processed buffers.
func (s Source) BuffersProcessed() int32 {
	return getSourcei(s, paramBuffersProcessed)
}

// OffsetSeconds returns the current playback position of the source in seconds.
func (s Source) OffsetSeconds() float32 {
	return getSourcef(s, paramSecOffset)
}

// OffsetSeconds returns the current playback position of the source in seconds.
func (s Source) SetOffsetSeconds(seconds float32) {
	setSourcef(s, paramSecOffset, seconds)
}

// OffsetSample returns the sample offset of the current playback position.
func (s Source) OffsetSample() float32 {
	return getSourcef(s, paramSampleOffset)
}

// OffsetSample returns the sample offset of the current playback position.
func (s Source) SetOffsetSample(samples float32) {
	setSourcef(s, paramSampleOffset, samples)
}

// OffsetByte returns the byte offset of the current playback position.
func (s Source) OffsetByte() int32 {
	return getSourcei(s, paramByteOffset)
}

// OffsetSample returns the sample offset of the current playback position.
func (s Source) SetOffsetBytes(bytes int32) {
	setSourcei(s, paramByteOffset, bytes)
}

func getSourcei(s Source, param int) int32 {
	return alGetSourcei(s, param)
}

func getSourcef(s Source, param int) float32 {
	return alGetSourcef(s, param)
}

func getSourcefv(s Source, param int, v []float32) {
	alGetSourcefv(s, param, v)
}

func setSourcei(s Source, param int, v int32) {
	alSourcei(s, param, v)
}

func setSourcef(s Source, param int, v float32) {
	alSourcef(s, param, v)
}

func setSourcefv(s Source, param int, v []float32) {
	alSourcefv(s, param, v)
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

// ListenerGain returns the total gain applied to the final mix.
func ListenerGain() float32 {
	return getListenerf(paramGain)
}

// ListenerPosition returns the position of the listener.
func ListenerPosition() Vector {
	v := Vector{}
	getListenerfv(paramPosition, v[:])
	return v
}

// ListenerVelocity returns the velocity of the listener.
func ListenerVelocity() Vector {
	v := Vector{}
	getListenerfv(paramVelocity, v[:])
	return v
}

// ListenerOrientation returns the orientation of the listener.
func ListenerOrientation() Orientation {
	v := make([]float32, 6)
	getListenerfv(paramOrientation, v)
	return orientationFromSlice(v)
}

// SetListenerGain sets the total gain that will be applied to the final mix.
func SetListenerGain(v float32) {
	setListenerf(paramGain, v)
}

// SetListenerPosition sets the position of the listener.
func SetListenerPosition(v Vector) {
	setListenerfv(paramPosition, v[:])
}

// SetListenerVelocity sets the velocity of the listener.
func SetListenerVelocity(v Vector) {
	setListenerfv(paramVelocity, v[:])
}

// SetListenerOrientation sets the orientation of the listener.
func SetListenerOrientation(v Orientation) {
	setListenerfv(paramOrientation, v.slice())
}

func getListenerf(param int) float32 {
	return alGetListenerf(param)
}

func getListenerfv(param int, v []float32) {
	alGetListenerfv(param, v)
}

func setListenerf(param int, v float32) {
	alListenerf(param, v)
}

func setListenerfv(param int, v []float32) {
	alListenerfv(param, v)
}

// A buffer represents a chunk of PCM audio data that could be buffered to an audio
// source. A single buffer could be shared between multiple sources.
type Buffer uint32

// GenBuffers generates n new buffers. The generated buffers should be deleted
// once they are no longer in use.
func GenBuffers(n int) []Buffer {
	return alGenBuffers(n)
}

// DeleteBuffers deletes the buffers.
func DeleteBuffers(buffer ...Buffer) {
	alDeleteBuffers(buffer)
}

func getBufferi(b Buffer, param int) int32 {
	return alGetBufferi(b, param)
}

// Frequency returns the frequency of the buffer data in Hertz (Hz).
func (b Buffer) Frequency() int32 {
	return getBufferi(b, paramFreq)
}

// Bits return the number of bits used to represent a sample.
func (b Buffer) Bits() int32 {
	return getBufferi(b, paramBits)
}

// Channels return the number of the audio channels.
func (b Buffer) Channels() int32 {
	return getBufferi(b, paramChannels)
}

// Size returns the size of the data.
func (b Buffer) Size() int32 {
	return getBufferi(b, paramSize)
}

// BufferData buffers PCM data to the current buffer.
func (b Buffer) BufferData(format uint32, data []byte, freq int32) {
	alBufferData(b, format, data, freq)
}

// Valid returns true if the buffer exists and is valid.
func (b Buffer) Valid() bool {
	return alIsBuffer(b)
}
