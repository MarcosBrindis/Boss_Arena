package entities

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/MarcosBrindis/boss-arena-go/internal/config"
	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
	"github.com/MarcosBrindis/boss-arena-go/internal/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Boss representa al jefe final
type Boss struct {
	// Posición y física
	Position utils.Vector2
	Velocity utils.Vector2
	Size     utils.Vector2

	// Estado
	State       BossState
	Phase       BossPhase
	FacingRight bool
	IsOnGround  bool

	// Combate
	Health         int
	MaxHealth      int
	Damage         int
	AttackCooldown int
	AttackRange    float64

	// IA
	Target           *Player
	AggroRange       float64
	AttackDelay      int
	DecisionTimer    int
	NextAction       BossState
	ConsecutivePogos int

	// Ataques especiales
	SlamCooldown    int
	ChargeCooldown  int
	RoarCooldown    int
	SlamDuration    int
	ChargeDuration  int
	RoarDuration    int
	ChargeSpeed     float64
	ChargeDirection utils.Vector2

	// Stun
	StunDuration int
	StunTimeLeft int

	// Transición de fase
	TransitionTimer int
	IsInvulnerable  bool

	// Referencias
	arena *world.Arena
	rng   *rand.Rand

	// Configuración
	config BossConfig

	// Colores
	bodyColor   color.RGBA
	accentColor color.RGBA

	// Disparo (Módulo 7)
	ShootCooldown  int
	ShootDelay     int
	WantsToShoot   bool
	ProjectileType int // Tipo de proyectil a disparar
}

// BossConfig contiene la configuración del boss
type BossConfig struct {
	// Movimiento
	WalkSpeed   float64
	ChargeSpeed float64
	JumpForce   float64

	// Combate
	AttackDamage   int
	AttackRange    float64
	AttackCooldown int

	// Ataques especiales (cooldowns en frames)
	SlamCooldown int
	SlamDuration int
	SlamDamage   int
	SlamRadius   float64

	ChargeCooldown int
	ChargeDuration int
	ChargeDamage   int

	RoarCooldown int
	RoarDuration int
	RoarStunTime int
	RoarRange    float64

	// IA
	AggroRange    float64
	DecisionDelay int

	// Física
	Gravity      float64
	MaxFallSpeed float64
	Friction     float64
}

// DefaultBossConfig retorna la configuración por defecto
func DefaultBossConfig() BossConfig {
	return BossConfig{
		// Movimiento
		WalkSpeed:   2.0,
		ChargeSpeed: 10.0,
		JumpForce:   10.0,

		// Combate
		AttackDamage:   15,
		AttackRange:    80.0,
		AttackCooldown: 60, // 1 segundo

		// Slam (golpe en el suelo)
		SlamCooldown: 180, // 3 segundos
		SlamDuration: 30,  // 0.5 segundos
		SlamDamage:   25,
		SlamRadius:   150.0,

		// Charge (carga)
		ChargeCooldown: 240, // 4 segundos
		ChargeDuration: 60,  // 1 segundo
		ChargeDamage:   30,

		// Roar (rugido)
		RoarCooldown: 300, // 5 segundos
		RoarDuration: 45,  // 0.75 segundos
		RoarStunTime: 60,  // 1 segundo de stun
		RoarRange:    200.0,

		// IA
		AggroRange:    400.0,
		DecisionDelay: 30, // Decide cada 0.5 segundos

		// Física
		Gravity:      0.6,
		MaxFallSpeed: 12.0,
		Friction:     0.9,
	}
}

// NewBoss crea un nuevo boss
func NewBoss(x, y float64, arena *world.Arena) *Boss {
	cfg := DefaultBossConfig()

	return &Boss{
		Position: utils.NewVector2(x, y),
		Velocity: utils.Zero(),
		Size:     utils.NewVector2(100, 120), // Boss más grande que el jugador

		State:       BossStateIdle,
		Phase:       Phase1,
		FacingRight: false, // Empieza mirando a la izquierda

		Health:    1000,
		MaxHealth: 1000,
		Damage:    15,

		AggroRange:  cfg.AggroRange,
		AttackRange: cfg.AttackRange,

		ConsecutivePogos: 0,

		arena:  arena,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
		config: cfg,

		bodyColor:   color.RGBA{255, 69, 0, 255}, // Fase 1
		accentColor: color.RGBA{255, 140, 0, 255},

		ShootCooldown:  0,     // NUEVO
		ShootDelay:     0,     // NUEVO
		WantsToShoot:   false, // NUEVO
		ProjectileType: 0,     // NUEVO
	}
}

