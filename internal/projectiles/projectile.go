package projectiles

import (
	"image/color"

	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ProjectileType representa el tipo de proyectil
type ProjectileType int

const (
	ProjectilePlayerBasic   ProjectileType = iota // Proyectil básico del jugador
	ProjectilePlayerCharged                       // Proyectil cargado (más daño)
	ProjectileBossFireball                        // Bola de fuego del boss
	ProjectileBossMissile                         // Misil del boss (persigue)
)

// Projectile representa un proyectil en el juego
type Projectile struct {
	// Identificación
	ID   int
	Type ProjectileType

	// Física
	Position utils.Vector2
	Velocity utils.Vector2
	Size     utils.Vector2

	// Propiedades
	Damage   int
	Speed    float64
	Lifetime int // Frames antes de auto-destruirse
	Age      int // Frames transcurridos
	IsActive bool

	// Propietario
	Owner string // "player" o "boss"

	// Visual
	Color color.RGBA

	// Comportamiento especial
	IsHoming    bool           // Si persigue al objetivo
	Target      *utils.Vector2 // Objetivo para proyectiles homing
	HomingForce float64        // Fuerza de persecución
}

// NewProjectile crea un nuevo proyectil (factory function)
func NewProjectile(id int, projectileType ProjectileType, position, direction utils.Vector2, owner string) *Projectile {
	p := &Projectile{
		ID:       id,
		Type:     projectileType,
		Position: position,
		Owner:    owner,
		Age:      0,
		IsActive: true,
	}

	// Configurar según tipo
	switch projectileType {
	case ProjectilePlayerBasic:
		p.Speed = 12.0
		p.Damage = 15
		p.Lifetime = 180 // 3 segundos @ 60 FPS
		p.Size = utils.NewVector2(8, 8)
		p.Color = color.RGBA{0, 200, 255, 255} // Cyan
		p.IsHoming = false

	case ProjectilePlayerCharged:
		p.Speed = 10.0
		p.Damage = 30
		p.Lifetime = 240
		p.Size = utils.NewVector2(12, 12)
		p.Color = color.RGBA{255, 255, 0, 255} // Amarillo
		p.IsHoming = false

	case ProjectileBossFireball:
		p.Speed = 8.0
		p.Damage = 20
		p.Lifetime = 300
		p.Size = utils.NewVector2(16, 16)
		p.Color = color.RGBA{255, 69, 0, 255} // Rojo fuego
		p.IsHoming = false

	case ProjectileBossMissile:
		p.Speed = 6.0
		p.Damage = 25
		p.Lifetime = 360
		p.Size = utils.NewVector2(10, 10)
		p.Color = color.RGBA{255, 0, 255, 255} // Magenta
		p.IsHoming = true
		p.HomingForce = 0.3
	}

	// Calcular velocidad inicial
	p.Velocity = direction.Normalize().Mul(p.Speed)

	return p
}

// Update actualiza el proyectil
func (p *Projectile) Update() {
	if !p.IsActive {
		return
	}

	// Incrementar edad
	p.Age++

	// Auto-destruir si excede lifetime
	if p.Age >= p.Lifetime {
		p.IsActive = false
		return
	}

	// Comportamiento homing (perseguir objetivo)
	if p.IsHoming && p.Target != nil {
		p.updateHoming()
	}

	// Aplicar movimiento
	p.Position = p.Position.Add(p.Velocity)

	// Verificar límites de pantalla
	p.checkBounds()
}

// updateHoming actualiza la trayectoria para perseguir al objetivo
func (p *Projectile) updateHoming() {
	// Dirección hacia el objetivo
	toTarget := p.Target.Sub(p.Position)
	distance := toTarget.Length()

	if distance > 0 {
		// Normalizar dirección
		direction := toTarget.Normalize()

		// Aplicar fuerza de persecución
		steering := direction.Mul(p.HomingForce)
		p.Velocity = p.Velocity.Add(steering)

		// Limitar velocidad máxima
		if p.Velocity.Length() > p.Speed {
			p.Velocity = p.Velocity.Normalize().Mul(p.Speed)
		}
	}
}

// checkBounds verifica límites de pantalla
func (p *Projectile) checkBounds() {
	// Destruir si sale de la pantalla
	if p.Position.X < -50 || p.Position.X > 1280+50 ||
		p.Position.Y < -50 || p.Position.Y > 720+50 {
		p.IsActive = false
	}
}

// GetHitbox retorna el rectángulo de colisión
func (p *Projectile) GetHitbox() utils.Rectangle {
	return utils.NewRectangle(
		p.Position.X-p.Size.X/2,
		p.Position.Y-p.Size.Y/2,
		p.Size.X,
		p.Size.Y,
	)
}

// Draw dibuja el proyectil
func (p *Projectile) Draw(screen *ebiten.Image) {
	if !p.IsActive {
		return
	}

	// Dibujar círculo
	vector.DrawFilledCircle(
		screen,
		float32(p.Position.X),
		float32(p.Position.Y),
		float32(p.Size.X/2),
		p.Color,
		false,
	)

	// Borde
	vector.StrokeCircle(
		screen,
		float32(p.Position.X),
		float32(p.Position.Y),
		float32(p.Size.X/2),
		2,
		color.White,
		false,
	)

	// Trail para proyectiles rápidos
	if p.Speed > 10 {
		trailColor := p.Color
		trailColor.A = 100
		vector.DrawFilledCircle(
			screen,
			float32(p.Position.X-p.Velocity.X*0.5),
			float32(p.Position.Y-p.Velocity.Y*0.5),
			float32(p.Size.X/3),
			trailColor,
			false,
		)
	}
}

// Reset resetea el proyectil para reutilización (pooling)
func (p *Projectile) Reset() {
	p.Position = utils.Zero()
	p.Velocity = utils.Zero()
	p.Age = 0
	p.IsActive = false
	p.Target = nil
}
