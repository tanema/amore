package gfx

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/tanema/amore/gfx/gl"
	"github.com/tanema/amore/mth"
	"github.com/tanema/amore/mth/rand"
)

type (
	// particle represents a single particle in the system.
	particle struct {
		lifetime               float32
		life                   float32
		position               mgl32.Vec2
		origin                 mgl32.Vec2
		velocity               mgl32.Vec2
		linearAcceleration     mgl32.Vec2
		radialAcceleration     float32
		tangentialAcceleration float32
		linearDamping          float32
		size                   float32
		sizeOffset             float32
		sizeIntervalSize       float32
		rotation               float32 // Amount of rotation applied to the final angle.
		angle                  float32
		spinStart              float32
		spinEnd                float32
		color                  *Color
		quadIndex              int
	}
	// ParticleSystem generates and emits particles at set speeds and rotations.
	ParticleSystem struct {
		particles                 []*particle
		texture                   iTexture
		insertMode                ParticleInsertion
		areaSpreadDistribution    ParticleDistribution
		active                    bool
		maxParticles              int
		emissionRate              float32
		emitCounter               float32
		position                  mgl32.Vec2
		prevPosition              mgl32.Vec2
		areaSpread                mgl32.Vec2
		lifetime                  float32
		life                      float32
		particleLifeMin           float32
		particleLifeMax           float32
		direction                 float32
		spread                    float32
		speedMin                  float32
		speedMax                  float32
		linearAccelerationMin     mgl32.Vec2
		linearAccelerationMax     mgl32.Vec2
		radialAccelerationMin     float32
		radialAccelerationMax     float32
		tangentialAccelerationMin float32
		tangentialAccelerationMax float32
		linearDampingMin          float32
		linearDampingMax          float32
		sizes                     []float32
		sizeVariation             float32
		rotationMin               float32
		rotationMax               float32
		spinStart                 float32
		spinEnd                   float32
		spinVariation             float32
		offset                    mgl32.Vec2
		defaultOffset             bool
		colors                    []*Color
		quads                     []Quad
		relativeRotation          bool
		quadIndices               *quadIndices
	}
)

// Release will release all the gl objects associates with the system and clean
// up the memory
func (system *ParticleSystem) Release() {
	system.quadIndices.Release()
}

// calculate_variation is used to calculate the variation in starting spin on
// each particle
func calculate_variation(inner, outer, v float32) float32 {
	low := inner - (outer/2.0)*v
	high := inner + (outer/2.0)*v
	r := rand.Rand()
	return low*(1-r) + high*r
}

// NewParticleSystem will create a new system where the texture given (generally
// an image) will be how each particle is represented and the size is the maximum
// number of particles displayed at a time.
func NewParticleSystem(texture iTexture, size int) *ParticleSystem {
	if size == 0 || size > MAX_PARTICLES {
		panic("Invalid ParticleSystem size.")
	}
	new_ps := &ParticleSystem{
		texture:                texture,
		active:                 true,
		insertMode:             INSERT_MODE_TOP,
		areaSpreadDistribution: DISTRIBUTION_NONE,
		lifetime:               -1,
		offset:                 mgl32.Vec2{float32(texture.GetWidth()) * 0.5, float32(texture.GetHeight()) * 0.5},
		defaultOffset:          true,
		colors:                 []*Color{&Color{1.0, 1.0, 1.0, 1.0}},
		sizes:                  []float32{1.0},
		particles:              []*particle{},
		maxParticles:           size,
		quadIndices:            newQuadIndices(size),
	}

	return new_ps
}

// resetOffset is called after that texture is set or a quad is set so that each
// particle offset is coming from the proper position.
func (system *ParticleSystem) resetOffset() {
	if len(system.quads) == 0 {
		system.offset = mgl32.Vec2{float32(system.texture.GetWidth()) * 0.5, float32(system.texture.GetHeight()) * 0.5}
	} else {
		x, y, _, _ := system.quads[0].GetViewport()
		system.offset = mgl32.Vec2{float32(x) * 0.5, float32(y) * 0.5}
	}
}

