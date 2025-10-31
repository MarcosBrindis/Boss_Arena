// internal/core/game.go
package core

import (
	"fmt"
	"image/color"
	"time"

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
	boss   *entities.Boss // ← NUEVO

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
		200, // X: Empieza a la izquierda
		300,
		controller,
		arena,
	)

	// Crear boss (NUEVO)
	boss := entities.NewBoss(
		1000, // X: Empieza a la derecha
		300,
		arena,
	)
	boss.SetTarget(player) // El boss apunta al jugador

	return &Game{
		config:     cfg,
		controller: controller,
		arena:      arena,
		player:     player,
		boss:       boss, // ← NUEVO

		state:      StatePlaying,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
	}
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

	// Resetear cooldowns del boss
	g.boss.AttackCooldown = 0
	g.boss.SlamCooldown = 0
	g.boss.ChargeCooldown = 0
	g.boss.RoarCooldown = 0

	// Actualizar colores del boss según fase 1

	// Volver a estado jugando
	g.state = StatePlaying
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

	// Actualizar arena
	g.arena.Update()

	// Actualizar jugador
	g.player.Update()

	// Actualizar boss (NUEVO)
	g.boss.Update()

	// Detectar colisiones jugador-boss (NUEVO)
	g.checkPlayerBossCollisions()

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

	// Dibujar según el estado actual
	switch g.state {
	case StateMainMenu:
		g.drawMainMenu(screen)
	case StatePlaying:
		g.drawPlaying(screen)
	case StatePaused:
		g.drawPaused(screen)
	case StateGameOver:
		g.drawGameOver(screen)
	case StateVictory:
		g.drawVictory(screen)
	}

	// Dibujar información de debug
	if g.config.ShowDebugInfo {
		g.drawDebugInfo(screen)
	}

	g.drawDuration = time.Since(start)
}

// Layout define el tamaño lógico de la pantalla
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
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
// SISTEMA DE COLISIONES JUGADOR-BOSS (NUEVO)
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
		damage := g.player.GetAttackDamage()
		if g.boss.TakeDamage(damage) {
			g.controller.Vibrate(100, 0.5)
		}
	}

	// Down Air Attack (MEJORADO)
	downAirHitbox := g.player.GetDownAirAttackHitbox()
	if downAirHitbox != nil && downAirHitbox.Intersects(bossHurtbox) {
		damage := g.player.GetDownAirAttackDamage()
		if g.boss.TakeDamage(damage) {
			// Feedback más fuerte
			g.controller.Vibrate(150, 0.6)

			// POGO EFFECT MEJORADO: Rebote más alto
			g.player.Velocity.Y = -13 // Era -10, ahora igual al salto normal
			g.player.State = entities.StateJumping
			g.player.AttackTimeLeft = 0
			g.player.JumpCount = 1 // Resetear contador de saltos (permite doble salto después)

			// Recuperar stamina como recompensa
			g.player.Stamina += 10
			if g.player.Stamina > g.player.MaxStamina {
				g.player.Stamina = g.player.MaxStamina
			}

			// ========================================================
			// CONTADOR DE POGOS CONSECUTIVOS (NUEVO)
			// ========================================================

			// Incrementar contador de pogos del boss
			g.boss.ConsecutivePogos++

			// Si recibió 3 pogos consecutivos, hacer Slam como contramedida
			if g.boss.ConsecutivePogos >= 3 {
				// Forzar Slam si no está en cooldown
				if g.boss.SlamCooldown == 0 && g.boss.IsOnGround {
					g.boss.NextAction = entities.BossStateSlam
					g.boss.DecisionTimer = 0 // Ejecutar inmediatamente
				}
				g.boss.ConsecutivePogos = 0 // Resetear contador
			}
		}
	}
}

