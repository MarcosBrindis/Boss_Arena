package entities

import (
	"image/color"

	"github.com/MarcosBrindis/boss-arena-go/internal/config"
	"github.com/MarcosBrindis/boss-arena-go/internal/input"
	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
	"github.com/MarcosBrindis/boss-arena-go/internal/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Player representa al jugador
type Player struct {
	// Posición y física
	Position utils.Vector2
	Velocity utils.Vector2
	Size     utils.Vector2

	// Estado
	State          PlayerState
	FacingRight    bool
	IsOnGround     bool
	IsTouchingWall bool
	WallSide       int // -1 = izquierda, 1 = derecha, 0 = ninguna

	// Salto
	JumpCount        int
	MaxJumps         int
	CoyoteTimeLeft   int
	WallJumpCooldown int

	// Dash
	CanDash       bool
	DashTimeLeft  int
	DashCooldown  int
	DashDirection utils.Vector2

	// Combate
	AttackTimeLeft int
	ComboCount     int
	ComboTimeLeft  int

	// Stats
	Health     int
	MaxHealth  int
	Stamina    float64
	MaxStamina float64

	// Referencias
	controller *input.Controller
	arena      *world.Arena

	// Configuración
	config PlayerConfig

	// Colores (para el placeholder antes de sprites)
	colorPrimary   color.RGBA
	colorSecondary color.RGBA
}

// PlayerConfig contiene la configuración del jugador
type PlayerConfig struct {
	// Movimiento
	MoveSpeed    float64
	MaxSpeed     float64
	Acceleration float64
	Deceleration float64
	AirControl   float64

	// Salto
	JumpForce        float64
	DoubleJumpForce  float64
	WallJumpForceX   float64
	WallJumpForceY   float64
	CoyoteTimeFrames int
	JumpBufferFrames int

	// Wall sliding
	WallSlideSpeed  float64
	WallStickFrames int

	// Dash
	DashSpeed    float64
	DashDuration int // Frames
	DashCooldown int // Frames

	// Combate
	AttackDuration int
	ComboDuration  int
	MaxCombo       int

	// Física
	Gravity        float64
	MaxFallSpeed   float64
	GroundFriction float64
	AirFriction    float64
}

// DefaultPlayerConfig retorna la configuración por defecto
func DefaultPlayerConfig() PlayerConfig {
	return PlayerConfig{
		// Movimiento
		MoveSpeed:    6.0,
		MaxSpeed:     8.0,
		Acceleration: 0.8,
		Deceleration: 0.9,
		AirControl:   0.6,

		// Salto
		JumpForce:        12.0,
		DoubleJumpForce:  10.0,
		WallJumpForceX:   10.0,
		WallJumpForceY:   12.0,
		CoyoteTimeFrames: 6,
		JumpBufferFrames: 5,

		// Wall sliding
		WallSlideSpeed:  1.5,
		WallStickFrames: 10,

		// Dash
		DashSpeed:    15.0,
		DashDuration: 10, // ~166ms a 60fps
		DashCooldown: 30, // ~500ms a 60fps

		// Combate
		AttackDuration: 15,
		ComboDuration:  30,
		MaxCombo:       3,

		// Física
		Gravity:        0.6,
		MaxFallSpeed:   12.0,
		GroundFriction: 0.85,
		AirFriction:    0.98,
	}
}

// NewPlayer crea un nuevo jugador
func NewPlayer(x, y float64, controller *input.Controller, arena *world.Arena) *Player {
	cfg := DefaultPlayerConfig()

	return &Player{
		Position: utils.NewVector2(x, y),
		Velocity: utils.Zero(),
		Size:     utils.NewVector2(40, 60), // 40x60 pixels

		State:       StateIdle,
		FacingRight: true,
		MaxJumps:    2, // Salto doble

		Health:     100,
		MaxHealth:  100,
		Stamina:    100,
		MaxStamina: 100,

		CanDash: true,

		controller: controller,
		arena:      arena,
		config:     cfg,

		colorPrimary:   config.ColorHeroPrimary,
		colorSecondary: config.ColorHeroSecondary,
	}
}

// Update actualiza el jugador
func (p *Player) Update() {
	// No actualizar si está muerto
	if p.State == StateDead {
		return
	}

	// 1. Actualizar temporizadores
	p.updateTimers()

	// 2. Verificar colisiones ANTES de todo (IMPORTANTE)
	p.updateCollisionState()

	// 3. Actualizar estado
	p.updateState()

	// 4. Procesar input
	p.handleInput()

	// 5. Aplicar física
	p.applyPhysics()

	// 6. Aplicar movimiento
	p.applyMovement()

	// 7. Regenerar stamina
	p.regenerateStamina()
}

// updateTimers actualiza todos los temporizadores
func (p *Player) updateTimers() {
	// Coyote time
	if p.CoyoteTimeLeft > 0 {
		p.CoyoteTimeLeft--
	}

	// Wall jump cooldown
	if p.WallJumpCooldown > 0 {
		p.WallJumpCooldown--
	}

	// Dash
	if p.DashTimeLeft > 0 {
		p.DashTimeLeft--
	}
	if p.DashCooldown > 0 {
		p.DashCooldown--
		if p.DashCooldown == 0 {
			p.CanDash = true
		}
	}

	// Ataque
	if p.AttackTimeLeft > 0 {
		p.AttackTimeLeft--
		if p.AttackTimeLeft == 0 {
			p.State = StateIdle
		}
	}

	// Combo
	if p.ComboTimeLeft > 0 {
		p.ComboTimeLeft--
		if p.ComboTimeLeft == 0 {
			p.ComboCount = 0
		}
	}
}

// updateCollisionState verifica colisiones con el mundo
func (p *Player) updateCollisionState() {
	hitbox := p.GetHitbox()

	// Verificar suelo
	wasOnGround := p.IsOnGround
	p.IsOnGround = p.arena.IsOnGround(hitbox)

	// Si acaba de tocar el suelo
	if p.IsOnGround && !wasOnGround {
		p.JumpCount = 0
		p.CanDash = true
		p.CoyoteTimeLeft = 0

		// Vibración sutil al aterrizar
		if p.Velocity.Y > 5 {
			p.controller.Vibrate(50, 0.2)
		}
	}

	// Si acaba de dejar el suelo (coyote time)
	if !p.IsOnGround && wasOnGround && p.Velocity.Y >= 0 {
		p.CoyoteTimeLeft = p.config.CoyoteTimeFrames
	}

	// Verificar paredes
	p.IsTouchingWall, p.WallSide = p.arena.IsTouchingWall(hitbox)
}

// updateState actualiza el estado del jugador
func (p *Player) updateState() {
	// Estados que terminan por tiempo
	if p.State == StateDashing && p.DashTimeLeft <= 0 {
		// Dash terminado, volver a estado normal
		p.State = StateIdle
		return
	}

	if p.State == StateAttacking && p.AttackTimeLeft <= 0 {
		// Ataque terminado
		p.State = StateIdle
		return
	}

	// Estados que no se pueden interrumpir mientras están activos
	if p.State == StateDashing || p.State == StateAttacking || p.State == StateHurt {
		return
	}

	// Determinar nuevo estado basado en condiciones
	var newState PlayerState

	if !p.IsOnGround {
		// En el aire
		if p.IsTouchingWall && p.Velocity.Y > 0 {
			// Tocando pared y cayendo = wall sliding
			newState = StateWallSliding
		} else if p.Velocity.Y < 0 {
			// Subiendo = jumping
			newState = StateJumping
		} else {
			// Cayendo = falling
			newState = StateFalling
		}
	} else {
		// En el suelo
		if utils.Abs(p.Velocity.X) > 0.5 {
			newState = StateWalking
		} else {
			newState = StateIdle
		}
	}

	// Transicionar si es válido
	if p.State.CanTransitionTo(newState) {
		p.State = newState
	}
}

// GetHitbox retorna el rectángulo de colisión
func (p *Player) GetHitbox() utils.Rectangle {
	return utils.NewRectangle(
		p.Position.X-p.Size.X/2,
		p.Position.Y-p.Size.Y/2,
		p.Size.X,
		p.Size.Y,
	)
}

// GetHurtbox retorna el rectángulo de daño (un poco más pequeño)
func (p *Player) GetHurtbox() utils.Rectangle {
	margin := 5.0
	return utils.NewRectangle(
		p.Position.X-p.Size.X/2+margin,
		p.Position.Y-p.Size.Y/2+margin,
		p.Size.X-margin*2,
		p.Size.Y-margin*2,
	)
}

// regenerateStamina regenera stamina con el tiempo
func (p *Player) regenerateStamina() {
	if p.Stamina < p.MaxStamina {
		// Regenerar más rápido en el suelo
		regenRate := 0.5
		if p.IsOnGround {
			regenRate = 1.0
		}

		p.Stamina += regenRate
		if p.Stamina > p.MaxStamina {
			p.Stamina = p.MaxStamina
		}
	}
}

// TakeDamage recibe daño
func (p *Player) TakeDamage(damage int, knockback utils.Vector2) {
	if p.State == StateDead || p.State == StateDashing {
		return
	}

	p.Health -= damage
	if p.Health <= 0 {
		p.Health = 0
		p.Die()
		return
	}

	// Knockback
	p.Velocity = knockback
	p.State = StateHurt

	// Vibración al recibir daño
	p.controller.Vibrate(200, 0.6)
}

// Die mata al jugador
func (p *Player) Die() {
	p.State = StateDead
	p.Velocity = utils.Zero()
	// TODO: Animación de muerte
}

// Draw dibuja al jugador (placeholder hasta tener sprites)
func (p *Player) Draw(screen *ebiten.Image) {
	hitbox := p.GetHitbox()

	// Color según estado
	bodyColor := p.colorPrimary
	switch p.State {
	case StateDashing:
		bodyColor = config.ColorHeroDash
	case StateAttacking:
		bodyColor = color.RGBA{255, 100, 100, 255}
	case StateHurt:
		bodyColor = color.RGBA{255, 0, 0, 255}
	case StateWallSliding:
		bodyColor = color.RGBA{100, 200, 255, 255}
	}

	// Dibujar cuerpo
	vector.DrawFilledRect(
		screen,
		float32(hitbox.X),
		float32(hitbox.Y),
		float32(hitbox.Width),
		float32(hitbox.Height),
		bodyColor,
		false,
	)

	// Dibujar borde
	vector.StrokeRect(
		screen,
		float32(hitbox.X),
		float32(hitbox.Y),
		float32(hitbox.Width),
		float32(hitbox.Height),
		2,
		p.colorSecondary,
		false,
	)

	// Indicador de dirección (flecha)
	p.drawDirectionIndicator(screen)

	// Indicador de dash disponible
	if p.CanDash {
		p.drawDashIndicator(screen)
	}
}

// drawDirectionIndicator dibuja una flecha indicando la dirección
func (p *Player) drawDirectionIndicator(screen *ebiten.Image) {
	arrowSize := float32(10)
	centerX := float32(p.Position.X)
	centerY := float32(p.Position.Y - 10)

	if p.FacingRight {
		// Flecha derecha
		vector.StrokeLine(screen, centerX, centerY, centerX+arrowSize, centerY, 2, color.White, false)
		vector.StrokeLine(screen, centerX+arrowSize, centerY, centerX+arrowSize-5, centerY-5, 2, color.White, false)
		vector.StrokeLine(screen, centerX+arrowSize, centerY, centerX+arrowSize-5, centerY+5, 2, color.White, false)
	} else {
		// Flecha izquierda
		vector.StrokeLine(screen, centerX, centerY, centerX-arrowSize, centerY, 2, color.White, false)
		vector.StrokeLine(screen, centerX-arrowSize, centerY, centerX-arrowSize+5, centerY-5, 2, color.White, false)
		vector.StrokeLine(screen, centerX-arrowSize, centerY, centerX-arrowSize+5, centerY+5, 2, color.White, false)
	}
}

// drawDashIndicator dibuja un indicador de dash disponible
func (p *Player) drawDashIndicator(screen *ebiten.Image) {
	x := float32(p.Position.X - 20)
	y := float32(p.Position.Y - 35)

	vector.DrawFilledCircle(screen, x, y, 3, color.RGBA{255, 255, 0, 255}, false)
}
