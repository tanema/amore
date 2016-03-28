package webal

type Cone struct {
	InnerAngle  int32
	OuterAngle  int32
	OuterVolume float32
}

type Orientation struct {
	Forward [3]float32
	Up      [3]float32
}