// SetBufferSize resizes the maximum amount of particles the system can create
// simultaneously.
func (system *ParticleSystem) SetBufferSize(size int) {
	if size == 0 || size > MAX_PARTICLES {
		panic("Invalid buffer size")
	}
	if len(system.particles) > int(size) {
		system.particles = system.particles[:int(size)]
	}
	system.maxParticles = size
	system.quadIndices.Release()
	system.quadIndices = newQuadIndices(size)
	system.life = system.lifetime
	system.emitCounter = 0
}

// GetBufferSize returns the maximum number of particles the system can create
// simultaneously.
func (system *ParticleSystem) GetBufferSize() int {
	return system.maxParticles
}

// addParticle handles when the system creates a new particle and it's insertion
// into the particle slice.
func (system *ParticleSystem) addParticle(t float32) {
	if system.IsFull() {
		return
	}
	// Gets a free particle and updates the allocation pointer.
	p := system.initParticle(t)
	switch system.insertMode {
	case INSERT_MODE_TOP:
		system.particles = append([]*particle{p}, system.particles...)
	case INSERT_MODE_BOTTOM:
		system.particles = append(system.particles, p)
	case INSERT_MODE_RANDOM:
		i := int(rand.RandMax(float32(system.maxParticles - 1)))
		system.particles = append(system.particles[:i], append([]*particle{p}, system.particles[i:]...)...)
	}
}

// initParticle creates a new particle to be inserted into the system with all the
// default values from the system.
func (system *ParticleSystem) initParticle(t float32) *particle {
	// Linearly interpolate between the previous and current emitter position.
	pos := system.prevPosition.Add(system.position.Sub(system.prevPosition).Mul(t))

	p := &particle{
		position: pos,
	}

	if system.particleLifeMin == system.particleLifeMax {
		p.life = system.particleLifeMin
	} else {
		p.life = rand.RandRange(system.particleLifeMin, system.particleLifeMax)
	}
	p.lifetime = p.life

	switch system.areaSpreadDistribution {
	case DISTRIBUTION_UNIFORM:
		p.position[0] += rand.RandRange(-system.areaSpread[0], system.areaSpread[0])
		p.position[1] += rand.RandRange(-system.areaSpread[1], system.areaSpread[1])
	case DISTRIBUTION_NORMAL:
		p.position[0] += rand.RandomNormal(system.areaSpread[0])
		p.position[1] += rand.RandomNormal(system.areaSpread[1])
	case DISTRIBUTION_NONE:
		//done
	}

	p.origin = pos

	speed := rand.RandRange(system.speedMin, system.speedMax)
	dir := rand.RandRange(system.direction-system.spread/2.0, system.direction+system.spread/2.0)
	p.velocity = mgl32.Vec2{mth.Cos(dir), mth.Sin(dir)}.Mul(speed)

	p.linearAcceleration[0] = rand.RandRange(system.linearAccelerationMin[0], system.linearAccelerationMax[0])
	p.linearAcceleration[1] = rand.RandRange(system.linearAccelerationMin[1], system.linearAccelerationMax[1])

	p.radialAcceleration = rand.RandRange(system.radialAccelerationMin, system.radialAccelerationMax)
	p.tangentialAcceleration = rand.RandRange(system.tangentialAccelerationMin, system.tangentialAccelerationMax)
	p.linearDamping = rand.RandRange(system.linearDampingMin, system.linearDampingMax)
	p.sizeOffset = rand.RandMax(system.sizeVariation) // time offset for size change
	p.sizeIntervalSize = (1.0 - rand.RandMax(system.sizeVariation)) - p.sizeOffset
	p.size = system.sizes[int(p.sizeOffset-0.5)*(len(system.sizes)-1)]
	p.spinStart = calculate_variation(system.spinStart, system.spinEnd, system.spinVariation)
	p.spinEnd = calculate_variation(system.spinEnd, system.spinStart, system.spinVariation)
	p.rotation = rand.RandRange(system.rotationMin, system.rotationMax)
	p.angle = p.rotation

	if system.relativeRotation {
		p.angle += mth.Atan2(p.velocity[1], p.velocity[0])
	}
	p.color = system.colors[0]
	p.quadIndex = 0

	return p
}

