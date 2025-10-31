package entities

import (
	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
)

// handleInput procesa el input del jugador
func (p *Player) handleInput() {
	// No procesar input si está en estados bloqueados
	if p.State == StateDashing || p.State == StateAttacking || p.State == StateHurt || p.State == StateDead {
		return
	}

	// Movimiento horizontal
	p.handleHorizontalInput()

	// Salto
	p.handleJumpInput()

	// Dash
	p.handleDashInput()

	// Ataque (lo implementaremos en player_combat.go)
	p.handleAttackInput()
}

// handleHorizontalInput maneja el movimiento horizontal
func (p *Player) handleHorizontalInput() {
	inputX := p.controller.GetHorizontalAxis()

	if inputX != 0 {
		// Actualizar dirección
		p.FacingRight = inputX > 0

		// Aplicar aceleración
		control := p.config.Acceleration
		if !p.IsOnGround {
			control *= p.config.AirControl
		}

		targetSpeed := inputX * p.config.MoveSpeed
		p.Velocity.X += (targetSpeed - p.Velocity.X) * control

		// Limitar velocidad máxima
		if utils.Abs(p.Velocity.X) > p.config.MaxSpeed {
			p.Velocity.X = utils.Sign(p.Velocity.X) * p.config.MaxSpeed
		}
	}
}

// handleJumpInput maneja el input de salto
func (p *Player) handleJumpInput() {
	// Detectar presión de salto
	if p.controller.IsJumpPressed() {
		p.tryJump()
	}

	// Salto variable (soltar botón = caer más rápido)
	if !p.controller.IsJumpHeld() && p.Velocity.Y < 0 {
		p.Velocity.Y *= 0.5
	}
}

// tryJump intenta realizar un salto
func (p *Player) tryJump() {
	// =========================================================================
	// PRIORIDAD 1: WALL CLIMB O WALL JUMP
	// =========================================================================

	if p.IsTouchingWall && !p.IsOnGround && p.WallJumpCooldown == 0 {
		// Verificar dirección del input
		inputX := p.controller.GetHorizontalAxis()

		// Si presiona HACIA la pared = wall climb (escalar)
		if (p.WallSide == -1 && inputX < -0.3) || (p.WallSide == 1 && inputX > 0.3) {
			// Verificar stamina para wall climb
			if p.Stamina >= 15 {
				p.performWallClimb()
				return
			} else {
				p.controller.ConsumeJumpBuffer()
				p.controller.Vibrate(50, 0.1)
				return
			}
		} else {
			// Si presiona LEJOS de la pared = wall jump tradicional
			if p.Stamina >= 10 {
				p.performWallJump()
				return
			} else {
				p.controller.ConsumeJumpBuffer()
				p.controller.Vibrate(50, 0.1)
				return
			}
		}
	}

	// =========================================================================
	// PRIORIDAD 2: SALTO DESDE EL SUELO (o coyote time)
	// =========================================================================

	if p.IsOnGround || p.CoyoteTimeLeft > 0 {
		p.performJump(p.config.JumpForce)
		p.JumpCount = 1
		p.CoyoteTimeLeft = 0
		return
	}

	// =========================================================================
	// PRIORIDAD 3: DOBLE SALTO (si hay stamina)
	// =========================================================================

	if p.JumpCount < p.MaxJumps && !p.IsTouchingWall {
		if p.Stamina >= 20 {
			p.performJump(p.config.DoubleJumpForce)
			p.JumpCount++
			p.Stamina -= 20
			return
		} else {
			p.controller.ConsumeJumpBuffer()
			p.controller.Vibrate(50, 0.1)
			return
		}
	}

	// =========================================================================
	// NO SE PUDO SALTAR
	// =========================================================================

	p.controller.ConsumeJumpBuffer()
}

// performJump realiza un salto
func (p *Player) performJump(force float64) {
	p.Velocity.Y = -force
	p.State = StateJumping

	// Vibración sutil al saltar
	p.controller.Vibrate(30, 0.15)
}

// performWallJump realiza un wall jump ALEJÁNDOSE de la pared
func (p *Player) performWallJump() {
	// Saltar en dirección opuesta a la pared
	direction := float64(-p.WallSide)

	p.Velocity.X = direction * p.config.WallJumpForceX
	p.Velocity.Y = -p.config.WallJumpForceY

	p.State = StateWallJumping
	p.JumpCount = 1
	p.WallJumpCooldown = p.config.WallStickFrames

	// Cambiar dirección
	p.FacingRight = direction > 0

	// Consumir stamina
	p.Stamina -= 10

	// Vibración al wall jump
	p.controller.Vibrate(50, 0.25)
}

