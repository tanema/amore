// +build !js

package al

import (
	"github.com/tanema/amore/audio/al/openal"
)

func Init() error {
	return openal.OpenDevice()
}

// DistanceModel returns the distance model.
func DistanceModel() int32 {
	return openal.DistanceModel()
}

// SetDistanceModel sets the distance model.
func SetDistanceModel(v int32) {
	openal.SetDistanceModel(v)
}

// DopplerFactor returns the doppler factor.
func DopplerFactor() float32 {
	return openal.DopplerFactor()
}

// SetDopplerFactor sets the doppler factor.
func SetDopplerFactor(v float32) {
	openal.SetDopplerFactor(v)
}

// DopplerVelocity returns the doppler velocity.
func DopplerVelocity() float32 {
	return openal.DopplerVelocity()
}

// SetDopplerVelocity sets the doppler velocity.
func SetDopplerVelocity(v float32) {
	openal.SetDopplerVelocity(v)
}

// SpeedOfSound is the speed of sound in meters per second (m/s).
func SpeedOfSound() float32 {
	return openal.SpeedOfSound()
}

// SetSpeedOfSound sets the speed of sound, its unit should be meters per second (m/s).
func SetSpeedOfSound(v float32) {
	openal.SetSpeedOfSound(v)
}

// Vendor returns the vendor.
func Vendor() string {
	return openal.Vendor()
}

// Version returns the version string.
func Version() string {
	return openal.Version()
}

// Error returns the most recently generated error.
func Error() int32 {
	return openal.Error()
}

// GenSources generates n new sources. These sources should be deleted
// once they are not in use.
func CreateSource() Source {
	return Source{openal.GenSources(1)[0]}
}

// PlaySources plays the sources.
func PlaySource(source Source) {
	openal.PlaySources(source.Source)
}

// PauseSources pauses the sources.
func PauseSource(source Source) {
	openal.PauseSources(source.Source)
}

// StopSources stops the sources.
func StopSource(source Source) {
	openal.StopSources(source.Source)
}

// RewindSources rewinds the sources to their beginning positions.
func RewindSource(source Source) {
	openal.RewindSources(source.Source)
}

// DeleteSources deletes the sources.
func DeleteSource(source Source) {
	openal.DeleteSources(source.Source)
}

// ListenerGain returns the total gain applied to the final mix.
func ListenerGain() float32 {
	return openal.ListenerGain()
}

// ListenerPosition returns the position of the listener.
func ListenerPosition() [3]float32 {
	return openal.ListenerPosition()
}

// ListenerVelocity returns the velocity of the listener.
func ListenerVelocity() [3]float32 {
	return openal.ListenerVelocity()
}

// ListenerOrientation returns the orientation of the listener.
func ListenerOrientation() Orientation {
	return Orientation{openal.ListenerOrientation()}
}

// SetListenerGain sets the total gain that will be applied to the final mix.
func SetListenerGain(v float32) {
	openal.SetListenerGain(v)
}

// SetListenerPosition sets the position of the listener.
func SetListenerPosition(v [3]float32) {
	openal.SetListenerPosition(v)
}

// SetListenerVelocity sets the velocity of the listener.
func SetListenerVelocity(v [3]float32) {
	openal.SetListenerVelocity(v)
}

// SetListenerOrientation sets the orientation of the listener.
func SetListenerOrientation(fx, fy, fz, ux, uy, uz float32) {
	openal.SetListenerOrientation(openal.Orientation{
		Forward: [3]float32{fx, fy, fz},
		Up:      [3]float32{ux, uy, uz},
	})
}

// GenBuffers generates n new buffers. The generated buffers should be deleted
// once they are no longer in use.
func CreateBuffer() Buffer {
	return Buffer{openal.GenBuffers(1)[0]}
}

// DeleteBuffers deletes the buffers.
func DeleteBuffer(buffer Buffer) {
	openal.DeleteBuffers(buffer.Buffer)
}

// UnqueueBuffers removes the specified buffers from the buffer queue.
func (s Source) UnqueueBuffer() Buffer {
	return Buffer{s.Source.UnqueueBuffer()}
}

// QueueBuffers adds the buffers to the buffer queue.
func (s Source) QueueBuffer(buffer Buffer) {
	s.Source.QueueBuffers(buffer.Buffer)
}

func (s Source) SetCone(c Cone) {
	s.Source.SetCone(openal.Cone(c))
}

func (s Source) SetBuffer(b Buffer) {
	s.Source.SetBuffer(b.Buffer)
}
