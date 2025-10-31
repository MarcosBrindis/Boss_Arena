// internal/core/game.go
package core

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/MarcosBrindis/boss-arena-go/internal/input" // ← NUEVO IMPORT
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game implementa la interfaz ebiten.Game
type Game struct {
	// Configuración
	config *Config

	// Input System (NUEVO)
	controller *input.Controller

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

	// Control de input global (sin sleep)
	escapeKeyPressedLastFrame bool
	f11KeyPressedLastFrame    bool
	f3KeyPressedLastFrame     bool
}

// NewGame crea una nueva instancia del juego
func NewGame() *Game {
	cfg := DefaultConfig()

	return &Game{
		config: cfg,
		// Inicializar controller con configuración
		controller: input.NewController(
			cfg.GamepadDeadzone,
			cfg.JumpBufferFrames,
			cfg.CoyoteTimeFrames,
		),
		state:      StatePlaying,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
	}
}

// Update actualiza la lógica del juego (llamado 60 veces por segundo)
func (g *Game) Update() error {
	start := time.Now()

	// Actualizar controller (NUEVO)
	g.controller.Update()

	// Calcular delta time
	now := time.Now()
	g.deltaTime = now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	// Incrementar contador de frames
	g.frame++

	// Calcular TPS/FPS cada segundo
	if g.frame%60 == 0 {
		g.tps = ebiten.ActualTPS()
		g.fps = ebiten.ActualFPS()
	}

	// Manejar input global (pausa, salir, etc)
	if err := g.handleGlobalInput(); err != nil {
		return err
	}

	// Si está pausado, no actualizar lógica del juego
	if g.isPaused {
		g.updateDuration = time.Since(start)
		return nil
	}

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

// Draw dibuja el juego en pantalla (llamado cada frame)
func (g *Game) Draw(screen *ebiten.Image) {
	start := time.Now()

	// Limpiar pantalla con color de fondo
	screen.Fill(ColorBackground)

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

	// Dibujar información de debug si está habilitada
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
	// ESC para pausar/despausar (sin sleep, instantáneo)
	escapePressed := ebiten.IsKeyPressed(ebiten.KeyEscape)
	if escapePressed && !g.escapeKeyPressedLastFrame {
		g.isPaused = !g.isPaused
	}
	g.escapeKeyPressedLastFrame = escapePressed

	// F11 para fullscreen
	f11Pressed := ebiten.IsKeyPressed(ebiten.KeyF11)
	if f11Pressed && !g.f11KeyPressedLastFrame {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	g.f11KeyPressedLastFrame = f11Pressed

	// F3 para toggle debug info
	f3Pressed := ebiten.IsKeyPressed(ebiten.KeyF3)
	if f3Pressed && !g.f3KeyPressedLastFrame {
		g.config.ShowDebugInfo = !g.config.ShowDebugInfo
	}
	g.f3KeyPressedLastFrame = f3Pressed

	return nil
}

func (g *Game) updateMainMenu() {
	// TODO: Implementar en módulos posteriores
}

func (g *Game) updatePlaying() {
	// Aquí probaremos el Input System
	// TODO: En Módulo 4 implementaremos el jugador completo
}

func (g *Game) updatePaused() {
	// TODO: Implementar en módulos posteriores
}

func (g *Game) updateGameOver() {
	// TODO: Implementar en módulos posteriores
}

func (g *Game) updateVictory() {
	// TODO: Implementar en módulos posteriores
}

// ============================================================================
// MÉTODOS DE DRAW POR ESTADO
// ============================================================================

func (g *Game) drawMainMenu(screen *ebiten.Image) {
	// TODO: Implementar en módulos posteriores
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	// Mensaje actualizado con controles correctos
	msg := "🎮 Módulo 2: Input System funcionando!\n\n"
	msg += "⌨️  Prueba los controles:\n"
	msg += "  WASD/Arrows = Mover\n"
	msg += "  Space/W     = Saltar\n"
	msg += "  Z/J         = Atacar\n"
	msg += "  X/K/SHIFT   = Dash\n" // ← ACTUALIZADO
	msg += "  C/L         = Especial\n\n"

	if g.controller.IsGamepadConnected() {
		msg += "🎮 Gamepad conectado!\n"
		msg += fmt.Sprintf("   %s\n\n", g.controller.GetGamepadName())
		msg += "  ✕ (Cross)    = Saltar\n"
		msg += "  ⬜ (Square)   = Atacar\n"
		msg += "  ⚪ (Circle)   = Dash\n"
		msg += "  R2           = Dash\n" // ← NUEVO
		msg += "  🔺 (Triangle) = Especial\n"
	} else {
		msg += "🎮 Conecta un gamepad PS5 para probarlo\n"
	}

	msg += "\n✅ Esperando Módulo 3 (Arena)..."

	bounds := screen.Bounds()
	x := float64(bounds.Dx()/2 - 220)
	y := 100.0

	ebitenutil.DebugPrintAt(screen, msg, int(x), int(y))

	// Dibujar cuadrado controlable (demo de input)
	g.drawControllableSquare(screen)

	// Dibujar indicadores de input
	g.drawInputIndicators(screen)
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
	// TODO: Implementar en módulos posteriores
}

func (g *Game) drawVictory(screen *ebiten.Image) {
	// TODO: Implementar en módulos posteriores
}

// ============================================================================
// DEBUG INFO
// ============================================================================

func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	// Información de input (NUEVO)
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
			"Delta: %.4fs\n"+
			"State: %s\n"+
			"Paused: %v\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"INPUT:\n"+
			"Method: %s\n"+
			"Horizontal: %.2f\n"+
			"Vertical: %.2f\n"+
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
		g.deltaTime,
		g.getStateName(),
		g.isPaused,
		inputMethod,
		g.controller.GetHorizontalAxis(),
		g.controller.GetVerticalAxis(),
	)

	debugBg := ebiten.NewImage(350, 320)
	debugBg.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(debugBg, nil)

	ebitenutil.DebugPrint(screen, debugText)
}

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