// SetTexture will change the texture for the particles being inserted into the system.
func (system *ParticleSystem) SetTexture(tex iTexture) {
	system.texture = tex
	if system.defaultOffset {
		system.resetOffset()
	}
}

// GetTexture will return the texture currently being used by this particle system.
func (system *ParticleSystem) GetTexture() iTexture {
	return system.texture
}

// SetInsertMode will set the ParticleInsertion on the system and if the particles
// will appear on top, bottom, or randomly in the stack.
func (system *ParticleSystem) SetInsertMode(mode ParticleInsertion) {
	system.insertMode = mode
}

// GetInsertMode will return the current ParticleInsertion of the system.
func (system *ParticleSystem) GetInsertMode() ParticleInsertion {
	return system.insertMode
}

// SetEmissionRate will set the rate in which particles are generated by the system.
// you will the system seems to pause if your rate is high and your buffer size is
// low because it will reach the max size very quickly.
func (system *ParticleSystem) SetEmissionRate(rate float32) {
	if rate < 0.0 {
		panic("Invalid emission rate")
	}
	system.emissionRate = rate
}

// GetEmissionRate will return the rate in which particles are generated by the system.
func (system *ParticleSystem) GetEmissionRate() float32 {
	return system.emissionRate
}

// SetEmitterLifetime will set how long the system should continue to emit
// particles in seconds. If -1 then it emits particles forever.
func (system *ParticleSystem) SetEmitterLifetime(life float32) {
	system.life = life
	system.lifetime = life
}

// GetEmitterLifetime will return the lifetime of the system in seconds. Default
// lifetime on a system is -1
func (system *ParticleSystem) GetEmitterLifetime() float32 {
	return system.lifetime
}

// SetParticleLifetime will set how long in seconds each particle is to remain on
// screen before being destroys and a new particle emitted.
func (system *ParticleSystem) SetParticleLifetime(min, max float32) {
	system.particleLifeMin = min
	system.particleLifeMax = max
}

// GetParticleLifetime will return the range, min, max, life time in seconds that
// each particle will have when they are emitted.
func (system *ParticleSystem) GetParticleLifetime() (float32, float32) {
	return system.particleLifeMin, system.particleLifeMax
}

// SetPosition sets the position the system will be rendered at. By using MoveTo
// the system can take into account the movement of the system when emitting
// particles. This is why you can use the system position insteam of drawing
// at this position.
func (system *ParticleSystem) SetPosition(x, y float32) {
	system.position = mgl32.Vec2{x, y}
	system.prevPosition = system.position
}

// GetPosition returns the position of the system.
func (system *ParticleSystem) GetPosition() (float32, float32) {
	return system.position[0], system.position[1]
}

// MoveTo moves the position of the emitter. This results in smoother particle
// spawning behaviour than if SetPosition is used every frame.
func (system *ParticleSystem) MoveTo(x, y float32) {
	system.position = mgl32.Vec2{x, y}
}

// SetAreaSpread sets area-based spawn parameters for the particles. Newly created
// particles will spawn in an area around the emitter based on the parameters to
// this function. x and y are the maximum spawn distance from the center of the
// emitter
func (system *ParticleSystem) SetAreaSpread(distribution ParticleDistribution, x, y float32) {
	system.areaSpread = mgl32.Vec2{x, y}
	system.areaSpreadDistribution = distribution
}

// GetAreaSpreadDistrobution will return the particle distrobution of the system.
func (system *ParticleSystem) GetAreaSpreadDistribution() ParticleDistribution {
	return system.areaSpreadDistribution
}