// SetTarget establece el objetivo del boss
func (b *Boss) SetTarget(player *Player) {
	b.Target = player
}

// Update actualiza el boss
func (b *Boss) Update() {
	// No actualizar si está muerto
	if b.State == BossStateDead {
		return
	}

	// Actualizar temporizadores
	b.updateTimers()

	// Verificar colisiones
	b.updateCollisionState()

	// Verificar cambio de fase
	b.updatePhase()

	// Procesar IA
	b.updateAI()

	// Aplicar física
	b.applyPhysics()

	// Aplicar movimiento
	b.applyMovement()

	// Actualizar estado
	b.updateState()
}

// updateTimers actualiza todos los temporizadores
func (b *Boss) updateTimers() {
	// Cooldowns de ataques
	if b.AttackCooldown > 0 {
		b.AttackCooldown--
	}
	if b.SlamCooldown > 0 {
		b.SlamCooldown--
	}
	if b.ChargeCooldown > 0 {
		b.ChargeCooldown--
	}
	if b.RoarCooldown > 0 {
		b.RoarCooldown--
	}

	// Duración de ataques
	if b.SlamDuration > 0 {
		b.SlamDuration--
		if b.SlamDuration == 0 {
			b.State = BossStateIdle
		}
	}
	if b.ChargeDuration > 0 {
		b.ChargeDuration--
		if b.ChargeDuration == 0 {
			b.State = BossStateIdle
			b.Velocity.X = 0
		}
	}
	if b.RoarDuration > 0 {
		b.RoarDuration--
		if b.RoarDuration == 0 {
			b.State = BossStateIdle
		}
	}

	// Stun
	if b.StunTimeLeft > 0 {
		b.StunTimeLeft--
		if b.StunTimeLeft == 0 {
			b.State = BossStateIdle
		}
	}

	// Transición de fase
	if b.TransitionTimer > 0 {
		b.TransitionTimer--
		if b.TransitionTimer == 0 {
			b.IsInvulnerable = false
			b.State = BossStateIdle
		}
	}

	// Decisión de IA
	if b.DecisionTimer > 0 {
		b.DecisionTimer--
	}

	// Delay de ataque
	if b.AttackDelay > 0 {
		b.AttackDelay--
	}
	// Cooldown de disparo (NUEVO)
	if b.ShootCooldown > 0 {
		b.ShootCooldown--
	}
	if b.ShootDelay > 0 {
		b.ShootDelay--
	}
}

// updateCollisionState verifica colisiones con el mundo
func (b *Boss) updateCollisionState() {
	hitbox := b.GetHitbox()
	b.IsOnGround = b.arena.IsOnGround(hitbox)

	// ========================================================
	// Resetear contador de pogos si el jugador toca el suelo
	// ========================================================
	if b.Target != nil && b.Target.IsOnGround {
		b.ConsecutivePogos = 0
	}
}

// updatePhase verifica y cambia de fase según la vida
func (b *Boss) updatePhase() {
	healthPercent := float64(b.Health) / float64(b.MaxHealth)

	var newPhase BossPhase

	if healthPercent > 0.66 {
		newPhase = Phase1
	} else if healthPercent > 0.33 {
		newPhase = Phase2
	} else {
		newPhase = Phase3
	}

	// Cambio de fase
	if newPhase != b.Phase {
		b.Phase = newPhase
		b.startPhaseTransition()
	}

	// Actualizar color según fase
	r, g, bl := b.Phase.GetColor()
	b.bodyColor = color.RGBA{r, g, bl, 255}
}

// startPhaseTransition inicia la transición de fase
func (b *Boss) startPhaseTransition() {
	b.State = BossStateTransition
	b.TransitionTimer = 90 // 1.5 segundos
	b.IsInvulnerable = true
	b.Velocity = utils.Zero()

	// Resetear cooldowns en cambio de fase
	b.SlamCooldown = 0
	b.ChargeCooldown = 0
	b.RoarCooldown = 0
}

// updateState actualiza el estado del boss
func (b *Boss) updateState() {
	// Estados que no se pueden interrumpir
	if b.State == BossStateSlam ||
		b.State == BossStateCharge ||
		b.State == BossStateRoar ||
		b.State == BossStateStunned ||
		b.State == BossStateTransition {
		return
	}

	// Determinar nuevo estado
	if !b.IsOnGround {
		if b.Velocity.Y < 0 {
			b.State = BossStateJumping
		} else {
			b.State = BossStateFalling
		}
	} else {
		if utils.Abs(b.Velocity.X) > 0.5 {
			b.State = BossStateWalking
		} else {
			b.State = BossStateIdle
		}
	}
}

