package entities

import "github.com/MarcosBrindis/boss-arena-go/internal/utils"

// updateAI actualiza la inteligencia artificial del boss
func (b *Boss) updateAI() {
	// No procesar IA en estos estados
	if b.State == BossStateSlam ||
		b.State == BossStateCharge ||
		b.State == BossStateRoar ||
		b.State == BossStateStunned ||
		b.State == BossStateTransition ||
		b.State == BossStateDead {
		return
	}

	// No hay objetivo
	if b.Target == nil {
		return
	}

	// Actualizar dirección hacia el jugador
	b.updateFacingDirection()

	// Tomar decisiones cada cierto tiempo
	if b.DecisionTimer <= 0 {
		b.makeDecision()
		b.DecisionTimer = b.config.DecisionDelay
	}

	// Ejecutar acción actual
	b.executeAction()
}

// updateFacingDirection actualiza la dirección hacia el jugador
func (b *Boss) updateFacingDirection() {
	if b.Target == nil {
		return
	}

	b.FacingRight = b.Target.Position.X > b.Position.X
}

// makeDecision decide la próxima acción del boss
// makeDecision decide la próxima acción del boss
func (b *Boss) makeDecision() {
	if b.Target == nil {
		return
	}

	// Distancia 2D completa (incluye altura)
	distanceToPlayer := b.Position.Distance(b.Target.Position)

	// Distancia solo horizontal (para algunos ataques específicos)
	horizontalDistance := utils.Abs(b.Position.X - b.Target.Position.X)

	// Usar distancia horizontal para decisiones de movimiento
	_ = horizontalDistance // Por ahora no se usa, pero está disponible

	// Lista de ataques disponibles
	availableAttacks := []BossState{}

	// Ataque básico siempre disponible si está en rango
	if distanceToPlayer <= b.config.AttackRange && b.AttackCooldown == 0 {
		availableAttacks = append(availableAttacks, BossStateAttacking)
	}

	// Slam disponible
	if b.SlamCooldown == 0 && distanceToPlayer <= b.config.SlamRadius {
		availableAttacks = append(availableAttacks, BossStateSlam)
	}

	// Charge disponible (si el jugador está lejos)
	if b.ChargeCooldown == 0 && distanceToPlayer > 150 && distanceToPlayer < 400 {
		availableAttacks = append(availableAttacks, BossStateCharge)
	}

	// Roar disponible
	if b.RoarCooldown == 0 && distanceToPlayer <= b.config.RoarRange {
		availableAttacks = append(availableAttacks, BossStateRoar)
	}

	// Elegir ataque según fase
	if len(availableAttacks) > 0 {
		// Fase 3 = más agresivo (prioriza ataques especiales)
		if b.Phase == Phase3 && len(availableAttacks) > 1 {
			// Elegir el ataque más fuerte disponible
			for _, attack := range availableAttacks {
				if attack == BossStateSlam || attack == BossStateCharge || attack == BossStateRoar {
					b.NextAction = attack
					return
				}
			}
		}

		// Elegir ataque aleatorio
		b.NextAction = availableAttacks[b.rng.Intn(len(availableAttacks))]
	} else {
		// No hay ataques disponibles, acercarse al jugador
		if distanceToPlayer > b.config.AttackRange {
			b.NextAction = BossStateWalking
		} else {
			b.NextAction = BossStateIdle
		}
	}
}

// executeAction ejecuta la acción decidida
func (b *Boss) executeAction() {
	switch b.NextAction {
	case BossStateWalking:
		b.walkTowardsPlayer()
	case BossStateAttacking:
		b.performBasicAttack()
	case BossStateSlam:
		b.performSlam()
	case BossStateCharge:
		b.performCharge()
	case BossStateRoar:
		b.performRoar()
	}
}

// walkTowardsPlayer camina hacia el jugador
func (b *Boss) walkTowardsPlayer() {
	if b.Target == nil || !b.IsOnGround {
		return
	}

	direction := 1.0
	if b.Target.Position.X < b.Position.X {
		direction = -1.0
	}

	// Velocidad aumenta con las fases
	speed := b.config.WalkSpeed
	if b.Phase == Phase2 {
		speed *= 1.3
	} else if b.Phase == Phase3 {
		speed *= 1.6
	}

	b.Velocity.X = direction * speed
}
