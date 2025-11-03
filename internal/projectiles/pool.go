package projectiles

import (
	"sync"

	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
)

// ProjectilePool implementa Object Pooling con channels
type ProjectilePool struct {
	// Pool como channel buffered
	pool chan *Projectile

	// Contador de IDs
	nextID int
	idMu   sync.Mutex

	// Estadísticas
	created int
	reused  int
	statsMu sync.Mutex
}

// NewProjectilePool crea un nuevo pool de proyectiles
func NewProjectilePool(size int) *ProjectilePool {
	return &ProjectilePool{
		pool:   make(chan *Projectile, size),
		nextID: 1,
	}
}

// Get obtiene un proyectil del pool (reutiliza o crea nuevo)
func (pp *ProjectilePool) Get(
	projectileType ProjectileType,
	position, direction utils.Vector2,
	owner string,
) *Projectile {
	var projectile *Projectile

	// Intentar obtener del pool
	select {
	case projectile = <-pp.pool:
		// Reutilizar proyectil existente
		pp.statsMu.Lock()
		pp.reused++
		pp.statsMu.Unlock()

	default:
		// Pool vacío, crear nuevo
		pp.idMu.Lock()
		id := pp.nextID
		pp.nextID++
		pp.idMu.Unlock()

		projectile = NewProjectile(id, projectileType, position, direction, owner)

		pp.statsMu.Lock()
		pp.created++
		pp.statsMu.Unlock()
	}

	// Reconfigurar el proyectil
	projectile.Type = projectileType
	projectile.Position = position
	projectile.Owner = owner
	projectile.Age = 0
	projectile.IsActive = true

	// Aplicar configuración según tipo
	switch projectileType {
	case ProjectilePlayerBasic:
		projectile.Speed = 12.0
		projectile.Damage = 15
		projectile.Lifetime = 180
		projectile.Size = utils.NewVector2(8, 8)
		projectile.Color.R, projectile.Color.G, projectile.Color.B = 0, 200, 255
		projectile.IsHoming = false

	case ProjectilePlayerCharged:
		projectile.Speed = 10.0
		projectile.Damage = 30
		projectile.Lifetime = 240
		projectile.Size = utils.NewVector2(12, 12)
		projectile.Color.R, projectile.Color.G, projectile.Color.B = 255, 255, 0
		projectile.IsHoming = false

	case ProjectileBossFireball:
		projectile.Speed = 8.0
		projectile.Damage = 20
		projectile.Lifetime = 300
		projectile.Size = utils.NewVector2(16, 16)
		projectile.Color.R, projectile.Color.G, projectile.Color.B = 255, 69, 0
		projectile.IsHoming = false

	case ProjectileBossMissile:
		projectile.Speed = 6.0
		projectile.Damage = 25
		projectile.Lifetime = 360
		projectile.Size = utils.NewVector2(10, 10)
		projectile.Color.R, projectile.Color.G, projectile.Color.B = 255, 0, 255
		projectile.IsHoming = true
		projectile.HomingForce = 0.3
	}

	projectile.Color.A = 255
	projectile.Velocity = direction.Normalize().Mul(projectile.Speed)

	return projectile
}

// Put devuelve un proyectil al pool
func (pp *ProjectilePool) Put(projectile *Projectile) {
	if projectile == nil {
		return
	}

	// Resetear el proyectil
	projectile.Reset()

	// Intentar devolver al pool (non-blocking)
	select {
	case pp.pool <- projectile:
		// Devuelto exitosamente
	default:
		// Pool lleno, descartar (será recolectado por GC)
	}
}

// GetStats retorna estadísticas del pool
func (pp *ProjectilePool) GetStats() (created, reused int) {
	pp.statsMu.Lock()
	defer pp.statsMu.Unlock()

	return pp.created, pp.reused
}

// Clear limpia el pool
func (pp *ProjectilePool) Clear() {
	// Vaciar el canal
	for {
		select {
		case <-pp.pool:
			// Descartar proyectil
		default:
			return
		}
	}
}
