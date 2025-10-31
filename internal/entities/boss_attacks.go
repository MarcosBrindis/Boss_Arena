// internal/entities/boss_attacks.go
package entities

import (
	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
)

// performBasicAttack realiza el ataque básico
func (b *Boss) performBasicAttack() {
	if b.AttackCooldown > 0 {
		return
	}

	b.State = BossStateAttacking
	b.AttackCooldown = b.config.AttackCooldown
	b.AttackDelay = 15 // 0.25 segundos antes de hacer daño

	// Pequeño impulso hacia adelante
	if b.FacingRight {
		b.Velocity.X = 3
	} else {
		b.Velocity.X = -3
	}
}

// performSlam realiza el ataque Slam
func (b *Boss) performSlam() {
	if b.SlamCooldown > 0 {
		return
	}

	b.State = BossStateSlam
	b.SlamDuration = b.config.SlamDuration
	b.SlamCooldown = b.config.SlamCooldown
	b.Velocity = utils.Zero()

	// TODO: En el próximo módulo, crear shockwave
}

// performCharge realiza el ataque Charge
func (b *Boss) performCharge() {
	if b.ChargeCooldown > 0 || b.Target == nil {
		return
	}

	b.State = BossStateCharge
	b.ChargeDuration = b.config.ChargeDuration
	b.ChargeCooldown = b.config.ChargeCooldown

	// Dirección hacia el jugador
	direction := b.Target.Position.Sub(b.Position).Normalize()
	b.ChargeDirection = direction

	// Velocidad aumenta con las fases
	b.ChargeSpeed = b.config.ChargeSpeed
	if b.Phase == Phase2 {
		b.ChargeSpeed *= 1.2
	} else if b.Phase == Phase3 {
		b.ChargeSpeed *= 1.5
	}
}

// performRoar realiza el ataque Roar
func (b *Boss) performRoar() {
	if b.RoarCooldown > 0 {
		return
	}

	b.State = BossStateRoar
	b.RoarDuration = b.config.RoarDuration
	b.RoarCooldown = b.config.RoarCooldown
	b.Velocity = utils.Zero()

	// Aplicar stun al jugador si está en rango
	if b.Target != nil {
		distance := b.Position.Distance(b.Target.Position)
		if distance <= b.config.RoarRange {
			b.stunPlayer()
		}
	}
}

// stunPlayer aturde al jugador
func (b *Boss) stunPlayer() {
	if b.Target == nil {
		return
	}

	// TODO: Implementar en el Módulo 6 (Combat System)
	// Por ahora solo guardamos que debe aplicarse
}

// GetAttackHitbox retorna el hitbox del ataque actual
func (b *Boss) GetAttackHitbox() *utils.Rectangle {
	if b.State != BossStateAttacking || b.AttackDelay > 0 {
		return nil
	}

	// Hitbox extendido en la dirección que mira
	attackReach := 50.0
	attackWidth := 60.0
	attackHeight := 80.0

	var attackX float64
	if b.FacingRight {
		attackX = b.Position.X + b.Size.X/2
	} else {
		attackX = b.Position.X - b.Size.X/2 - attackWidth
	}

	rect := utils.NewRectangle(
		attackX,
		b.Position.Y-attackHeight/2,
		attackWidth+attackReach,
		attackHeight,
	)

	return &rect
}

// GetSlamHitbox retorna el hitbox del slam
func (b *Boss) GetSlamHitbox() *utils.Rectangle {
	if b.State != BossStateSlam || b.SlamDuration > 15 {
		return nil
	}

	// Área circular alrededor del boss
	rect := utils.NewRectangle(
		b.Position.X-b.config.SlamRadius,
		b.Position.Y-b.config.SlamRadius,
		b.config.SlamRadius*2,
		b.config.SlamRadius*2,
	)

	return &rect
}

// GetChargeHitbox retorna el hitbox del charge
func (b *Boss) GetChargeHitbox() *utils.Rectangle {
	if b.State != BossStateCharge {
		return nil
	}

	// NOTA: El jugador puede hacer pogo sobre el boss durante charge
	// Esto es intencional - alto riesgo, alta recompensa
	hitbox := b.GetHitbox()
	return &hitbox
}
