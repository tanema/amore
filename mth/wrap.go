package mth

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

func Abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
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