// performWallClimb permite escalar la pared saltando (cuesta stamina)
func (p *Player) performWallClimb() {
	// Consumir stamina
	p.Stamina -= 15

	// Salto más alto que el normal (para escalar)
	// Pero sin empuje horizontal (se queda en la pared)
	wallClimbForce := p.config.JumpForce * 1.2 // 20% más alto

	p.Velocity.Y = -wallClimbForce
	p.Velocity.X *= 0.3 // Reducir velocidad horizontal para quedarse cerca de la pared

	p.State = StateJumping
	p.JumpCount = 1
	p.WallJumpCooldown = 5 // Cooldown corto entre saltos

	// Vibración al wall climb
	p.controller.Vibrate(40, 0.2)
}

// handleDashInput maneja el input de dash
func (p *Player) handleDashInput() {
	if !p.controller.IsDashPressed() {
		return
	}

	if !p.CanDash || p.Stamina < 25 {
		return
	}

	p.performDash()
}

// performDash realiza un dash
func (p *Player) performDash() {
	// Dirección del dash
	inputX := p.controller.GetHorizontalAxis()
	inputY := p.controller.GetVerticalAxis()

	var dashDir utils.Vector2

	if inputX != 0 || inputY != 0 {
		// Dash en dirección del input
		dashDir = utils.NewVector2(inputX, inputY).Normalize()
	} else {
		// Dash en dirección que mira el jugador
		if p.FacingRight {
			dashDir = utils.Right()
		} else {
			dashDir = utils.Left()
		}
	}

	// Aplicar dash
	p.Velocity = dashDir.Mul(p.config.DashSpeed)
	p.DashDirection = dashDir
	p.DashTimeLeft = p.config.DashDuration
	p.DashCooldown = p.config.DashCooldown
	p.CanDash = false
	p.State = StateDashing

	// Consumir stamina
	p.Stamina -= 25

	// Vibración al dashear
	p.controller.Vibrate(100, 0.3)
}

// applyPhysics aplica física al jugador
func (p *Player) applyPhysics() {
	// Durante dash, mantener velocidad constante (no aplicar física)
	if p.State == StateDashing {
		// Mantener la velocidad del dash sin cambios
		p.Velocity = p.DashDirection.Mul(p.config.DashSpeed)
		return
	}

	// Gravedad
	if !p.IsOnGround {
		// Wall sliding reduce la gravedad
		if p.State == StateWallSliding {
			p.Velocity.Y += p.config.Gravity * 0.3
			if p.Velocity.Y > p.config.WallSlideSpeed {
				p.Velocity.Y = p.config.WallSlideSpeed
			}
		} else {
			p.Velocity.Y += p.config.Gravity
			if p.Velocity.Y > p.config.MaxFallSpeed {
				p.Velocity.Y = p.config.MaxFallSpeed
			}
		}
	}

	// Fricción
	if p.IsOnGround {
		p.Velocity.X *= p.config.GroundFriction
		if utils.Abs(p.Velocity.X) < 0.1 {
			p.Velocity.X = 0
		}
	} else {
		p.Velocity.X *= p.config.AirFriction
	}
}

// applyMovement aplica el movimiento con colisiones
func (p *Player) applyMovement() {
	// =========================================================================
	// MOVIMIENTO HORIZONTAL (separado del vertical)
	// =========================================================================

	newX := p.Position.X + p.Velocity.X

	// Crear hitbox en nueva posición X
	testRect := utils.NewRectangle(
		newX-p.Size.X/2,
		p.Position.Y-p.Size.Y/2,
		p.Size.X,
		p.Size.Y,
	)

	// Verificar colisión horizontal
	collides, _ := p.arena.CheckCollision(testRect)
	if !collides {
		// No hay colisión, mover
		p.Position.X = newX
	} else {
		// Hay colisión, detener velocidad
		p.Velocity.X = 0
	}

	// =========================================================================
	// MOVIMIENTO VERTICAL (separado del horizontal)
	// =========================================================================

	newY := p.Position.Y + p.Velocity.Y

	// Crear hitbox en nueva posición Y
	testRect = utils.NewRectangle(
		p.Position.X-p.Size.X/2,
		newY-p.Size.Y/2,
		p.Size.X,
		p.Size.Y,
	)

	// Verificar colisión vertical
	collides, _ = p.arena.CheckCollision(testRect)
	if !collides {
		// No hay colisión, mover
		p.Position.Y = newY
	} else {
		// Hay colisión, detener velocidad
		p.Velocity.Y = 0
	}

	// =========================================================================
	// LÍMITE DE SEGURIDAD: No caer demasiado abajo
	// =========================================================================

	floorY := p.arena.GetFloorY()
	if p.Position.Y > floorY+100 {
		// Resetear posición si cae muy abajo
		p.Position.X = 640
		p.Position.Y = 300
		p.Velocity = utils.Zero()
	}

	// Límites laterales
	margin := p.Size.X / 2
	if p.Position.X < margin+60 {
		p.Position.X = margin + 60
		p.Velocity.X = 0
	}
	if p.Position.X > 1280-margin-60 {
		p.Position.X = 1280 - margin - 60
		p.Velocity.X = 0
	}

	// Límite superior
	if p.Position.Y < p.Size.Y/2 {
		p.Position.Y = p.Size.Y / 2
		p.Velocity.Y = 0
	}
}
