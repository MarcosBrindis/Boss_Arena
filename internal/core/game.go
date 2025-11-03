package core

import (
	"fmt"
	"image/color"
	"log"
	"sync"
	"time"

	"github.com/MarcosBrindis/boss-arena-go/internal/audio"
	"github.com/MarcosBrindis/boss-arena-go/internal/combat"
	"github.com/MarcosBrindis/boss-arena-go/internal/effects"
	"github.com/MarcosBrindis/boss-arena-go/internal/entities"
	"github.com/MarcosBrindis/boss-arena-go/internal/input"
	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
	"github.com/MarcosBrindis/boss-arena-go/internal/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Game implementa la interfaz ebiten.Game
type Game struct {
	// Configuración
	config *Config

	// Input System
	controller *input.Controller

	// World
	arena *world.Arena

	// Entities
	player *entities.Player
	boss   *entities.Boss

	// Combat System (NUEVO)
	eventSystem   *combat.EventSystem
	damageCalc    *combat.DamageCalculator
	effectManager *combat.EffectManager

	// Visual Effects (NUEVO)
	particleSystem *effects.ParticleSystem
	screenShake    *effects.ScreenShake
	hitStop        *effects.HitStop

	// Audio System (NUEVO)
	soundSystem *audio.SoundSystem

	// Estado del juego
	state     GameState
	frame     uint64
	startTime time.Time

	// Performance tracking
	updateDuration time.Duration
	drawDuration   time.Duration
	tps            float64
	fps            float64

	// Control de tiempo
	lastUpdate time.Time
	deltaTime  float64

	// Banderas
	isPaused bool

	// Control de input global
	escapeKeyPressedLastFrame bool
	f11KeyPressedLastFrame    bool
	f3KeyPressedLastFrame     bool
}

// NewGame crea una nueva instancia del juego
func NewGame() *Game {
	cfg := DefaultConfig()

	// Crear arena
	arena := world.NewArena(ScreenWidth, ScreenHeight)

	// Crear controller
	controller := input.NewController(
		cfg.GamepadDeadzone,
		cfg.JumpBufferFrames,
		cfg.CoyoteTimeFrames,
	)

	// Crear jugador
	player := entities.NewPlayer(
		200,
		300,
		controller,
		arena,
	)

	// Crear boss
	boss := entities.NewBoss(
		1000,
		300,
		arena,
	)
	boss.SetTarget(player)

	// ========================================================================
	// CREAR SISTEMAS DE COMBATE (NUEVO)
	// ========================================================================

	// Event System (con buffer de 100 eventos)
	eventSystem := combat.NewEventSystem(100)
	eventSystem.Start() // ← Inicia la goroutine

	// Damage Calculator
	damageCalc := combat.NewDamageCalculator()

	// Effect Manager
	effectManager := combat.NewEffectManager(50)

	// Particle System
	particleSystem := effects.NewParticleSystem(200)

	// Screen Shake
	screenShake := effects.NewScreenShake()

	// Hit Stop
	hitStop := effects.NewHitStop()

	// Sound System
	soundSystem := audio.NewSoundSystem()

	game := &Game{
		config:     cfg,
		controller: controller,
		arena:      arena,
		player:     player,
		boss:       boss,

		// Combat Systems (NUEVO)
		eventSystem:    eventSystem,
		damageCalc:     damageCalc,
		effectManager:  effectManager,
		particleSystem: particleSystem,
		screenShake:    screenShake,
		hitStop:        hitStop,
		soundSystem:    soundSystem,

		state:      StatePlaying,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
	}

	// ========================================================================
	// REGISTRAR LISTENERS DE EVENTOS (NUEVO)
	// ========================================================================
	game.setupEventListeners()

	return game
}

// setupEventListeners registra listeners para eventos de combate
func (g *Game) setupEventListeners() {
	// Listener: Cuando se hace daño
	g.eventSystem.AddListener(combat.EventDamageDealt, func(event combat.CombatEvent) {
		// Spawn partículas en la posición del impacto
		particleColor := color.RGBA{255, 100, 100, 255}
		if event.IsCritical {
			particleColor = color.RGBA{255, 255, 0, 255} // Amarillo para críticos
		}
		g.particleSystem.Emit(event.Position, 5, particleColor)

		// Screen shake según el daño
		shakeIntensity := float64(event.Damage) / 10.0
		g.screenShake.Start(shakeIntensity, 5)

		// Hit stop para críticos
		if event.IsCritical {
			g.hitStop.Start(3)
		}

		// Sonido
		g.soundSystem.PlaySound(audio.SoundHit)
	})

	// Listener: Cuando aumenta el combo
	g.eventSystem.AddListener(combat.EventComboIncreased, func(event combat.CombatEvent) {
		// Efecto visual de combo
		g.effectManager.SpawnEffect(
			combat.EffectSlash,
			event.Position,
			color.RGBA{255, 255, 0, 255},
		)
	})

	// Listener: Cuando mata al boss
	g.eventSystem.AddListener(combat.EventKill, func(event combat.CombatEvent) {
		if event.Target == "boss" {
			// Explosión grande
			g.particleSystem.Emit(event.Position, 30, color.RGBA{255, 140, 0, 255})
			g.screenShake.Start(20, 30)
			g.soundSystem.PlaySound(audio.SoundExplosion)
		}
	})
}

// Update actualiza la lógica del juego
func (g *Game) Update() error {
	start := time.Now()

	// Actualizar controller
	g.controller.Update()

	// Calcular delta time
	now := time.Now()
	g.deltaTime = now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	// Incrementar contador de frames
	g.frame++

	// Calcular TPS/FPS
	if g.frame%60 == 0 {
		g.tps = ebiten.ActualTPS()
		g.fps = ebiten.ActualFPS()
	}

	// Manejar input global
	if err := g.handleGlobalInput(); err != nil {
		return err
	}

	// Si está pausado, no actualizar lógica
	if g.isPaused {
		g.updateDuration = time.Since(start)
		return nil
	}

	// ========================================================================
	// HIT STOP: Si está activo, congelar el juego
	// ========================================================================
	g.hitStop.Update()
	if g.hitStop.ShouldFreeze() {
		g.updateDuration = time.Since(start)
		return nil // No actualizar nada durante freeze frame
	}

	// ========================================================================
	// ACTUALIZAR SISTEMAS
	// ========================================================================

	// Actualizar arena
	g.arena.Update()

	// Actualizar jugador
	g.player.Update()

	// Actualizar boss
	g.boss.Update()

	// Detectar colisiones jugador-boss
	g.checkPlayerBossCollisions()

	// Actualizar efectos visuales (NUEVO)
	g.effectManager.Update()
	g.particleSystem.Update()
	g.screenShake.Update()

	// Actualizar según el estado actual
	switch g.state {
	case StateMainMenu:
		g.updateMainMenu()
	case StatePlaying:
		g.updatePlaying()
	case StatePaused:
		g.updatePaused()
	case StateGameOver:
		g.updateGameOver()
	case StateVictory:
		g.updateVictory()
	}

	g.updateDuration = time.Since(start)
	return nil
}

// Draw dibuja el juego en pantalla
func (g *Game) Draw(screen *ebiten.Image) {
	start := time.Now()

	// ========================================================================
	// APLICAR SCREEN SHAKE (NUEVO)
	// ========================================================================
	shakeOffset := g.screenShake.GetOffset()

	// Crear imagen temporal con offset
	tempScreen := screen

	// Si hay shake, aplicar offset
	if g.screenShake.IsActive() {
		// Nota: En una implementación completa, aplicaríamos el offset
		// a una cámara. Por ahora es visual placeholder.
		_ = shakeOffset
	}

	// Dibujar según el estado actual
	switch g.state {
	case StateMainMenu:
		g.drawMainMenu(tempScreen)
	case StatePlaying:
		g.drawPlaying(tempScreen)
	case StatePaused:
		g.drawPaused(tempScreen)
	case StateGameOver:
		g.drawGameOver(tempScreen)
	case StateVictory:
		g.drawVictory(tempScreen)
	}

	// Dibujar información de debug
	if g.config.ShowDebugInfo {
		g.drawDebugInfo(tempScreen)
	}

	g.drawDuration = time.Since(start)
}

// Layout define el tamaño lógico de la pantalla
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// RestartGame reinicia el juego a su estado inicial
func (g *Game) RestartGame() {
	// Resetear jugador
	g.player.Position = utils.NewVector2(200, 300)
	g.player.Velocity = utils.Zero()
	g.player.Health = g.player.MaxHealth
	g.player.Stamina = g.player.MaxStamina
	g.player.State = entities.StateIdle
	g.player.CanDash = true
	g.player.JumpCount = 0

	// Resetear boss
	g.boss.Position = utils.NewVector2(1000, 300)
	g.boss.Velocity = utils.Zero()
	g.boss.Health = g.boss.MaxHealth
	g.boss.State = entities.BossStateIdle
	g.boss.Phase = entities.Phase1
	g.boss.IsInvulnerable = false
	g.boss.ConsecutivePogos = 0

	// Resetear cooldowns del boss
	g.boss.AttackCooldown = 0
	g.boss.SlamCooldown = 0
	g.boss.ChargeCooldown = 0
	g.boss.RoarCooldown = 0

	// Actualizar colores del boss según fase 1
	g.boss.UpdateColor()

	// Limpiar efectos (NUEVO)
	g.particleSystem.Clear()
	g.effectManager.Clear()
	g.screenShake.Stop()

	// Resetear estadísticas (NUEVO)
	g.eventSystem.ResetStats()

	// Volver a estado jugando
	g.state = StatePlaying
}

// Cleanup limpia recursos al cerrar el juego
func (g *Game) Cleanup() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		g.eventSystem.Stop()
	}()

	// Esperar que todas las goroutines terminen
	wg.Wait()

	log.Println("✅ Todas las goroutines cerradas correctamente")
}