// GetAreaSpreadParameters will return the maximum spawn distance from the emitter
// for uniform distribution, or the standard deviation for normal distribution.
func (system *ParticleSystem) GetAreaSpreadParameters() (float32, float32) {
	return system.areaSpread[0], system.areaSpread[1]
}

// SetDirection sets the direction in radians the particles will be emitted in.
func (system *ParticleSystem) SetDirection(direction float32) {
	system.direction = direction
}

// GetDirection will return the direction in radians that the particles will be
// emitted in.
func (system *ParticleSystem) GetDirection() float32 {
	return system.direction
}

// SetSpread sets the spread in radian from the direction being emitted.
func (system *ParticleSystem) SetSpread(spread float32) {
	system.spread = spread
}

// GetSpread returns the spread in radians that the particle are emitted in.
func (system *ParticleSystem) GetSpread() float32 {
	return system.spread
}

// SetSpeed will set a range of speeds that the particles will be emitted in.
func (system *ParticleSystem) SetSpeed(min, max float32) {
	system.speedMin = min
	system.speedMax = max
}

// GetSpeed will return the min, max range of speeds that the particles will be
// emitted in.
func (system *ParticleSystem) GetSpeed() (float32, float32) {
	return system.speedMin, system.speedMax
}

// SetLinearAcceleration Sets a linear acceleration along the x and y axes, for
//particles.  Every particle created will accelerate along the x and y axes between
// xmin,ymin and xmax,ymax.
func (system *ParticleSystem) SetLinearAcceleration(xmin, ymin, xmax, ymax float32) {
	system.linearAccelerationMin = mgl32.Vec2{xmin, ymin}
	system.linearAccelerationMax = mgl32.Vec2{xmax, ymax}
}

// GetLinearAcceleration will return the min and max ranges od linear acceleration
// on each particle.
func (system *ParticleSystem) GetLinearAcceleration() (xmin, ymin, xmax, ymax float32) {
	return system.linearAccelerationMin[0], system.linearAccelerationMin[1], system.linearAccelerationMax[0], system.linearAccelerationMax[1]
}

// SetRadialAcceleration Sets a radial acceleration along the x and y axes, for
//particles.  Every particle created will accelerate along the x and y axes between
// xmin,ymin and xmax,ymax.
func (system *ParticleSystem) SetRadialAcceleration(min, max float32) {
	system.radialAccelerationMin = min
	system.radialAccelerationMax = max
}

// GetRadialAcceleration will return the range min, max of the radial acceleration
// on each particle.
func (system *ParticleSystem) GetRadialAcceleration() (min, max float32) {
	return system.radialAccelerationMin, system.radialAccelerationMax
}

// SetTangentialAcceleration sets a tangential acceleration along the x and y axes, for
// particles.  Every particle created will accelerate along the x and y axes between
// xmin,ymin and xmax,ymax.
func (system *ParticleSystem) SetTangentialAcceleration(min, max float32) {
	system.tangentialAccelerationMin = min
	system.tangentialAccelerationMax = max
}

// GetTangentialAcceleration will return the range min, max of the tangential acceleration
// on each particle.
func (system *ParticleSystem) GetTangentialAcceleration() (min, max float32) {
	return system.tangentialAccelerationMin, system.tangentialAccelerationMax
}

// SetLinearDamping sets the range of linear damping (constant deceleration) for particles.
func (system *ParticleSystem) SetLinearDamping(min, max float32) {
	system.linearDampingMin = min
	system.linearDampingMax = max
}

// GetLinearDamping returns the range min, max of linear deceleration for particles.
func (system *ParticleSystem) GetLinearDamping() (min, max float32) {
	return system.linearDampingMin, system.linearDampingMax
}

// SetSize sets the size of the particles in pixels over its lifetime. If only
// one is provided it will stay that size, if multiple are give it will interpolate
// between the sizes over time. A max of 8 sizes
func (system *ParticleSystem) SetSize(size ...float32) {
	if size == nil || len(size) == 0 || len(size) > 8 {
		return
	}
	system.sizes = size
}

