package gfx

import (
	"math"
	"math/rand"
	"time"
)

var rng = newRandomGenerator()

type randomGenerator struct {
	rng              *rand.Rand
	seed             int64
	lastRandomnormal float64
}

func newRandomGenerator() *randomGenerator {
	seed := time.Now().UnixNano()

	return &randomGenerator{
		rng:  rand.New(rand.NewSource(seed)),
		seed: seed,
	}
}

func (generator *randomGenerator) rand() float32 {
	return generator.rng.Float32()
}

func (generator *randomGenerator) randMax(max float32) float32 {
	return generator.rand() * max
}

func (generator *randomGenerator) randRange(min, max float32) float32 {
	return generator.rand()*(max-min) + min
}

func (generator *randomGenerator) randomNormal(stddev float32) float32 {
	// use cached number if possible
	if generator.lastRandomnormal != math.Inf(1) {
		r := float32(generator.lastRandomnormal)
		generator.lastRandomnormal = math.Inf(1)
		return r * stddev
	}
	r := math.Sqrt(-2.0 * math.Log(1.0-float64(generator.rand())))
	phi := 2.0 * math.Pi * (1.0 - float64(generator.rand()))
	generator.lastRandomnormal = r * math.Cos(phi)
	return float32(r * math.Sin(phi) * float64(stddev))
}

func (generator *randomGenerator) setSeed(seed int64) {
	generator.seed = seed
	generator.rng.Seed(seed)
}

func (generator *randomGenerator) getSeed() int64 {
	return generator.seed
}

// Rand will return a random number between 0 and 1
func Rand() float32 { return rng.rand() }

// RandMax will return a random number between 0 and max
func RandMax(max float32) float32 { return rng.randMax(max) }

// RandRange will return a random number between min and max
func RandRange(min, max float32) float32 { return rng.randRange(min, max) }

// RandomNormal does a Box–Muller transform with the standard deviation provided.
// https://en.wikipedia.org/wiki/Box%E2%80%93Muller_transform
func RandomNormal(stddev float32) float32 { return rng.randomNormal(stddev) }

// SetSeed will set the seed for random number generation
func SetSeed(seed int64) { rng.setSeed(seed) }

// GetSeed will return the seed currently used for random number generation
func GetSeed() int64 { return rng.getSeed() }