// ============================================================================
// MÉTODOS DE UPDATE POR ESTADO
// ============================================================================

func (g *Game) handleGlobalInput() error {
	escapePressed := ebiten.IsKeyPressed(ebiten.KeyEscape)
	if escapePressed && !g.escapeKeyPressedLastFrame {
		g.isPaused = !g.isPaused
	}
	g.escapeKeyPressedLastFrame = escapePressed

	f11Pressed := ebiten.IsKeyPressed(ebiten.KeyF11)
	if f11Pressed && !g.f11KeyPressedLastFrame {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	g.f11KeyPressedLastFrame = f11Pressed

	f3Pressed := ebiten.IsKeyPressed(ebiten.KeyF3)
	if f3Pressed && !g.f3KeyPressedLastFrame {
		g.config.ShowDebugInfo = !g.config.ShowDebugInfo
	}
	g.f3KeyPressedLastFrame = f3Pressed

	return nil
}

func (g *Game) updateMainMenu() {
	// TODO
}

func (g *Game) updatePlaying() {
	// Verificar victoria
	if g.boss.State == entities.BossStateDead && g.state != StateVictory {
		g.state = StateVictory

		// Emitir evento de victoria
		g.eventSystem.EmitEvent(combat.CombatEvent{
			Type:     combat.EventKill,
			Target:   "boss",
			Attacker: "player",
			Position: g.boss.Position,
		})
	}

	// Verificar derrota
	if g.player.State == entities.StateDead && g.state != StateGameOver {
		g.state = StateGameOver
	}
}

func (g *Game) updatePaused() {
	// TODO
}

func (g *Game) updateGameOver() {
	// Reiniciar con R (teclado) o Start (gamepad)
	if ebiten.IsKeyPressed(ebiten.KeyR) || g.controller.IsSpecialPressed() {
		g.RestartGame()
	}
}

func (g *Game) updateVictory() {
	// Reiniciar con R (teclado) o Start (gamepad)
	if ebiten.IsKeyPressed(ebiten.KeyR) || g.controller.IsSpecialPressed() {
		g.RestartGame()
	}
}

// ============================================================================
// SISTEMA DE COLISIONES JUGADOR-BOSS
// ============================================================================

func (g *Game) checkPlayerBossCollisions() {
	// No verificar si alguno está muerto
	if g.player.State == entities.StateDead || g.boss.State == entities.BossStateDead {
		return
	}

	// 1. Verificar si el jugador golpea al boss
	g.checkPlayerAttacksBoss()

	// 2. Verificar si el boss golpea al jugador
	g.checkBossAttacksPlayer()
}

// checkPlayerAttacksBoss verifica si el jugador está atacando al boss
func (g *Game) checkPlayerAttacksBoss() {
	bossHurtbox := g.boss.GetHurtbox()

	// Ataque normal
	attackHitbox := g.player.GetAttackHitbox()
	if attackHitbox != nil && attackHitbox.Intersects(bossHurtbox) {
		baseDamage := g.player.GetAttackDamage()

		// Calcular daño con sistema nuevo
		isCritical := g.damageCalc.RollCritical(0.15) // 15% de crítico
		comboMultiplier := 1.0 + float64(g.player.ComboCount)*0.1
		damage := g.damageCalc.CalculateDamage(baseDamage, combat.DamagePhysical, isCritical, comboMultiplier)

		if g.boss.TakeDamage(damage) {
			g.controller.Vibrate(100, 0.5)

			// Emitir evento de daño
			g.eventSystem.EmitEvent(combat.CombatEvent{
				Type:       combat.EventDamageDealt,
				Damage:     damage,
				Position:   g.boss.Position,
				Attacker:   "player",
				Target:     "boss",
				IsCritical: isCritical,
				ComboCount: g.player.ComboCount,
			})
		}
	}

	// Down Air Attack (MEJORADO)
	downAirHitbox := g.player.GetDownAirAttackHitbox()
	if downAirHitbox != nil && downAirHitbox.Intersects(bossHurtbox) {
		baseDamage := g.player.GetDownAirAttackDamage()

		// Calcular daño
		isCritical := g.damageCalc.RollCritical(0.25) // 25% de crítico para pogo
		damage := g.damageCalc.CalculateDamage(baseDamage, combat.DamagePhysical, isCritical, 1.0)

		if g.boss.TakeDamage(damage) {
			// Feedback más fuerte
			g.controller.Vibrate(150, 0.6)

			// POGO EFFECT MEJORADO
			g.player.Velocity.Y = -13
			g.player.State = entities.StateJumping
			g.player.AttackTimeLeft = 0
			g.player.JumpCount = 1

			// Recuperar stamina
			g.player.Stamina += 10
			if g.player.Stamina > g.player.MaxStamina {
				g.player.Stamina = g.player.MaxStamina
			}

			// Contador de pogos consecutivos
			g.boss.ConsecutivePogos++

			if g.boss.ConsecutivePogos >= 3 {
				if g.boss.SlamCooldown == 0 && g.boss.IsOnGround {
					g.boss.NextAction = entities.BossStateSlam
					g.boss.DecisionTimer = 0
				}
				g.boss.ConsecutivePogos = 0
			}

			// Emitir evento (NUEVO)
			g.eventSystem.EmitEvent(combat.CombatEvent{
				Type:       combat.EventDamageDealt,
				Damage:     damage,
				Position:   g.boss.Position,
				Attacker:   "player",
				Target:     "boss",
				IsCritical: isCritical,
			})
		}
	}
}

// checkBossAttacksPlayer verifica si el boss está atacando al jugador
func (g *Game) checkBossAttacksPlayer() {
	playerHurtbox := g.player.GetHurtbox()

	// Verificar ataque básico
	attackHitbox := g.boss.GetAttackHitbox()
	if attackHitbox != nil && attackHitbox.Intersects(playerHurtbox) {
		direction := g.player.Position.Sub(g.boss.Position).Normalize()
		knockback := direction.Mul(8)

		g.player.TakeDamage(g.boss.Damage, knockback)

		// Emitir evento
		g.eventSystem.EmitEvent(combat.CombatEvent{
			Type:     combat.EventDamageDealt,
			Damage:   g.boss.Damage,
			Position: g.player.Position,
			Attacker: "boss",
			Target:   "player",
		})
	}

	// Verificar Slam
	slamHitbox := g.boss.GetSlamHitbox()
	if slamHitbox != nil && slamHitbox.Intersects(playerHurtbox) {
		direction := g.player.Position.Sub(g.boss.Position).Normalize()
		knockback := direction.Mul(12)

		g.player.TakeDamage(g.boss.Damage*2, knockback)

		// Emitir evento
		g.eventSystem.EmitEvent(combat.CombatEvent{
			Type:     combat.EventDamageDealt,
			Damage:   g.boss.Damage * 2,
			Position: g.player.Position,
			Attacker: "boss",
			Target:   "player",
		})
	}

	// Verificar Charge
	chargeHitbox := g.boss.GetChargeHitbox()
	if chargeHitbox != nil && chargeHitbox.Intersects(playerHurtbox) {
		direction := g.boss.ChargeDirection
		knockback := direction.Mul(8)

		g.player.TakeDamage(g.boss.Damage*2, knockback)

		// Emitir evento
		g.eventSystem.EmitEvent(combat.CombatEvent{
			Type:     combat.EventDamageDealt,
			Damage:   g.boss.Damage * 2,
			Position: g.player.Position,
			Attacker: "boss",
			Target:   "player",
		})
	}

	// Verificar contacto
	if g.boss.State != entities.BossStateDead &&
		g.boss.State != entities.BossStateStunned &&
		g.boss.State != entities.BossStateTransition {

		bossHitbox := g.boss.GetHitbox()

		if bossHitbox.Intersects(playerHurtbox) {
			contactDamage := 5

			direction := g.player.Position.Sub(g.boss.Position).Normalize()
			knockback := direction.Mul(6)

			g.player.TakeDamage(contactDamage, knockback)

			// Emitir evento
			g.eventSystem.EmitEvent(combat.CombatEvent{
				Type:     combat.EventDamageDealt,
				Damage:   contactDamage,
				Position: g.player.Position,
				Attacker: "boss",
				Target:   "player",
			})
		}
	}
}

// ============================================================================
// MÉTODOS DE DRAW POR ESTADO
// ============================================================================

func (g *Game) drawMainMenu(screen *ebiten.Image) {
	// TODO
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	// 1. Dibujar arena
	g.arena.Draw(screen)

	// 2. Dibujar boss
	g.boss.Draw(screen)

	// 3. Dibujar jugador
	g.player.Draw(screen)

	// 4. Dibujar efectos visuales (NUEVO)
	g.drawVisualEffects(screen)

	// 5. Debug: Hitboxes
	if g.config.ShowDebugInfo {
		g.drawDebugHitboxes(screen)
	}

	// 6. Mensaje
	msg := "⚔️ Módulo 6: Combat System!\n\n"
	msg += "🎮 Controles:\n"
	msg += "  WASD/Stick = Mover\n"
	msg += "  Space/✕    = Saltar\n"
	msg += "  Z/⬜        = Atacar\n"
	msg += "  X/⚪/R2     = Dash\n"
	msg += "  Down+Z     = Pogo\n\n"
	msg += "✨ NUEVO:\n"
	msg += "  💥 Partículas de impacto\n"
	msg += "  🌟 Screen shake\n"
	msg += "  ⏸️  Hit stop (freeze frame)\n"
	msg += "  🎯 Sistema de críticos\n"
	msg += "  📊 Estadísticas en tiempo real\n"
	msg += "  🔥 Eventos con concurrencia"

	ebitenutil.DebugPrintAt(screen, msg, 20, 20)

	// 7. HUD
	g.drawPlayerHUD(screen)
	g.drawBossHUD(screen)
	g.drawStatsHUD(screen)
}

func (g *Game) drawPaused(screen *ebiten.Image) {
	g.drawPlaying(screen)

	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(overlay, nil)

	msg := "⏸️  PAUSA\n\nPresiona ESC para continuar"
	ebitenutil.DebugPrintAt(screen, msg, ScreenWidth/2-100, ScreenHeight/2-20)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	g.drawPlaying(screen)

	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 200})
	screen.DrawImage(overlay, nil)

	stats := g.eventSystem.GetStats()

	msg := fmt.Sprintf(
		"💀 GAME OVER\n\n"+
			"Daño hecho: %d\n"+
			"Daño recibido: %d\n"+
			"Combo máximo: %d\n"+
			"Críticos: %d\n\n"+
			"Presiona R (teclado) o\n"+
			"△/Y (gamepad) para reintentar",
		stats.PlayerDamageDealt,
		stats.PlayerDamageTaken,
		stats.HighestCombo,
		stats.CriticalHits,
	)

	ebitenutil.DebugPrintAt(screen, msg, ScreenWidth/2-120, ScreenHeight/2-60)
}

