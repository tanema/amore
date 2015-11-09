package audio

import (
	"golang.org/x/mobile/exp/audio/al"
)

func init() {
	if err := al.OpenDevice(); err != nil {
		panic(err)
	}
	createPool()
}

//	Creates a new Source from a file, SoundData, or Decoder
func NewSource(filepath string) {
}

//	Returns the distance attenuation model
func GetDistanceModel() {}

//	Gets the global scale factor for doppler effects
func GetDopplerScale() {}

//	Returns the orientation of the listener
func GetOrientation() {}

//	Returns the position of the listener
func GetPosition() {}

//	Gets the current number of simultaneously playing sources
func GetSourceCount() {}

//	Returns the velocity of the listener
func GetVelocity() {}

//	Returns the master volume
func GetVolume() {}

//	Pauses all audio
func Pause() {}

//	Plays the specified Source
func Play() {}

//	Resumes all audio
func Resume() {}

//	Rewinds all playing audio
func Rewind() {}

//	Sets the distance attenuation model
func SetDistanceModel() {}

//	Sets a global scale factor for doppler effects
func SetDopplerScale() {}

//	Sets the orientation of the listener
func SetOrientation() {}

//	Sets the position of the listener
func SetPosition() {}

//	Sets the velocity of the listener
func SetVelocity() {}

//	Sets the master volume
func SetVolume() {}

//	Stops currently played sources.
func Stop() {}
