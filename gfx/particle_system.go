package gfx

import (
	"math"

	//"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	MAX_PARTICLES = math.MaxInt32 / 4
)

type (
	// Represents a single particle.
	Particle struct {
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
		color                  Color
		quadIndex              int
	}
	ParticleSystem struct {
		particles                 []*Particle
		texture                   iTexture
		active                    bool
		insertMode                ParticleInsertion
		maxParticles              uint32
		activeParticles           uint32
		emissionRate              float32
		emitCounter               float32
		position                  mgl32.Vec2
		prevPosition              mgl32.Vec2
		areaSpreadDistribution    ParticleDistribution
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
		colors                    []Color
		quads                     []Quad
		relativeRotation          bool

		// array of transformed vertex data for all particles, for drawing
		particleVerts []mgl32.Vec2
		// Vertex index buffer.
		quadIndices *quadIndices
	}
)

func calculate_variation(inner, outer, v float32) float32 {
	low := inner - (outer/2.0)*v
	high := inner + (outer/2.0)*v
	r := rng.Rand()
	return low*(1-r) + high*r
}

func NewParticleSystem(texture iTexture, size uint32) *ParticleSystem {
	if size == 0 || size > MAX_PARTICLES {
		panic("Invalid ParticleSystem size.")
	}
	return &ParticleSystem{
		texture:                texture,
		active:                 true,
		insertMode:             INSERT_MODE_TOP,
		areaSpreadDistribution: DISTRIBUTION_NONE,
		lifetime:               -1,
		offset:                 mgl32.Vec2{float32(texture.GetWidth()) * 0.5, float32(texture.GetHeight()) * 0.5},
		defaultOffset:          true,
		colors:                 []Color{Color{1.0, 1.0, 1.0, 1.0}},
		sizes:                  []float32{1.0},
		particles:              make([]*Particle, size),
		maxParticles:           size,
	}
}

func (system *ParticleSystem) resetOffset() {
	if len(system.quads) == 0 {
		system.offset = mgl32.Vec2{float32(system.texture.GetWidth()) * 0.5, float32(system.texture.GetHeight()) * 0.5}
	} else {
		x, y, _, _ := system.quads[0].GetViewport()
		system.offset = mgl32.Vec2{x * 0.5, y * 0.5}
	}
}

func (system *ParticleSystem) createBuffers(size int) {
	system.particles = make([]*Particle, size)
	system.particleVerts = make([]mgl32.Vec2, size)
}

func (system *ParticleSystem) setBufferSize(size uint32) {
	if size == 0 || size > MAX_PARTICLES {
		panic("Invalid buffer size")
	}
	system.maxParticles = size
	system.createBuffers(int(size))
	system.quadIndices = newQuadIndices(int(size))
	system.reset()
}

func (system *ParticleSystem) getBufferSize() uint32 {
	return system.maxParticles
}

func (system *ParticleSystem) addParticle(t float32) {
	if system.isFull() {
		return
	}
	// Gets a free particle and updates the allocation pointer.
	p := system.initParticle(t)
	switch system.insertMode {
	case INSERT_MODE_TOP:
		system.particles = append([]*Particle{p}, system.particles...)
	case INSERT_MODE_BOTTOM:
		system.particles = append(system.particles, p)
	case INSERT_MODE_RANDOM:
		i := int(rng.RandMax(float32(system.maxParticles - 1)))
		system.particles = append(system.particles[:i], append([]*Particle{p}, system.particles[i:]...)...)
	}
	system.activeParticles++
}