func (g *Game) drawVictory(screen *ebiten.Image) {
	g.drawPlaying(screen)

	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{255, 215, 0, 100})
	screen.DrawImage(overlay, nil)

	stats := g.eventSystem.GetStats()

	msg := fmt.Sprintf(
		"🏆 ¡VICTORIA!\n\n"+
			"Daño total: %d\n"+
			"Daño recibido: %d\n"+
			"Combo máximo: %d\n"+
			"Críticos: %d\n"+
			"Precisión: %.1f%%\n\n"+
			"Módulo 6 completado 🎉\n\n"+
			"Presiona R (teclado) o\n"+
			"△/Y (gamepad) para jugar otra vez",
		stats.PlayerDamageDealt,
		stats.PlayerDamageTaken,
		stats.HighestCombo,
		stats.CriticalHits,
		float64(stats.PlayerAttacksLanded)/float64(stats.PlayerAttacksLanded+stats.PlayerAttacksMissed)*100,
	)

	ebitenutil.DebugPrintAt(screen, msg, ScreenWidth/2-130, ScreenHeight/2-80)
}

// ============================================================================
// DIBUJO DE EFECTOS VISUALES
// ============================================================================

func (g *Game) drawVisualEffects(screen *ebiten.Image) {
	// Dibujar partículas
	for _, particle := range g.particleSystem.GetParticles() {
		vector.DrawFilledCircle(
			screen,
			float32(particle.Position.X),
			float32(particle.Position.Y),
			float32(particle.Size),
			particle.Color,
			false,
		)
	}

	// Dibujar efectos de combate
	for _, effect := range g.effectManager.GetActiveEffects() {
		alpha := uint8(255 * (1.0 - float64(effect.Age)/float64(effect.Lifetime)))
		effectColor := effect.Color
		effectColor.A = alpha

		switch effect.Type {
		case combat.EffectHitSpark:
			vector.DrawFilledCircle(
				screen,
				float32(effect.Position.X),
				float32(effect.Position.Y),
				float32(effect.Size),
				effectColor,
				false,
			)

		case combat.EffectSlash:
			vector.StrokeRect(
				screen,
				float32(effect.Position.X-effect.Size/2),
				float32(effect.Position.Y-effect.Size/2),
				float32(effect.Size),
				float32(effect.Size),
				3,
				effectColor,
				false,
			)
		}
	}
}

