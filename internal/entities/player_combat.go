// internal/entities/player_combat.go
package entities

import (
	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
)

// handleAttackInput maneja el input de ataque
func (p *Player) handleAttackInput() {
	// Primero verificar Down Air Attack (prioridad)
	p.handleDownAirAttackInput()

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

	attackReach := 40.0
	attackWidth := 50.0
	attackHeight := 60.0

	var attackX float64
	if p.FacingRight {
		// Atacando a la derecha: empieza desde el borde derecho del jugador
		attackX = p.Position.X + p.Size.X/2
	} else {
		// Atacando a la izquierda: empieza desde attackReach más allá del borde izquierdo
		attackX = p.Position.X - p.Size.X/2 - attackWidth - attackReach // ← ARREGLADO
	}

	rect := utils.NewRectangle(
		attackX,
		p.Position.Y-attackHeight/2,
		attackWidth+attackReach, // Ancho total del ataque
		attackHeight,
	)

	return &rect
}

// GetAttackDamage retorna el daño del ataque actual
func (p *Player) GetAttackDamage() int {
	baseDamage := 10
	return baseDamage * p.ComboCount
}

// handleDownAirAttackInput maneja el ataque hacia abajo en el aire
func (p *Player) handleDownAirAttackInput() {
	// Solo en el aire
	if p.IsOnGround {
		return
	}

	// Permitir cancelar otros ataques aéreos (MEJORADO)
	if p.State == StateDashing {
		return // No cancelar dash
	}

	// Verificar input: Down + Attack
	if !p.controller.IsAttackPressed() {
		return
	}

	verticalInput := p.controller.GetVerticalAxis()
	if verticalInput <= 0.5 { // No está presionando down
		return
	}

	// Verificar stamina
	staminaCost := 15.0
	if p.Stamina < staminaCost {
		p.controller.Vibrate(50, 0.1)
		return
	}

	p.performDownAirAttack()
}

// performDownAirAttack realiza el ataque hacia abajo
func (p *Player) performDownAirAttack() {
	// Consumir stamina
	p.Stamina -= 15

	p.State = StateDownAirAttack
	p.AttackTimeLeft = 20 // Duración del ataque

	// Impulso hacia abajo (para el pogo effect)
	p.Velocity.Y = 8 // Caída rápida

	// Vibración
	p.controller.Vibrate(100, 0.4)
}

// GetDownAirAttackHitbox retorna el hitbox del ataque hacia abajo
func (p *Player) GetDownAirAttackHitbox() *utils.Rectangle {
	if p.State != StateDownAirAttack {
		return nil
	}

	// Hitbox debajo del jugador
	attackWidth := 50.0
	attackHeight := 40.0

	rect := utils.NewRectangle(
		p.Position.X-attackWidth/2,
		p.Position.Y+p.Size.Y/2-10, // Justo debajo
		attackWidth,
		attackHeight,
	)

	return &rect
}

// GetDownAirAttackDamage retorna el daño del ataque hacia abajo
func (p *Player) GetDownAirAttackDamage() int {
	return 20 // Más daño que ataque normal
}
