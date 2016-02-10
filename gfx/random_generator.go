package gfx

import (
	"math"
	"math/rand"
	"time"
)

var rng = newRandomGenerator()

type RandomGenerator struct {
	rng               *rand.Rand
	seed              int64
	last_randomnormal float64
}

func newRandomGenerator() *RandomGenerator {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)

	return &RandomGenerator{
		rng:  rand.New(source),
		seed: seed,
	}
}

func (generator *RandomGenerator) Rand() float32 {
	return generator.rng.Float32()
}

func (generator *RandomGenerator) RandMax(max float32) float32 {
	return generator.Rand() * max
}

func (generator *RandomGenerator) RandRange(min, max float32) float32 {
	return generator.Rand()*(max-min) + min
}

// Box–Muller transform
func (generator *RandomGenerator) RandomNormal(stddev float32) float32 {
	// use cached number if possible
	if generator.last_randomnormal != math.Inf(1) {
		r := float32(generator.last_randomnormal)
		generator.last_randomnormal = math.Inf(1)
		return r * stddev
	}
	r := math.Sqrt(-2.0 * math.Log(1.0-float64(generator.Rand())))
	phi := 2.0 * math.Pi * (1.0 - float64(generator.Rand()))
	generator.last_randomnormal = r * math.Cos(phi)
	return float32(r * math.Sin(phi) * float64(stddev))
}

func (generator *RandomGenerator) SetSeed(seed int64) {
	generator.seed = seed
	generator.rng.Seed(seed)
}

func (generator *RandomGenerator) GetSeed() int64 {
	return generator.seed
}