// ============================================================================
// HUD
// ============================================================================

func (g *Game) drawPlayerHUD(screen *ebiten.Image) {
	hudX := float32(20)
	hudY := float32(ScreenHeight - 100)

	// Fondo del HUD
	hudBg := ebiten.NewImage(300, 80)
	hudBg.Fill(color.RGBA{0, 0, 0, 150})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(hudX), float64(hudY))
	screen.DrawImage(hudBg, op)

	// Barra de salud
	g.drawBar(screen, hudX+10, hudY+10, 280, 20,
		float64(g.player.Health)/float64(g.player.MaxHealth),
		g.getHealthColor(float64(g.player.Health)/float64(g.player.MaxHealth)),
		"PLAYER HP")

	// Barra de stamina
	staminaPercent := g.player.Stamina / g.player.MaxStamina
	staminaColor := color.RGBA{0, 200, 255, 255}

	if staminaPercent < 0.3 {
		staminaColor = color.RGBA{255, 100, 100, 255}
	} else if staminaPercent < 0.5 {
		staminaColor = color.RGBA{255, 200, 0, 255}
	}

	g.drawBar(screen, hudX+10, hudY+40, 280, 15,
		staminaPercent,
		staminaColor,
		"STAMINA")

	// Info de combo
	if g.player.ComboCount > 0 {
		comboText := fmt.Sprintf("COMBO x%d", g.player.ComboCount)
		ebitenutil.DebugPrintAt(screen, comboText, int(hudX+200), int(hudY+60))
	}

	// Advertencia de stamina baja
	if staminaPercent < 0.2 {
		warningText := "⚠️ STAMINA BAJA"
		ebitenutil.DebugPrintAt(screen, warningText, int(hudX+10), int(hudY+60))
	}
}