// applyPhysics aplica física al boss
func (b *Boss) applyPhysics() {
	// Durante charge, mantener velocidad constante
	if b.State == BossStateCharge {
		b.Velocity = b.ChargeDirection.Mul(b.ChargeSpeed)
		return
	}

	// Durante stun, no aplicar física horizontal
	if b.State == BossStateStunned {
		b.Velocity.X = 0
	}

	// ========================================================================
	// GRAVEDAD (con límites)
	// ========================================================================
	if !b.IsOnGround {
		b.Velocity.Y += b.config.Gravity
		if b.Velocity.Y > b.config.MaxFallSpeed {
			b.Velocity.Y = b.config.MaxFallSpeed
		}

		// IMPORTANTE: Limitar velocidad hacia arriba también (NUEVO)
		if b.Velocity.Y < -b.config.MaxFallSpeed {
			b.Velocity.Y = -b.config.MaxFallSpeed
		}
	} else {
		// Si está en el suelo, resetear velocidad Y
		if b.Velocity.Y > 0 {
			b.Velocity.Y = 0
		}
	}

	// ========================================================================
	// FRICCIÓN
	// ========================================================================
	if b.IsOnGround && b.State != BossStateCharge {
		b.Velocity.X *= b.config.Friction
		if utils.Abs(b.Velocity.X) < 0.1 {
			b.Velocity.X = 0
		}
	}
}

// applyMovement aplica el movimiento con colisiones
func (b *Boss) applyMovement() {
	// ========================================================================
	// MOVIMIENTO HORIZONTAL
	// ========================================================================
	newX := b.Position.X + b.Velocity.X
	testRect := utils.NewRectangle(
		newX-b.Size.X/2,
		b.Position.Y-b.Size.Y/2,
		b.Size.X,
		b.Size.Y,
	)

	collides, _ := b.arena.CheckCollision(testRect)
	if !collides {
		b.Position.X = newX
	} else {
		b.Velocity.X = 0
		// Si está en charge y choca, detenerlo
		if b.State == BossStateCharge {
			b.ChargeDuration = 0
		}
	}

	// ========================================================================
	// MOVIMIENTO VERTICAL
	// ========================================================================
	newY := b.Position.Y + b.Velocity.Y
	testRect = utils.NewRectangle(
		b.Position.X-b.Size.X/2,
		newY-b.Size.Y/2,
		b.Size.X,
		b.Size.Y,
	)

	collides, _ = b.arena.CheckCollision(testRect)
	if !collides {
		b.Position.Y = newY
	} else {
		b.Velocity.Y = 0
	}

	// ========================================================================
	// LÍMITES DE SEGURIDAD (MEJORADOS)
	// ========================================================================

	// Límite inferior (piso)
	floorY := b.arena.GetFloorY()
	if b.Position.Y > floorY+100 {
		b.Position.Y = 300
		b.Velocity = utils.Zero()
	}

	// Límite superior (TECHO - NUEVO) ← IMPORTANTE
	ceilingY := b.Size.Y / 2
	if b.Position.Y < ceilingY {
		b.Position.Y = ceilingY
		b.Velocity.Y = 0 // Detener movimiento hacia arriba
	}

	// Límites laterales
	margin := b.Size.X / 2
	if b.Position.X < margin+60 {
		b.Position.X = margin + 60
		b.Velocity.X = 0
		// Detener charge si choca con pared
		if b.State == BossStateCharge {
			b.ChargeDuration = 0
			b.State = BossStateIdle
		}
	}
	if b.Position.X > 1280-margin-60 {
		b.Position.X = 1280 - margin - 60
		b.Velocity.X = 0
		// Detener charge si choca con pared
		if b.State == BossStateCharge {
			b.ChargeDuration = 0
			b.State = BossStateIdle
		}
	}
}

// GetHitbox retorna el rectángulo de colisión
func (b *Boss) GetHitbox() utils.Rectangle {
	return utils.NewRectangle(
		b.Position.X-b.Size.X/2,
		b.Position.Y-b.Size.Y/2,
		b.Size.X,
		b.Size.Y,
	)
}

