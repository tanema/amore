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
package openal

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
	Forward [3]float32
	// Up vector represents the rotation of the object.
	Up [3]float32
}

func orientationFromSlice(v []float32) Orientation {
	return Orientation{
		Forward: [3]float32{v[0], v[1], v[2]},
		Up:      [3]float32{v[3], v[4], v[5]},
	}
}

func (v Orientation) slice() []float32 {
	return []float32{v.Forward[0], v.Forward[1], v.Forward[2], v.Up[0], v.Up[1], v.Up[2]}
}

// DistanceModel returns the distance model.
func DistanceModel() int32 {
	return alGetInteger(paramDistanceModel)
}

// SetDistanceModel sets the distance model.
func SetDistanceModel(v int32) {
	alDistanceModel(v)
}

// DopplerFactor returns the doppler factor.
func DopplerFactor() float32 {
	return alGetFloat(paramDopplerFactor)
}

// SetDopplerFactor sets the doppler factor.
func SetDopplerFactor(v float32) {
	alDopplerFactor(v)
}

// DopplerVelocity returns the doppler velocity.
func DopplerVelocity() float32 {
	return alGetFloat(paramDopplerVelocity)
}

// SetDopplerVelocity sets the doppler velocity.
func SetDopplerVelocity(v float32) {
	alDopplerVelocity(v)
}

// SpeedOfSound is the speed of sound in meters per second (m/s).
func SpeedOfSound() float32 {
	return alGetFloat(paramSpeedOfSound)
}

// SetSpeedOfSound sets the speed of sound, its unit should be meters per second (m/s).
func SetSpeedOfSound(v float32) {
	alSpeedOfSound(v)
}

// Vendor returns the vendor.
func Vendor() string {
	return alGetString(paramVendor)
}

// Version returns the version string.
func Version() string {
	return alGetString(paramVersion)
}

// Renderer returns the renderer information.
func Renderer() string {
	return alGetString(paramRenderer)
}

// Extensions returns the enabled extensions.
func Extensions() string {
	return alGetString(paramExtensions)
}

// Error returns the most recently generated error.
func Error() int32 {
	return alGetError()
}

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

// ListenerGain returns the total gain applied to the final mix.
func ListenerGain() float32 {
	return alGetListenerf(paramGain)
}

// ListenerPosition returns the position of the listener.
func ListenerPosition() [3]float32 {
	v := [3]float32{}
	alGetListenerfv(paramPosition, v[:])
	return v
}

// ListenerVelocity returns the velocity of the listener.
func ListenerVelocity() [3]float32 {
	v := [3]float32{}
	alGetListenerfv(paramVelocity, v[:])
	return v
}

// ListenerOrientation returns the orientation of the listener.
func ListenerOrientation() Orientation {
	v := make([]float32, 6)
	alGetListenerfv(paramOrientation, v)
	return orientationFromSlice(v)
}

// SetListenerGain sets the total gain that will be applied to the final mix.
func SetListenerGain(v float32) {
	alListenerf(paramGain, v)
}

// SetListenerPosition sets the position of the listener.
func SetListenerPosition(v [3]float32) {
	alListenerfv(paramPosition, v[:])
}

// SetListenerVelocity sets the velocity of the listener.
func SetListenerVelocity(v [3]float32) {
	alListenerfv(paramVelocity, v[:])
}

// SetListenerOrientation sets the orientation of the listener.
func SetListenerOrientation(v Orientation) {
	alListenerfv(paramOrientation, v.slice())
}

// GenBuffers generates n new buffers. The generated buffers should be deleted
// once they are no longer in use.
func GenBuffers(n int) []Buffer {
	return alGenBuffers(n)
}

// DeleteBuffers deletes the buffers.
func DeleteBuffers(buffer ...Buffer) {
	alDeleteBuffers(buffer)
}
