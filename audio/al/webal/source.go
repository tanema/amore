package webal

type Source struct{}

func (s Source) IsValid() bool                    { return false }
func (s Source) Gain() float32                    { return 0 }
func (s Source) SetGain(v float32)                {}
func (s Source) Pitch() float32                   { return 0 }
func (s Source) SetPitch(p float32)               {}
func (s Source) Rolloff() float32                 { return 0 }
func (s Source) SetRolloff(roll_off float32)      {}
func (s Source) ReferenceDistance() float32       { return 0 }
func (s Source) SetReferenceDistance(dis float32) {}
func (s Source) MaxDistance() float32             { return 0 }
func (s Source) SetMaxDistance(dis float32)       {}
func (s Source) Looping() bool                    { return false }
func (s Source) SetLooping(should_loop bool)      {}
func (s Source) Relative() bool                   { return false }
func (s Source) SetRelative(is_relative bool)     {}
func (s Source) MinGain() float32                 { return 0 }
func (s Source) SetMinGain(v float32)             {}
func (s Source) MaxGain() float32                 { return 0 }
func (s Source) SetMaxGain(v float32)             {}
func (s Source) Position() [3]float32             { return [3]float32{} }
func (s Source) SetPosition(v [3]float32)         {}
func (s Source) Direction() [3]float32            { return [3]float32{} }
func (s Source) SetDirection(v [3]float32)        {}
func (s Source) Cone() Cone                       { return Cone{} }
func (s Source) SetCone(c Cone)                   {}
func (s Source) Velocity() [3]float32             { return [3]float32{} }
func (s Source) SetVelocity(v [3]float32)         {}
func (s Source) Orientation() Orientation         { return Orientation{} }
func (s Source) SetOrientation(o Orientation)     {}
func (s Source) State() int32                     { return 0 }
func (s Source) SetBuffer(b Buffer)               {}
func (s Source) Buffer() Buffer                   { return Buffer{} }
func (s Source) ClearBuffers()                    {}
func (s Source) BuffersQueued() int32             { return 0 }
func (s Source) BuffersProcessed() int32          { return 0 }
func (s Source) OffsetSeconds() float32           { return 0 }
func (s Source) SetOffsetSeconds(seconds float32) {}
func (s Source) OffsetSample() float32            { return 0 }
func (s Source) SetOffsetSample(samples float32)  {}
func (s Source) OffsetByte() int32                { return 0 }
func (s Source) SetOffsetBytes(bytes int32)       {}
func (s Source) QueueBuffers(buffer ...Buffer)    {}
func (s Source) UnqueueBuffer() Buffer            { return Buffer{} }