// GetHurtbox retorna el rectángulo de daño
func (b *Boss) GetHurtbox() utils.Rectangle {
	margin := 10.0
	return utils.NewRectangle(
		b.Position.X-b.Size.X/2+margin,
		b.Position.Y-b.Size.Y/2+margin,
		b.Size.X-margin*2,
		b.Size.Y-margin*2,
	)
}

// TakeDamage recibe daño
func (b *Boss) TakeDamage(damage int) bool {
	// No recibir daño si es invulnerable
	if b.IsInvulnerable || b.State == BossStateTransition || b.State == BossStateDead {
		return false
	}

	b.Health -= damage
	if b.Health <= 0 {
		b.Health = 0
		b.Die()
		return true
	}

	// Pequeño knockback
	if b.Target != nil {
		direction := b.Position.Sub(b.Target.Position).Normalize()
		b.Velocity.X += direction.X * 2
	}

	return true
}

// Die mata al boss
func (b *Boss) Die() {
	b.State = BossStateDead
	b.Velocity = utils.Zero()
	// TODO: Animación de muerte
}

// Draw dibuja al boss (placeholder)
func (b *Boss) Draw(screen *ebiten.Image) {
	hitbox := b.GetHitbox()

	// Color según estado
	bodyColor := b.bodyColor
	switch b.State {
	case BossStateSlam:
		bodyColor = color.RGBA{255, 0, 0, 255}
	case BossStateCharge:
		bodyColor = color.RGBA{255, 100, 0, 255}
	case BossStateRoar:
		bodyColor = color.RGBA{255, 255, 0, 255}
	case BossStateStunned:
		bodyColor = color.RGBA{100, 100, 255, 255}
	case BossStateTransition:
		// Efecto de parpadeo
		if (b.TransitionTimer/5)%2 == 0 {
			bodyColor = color.RGBA{255, 255, 255, 255}
		}
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
		3,
		color.RGBA{255, 255, 255, 255},
		false,
	)

	// Indicador de dirección
	b.drawDirectionIndicator(screen)

	// Barra de vida individual
	b.drawHealthBar(screen)
}

// drawDirectionIndicator dibuja un indicador de dirección
func (b *Boss) drawDirectionIndicator(screen *ebiten.Image) {
	centerX := float32(b.Position.X)
	centerY := float32(b.Position.Y - 20)
	arrowSize := float32(15)

	arrowColor := color.RGBA{255, 0, 0, 255}

	if b.FacingRight {
		vector.StrokeLine(screen, centerX, centerY, centerX+arrowSize, centerY, 3, arrowColor, false)
		vector.StrokeLine(screen, centerX+arrowSize, centerY, centerX+arrowSize-7, centerY-7, 3, arrowColor, false)
		vector.StrokeLine(screen, centerX+arrowSize, centerY, centerX+arrowSize-7, centerY+7, 3, arrowColor, false)
	} else {
		vector.StrokeLine(screen, centerX, centerY, centerX-arrowSize, centerY, 3, arrowColor, false)
		vector.StrokeLine(screen, centerX-arrowSize, centerY, centerX-arrowSize+7, centerY-7, 3, arrowColor, false)
		vector.StrokeLine(screen, centerX-arrowSize, centerY, centerX-arrowSize+7, centerY+7, 3, arrowColor, false)
	}
}

// drawHealthBar dibuja la barra de vida sobre el boss
func (b *Boss) drawHealthBar(screen *ebiten.Image) {
	barWidth := float32(b.Size.X)
	barHeight := float32(8)
	barX := float32(b.Position.X - b.Size.X/2)
	barY := float32(b.Position.Y - b.Size.Y/2 - 15)

	// Fondo
	vector.DrawFilledRect(screen, barX, barY, barWidth, barHeight, color.RGBA{50, 50, 50, 255}, false)

	// Relleno según salud
	healthPercent := float32(b.Health) / float32(b.MaxHealth)
	fillWidth := barWidth * healthPercent

	var fillColor color.RGBA
	if healthPercent > 0.66 {
		fillColor = config.ColorHPBarFull
	} else if healthPercent > 0.33 {
		fillColor = config.ColorHPBarMid
	} else {
		fillColor = config.ColorHPBarLow
	}

	vector.DrawFilledRect(screen, barX, barY, fillWidth, barHeight, fillColor, false)

	// Borde
	vector.StrokeRect(screen, barX, barY, barWidth, barHeight, 1, color.White, false)
}

// UpdateColor actualiza el color del boss según su fase actual
func (b *Boss) UpdateColor() {
	r, g, bl := b.Phase.GetColor()
	b.bodyColor = color.RGBA{r, g, bl, 255}
}
