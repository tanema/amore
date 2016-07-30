package rand

import (
	"math"
	"math/rand"
	"time"
)

var (
	seed              = time.Now().UnixNano()
	source            = rand.NewSource(seed)
	rng               = rand.New(source)
	last_randomnormal float64
)

func Rand() float32 {
	return rng.Float32()
}

func RandMax(max float32) float32 {
	return Rand() * max
}

func RandRange(min, max float32) float32 {
	return Rand()*(max-min) + min
}

// Box–Muller transform
func RandomNormal(stddev float32) float32 {
	// use cached number if possible
	if last_randomnormal != math.Inf(1) {
		r := float32(last_randomnormal)
		last_randomnormal = math.Inf(1)
		return r * stddev
	}
	r := math.Sqrt(-2.0 * math.Log(1.0-float64(Rand())))
	phi := 2.0 * math.Pi * (1.0 - float64(Rand()))
	last_randomnormal = r * math.Cos(phi)
	return float32(r * math.Sin(phi) * float64(stddev))
}

func SetSeed(s int64) {
	seed = s
	rng.Seed(seed)
}

func GetSeed() int64 {
	return seed
}