func (g *Game) drawBossHUD(screen *ebiten.Image) {
	// Barra grande en la parte superior
	barWidth := float32(600)
	barHeight := float32(30)
	barX := float32(ScreenWidth/2) - barWidth/2
	barY := float32(10)

	// Fondo
	hudBg := ebiten.NewImage(int(barWidth+20), int(barHeight+40))
	hudBg.Fill(color.RGBA{0, 0, 0, 180})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(barX-10), float64(barY-10))
	screen.DrawImage(hudBg, op)

	// Nombre del boss
	bossName := fmt.Sprintf("🐉 TITAN - %s", g.boss.Phase.String())
	ebitenutil.DebugPrintAt(screen, bossName, int(barX), int(barY))

	// Barra de vida
	healthPercent := float64(g.boss.Health) / float64(g.boss.MaxHealth)
	g.drawBar(screen, barX, barY+20, barWidth, barHeight,
		healthPercent,
		g.getHealthColor(healthPercent),
		"")

	// Texto de HP
	hpText := fmt.Sprintf("%d / %d", g.boss.Health, g.boss.MaxHealth)
	ebitenutil.DebugPrintAt(screen, hpText, int(barX+barWidth/2-30), int(barY+27))
}

// drawStatsHUD dibuja estadísticas en tiempo real (NUEVO)
func (g *Game) drawStatsHUD(screen *ebiten.Image) {
	stats := g.eventSystem.GetStats()

	hudX := float32(ScreenWidth - 250)
	hudY := float32(ScreenHeight - 150)

	// Fondo
	hudBg := ebiten.NewImage(230, 140)
	hudBg.Fill(color.RGBA{0, 0, 0, 150})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(hudX), float64(hudY))
	screen.DrawImage(hudBg, op)

	// Estadísticas
	statsText := fmt.Sprintf(
		"📊 ESTADÍSTICAS\n\n"+
			"Daño hecho: %d\n"+
			"Golpes: %d\n"+
			"Combo máx: %d\n"+
			"Críticos: %d\n"+
			"Eventos: %d",
		stats.PlayerDamageDealt,
		stats.PlayerAttacksLanded,
		stats.HighestCombo,
		stats.CriticalHits,
		stats.TotalEvents,
	)

	ebitenutil.DebugPrintAt(screen, statsText, int(hudX+10), int(hudY+10))
}

