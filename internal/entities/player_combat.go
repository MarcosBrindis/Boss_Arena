// internal/entities/player_combat.go
package entities

import (
	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
)

// handleAttackInput maneja el input de ataque
func (p *Player) handleAttackInput() {
	if !p.controller.IsAttackPressed() {
		return
	}

	// No atacar si no está en estado válido
	if p.State != StateIdle && p.State != StateWalking && p.State != StateJumping && p.State != StateFalling {
		return
	}

	p.performAttack()
}

// performAttack realiza un ataque
func (p *Player) performAttack() {
	// Verificar stamina
	staminaCost := 10.0 + float64(p.ComboCount)*5.0 // Cada combo cuesta más

	if p.Stamina < staminaCost {
		// No hay stamina suficiente
		p.controller.Vibrate(50, 0.1)
		return
	}

	// Consumir stamina
	p.Stamina -= staminaCost

	// Incrementar combo si está dentro del tiempo
	if p.ComboTimeLeft > 0 && p.ComboCount < p.config.MaxCombo {
		p.ComboCount++
	} else {
		p.ComboCount = 1
	}

	p.ComboTimeLeft = p.config.ComboDuration
	p.AttackTimeLeft = p.config.AttackDuration
	p.State = StateAttacking

	// Pequeño impulso hacia adelante al atacar
	if p.IsOnGround {
		pushForce := 2.0
		if p.FacingRight {
			p.Velocity.X += pushForce
		} else {
			p.Velocity.X -= pushForce
		}
	}

	// Vibración al atacar (más fuerte según el combo)
	vibrateStrength := 0.3 + float64(p.ComboCount)*0.1
	p.controller.Vibrate(80, vibrateStrength)
}

// GetAttackHitbox retorna el hitbox del ataque
func (p *Player) GetAttackHitbox() *utils.Rectangle {
	if p.State != StateAttacking {
		return nil
	}

	// Hitbox extendido en la dirección que mira
	attackReach := 30.0
	attackWidth := 40.0
	attackHeight := 50.0

	var attackX float64
	if p.FacingRight {
		attackX = p.Position.X + p.Size.X/2
	} else {
		attackX = p.Position.X - p.Size.X/2 - attackWidth
	}

	rect := utils.NewRectangle(
		attackX,
		p.Position.Y-attackHeight/2,
		attackWidth+attackReach,
		attackHeight,
	)

	return &rect
}

// GetAttackDamage retorna el daño del ataque actual
func (p *Player) GetAttackDamage() int {
	baseDamage := 10
	return baseDamage * p.ComboCount
}