// GetSizes will return all the sizes set for each particle.
func (system *ParticleSystem) GetSizes() []float32 {
	return system.sizes
}

// SetSizeVariation sets the amount of size variation when the particle is emitted.
func (system *ParticleSystem) SetSizeVariation(variation float32) {
	system.sizeVariation = variation
}

// GetSizeVariation will return the set variation of sizes.
func (system *ParticleSystem) GetSizeVariation() float32 {
	return system.sizeVariation
}

// SetRotation sets a range min, max initial angle of the texture in radians.
func (system *ParticleSystem) SetRotation(min, max float32) {
	system.rotationMin = min
	system.rotationMax = max
}

// GetRotation will return the rotation range min, max for the texture
func (system *ParticleSystem) GetRotation() (min, max float32) {
	return system.rotationMin, system.rotationMax
}

// SetSpin sets the range min, max spin (radians per second).
func (system *ParticleSystem) SetSpin(start, end float32) {
	system.spinStart = start
	system.spinEnd = end
}

// GetSpin returns the range min, max spin (radians per second).
func (system *ParticleSystem) GetSpin() (start, end float32) {
	return system.spinStart, system.spinEnd
}

// SetSpinVariation sets the amount of spin variation (0 meaning no variation
// and 1 meaning full variation between start and end).
func (system *ParticleSystem) SetSpinVariation(variation float32) {
	system.spinVariation = variation
}

// GetSpinVariation returns the variation in spin of the particles
func (system *ParticleSystem) GetSpinVariation() float32 {
	return system.spinVariation
}

// SetOffset set the offset position which the particle sprite is rotated around.
// If this function is not used, the particles rotate around their center.
func (system *ParticleSystem) SetOffset(x, y float32) {
	system.offset = mgl32.Vec2{x, y}
	system.defaultOffset = false
}

// GetOffset will return the offset that each particle rotates at
func (system *ParticleSystem) GetOffset() (x, y float32) {
	return system.offset[0], system.offset[1]
}

// SetColor sets a series of colors to apply to the particle sprite. The particle
// system will interpolate between each color evenly over the particle's lifetime.
func (system *ParticleSystem) SetColor(newColors ...*Color) {
	if newColors == nil {
		system.colors = []*Color{&Color{1.0, 1.0, 1.0, 1.0}}
	} else {
		system.colors = newColors
	}
}

// GetColor will return all the colors set for each pixel. See Set color as to why
// this is an array
func (system *ParticleSystem) GetColor() []*Color {
	return system.colors
}

// SetQuads sets a series of Quads to use for the particle sprites. Particles will
// choose a Quad from the list based on the particle's current lifetime, allowing
// for the use of animated sprite sheets with ParticleSystems. If no quads are passed
// in it will clear the quads and use the whole texture again
func (system *ParticleSystem) SetQuads(newQuads ...Quad) {
	if newQuads == nil {
		system.quads = []Quad{}
	} else {
		system.quads = newQuads
		if system.defaultOffset {
			system.resetOffset()
		}
	}
}

// GetQuads will return the quads that are to be used by the system.
func (system *ParticleSystem) GetQuads() []Quad {
	return system.quads
}

// SetRelativeRotation sets whether particle angles and rotations are relative
// to their velocities. If enabled, particles are aligned to the angle of their
// velocities and rotate relative to that angle.
func (system *ParticleSystem) SetRelativeRotation(enable bool) {
	system.relativeRotation = enable
}

// GetRelativeRotation returns whether particle angles and rotations are relative
// to their velocities. If enabled, particles are aligned to the angle of their
// velocities and rotate relative to that angle.
func (system *ParticleSystem) HasRelativeRotation() bool {
	return system.relativeRotation
}

// GetCount will return the count of currently active particles
func (system *ParticleSystem) GetCount() int {
	return len(system.particles)
}

// Start will activate a stopped system.
func (system *ParticleSystem) Start() {
	system.active = true
}