func (g *Game) drawBar(screen *ebiten.Image, x, y, width, height float32, fill float64, col color.RGBA, label string) {
	// Fondo
	vector.DrawFilledRect(screen, x, y, width, height, color.RGBA{50, 50, 50, 255}, false)

	// Relleno
	fillWidth := float32(fill) * width
	vector.DrawFilledRect(screen, x, y, fillWidth, height, col, false)

	// Borde
	vector.StrokeRect(screen, x, y, width, height, 2, color.White, false)

	// Label
	if label != "" {
		ebitenutil.DebugPrintAt(screen, label, int(x), int(y-15))
	}
}

func (g *Game) getHealthColor(percent float64) color.RGBA {
	if percent > 0.66 {
		return color.RGBA{0, 255, 0, 255}
	} else if percent > 0.33 {
		return color.RGBA{255, 200, 0, 255}
	} else {
		return color.RGBA{255, 0, 0, 255}
	}
}

func (g *Game) drawDebugHitboxes(screen *ebiten.Image) {
	// Hitbox de ataque del jugador
	playerAttack := g.player.GetAttackHitbox()
	if playerAttack != nil {
		vector.StrokeRect(
			screen,
			float32(playerAttack.X),
			float32(playerAttack.Y),
			float32(playerAttack.Width),
			float32(playerAttack.Height),
			2,
			color.RGBA{0, 255, 0, 150},
			false,
		)
	}

	// Hitbox de ataque del boss
	bossAttack := g.boss.GetAttackHitbox()
	if bossAttack != nil {
		vector.StrokeRect(
			screen,
			float32(bossAttack.X),
			float32(bossAttack.Y),
			float32(bossAttack.Width),
			float32(bossAttack.Height),
			2,
			color.RGBA{255, 0, 0, 150},
			false,
		)
	}

	// Hitbox de slam
	slamHitbox := g.boss.GetSlamHitbox()
	if slamHitbox != nil {
		vector.StrokeRect(
			screen,
			float32(slamHitbox.X),
			float32(slamHitbox.Y),
			float32(slamHitbox.Width),
			float32(slamHitbox.Height),
			2,
			color.RGBA{255, 100, 0, 150},
			false,
		)
	}
}