func (system *ParticleSystem) initParticle(t float32) *Particle {
	pos := system.prevPosition.Add(system.position.Sub(system.prevPosition).Mul(t))
	p := &Particle{
		// Linearly interpolate between the previous and current emitter position.
		position: pos,
	}

	min := system.particleLifeMin
	max := system.particleLifeMax
	if min == max {
		p.life = min
	} else {
		p.life = rng.RandRange(min, max)
	}
	p.lifetime = p.life

	switch system.areaSpreadDistribution {
	case DISTRIBUTION_UNIFORM:
		p.position[0] += rng.RandRange(-system.areaSpread[0], system.areaSpread[0])
		p.position[1] += rng.RandRange(-system.areaSpread[1], system.areaSpread[1])
	case DISTRIBUTION_NORMAL:
		p.position[0] += rng.RandomNormal(system.areaSpread[0])
		p.position[1] += rng.RandomNormal(system.areaSpread[1])
	case DISTRIBUTION_NONE:
		//discard
	}

	p.origin = pos

	speed := rng.RandRange(system.speedMin, system.speedMax)
	dir := float64(rng.RandRange(system.direction-system.spread/2.0, system.direction+system.spread/2.0))
	p.velocity = mgl32.Vec2{float32(math.Cos(dir)), float32(math.Sin(dir))}.Mul(speed)

	p.linearAcceleration[0] = rng.RandRange(system.linearAccelerationMin[0], system.linearAccelerationMax[0])
	p.linearAcceleration[1] = rng.RandRange(system.linearAccelerationMin[1], system.linearAccelerationMax[1])

	p.radialAcceleration = rng.RandRange(system.radialAccelerationMin, system.radialAccelerationMax)
	p.tangentialAcceleration = rng.RandRange(system.tangentialAccelerationMin, system.tangentialAccelerationMax)
	p.linearDamping = rng.RandRange(system.linearDampingMin, system.linearDampingMax)
	p.sizeOffset = rng.RandMax(system.sizeVariation) // time offset for size change
	p.sizeIntervalSize = (1.0 - rng.RandMax(system.sizeVariation)) - p.sizeOffset
	p.size = system.sizes[int(p.sizeOffset-0.5)*(len(system.sizes)-1)]
	p.spinStart = calculate_variation(system.spinStart, system.spinEnd, system.spinVariation)
	p.spinEnd = calculate_variation(system.spinEnd, system.spinStart, system.spinVariation)
	p.rotation = rng.RandRange(system.rotationMin, system.rotationMax)
	p.angle = p.rotation

	if system.relativeRotation {
		p.angle += float32(math.Atan2(float64(p.velocity[1]), float64(p.velocity[0])))
	}
	p.color = system.colors[0]
	p.quadIndex = 0

	return p
}

func (system *ParticleSystem) removeParticle(p *Particle) {
	found := -1
	for i, particle := range system.particles {
		if particle == p {
			found = i
			break
		}
	}
	system.particles = append(system.particles[:found], system.particles[found+1:]...)
}

func (system *ParticleSystem) setTexture(tex iTexture) {
	system.texture = tex
	if system.defaultOffset {
		system.resetOffset()
	}
}

func (system *ParticleSystem) getTexture() iTexture {
	return system.texture
}

func (system *ParticleSystem) setInsertMode(mode ParticleInsertion) {
	system.insertMode = mode
}

func (system *ParticleSystem) getInsertMode() ParticleInsertion {
	return system.insertMode
}

func (system *ParticleSystem) setEmissionRate(rate float32) {
	if rate < 0.0 {
		panic("Invalid emission rate")
	}
	system.emissionRate = rate
}

func (system *ParticleSystem) getEmissionRate() float32 {
	return system.emissionRate
}

func (system *ParticleSystem) setEmitterLifetime(life float32) {
	system.life = life
	system.lifetime = life
}

func (system *ParticleSystem) getEmitterLifetime() float32 {
	return system.lifetime
}

func (system *ParticleSystem) setParticleLifetime(min, max float32) {
	system.particleLifeMin = min
	if max == 0 {
		system.particleLifeMax = min
	} else {
		system.particleLifeMax = max
	}
}

func (system *ParticleSystem) getParticleLifetime() (float32, float32) {
	return system.particleLifeMin, system.particleLifeMax
}

func (system *ParticleSystem) setPosition(x, y float32) {
	system.position = mgl32.Vec2{x, y}
	system.prevPosition = system.position
}

func (system *ParticleSystem) getPosition() (float32, float32) {
	return system.position[0], system.position[1]
}

func (system *ParticleSystem) moveTo(x, y float32) {
	system.position = mgl32.Vec2{x, y}
}

func (system *ParticleSystem) setAreaSpread(distribution ParticleDistribution, x, y float32) {
	system.areaSpread = mgl32.Vec2{x, y}
	system.areaSpreadDistribution = distribution
}

func (system *ParticleSystem) getAreaSpreadDistribution() ParticleDistribution {
	return system.areaSpreadDistribution
}

func (system *ParticleSystem) getAreaSpreadParameters() (float32, float32) {
	return system.areaSpread[0], system.areaSpread[1]
}

func (system *ParticleSystem) setDirection(direction float32) {
	system.direction = direction
}