// ============================================================================
// DEMO DE INPUT SYSTEM
// ============================================================================

// Posición del cuadrado controlable (demo)
var demoSquareX float64 = 640
var demoSquareY float64 = 450

func (g *Game) drawControllableSquare(screen *ebiten.Image) {
	// Mover el cuadrado con el input
	speed := 5.0
	demoSquareX += g.controller.GetHorizontalAxis() * speed
	demoSquareY += g.controller.GetVerticalAxis() * speed

	// Limitar a la pantalla
	if demoSquareX < 25 {
		demoSquareX = 25
	}
	if demoSquareX > ScreenWidth-25 {
		demoSquareX = ScreenWidth - 25
	}
	if demoSquareY < 25 {
		demoSquareY = 25
	}
	if demoSquareY > ScreenHeight-25 {
		demoSquareY = ScreenHeight - 25
	}

	// Color según acción
	squareColor := ColorHeroPrimary
	if g.controller.IsAttackHeld() {
		squareColor = color.RGBA{255, 0, 0, 255} // Rojo si ataca
	} else if g.controller.IsDashHeld() {
		squareColor = color.RGBA{255, 255, 0, 255} // Amarillo si dash
	}

	// Crear y dibujar cuadrado
	square := ebiten.NewImage(50, 50)
	square.Fill(squareColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(demoSquareX-25, demoSquareY-25)
	screen.DrawImage(square, op)

	// Texto
	ebitenutil.DebugPrintAt(
		screen,
		"▲ Muéveme con WASD o Stick",
		int(demoSquareX)-80,
		int(demoSquareY)+40,
	)
}

func (g *Game) drawInputIndicators(screen *ebiten.Image) {
	startX := 400.0
	startY := 550.0
	spacing := 120.0

	// Jump
	g.drawButton(screen, startX, startY, "JUMP", g.controller.IsJumpHeld())

	// Attack
	g.drawButton(screen, startX+spacing, startY, "ATTACK", g.controller.IsAttackHeld())

	// Dash
	g.drawButton(screen, startX+spacing*2, startY, "DASH", g.controller.IsDashHeld())

	// Special
	g.drawButton(screen, startX+spacing*3, startY, "SPECIAL", g.controller.IsSpecialHeld())
}

func (g *Game) drawButton(screen *ebiten.Image, x, y float64, label string, pressed bool) {
	buttonColor := color.RGBA{60, 60, 80, 255}
	if pressed {
		buttonColor = ColorHeroPrimary
	}

	button := ebiten.NewImage(100, 40)
	button.Fill(buttonColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(button, op)

	ebitenutil.DebugPrintAt(screen, label, int(x)+20, int(y)+15)
}

// ============================================================================
// ANIMACIÓN ORIGINAL (mantenida)
// ============================================================================

func (g *Game) drawAnimatedSquare(screen *ebiten.Image) {
	centerX := float64(ScreenWidth / 2)
	centerY := float64(ScreenHeight/2) + 100
	radius := 80.0

	angle := float64(g.frame) * 0.05

	squareX := int(centerX + radius*math.Cos(angle))
	squareY := int(centerY + radius*math.Sin(angle))

	square := ebiten.NewImage(50, 50)
	square.Fill(ColorHeroPrimary)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(squareX-25), float64(squareY-25))
	screen.DrawImage(square, op)

	ebitenutil.DebugPrintAt(
		screen,
		"▲ Cuadrado animado (verificando 60 FPS)",
		ScreenWidth/2-120,
		int(centerY+radius+30),
	)
}
