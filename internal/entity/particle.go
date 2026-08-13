package entity

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Particle represents a single visual FX particle.
type Particle struct {
	X, Y     float64
	VX, VY   float64
	Life     float64
	MaxLife  float64
	Size     float32
	Color    color.RGBA
	Active   bool
}

// ParticlePool manages pre-allocated particles.
type ParticlePool struct {
	particles []Particle
}

// NewParticlePool creates a particle pool with a fixed capacity.
func NewParticlePool(capacity int) *ParticlePool {
	return &ParticlePool{
		particles: make([]Particle, capacity),
	}
}

// Emit spawns a new particle into the pool.
func (p *ParticlePool) Emit(x, y, vx, vy float64, size float32, life float64, col color.RGBA) {
	for i := range p.particles {
		if !p.particles[i].Active {
			p.particles[i] = Particle{
				X:       x,
				Y:       y,
				VX:      vx,
				VY:      vy,
				Size:    size,
				Life:    life,
				MaxLife: life,
				Color:   col,
				Active:  true,
			}
			return
		}
	}
}

// EmitBurst spawns multiple particles in a burst pattern.
func (p *ParticlePool) EmitBurst(x, y float64, count int, col color.RGBA) {
	for i := 0; i < count; i++ {
		vx := (rand.Float64()*2.0 - 1.0) * 3.0
		vy := (rand.Float64()*2.0 - 1.0) * 3.0
		p.Emit(x, y, vx, vy, float32(rand.Float64()*3.0+2.0), 0.4+rand.Float64()*0.3, col)
	}
}

// Update advances all particles in the pool.
func (p *ParticlePool) Update() {
	for i := range p.particles {
		if p.particles[i].Active {
			p.particles[i].X += p.particles[i].VX
			p.particles[i].Y += p.particles[i].VY
			p.particles[i].Life -= 0.016
			if p.particles[i].Life <= 0 {
				p.particles[i].Active = false
			}
		}
	}
}

// Draw renders active particles onto the screen.
func (p *ParticlePool) Draw(screen *ebiten.Image, cameraX float64) {
	for i := range p.particles {
		pt := &p.particles[i]
		if pt.Active {
			screenX := float32(pt.X - cameraX)
			alpha := float32(pt.Life / pt.MaxLife)
			col := color.RGBA{
				R: pt.Color.R,
				G: pt.Color.G,
				B: pt.Color.B,
				A: uint8(float32(pt.Color.A) * alpha),
			}
			vector.DrawFilledCircle(screen, screenX, float32(pt.Y), pt.Size, col, false)
		}
	}
}
