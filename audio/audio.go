package audio

import (
	"github.com/tanema/amore/audio/al"
)

type DistanceModel int32

const (
	DISTANCE_NONE             DistanceModel = 0xD000
	DISTANCE_INVERSE          DistanceModel = 0xD001
	DISTANCE_INVERSE_CLAMPED  DistanceModel = 0xD002
	DISTANCE_LINEAR           DistanceModel = 0xD003
	DISTANCE_LINEAR_CLAMPED   DistanceModel = 0xD004
	DISTANCE_EXPONENT         DistanceModel = 0xD005
	DISTANCE_EXPONENT_CLAMPED DistanceModel = 0xD006
)

func init() {
	if err := al.OpenDevice(); err != nil {
		panic(err)
	}
	createPool()
}

//	Returns the distance attenuation model
func GetDistanceModel() DistanceModel {
	return DistanceModel(al.DistanceModel())
}

//	Gets the global scale factor for doppler effects
func GetDopplerScale() float32 {
	return al.DopplerFactor()
}

//	Returns the orientation of the listener
func GetOrientation() (float32, float32, float32, float32, float32, float32) {
	ori := al.ListenerOrientation()
	return ori.Forward[0], ori.Forward[1], ori.Forward[2], ori.Up[0], ori.Up[1], ori.Up[2]
}

//	Returns the position of the listener
func GetPosition() (float32, float32, float32) {
	pos := al.ListenerPosition()
	return pos[0], pos[1], pos[2]
}

//	Gets the current number of simultaneously playing sources
func GetSourceCount() int {
	return pool.GetSourceCount()
}

func GetMaxSources() int {
	return pool.GetMaxSources()
}

//	Returns the velocity of the listener
func GetVelocity() (float32, float32, float32) {
	vel := al.ListenerVelocity()
	return vel[0], vel[1], vel[2]
}

//	Returns the master volume
func GetVolume() float32 {
	return al.ListenerGain()
}

//	Pauses all audio
func Pause(source *Source) {
	if source == nil {
		pool.Pause(nil)
	} else {
		source.Pause()
	}
}

//	Plays the specified Source
func Play(source *Source) {
	if source == nil {
		pool.Resume(nil)
	} else {
		source.Play()
	}
}

//	Resumes all audio
func Resume(source *Source) {
	if source == nil {
		pool.Resume(nil)
	} else {
		source.Resume()
	}
}

//	Rewinds all playing audio
func Rewind(source *Source) {
	if source == nil {
		pool.Rewind(nil)
	} else {
		source.Rewind()
	}
}

//	Sets the distance attenuation model
func SetDistanceModel(model DistanceModel) {
	al.SetDistanceModel(int32(model))
}

//	Sets a global scale factor for doppler effects
func SetDopplerScale(scale float32) {
	if scale >= 0.0 {
		al.SetDopplerFactor(scale)
	}
}

//	Sets the orientation of the listener
func SetOrientation(fx, fy, fz, ux, uy, uz float32) {
	al.SetListenerOrientation(al.Orientation{
		Forward: al.Vector{fx, fy, fz},
		Up:      al.Vector{ux, uy, uz},
	})
}

//	Sets the position of the listener
func SetPosition(x, y, z float32) {
	al.SetListenerPosition(al.Vector{x, y, z})
}

//	Sets the velocity of the listener
func SetVelocity(x, y, z float32) {
	al.SetListenerVelocity(al.Vector{x, y, z})
}

//	Sets the master volume
func SetVolume(gain float32) {
	al.SetListenerGain(gain)
}

//	Stops currently played sources.
func Stop(source *Source) {
	if source == nil {
		pool.Stop(nil)
	} else {
		source.Stop()
	}
}
