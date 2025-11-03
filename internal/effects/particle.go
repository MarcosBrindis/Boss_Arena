package effects

import (
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
)

// Particle representa una partícula individual
type Particle struct {
	Position utils.Vector2
	Velocity utils.Vector2
	Color    color.RGBA
	Size     float64
	Lifetime int
	Age      int
	IsActive bool
}

// ParticleSystem maneja un sistema de partículas (THREAD-SAFE)
type ParticleSystem struct {
	particles    []*Particle
	maxParticles int
	rng          *rand.Rand
	mu           sync.Mutex
}

// NewParticleSystem crea un nuevo sistema de partículas
func NewParticleSystem(maxParticles int) *ParticleSystem {
	return &ParticleSystem{
		particles:    make([]*Particle, 0, maxParticles),
		maxParticles: maxParticles,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Emit emite partículas desde una posición (THREAD-SAFE)
func (ps *ParticleSystem) Emit(position utils.Vector2, count int, particleColor color.RGBA) {
	ps.mu.Lock()         // ← LOCK: Bloquear acceso
	defer ps.mu.Unlock() // ← UNLOCK: Desbloquear al salir

	for i := 0; i < count; i++ {
		if len(ps.particles) >= ps.maxParticles {
			break
		}

		// Velocidad aleatoria
		angle := ps.rng.Float64() * 2 * math.Pi
		speed := 2.0 + ps.rng.Float64()*4.0
		velocity := utils.Vector2{
			X: speed * math.Cos(angle),
			Y: speed * math.Sin(angle),
		}

		particle := &Particle{
			Position: position,
			Velocity: velocity,
			Color:    particleColor,
			Size:     2 + ps.rng.Float64()*4,
			Lifetime: 20 + ps.rng.Intn(20),
			Age:      0,
			IsActive: true,
		}

		ps.particles = append(ps.particles, particle)
	}
}

// Update actualiza todas las partículas (THREAD-SAFE)
func (ps *ParticleSystem) Update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i := len(ps.particles) - 1; i >= 0; i-- {
		p := ps.particles[i]

		if !p.IsActive {
			// Remover partícula inactiva
			ps.particles = append(ps.particles[:i], ps.particles[i+1:]...)
			continue
		}

		// Actualizar partícula
		p.Age++
		p.Position = p.Position.Add(p.Velocity)
		p.Velocity.Y += 0.2 // Gravedad

		// Fade out
		alpha := float64(p.Lifetime-p.Age) / float64(p.Lifetime)
		p.Color.A = uint8(255 * alpha)

		// Desactivar si cumplió su tiempo
		if p.Age >= p.Lifetime {
			p.IsActive = false
		}
	}
}

// GetParticles retorna una COPIA de las partículas activas (THREAD-SAFE)
func (ps *ParticleSystem) GetParticles() []*Particle {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Crear copia para evitar data races durante el dibujado
	particles := make([]*Particle, len(ps.particles))
	copy(particles, ps.particles)
	return particles
}

// Clear limpia todas las partículas (THREAD-SAFE)
func (ps *ParticleSystem) Clear() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.particles = ps.particles[:0]
}