// checkBossAttacksPlayer verifica si el boss está atacando al jugador
func (g *Game) checkBossAttacksPlayer() {
	playerHurtbox := g.player.GetHurtbox()

	// Verificar ataque básico
	attackHitbox := g.boss.GetAttackHitbox()
	if attackHitbox != nil && attackHitbox.Intersects(playerHurtbox) {
		// Calcular knockback
		direction := g.player.Position.Sub(g.boss.Position).Normalize()
		knockback := direction.Mul(8)

		g.player.TakeDamage(g.boss.Damage, knockback)
	}

	// Verificar Slam
	slamHitbox := g.boss.GetSlamHitbox()
	if slamHitbox != nil && slamHitbox.Intersects(playerHurtbox) {
		direction := g.player.Position.Sub(g.boss.Position).Normalize()
		knockback := direction.Mul(12) // Knockback más fuerte

		g.player.TakeDamage(g.boss.Damage*2, knockback)
	}

	// Verificar Charge
	chargeHitbox := g.boss.GetChargeHitbox()
	if chargeHitbox != nil && chargeHitbox.Intersects(playerHurtbox) {
		direction := g.boss.ChargeDirection
		knockback := direction.Mul(8)

		g.player.TakeDamage(g.boss.Damage*2, knockback)
	}

	// =========================================================================
	// DAÑO POR CONTACTO (NUEVO)
	// =========================================================================

	// Si el jugador toca al boss (y el boss no está muerto/aturdido/en transición)
	if g.boss.State != entities.BossStateDead &&
		g.boss.State != entities.BossStateStunned &&
		g.boss.State != entities.BossStateTransition {

		bossHitbox := g.boss.GetHitbox()

		if bossHitbox.Intersects(playerHurtbox) {
			// Daño pequeño por contacto (5 HP)
			contactDamage := 5

			// Knockback suave alejándose del boss
			direction := g.player.Position.Sub(g.boss.Position).Normalize()
			knockback := direction.Mul(6) // Knockback moderado

			g.player.TakeDamage(contactDamage, knockback)
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

	// 2. Dibujar boss (NUEVO - dibujarlo primero para que esté detrás)
	g.boss.Draw(screen)

	// 3. Dibujar jugador
	g.player.Draw(screen)

	// 4. Debug: Hitboxes (NUEVO)
	if g.config.ShowDebugInfo {
		g.drawDebugHitboxes(screen)
	}

	// 5. Mensaje actualizado
	msg := "🐉 Módulo 5: Boss Entity!\n\n"
	msg += "🎮 Controles:\n"
	msg += "  WASD/Stick = Mover\n"
	msg += "  Space/✕    = Saltar\n"
	msg += "  Z/⬜        = Atacar\n"
	msg += "  X/⚪/R2     = Dash\n\n"
	msg += "🐉 Boss:\n"
	msg += "  ✅ 3 Fases (color cambia)\n"
	msg += "  ✅ 4 Ataques diferentes\n"
	msg += "  ✅ IA reactiva\n"
	msg += "  ✅ Aumenta velocidad por fase\n\n"
	msg += "⚔️  ¡Derrota al Titan!"

	ebitenutil.DebugPrintAt(screen, msg, 20, 20)

	// 6. HUD del jugador y boss
	g.drawPlayerHUD(screen)
	g.drawBossHUD(screen) // ← NUEVO
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

	msg := "💀 GAME OVER\n\n"
	msg += "Has sido derrotado por el Titan\n\n"
	msg += "Presiona R (teclado) o\n"
	msg += "△/Y (gamepad) para reintentar"

	ebitenutil.DebugPrintAt(screen, msg, ScreenWidth/2-120, ScreenHeight/2-40)
}

func (g *Game) drawVictory(screen *ebiten.Image) {
	g.drawPlaying(screen)

	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{255, 215, 0, 100})
	screen.DrawImage(overlay, nil)

	msg := "🏆 ¡VICTORIA!\n\n"
	msg += "¡Derrotaste al Titan!\n\n"
	msg += "Módulo 5 completado 🎉\n\n"
	msg += "Presiona R (teclado) o\n"
	msg += "△/Y (gamepad) para jugar otra vez"

	ebitenutil.DebugPrintAt(screen, msg, ScreenWidth/2-130, ScreenHeight/2-50)
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

// drawBossHUD dibuja el HUD del boss (NUEVO)
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

// drawDebugHitboxes dibuja los hitboxes en modo debug (NUEVO)
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

	debugText := fmt.Sprintf(
		"🎮 %s %s\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"FPS: %.1f / TPS: %.1f\n"+
			"Frame: %d\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"PLAYER:\n"+
			"HP: %d/%d\n"+
			"State: %s\n"+
			"Pos: (%.0f, %.0f)\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"BOSS:\n"+
			"HP: %d/%d\n"+
			"Phase: %s\n"+
			"State: %s\n"+
			"Pos: (%.0f, %.0f)\n"+
			"Pogos: %d/3\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"INPUT: %s\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"F3: Toggle Debug\n"+
			"F11: Fullscreen\n"+
			"ESC: Pause",
		GameTitle,
		GameVersion,
		g.fps,
		g.tps,
		g.frame,
		g.player.Health,
		g.player.MaxHealth,
		g.player.State,
		g.player.Position.X,
		g.player.Position.Y,
		g.boss.Health,
		g.boss.MaxHealth,
		g.boss.Phase,
		g.boss.State,
		g.boss.Position.X,
		g.boss.Position.Y,
		g.boss.ConsecutivePogos,
		inputMethod,
	)

	debugBg := ebiten.NewImage(300, 450)
	debugBg.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(debugBg, nil)

	ebitenutil.DebugPrint(screen, debugText)
}
