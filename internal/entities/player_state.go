package entities

// PlayerState representa el estado actual del jugador
type PlayerState int

const (
	StateIdle PlayerState = iota
	StateWalking
	StateJumping
	StateFalling
	StateWallSliding
	StateWallJumping
	StateDashing
	StateAttacking
	StateDownAirAttack
	StateHurt
	StateDead
)

// String retorna el nombre del estado (para debug)
func (s PlayerState) String() string {
	switch s {
	case StateIdle:
		return "Idle"
	case StateWalking:
		return "Walking"
	case StateJumping:
		return "Jumping"
	case StateFalling:
		return "Falling"
	case StateWallSliding:
		return "WallSliding"
	case StateWallJumping:
		return "WallJumping"
	case StateDashing:
		return "Dashing"
	case StateAttacking:
		return "Attacking"
	case StateDownAirAttack:
		return "DownAirAttack"
	case StateHurt:
		return "Hurt"
	case StateDead:
		return "Dead"
	default:
		return "Unknown"
	}
}

// CanTransitionTo verifica si puede cambiar a otro estado
func (s PlayerState) CanTransitionTo(newState PlayerState) bool {
	// Reglas de transición de estados
	switch s {
	case StateDashing:
		// Durante dash, solo puede ir a hurt o dead
		return newState == StateHurt || newState == StateDead

	case StateAttacking:
		// Durante ataque, solo puede ir a hurt, dead, o terminar ataque
		return newState == StateHurt || newState == StateDead || newState == StateIdle

	case StateDead:
		// Muerto es estado final
		return false

	default:
		// Otros estados pueden transicionar libremente
		return true
	}
}