func (system *ParticleSystem) getDirection() float32 {
	return system.direction
}

func (system *ParticleSystem) setSpread(spread float32) {
	system.spread = spread
}

func (system *ParticleSystem) getSpread() float32 {
	return system.spread
}

func (system *ParticleSystem) setSpeed(min, max float32) {
	system.speedMin = min
	system.speedMax = max
}

func (system *ParticleSystem) getSpeed() (float32, float32) {
	return system.speedMin, system.speedMax
}

func (system *ParticleSystem) setLinearAcceleration(xmin, ymin, xmax, ymax float32) {
	system.linearAccelerationMin = mgl32.Vec2{xmin, ymin}
	system.linearAccelerationMax = mgl32.Vec2{xmax, ymax}
}

func (system *ParticleSystem) getLinearAcceleration() (xmin, ymin, xmax, ymax float32) {
	return system.linearAccelerationMin[0], system.linearAccelerationMin[1], system.linearAccelerationMax[0], system.linearAccelerationMax[1]
}

func (system *ParticleSystem) setRadialAcceleration(min, max float32) {
	system.radialAccelerationMin = min
	system.radialAccelerationMax = max
}

func (system *ParticleSystem) getRadialAcceleration() (min, max float32) {
	return system.radialAccelerationMin, system.radialAccelerationMax
}

func (system *ParticleSystem) setTangentialAcceleration(min, max float32) {
	system.tangentialAccelerationMin = min
	system.tangentialAccelerationMax = max
}

func (system *ParticleSystem) getTangentialAcceleration() (min, max float32) {
	return system.tangentialAccelerationMin, system.tangentialAccelerationMax
}

func (system *ParticleSystem) setLinearDamping(min, max float32) {
	system.linearDampingMin = min
	system.linearDampingMax = max
}

func (system *ParticleSystem) getLinearDamping() (min, max float32) {
	return system.linearDampingMin, system.linearDampingMax
}

func (system *ParticleSystem) setSize(size float32) {
	system.sizes = []float32{size}
}

func (system *ParticleSystem) setSizes(newSizes []float32) {
	system.sizes = newSizes
}

func (system *ParticleSystem) getSizes() []float32 {
	return system.sizes
}

func (system *ParticleSystem) setSizeVariation(variation float32) {
	system.sizeVariation = variation
}

func (system *ParticleSystem) getSizeVariation() float32 {
	return system.sizeVariation
}

func (system *ParticleSystem) setRotation(min, max float32) {
	system.rotationMin = min
	system.rotationMax = max
}

func (system *ParticleSystem) getRotation() (min, max float32) {
	return system.rotationMin, system.rotationMax
}

func (system *ParticleSystem) setSpin(start, end float32) {
	system.spinStart = start
	system.spinEnd = end
}

func (system *ParticleSystem) getSpin() (start, end float32) {
	return system.spinStart, system.spinEnd
}

func (system *ParticleSystem) setSpinVariation(variation float32) {
	system.spinVariation = variation
}

func (system *ParticleSystem) getSpinVariation() float32 {
	return system.spinVariation
}

func (system *ParticleSystem) setOffset(x, y float32) {
	system.offset = mgl32.Vec2{x, y}
	system.defaultOffset = false
}

func (system *ParticleSystem) getOffset() (x, y float32) {
	return system.offset[0], system.offset[1]
}

func (system *ParticleSystem) setColor(newColors []Color) {
	system.colors = newColors
}

func (system *ParticleSystem) getColor() []Color {
	return system.colors
}

func (system *ParticleSystem) setQuads(newQuads []Quad) {
	system.quads = newQuads
	if system.defaultOffset {
		system.resetOffset()
	}
}

func (system *ParticleSystem) clearQuads() {
	system.quads = []Quad{}
}

func (system *ParticleSystem) getQuads() []Quad {
	return system.quads
}

func (system *ParticleSystem) setRelativeRotation(enable bool) {
	system.relativeRotation = enable
}

func (system *ParticleSystem) hasRelativeRotation() bool {
	return system.relativeRotation
}

func (system *ParticleSystem) getCount() uint32 {
	return system.activeParticles
}

func (system *ParticleSystem) start() {
	system.active = true
}

func (system *ParticleSystem) stop() {
	system.active = false
	system.life = system.lifetime
	system.emitCounter = 0
}

func (system *ParticleSystem) pause() {
	system.active = false
}

