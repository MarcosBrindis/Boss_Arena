// internal/ai/dodge.go
package ai

import (
	"math"

	"github.com/MarcosBrindis/boss-arena-go/internal/projectiles"
	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
)

// DodgeSystem maneja la lógica de esquiva del boss
type DodgeSystem struct {
	// Configuración
	detectionRadius float64
	reactionTime    int // Frames antes de reaccionar
	dodgeSpeed      float64

	// Estado
	isDetectingThreat bool
	threatPosition    utils.Vector2
	dodgeTimer        int
}

// NewDodgeSystem crea un nuevo sistema de esquiva
func NewDodgeSystem() *DodgeSystem {
	return &DodgeSystem{
		detectionRadius: 200.0,
		reactionTime:    15, // 0.25 segundos @ 60 FPS
		dodgeSpeed:      5.0,
	}
}

// ShouldDodge verifica si el boss debe esquivar un proyectil
func (ds *DodgeSystem) ShouldDodge(
	bossPosition utils.Vector2,
	projectileList []*projectiles.Projectile,
) (bool, utils.Vector2) {
	// Buscar proyectiles cercanos del jugador
	for _, proj := range projectileList {
		if proj.Owner != "player" || !proj.IsActive {
			continue
		}

		// Calcular distancia al proyectil
		distance := bossPosition.Distance(proj.Position)

		// Si está dentro del radio de detección
		if distance <= ds.detectionRadius {
			// Predecir si el proyectil va a impactar
			willHit := ds.predictImpact(bossPosition, proj)

			if willHit {
				// Calcular dirección de esquiva (perpendicular al proyectil)
				dodgeDirection := ds.calculateDodgeDirection(bossPosition, proj)
				return true, dodgeDirection
			}
		}
	}

	return false, utils.Zero()
}

// predictImpact predice si un proyectil va a impactar al boss
func (ds *DodgeSystem) predictImpact(bossPosition utils.Vector2, proj *projectiles.Projectile) bool {
	// Vector del proyectil al boss
	toBoss := bossPosition.Sub(proj.Position)

	// Normalizar velocidad del proyectil
	projDirection := proj.Velocity.Normalize()

	// Producto punto para ver si el proyectil apunta hacia el boss
	dot := toBoss.Dot(projDirection)

	// Si el proyectil va hacia el boss y está cerca
	return dot > 0 && toBoss.Length() < ds.detectionRadius
}

// calculateDodgeDirection calcula la dirección de esquiva
func (ds *DodgeSystem) calculateDodgeDirection(bossPosition utils.Vector2, proj *projectiles.Projectile) utils.Vector2 {
	// Vector perpendicular a la velocidad del proyectil
	perpendicular := utils.Vector2{
		X: -proj.Velocity.Y,
		Y: proj.Velocity.X,
	}

	// Normalizar y aplicar velocidad de esquiva
	dodgeDir := perpendicular.Normalize().Mul(ds.dodgeSpeed)

	// Verificar que no se salga de los límites
	futurePos := bossPosition.Add(dodgeDir.Mul(10)) // Predecir 10 frames adelante

	// Si se sale por la derecha, esquivar a la izquierda
	if futurePos.X > 1200 {
		dodgeDir.X = -math.Abs(dodgeDir.X)
	}

	// Si se sale por la izquierda, esquivar a la derecha
	if futurePos.X < 80 {
		dodgeDir.X = math.Abs(dodgeDir.X)
	}

	return dodgeDir
}
