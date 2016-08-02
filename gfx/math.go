package gfx

import (
	"math"
)

const Pi = float32(math.Pi)
const MaxInt32 = math.MaxInt32
const MaxUint16 = math.MaxUint16
const MaxUint32 = math.MaxUint32

func Maxi(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

func Mini(x, y int) int {
	return int(math.Min(float64(x), float64(y)))
}

func Maxi32(x, y int32) int32 {
	return int32(math.Max(float64(x), float64(y)))
}

func Mini32(x, y int32) int32 {
	return int32(math.Min(float64(x), float64(y)))
}

func Max(x, y float32) float32 {
	return float32(math.Max(float64(x), float64(y)))
}

func Min(x, y float32) float32 {
	return float32(math.Min(float64(x), float64(y)))
}

func Abs(a float32) float32 {
	if a < 0 {
		return -a
	} else if a == 0 {
		return 0
	}

	return a
}

func Clamp(a, low, high float32) float32 {
	if a < low {
		return low
	} else if a > high {
		return high
	}

	return a
}

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Atan2(x, y float32) float32 {
	return float32(math.Atan2(float64(x), float64(y)))
}

func Round(v float32, precision int) float32 {
	p := float64(precision)
	t := float64(v) * math.Pow(10, p)
	if t > 0 {
		return float32(math.Floor(t+0.5) / math.Pow(10, p))
	}
	return float32(math.Ceil(t-0.5) / math.Pow(10, p))
}

// Converts degrees to radians
func DegToRad(angle float32) float32 {
	return angle * float32(math.Pi) / 180
}

// Converts radians to degrees
func RadToDeg(angle float32) float32 {
	return angle * 180 / float32(math.Pi)
}