// Stop will deactivate a system, stop emitting particles and restart the life of
// the system.
func (system *ParticleSystem) Stop() {
	system.active = false
	system.life = system.lifetime
	system.emitCounter = 0
}

// Pause will deactivate a system and stop emitting particles.
func (system *ParticleSystem) Pause() {
	system.active = false
}

// Emit will emits a burst of particles from the particle emitter.
func (system *ParticleSystem) Emit(num int) {
	if !system.active {
		return
	}

	num = mth.Mini(num, system.maxParticles-len(system.particles))

	for ; num > 0; num-- {
		system.addParticle(1.0)
	}
}

// IsActive will return if the system is activly emitting particles
func (system *ParticleSystem) IsActive() bool {
	return system.active
}

// IsPaused will return true if the system has been paused
func (system *ParticleSystem) IsPaused() bool {
	return !system.active && system.life < system.lifetime
}

// IsStopped will return if the system is not active or has reached the end of its
// life time.
func (system *ParticleSystem) IsStopped() bool {
	return !system.active && system.life >= system.lifetime
}

// IsEmpty will return true if there are no active particles
func (system *ParticleSystem) IsEmpty() bool {
	return len(system.particles) == 0
}

// IsFull will return true if the number of active particles has reached the max
// buffer size
func (system *ParticleSystem) IsFull() bool {
	return len(system.particles) == system.maxParticles
}

// Update should be called from the amore update function with the delta time,
// this updates the position, color and rotation of all the particles.
func (system *ParticleSystem) Update(dt float32) {
	if dt == 0.0 {
		return
	}

	// Traverse all particles and update.
	for i := len(system.particles) - 1; i >= 0; i-- {
		p := system.particles[i]
		// Decrease lifespan.
		p.life -= dt

		if p.life <= 0 {
			//remove particle
			if i == len(system.particles)-1 {
				system.particles = system.particles[:i]
			} else if i > -1 {
				system.particles = append(system.particles[:i], system.particles[i+1:]...)
			}
			continue
		}

		ppos := p.position
		// Get vector from particle center to particle.
		radial := ppos.Sub(p.origin)
		if radial.Len() > 0 {
			radial = radial.Normalize()
		}
		tangential := radial

		// Resize radial acceleration.
		radial = radial.Mul(p.radialAcceleration)

		// Calculate tangential acceleration.
		a := tangential[0]
		tangential[0] = -tangential[1]
		tangential[1] = a

		// Resize tangential.
		tangential = tangential.Mul(p.tangentialAcceleration)

		// Update velocity.
		p.velocity = p.velocity.Add(radial.Add(tangential).Add(p.linearAcceleration).Mul(dt))

		// Apply damping.
		p.velocity = p.velocity.Mul(1.0 / (1.0 + p.linearDamping*dt))

		// Modify position.
		ppos = ppos.Add(p.velocity.Mul(dt))

		p.position = ppos

		t := 1.0 - p.life/p.lifetime

		// Rotate.
		p.rotation += (p.spinStart*(1.0-t) + p.spinEnd*t) * dt

		p.angle = p.rotation
		if system.relativeRotation {
			p.angle += mth.Atan2(p.velocity[1], p.velocity[0])
		}

		// Change size according to given intervals:
		// i = 0       1       2      3          n-1
		//     |-------|-------|------|--- ... ---|
		// t = 0    1/(n-1)        3/(n-1)        1
		//
		// `s' is the interpolation variable scaled to the current
		// interval width, e.g. if n = 5 and t = 0.3, then the current
		// indices are 1,2 and s = 0.3 - 0.25 = 0.05
		s := p.sizeOffset + t*p.sizeIntervalSize // size variation
		s *= float32(len(system.sizes) - 1)      // 0 <= s < sizes.size()
		i := int(s)
		k := i
		if i != len(system.sizes)-1 {
			k += 1 // boundary check (prevents failing on t = 1.0f)
		}
		s -= float32(i) // transpose s to be in interval [0:1]: i <= s < i + 1 ~> 0 <= s < 1
		p.size = system.sizes[i]*(1.0-s) + system.sizes[k]*s

		// Update color according to given intervals (as above)
		s = t * float32(len(system.colors)-1)
		i = int(s)
		k = i
		if i != len(system.colors)-1 {
			k += 1 // boundary check (prevents failing on t = 1.0f)
		}
		s -= float32(i) // 0 <= s <= 1
		p.color = system.colors[i].Mul(1.0 - s).Add(system.colors[k])

		// Update the quad index.
		k = len(system.quads)
		if k > 0 {
			s = t * float32(k) // [0:numquads-1] (clamped below)
			if s > 0.0 {
				i = int(s)
			} else {
				i = 0
			}
			if i < k {
				p.quadIndex = i
			} else {
				p.quadIndex = k - 1
			}
		}
	}

	// Make some more particles.
	if system.active {
		rate := 1.0 / system.emissionRate // the amount of time between each particle emit
		system.emitCounter += dt
		total := system.emitCounter - rate

		for ; system.emitCounter > rate; system.emitCounter -= rate {
			system.addParticle(1.0 - (system.emitCounter-rate)/total)
		}

		system.life -= dt
		if system.lifetime != -1 && system.life < 0 {
			system.Stop()
		}
	}

	system.prevPosition = system.position
}

