// +build js

package al

func Init() error                                           { return nil }
func DistanceModel() int32                                  { return 0 }
func SetDistanceModel(v int32)                              {}
func DopplerFactor() float32                                { return 0 }
func SetDopplerFactor(v float32)                            {}
func DopplerVelocity() float32                              { return 0 }
func SetDopplerVelocity(v float32)                          {}
func SpeedOfSound() float32                                 { return 0 }
func SetSpeedOfSound(v float32)                             {}
func Vendor() string                                        { return "web audio api" }
func Version() string                                       { return "0" }
func Error() int32                                          { return 0 }
func CreateSource() Source                                  { return Source{} }
func PlaySource(source Source)                              {}
func PauseSource(source Source)                             {}
func StopSource(source Source)                              {}
func RewindSource(source Source)                            {}
func DeleteSource(source Source)                            {}
func ListenerGain() float32                                 { return 0 }
func ListenerPosition() [3]float32                          { return [3]float32{} }
func ListenerVelocity() [3]float32                          { return [3]float32{} }
func ListenerOrientation() Orientation                      { return Orientation{} }
func SetListenerGain(v float32)                             {}
func SetListenerPosition(v [3]float32)                      {}
func SetListenerVelocity(v [3]float32)                      {}
func SetListenerOrientation(fx, fy, fz, ux, uy, uz float32) {}
func CreateBuffer() Buffer                                  { return Buffer{} }
func DeleteBuffer(buffer Buffer)                            {}
func (s Source) UnqueueBuffer() Buffer                      { return Buffer{} }
func (s Source) QueueBuffer(buffer Buffer)                  {}
func (s Source) SetCone(c Cone)                             {}
func (s Source) SetBuffer(b Buffer)                         {}