// ============================================================================
// DEBUG INFO
// ============================================================================

func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	inputMethod := "Keyboard"
	if g.controller.IsGamepadConnected() {
		inputMethod = "Gamepad: " + g.controller.GetGamepadName()
	}

	stats := g.eventSystem.GetStats()

	debugText := fmt.Sprintf(
		"🎮 %s %s\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"FPS: %.1f / TPS: %.1f\n"+
			"Frame: %d\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"PLAYER:\n"+
			"HP: %d/%d\n"+
			"State: %s\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"BOSS:\n"+
			"HP: %d/%d\n"+
			"Phase: %s\n"+
			"State: %s\n"+
			"Pogos: %d/3\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"COMBAT:\n"+
			"Partículas: %d\n"+
			"Eventos: %d\n"+
			"Screen Shake: %v\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"INPUT: %s",
		GameTitle,
		GameVersion,
		g.fps,
		g.tps,
		g.frame,
		g.player.Health,
		g.player.MaxHealth,
		g.player.State,
		g.boss.Health,
		g.boss.MaxHealth,
		g.boss.Phase,
		g.boss.State,
		g.boss.ConsecutivePogos,
		len(g.particleSystem.GetParticles()),
		stats.TotalEvents,
		g.screenShake.IsActive(),
		inputMethod,
	)

	debugBg := ebiten.NewImage(300, 500)
	debugBg.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(debugBg, nil)

	ebitenutil.DebugPrint(screen, debugText)
}
