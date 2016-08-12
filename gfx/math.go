package gfx

import (
	"math"
)

const pi = float32(math.Pi)
const maxInt32 = math.MaxInt32
const maxUint16 = math.MaxUint16
const maxUint32 = math.MaxUint32

func maxi(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

func mini(x, y int) int {
	return int(math.Min(float64(x), float64(y)))
}

func maxi32(x, y int32) int32 {
	return int32(math.Max(float64(x), float64(y)))
}

func mini32(x, y int32) int32 {
	return int32(math.Min(float64(x), float64(y)))
}

func max(x, y float32) float32 {
	return float32(math.Max(float64(x), float64(y)))
}

func min(x, y float32) float32 {
	return float32(math.Min(float64(x), float64(y)))
}

func abs(a float32) float32 {
	if a < 0 {
		return -a
	} else if a == 0 {
		return 0
	}

	return a
}

func clamp(a, low, high float32) float32 {
	if a < low {
		return low
	} else if a > high {
		return high
	}

	return a
}

func cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func atan2(x, y float32) float32 {
	return float32(math.Atan2(float64(x), float64(y)))
}

func round(v float32, precision int) float32 {
	p := float64(precision)
	t := float64(v) * math.Pow(10, p)
	if t > 0 {
		return float32(math.Floor(t+0.5) / math.Pow(10, p))
	}
	return float32(math.Ceil(t-0.5) / math.Pow(10, p))
}

// Converts degrees to radians
func degToRad(angle float32) float32 {
	return angle * float32(math.Pi) / 180
}

// Converts radians to degrees
func radToDeg(angle float32) float32 {
	return angle * 180 / float32(math.Pi)
}
