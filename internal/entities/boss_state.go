package entities

// BossState representa el estado actual del boss
type BossState int

const (
	BossStateIdle BossState = iota
	BossStateWalking
	BossStateJumping
	BossStateFalling
	BossStateAttacking
	BossStateSlam       // Ataque de golpe en el suelo
	BossStateCharge     // Ataque de carga
	BossStateRoar       // Rugido (stun a jugador)
	BossStateShooting   // Disparo de proyectil
	BossStateStunned    // Aturdido (vulnerable)
	BossStateTransition // Transición de fase
	BossStateDead
)

// String retorna el nombre del estado (para debug)
func (s BossState) String() string {
	switch s {
	case BossStateIdle:
		return "Idle"
	case BossStateWalking:
		return "Walking"
	case BossStateJumping:
		return "Jumping"
	case BossStateFalling:
		return "Falling"
	case BossStateAttacking:
		return "Attacking"
	case BossStateSlam:
		return "Slam"
	case BossStateCharge:
		return "Charge"
	case BossStateRoar:
		return "Roar"
	case BossStateShooting:
		return "Shooting"
	case BossStateStunned:
		return "Stunned"
	case BossStateTransition:
		return "Transition"
	case BossStateDead:
		return "Dead"
	default:
		return "Unknown"
	}
}

// BossPhase representa la fase actual del boss
type BossPhase int

const (
	Phase1 BossPhase = iota // 100% - 66% HP (normal)
	Phase2                  // 66% - 33% HP (agresivo)
	Phase3                  // 33% - 0% HP (berserk)
)

// String retorna el nombre de la fase
func (p BossPhase) String() string {
	switch p {
	case Phase1:
		return "Phase 1 (Normal)"
	case Phase2:
		return "Phase 2 (Aggressive)"
	case Phase3:
		return "Phase 3 (Berserk)"
	default:
		return "Unknown"
	}
}

// GetColor retorna el color según la fase
func (p BossPhase) GetColor() (r, g, b uint8) {
	switch p {
	case Phase1:
		return 255, 69, 0 // Rojo fuego
	case Phase2:
		return 255, 99, 71 // Naranja
	case Phase3:
		return 255, 165, 0 // Amarillo
	default:
		return 255, 255, 255 // Blanco
	}
}