func (system *ParticleSystem) reset() {
	system.activeParticles = 0
	system.life = system.lifetime
	system.emitCounter = 0
}

func (system *ParticleSystem) emit(num uint32) {
	if !system.active {
		return
	}

	num = uint32(math.Min(float64(num), float64(system.maxParticles-system.activeParticles)))

	for ; num > 0; num-- {
		system.addParticle(1.0)
	}
}

func (system *ParticleSystem) isActive() bool {
	return system.active
}

func (system *ParticleSystem) isPaused() bool {
	return !system.active && system.life < system.lifetime
}

func (system *ParticleSystem) isStopped() bool {
	return !system.active && system.life >= system.lifetime
}

func (system *ParticleSystem) isEmpty() bool {
	return system.activeParticles == 0
}

func (system *ParticleSystem) isFull() bool {
	return system.activeParticles == system.maxParticles
}

func (system *ParticleSystem) update(dt float32) {
	if dt == 0.0 || system.activeParticles == 0 {
		return
	}

	// Traverse all particles and update.
	for _, p := range system.particles {
		// Decrease lifespan.
		p.life -= dt

		if p.life <= 0 {
			system.removeParticle(p)
		} else {
			ppos := p.position
			// Get vector from particle center to particle.
			radial := ppos.Sub(p.origin).Normalize()
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
				p.angle += float32(math.Atan2(float64(p.velocity[1]), float64(p.velocity[0])))
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
			p.color = *system.colors[i].Mul(1.0 - s).Add(system.colors[k].Mul(s))

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
			system.stop()
		}
	}

	system.prevPosition = system.position
}

func (system *ParticleSystem) Draw(args ...float32) {
	pCount := system.getCount()

	if pCount == 0 || system.particleVerts == nil {
		return
	}

	//textureVerts := system.texture.getVertices()
	//pVerts := system.particleVerts
	//useQuads := len(system.quads) == 0

	//// set the vertex data for each particle (transformation, texcoords, color)
	//for _, p := range system.particles {
	//if system.useQuads {
	//textureVerts = system.quads[p.quadIndex].getVertices()
	//}

	//// particle vertices are image vertices transformed by particle info
	//t := generateModelMatFromArgs([]float32{p.position[0], p.position[1], p.angle, p.size, p.size, system.offset[0], system.offset[1], 0.0, 0.0})
	//pVerts = matTransform(t, textureVerts, 4)

	//// set the texture coordinate and color data for particle vertices
	//for v := 0; v < 4; v++ {
	//pVerts[v].s = textureVerts[v].s
	//pVerts[v].t = textureVerts[v].t

	//// Particle colors are stored as floats (0-1) but vertex colors are
	//// unsigned bytes (0-255).
	//pVerts[v].r = p.color.r * 255
	//pVerts[v].g = p.color.g * 255
	//pVerts[v].b = p.color.b * 255
	//pVerts[v].a = p.color.a * 255
	//}

	//pVerts += 4
	//}

	//prepareDraw(nil)
	//bindTexture(system.texture.getHandle())

	//gl.EnableVertexAttribArray(ATTRIB_POS)
	//gl.EnableVertexAttribArray(ATTRIB_TEXCOORD)
	//gl.EnableVertexAttribArray(ATTRIB_COLOR)

	//gl.VertexAttribPointer(ATTRIB_COLOR, 4, gl.UNSIGNED_BYTE, false, sizeof(Vertex), &particleVerts[0].r)
	//gl.VertexAttribPointer(ATTRIB_POS, 2, gl.FLOAT, false, sizeof(Vertex), &particleVerts[0].x)
	//gl.VertexAttribPointer(ATTRIB_TEXCOORD, 2, gl.FLOAT, false, sizeof(Vertex), &particleVerts[0].s)

	//// We use a client-side index array instead of an Index Buffers, because
	//// at least one graphics driver (the one for Kepler nvidia GPUs in OS X
	//// 10.11) fails to render geometry if an index buffer is used with
	//// client-side vertex arrays.
	//gl.DrawElements(gl.TRIANGLES, pCount*6, gl.FLOAT, quadIndices.indices)

	//gl.DisableVertexAttribArray(ATTRIB_COLOR)
	//gl.DisableVertexAttribArray(ATTRIB_TEXCOORD)
	//gl.DisableVertexAttribArray(ATTRIB_POS)
}
