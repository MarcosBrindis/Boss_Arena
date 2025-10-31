// internal/core/game.go
package core

import (
	"fmt"
	"image/color"
	"time"

	"github.com/MarcosBrindis/boss-arena-go/internal/entities" // ← NUEVO
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

	// Entities (NUEVO)
	player *entities.Player

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

	// Crear jugador (NUEVO)
	player := entities.NewPlayer(
		float64(ScreenWidth/2),
		300,
		controller,
		arena,
	)

	return &Game{
		config:     cfg,
		controller: controller,
		arena:      arena,
		player:     player, // ← NUEVO

		state:      StatePlaying,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
	}
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

	// Actualizar jugador (NUEVO)
	g.player.Update()

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
	// La lógica ya se maneja en player.Update()
}

func (g *Game) updatePaused() {
	// TODO
}

func (g *Game) updateGameOver() {
	// TODO
}

func (g *Game) updateVictory() {
	// TODO
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

	// 2. Dibujar jugador (NUEVO)
	g.player.Draw(screen)

	// 2.5 DEBUG: Dibujar área de detección de paredes (NUEVO)
	if g.config.ShowDebugInfo {
		hitbox := g.player.GetHitbox()
		margin := 8.0

		// Área de detección izquierda (rojo)
		testRectLeft := utils.NewRectangle(
			hitbox.X-margin,
			hitbox.Y+5,
			hitbox.Width+margin,
			hitbox.Height-10,
		)
		vector.StrokeRect(
			screen,
			float32(testRectLeft.X),
			float32(testRectLeft.Y),
			float32(testRectLeft.Width),
			float32(testRectLeft.Height),
			2,
			color.RGBA{255, 0, 0, 150},
			false,
		)

		// Área de detección derecha (verde)
		testRectRight := utils.NewRectangle(
			hitbox.X,
			hitbox.Y+5,
			hitbox.Width+margin,
			hitbox.Height-10,
		)
		vector.StrokeRect(
			screen,
			float32(testRectRight.X),
			float32(testRectRight.Y),
			float32(testRectRight.Width),
			float32(testRectRight.Height),
			2,
			color.RGBA{0, 255, 0, 150},
			false,
		)
	}

	// 3. Dibujar hitbox de ataque (debug)
	if g.config.ShowDebugInfo {
		attackHitbox := g.player.GetAttackHitbox()
		if attackHitbox != nil {
			// Dibujar hitbox de ataque en rojo semi-transparente
			vector.StrokeRect(
				screen,
				float32(attackHitbox.X),
				float32(attackHitbox.Y),
				float32(attackHitbox.Width),
				float32(attackHitbox.Height),
				2,
				color.RGBA{255, 0, 0, 150},
				false,
			)
		}
	}

	// 4. Mensaje actualizado
	msg := "🎮 Módulo 4: Player completado!\n\n"
	msg += "Controles:\n"
	msg += "  WASD/Stick = Mover\n"
	msg += "  Space/✕    = Saltar (doble salto)\n"
	msg += "  Z/⬜        = Atacar (combo x3) 💥\n"
	msg += "  X/⚪/R2     = Dash 🚀\n\n"
	msg += "Mecánicas:\n"
	msg += "  ✅ Salto doble (20 stamina)\n"
	msg += "  ✅ Wall-climb (15 stamina/salto) 🧗\n"
	msg += "  ✅ Wall-slide (gratis)\n"
	msg += "  ✅ Dash 360° (25 stamina)\n"
	msg += "  ✅ Combo x3 (10/15/20 stamina)\n"
	msg += "  ✅ Coyote time\n"
	msg += "  ✅ Jump buffer\n\n"
	msg += "💡 Stamina se regenera más rápido en el suelo\n"
	msg += "⏭️  Esperando Módulo 5 (Boss)..."

	ebitenutil.DebugPrintAt(screen, msg, 20, 20)

	// 5. HUD del jugador
	g.drawPlayerHUD(screen)
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
	// TODO
}

func (g *Game) drawVictory(screen *ebiten.Image) {
	// TODO
}

// ============================================================================
// HUD DEL JUGADOR
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
		color.RGBA{255, 0, 0, 255},
		"HP")

	// Barra de stamina (cambia de color si está baja)
	staminaPercent := g.player.Stamina / g.player.MaxStamina
	staminaColor := color.RGBA{0, 200, 255, 255}

	if staminaPercent < 0.3 {
		// Stamina baja = color rojo
		staminaColor = color.RGBA{255, 100, 100, 255}
	} else if staminaPercent < 0.5 {
		// Stamina media = color amarillo
		staminaColor = color.RGBA{255, 200, 0, 255}
	}

	g.drawBar(screen, hudX+10, hudY+40, 280, 15,
		staminaPercent,
		staminaColor,
		"STAMINA")

	// Info de combo
	// Info de combo
	if g.player.ComboCount > 0 {
		comboText := fmt.Sprintf("COMBO x%d", g.player.ComboCount)

		// Dibujar texto de combo
		ebitenutil.DebugPrintAt(screen, comboText, int(hudX+200), int(hudY+60))
	}

	// Advertencia de stamina baja
	if staminaPercent < 0.2 {
		warningText := "⚠️ STAMINA BAJA"
		ebitenutil.DebugPrintAt(screen, warningText, int(hudX+10), int(hudY+60))
	}
}
func (g *Game) drawBar(screen *ebiten.Image, x, y, width, height float32, fill float64, col color.RGBA, label string) {
	// Fondo
	vector.DrawFilledRect(screen, x, y, width, height, color.RGBA{50, 50, 50, 255}, false)

	// Relleno
	fillWidth := float32(fill) * width
	vector.DrawFilledRect(screen, x, y, fillWidth, height, col, false)

	// Borde
	vector.StrokeRect(screen, x, y, width, height, 1, color.White, false)

	// Label
	ebitenutil.DebugPrintAt(screen, label, int(x), int(y-15))
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
			"Update: %.2fms\n"+
			"Draw: %.2fms\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"PLAYER:\n"+
			"State: %s\n"+
			"Pos: (%.0f, %.0f)\n"+
			"Vel: (%.1f, %.1f)\n"+
			"OnGround: %v\n"+
			"OnWall: %v\n"+
			"JumpCount: %d/%d\n"+
			"CanDash: %v\n"+
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
		float64(g.updateDuration.Microseconds())/1000.0,
		float64(g.drawDuration.Microseconds())/1000.0,
		g.player.State,
		g.player.Position.X,
		g.player.Position.Y,
		g.player.Velocity.X,
		g.player.Velocity.Y,
		g.player.IsOnGround,
		g.player.IsTouchingWall,
		g.player.JumpCount,
		g.player.MaxJumps,
		g.player.CanDash,
		inputMethod,
	)

	debugBg := ebiten.NewImage(300, 400)
	debugBg.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(debugBg, nil)

	ebitenutil.DebugPrint(screen, debugText)
}

/*
func (g *Game) getStateName() string {
	switch g.state {
	case StateMainMenu:
		return "MainMenu"
	case StatePlaying:
		return "Playing"
	case StatePaused:
		return "Paused"
	case StateGameOver:
		return "GameOver"
	case StateVictory:
		return "Victory"
	default:
		return "Unknown"
	}
}
*/