// Draw satisfies the Drawable interface. Inputs are as follows
// x, y, r, sx, sy, ox, oy, kx, ky
// x, y are position
// r is rotation
// sx, sy is the scale, if sy is not given sy will equal sx
// ox, oy are offset
// kx, ky are the shear. If ky is not given ky will equal kx
func (system *ParticleSystem) Draw(args ...float32) {
	if !system.active || len(system.particles) == 0 {
		return
	}

	useQuads := len(system.quads) > 0
	particleVerts := make([]float32, 32*len(system.particles))
	textureVerts := system.texture.getVerticies()
	for particle_index, p := range system.particles {
		if useQuads {
			textureVerts = system.quads[p.quadIndex].getVertices()
		}

		// particle vertices are image vertices transformed by particle info
		mat := generateModelMatFromArgs([]float32{
			p.position[0], p.position[1],
			p.angle,
			p.size, p.size,
			system.offset[0], system.offset[1],
		})

		pvi := 32 * particle_index
		for i := 0; i < 32; i += 8 {
			j := (i / 2)
			particleVerts[pvi+i+0] = (mat[0] * textureVerts[j+0]) + (mat[4] * textureVerts[j+1]) + mat[12]
			particleVerts[pvi+i+1] = (mat[1] * textureVerts[j+0]) + (mat[5] * textureVerts[j+1]) + mat[13]
			particleVerts[pvi+i+2] = textureVerts[j+2]
			particleVerts[pvi+i+3] = textureVerts[j+3]
			particleVerts[pvi+i+4] = p.color[0]
			particleVerts[pvi+i+5] = p.color[1]
			particleVerts[pvi+i+6] = p.color[2]
			particleVerts[pvi+i+7] = p.color[3]
		}
	}

	prepareDraw(generateModelMatFromArgs(args))
	bindTexture(system.texture.getHandle())

	useVertexAttribArrays(attribflag_pos | attribflag_texcoord | attribflag_color)
	gl.VertexAttribPointer(attrib_pos, 2, gl.FLOAT, false, 8*4, gl.Ptr(particleVerts))
	gl.VertexAttribPointer(attrib_texcoord, 2, gl.FLOAT, false, 8*4, gl.Ptr(&particleVerts[2]))
	gl.VertexAttribPointer(attrib_color, 4, gl.FLOAT, false, 8*4, gl.Ptr(&particleVerts[4]))

	// We use a client-side index array instead of an Index Buffers, because
	// at least one graphics driver (the one for Kepler nvidia GPUs in OS X
	// 10.11) fails to render geometry if an index buffer is used with
	// client-side vertex arrays.
	system.quadIndices.drawElementsLocal(gl.TRIANGLES, 0, len(system.particles))
}
