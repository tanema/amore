// Package audio is use for creating audio sources, managing/pooling resources,
// and playback of those audio sources.
package audio

import (
	"github.com/tanema/amore/audio/al"
)

// DistanceModel defines sound attenuation.
type DistanceModel int32

// Distance models to be set
const (
	DistanceNone            DistanceModel = 0xD000
	DistanceInverse         DistanceModel = 0xD001
	DistanceInverseClamped  DistanceModel = 0xD002
	DistanceLinear          DistanceModel = 0xD003
	DistanceLinearClamped   DistanceModel = 0xD004
	DistanceExponent        DistanceModel = 0xD005
	DistanceExponentClamped DistanceModel = 0xD006
)

// init will open the audio interface.
func init() {
	if err := al.OpenDevice(); err != nil {
		panic(err)
	}
}

// GetDistanceModel returns the distance attenuation model
func GetDistanceModel() DistanceModel {
	return DistanceModel(al.DistanceModel())
}

// GetDopplerScale gets the global scale factor for doppler effects
func GetDopplerScale() float32 {
	return al.DopplerFactor()
}

// GetOrientation returns the orientation vectors of the listener in the order of
// Forward{x,y,z} Up{x,y,z}
func GetOrientation() (float32, float32, float32, float32, float32, float32) {
	ori := al.ListenerOrientation()
	return ori.Forward[0], ori.Forward[1], ori.Forward[2], ori.Up[0], ori.Up[1], ori.Up[2]
}

// GetPosition returns the position of the listener x, y, z
func GetPosition() (float32, float32, float32) {
	pos := al.ListenerPosition()
	return pos[0], pos[1], pos[2]
}

// GetSourceCount gets the current number of simultaneously playing sources
func GetSourceCount() int {
	return pool.getSourceCount()
}

// GetVelocity returns the velocity of the listener
func GetVelocity() (float32, float32, float32) {
	vel := al.ListenerVelocity()
	return vel[0], vel[1], vel[2]
}

// GetVolume returns the master volume.
func GetVolume() float32 { return al.ListenerGain() }

// Pause pauses a specified source, if source is nil it will pause all
func Pause(source *Source) { pool.pause(source) }

// Play plays a specified source, if source is nil it will play all
func Play(source *Source) { pool.play(source) }

// Resume playes a specified source, if source is nil it will resume all
func Resume(source *Source) { pool.resume(source) }

// Rewind rewinds a specified source, if source is nil it will rewind all
func Rewind(source *Source) { pool.rewind(source) }

// Stop stops a specified source, if source is nil it will stop all
func Stop(source *Source) { pool.stop(source) }

// SetDistanceModel sets the distance attenuation model
func SetDistanceModel(model DistanceModel) {
	al.SetDistanceModel(int32(model))
}

// SetDopplerScale sets a global scale factor for doppler effects
func SetDopplerScale(scale float32) {
	if scale >= 0.0 {
		al.SetDopplerFactor(scale)
	}
}

// SetOrientation sets the orientation of the listener in the order
// Front{x,y,z} Up{x,y,z}
func SetOrientation(fx, fy, fz, ux, uy, uz float32) {
	al.SetListenerOrientation(al.Orientation{
		Forward: al.Vector{fx, fy, fz},
		Up:      al.Vector{ux, uy, uz},
	})
}

// SetPosition sets the position of the listener
func SetPosition(x, y, z float32) {
	al.SetListenerPosition(al.Vector{x, y, z})
}

// SetVelocity sets the velocity of the listener
func SetVelocity(x, y, z float32) {
	al.SetListenerVelocity(al.Vector{x, y, z})
}

// SetVolume sets the master volume
func SetVolume(gain float32) {
	al.SetListenerGain(gain)
}
