package core

import (
	"fmt"
	"image/color"
	"time"

	"github.com/MarcosBrindis/boss-arena-go/internal/input"
	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
	"github.com/MarcosBrindis/boss-arena-go/internal/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game implementa la interfaz ebiten.Game
type Game struct {
	// Configuración
	config *Config

	// Input System
	controller *input.Controller

	// World (NUEVO)
	arena *world.Arena

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

	return &Game{
		config: cfg,
		controller: input.NewController(
			cfg.GamepadDeadzone,
			cfg.JumpBufferFrames,
			cfg.CoyoteTimeFrames,
		),
		// Crear arena (NUEVO)
		arena: world.NewArena(ScreenWidth, ScreenHeight),

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

	// Calcular TPS/FPS cada segundo
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

	// Actualizar arena (NUEVO)
	g.arena.Update()

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

	// Limpiar pantalla (ya no es necesario, la arena dibuja el fondo)
	// screen.Fill(ColorBackground)

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
	// TODO: En Módulo 4 implementaremos el jugador
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
	// 1. Dibujar arena (NUEVO)
	g.arena.Draw(screen)

	// 2. Mensaje actualizado
	msg := "🏛️ Módulo 3: Arena completada!\n\n"
	msg += "✅ Paredes escalonadas estilo Mega Man X\n"
	msg += "✅ Fondo con efecto parallax\n"
	msg += "✅ Sistema de colisiones AABB\n"
	msg += "✅ Grid decorativo en el piso\n\n"
	msg += "⏭️  Esperando Módulo 4 (Player)..."

	// Dibujar en el centro superior
	ebitenutil.DebugPrintAt(screen, msg, ScreenWidth/2-200, 50)

	// 3. Cuadrado controlable con colisiones (NUEVO)
	g.drawControllableSquareWithCollisions(screen)

	// 4. Indicadores de input
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
	// TODO
}

func (g *Game) drawVictory(screen *ebiten.Image) {
	// TODO
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
// DEMO DE INPUT + COLISIONES (NUEVO)
// ============================================================================

var (
	demoSquareX  float64 = 640
	demoSquareY  float64 = 575 // 600 (piso) - 25 (mitad del cuadrado)
	demoVelocity utils.Vector2
)

func (g *Game) drawControllableSquareWithCollisions(screen *ebiten.Image) {
	squareSize := 50.0
	speed := 4.0
	gravity := 0.5
	maxFallSpeed := 10.0

	// Input horizontal
	inputX := g.controller.GetHorizontalAxis()

	// Aplicar input con fricción
	if inputX != 0 {
		demoVelocity.X = inputX * speed
	} else {
		demoVelocity.X *= 0.85
		if utils.Abs(demoVelocity.X) < 0.1 {
			demoVelocity.X = 0
		}
	}

	// Crear hitbox ANTES de mover
	squareRect := utils.NewRectangle(
		demoSquareX-squareSize/2,
		demoSquareY-squareSize/2,
		squareSize,
		squareSize,
	)

	// Verificar si está en suelo
	isOnGround := g.arena.IsOnGround(squareRect)

	// Aplicar gravedad SOLO si no está en suelo
	if !isOnGround {
		demoVelocity.Y += gravity
		if demoVelocity.Y > maxFallSpeed {
			demoVelocity.Y = maxFallSpeed
		}
	} else {
		demoVelocity.Y = 0
	}

	// =========================================================================
	// MOVIMIENTO HORIZONTAL CON COLISIONES
	// =========================================================================

	// Intentar mover horizontalmente
	newX := demoSquareX + demoVelocity.X

	// Crear hitbox en la nueva posición X
	testRect := utils.NewRectangle(
		newX-squareSize/2,
		demoSquareY-squareSize/2,
		squareSize,
		squareSize,
	)

	// Verificar colisión horizontal
	collidesX, _ := g.arena.CheckCollision(testRect)
	if !collidesX {
		// No hay colisión, permitir movimiento
		demoSquareX = newX
	} else {
		// Hay colisión, detener velocidad horizontal
		demoVelocity.X = 0
	}

	// =========================================================================
	// MOVIMIENTO VERTICAL CON COLISIONES
	// =========================================================================

	// Intentar mover verticalmente
	newY := demoSquareY + demoVelocity.Y

	// Crear hitbox en la nueva posición Y
	testRect = utils.NewRectangle(
		demoSquareX-squareSize/2,
		newY-squareSize/2,
		squareSize,
		squareSize,
	)

	// Verificar colisión vertical
	collidesY, _ := g.arena.CheckCollision(testRect)
	if !collidesY {
		// No hay colisión, permitir movimiento
		demoSquareY = newY
	} else {
		// Hay colisión, detener velocidad vertical
		demoVelocity.Y = 0
	}

	// =========================================================================
	// LÍMITES DE PANTALLA
	// =========================================================================

	floorY := g.arena.GetFloorY()

	// Límite inferior: Si cae muy abajo, resetear
	if demoSquareY > floorY+100 {
		demoSquareX = float64(ScreenWidth / 2)
		demoSquareY = 575
		demoVelocity = utils.Zero()
	}

	// Límites laterales: No salir de pantalla
	minX := squareSize/2 + 60 // Margen para la pared
	maxX := float64(ScreenWidth) - squareSize/2 - 60

	if demoSquareX < minX {
		demoSquareX = minX
		demoVelocity.X = 0
	}
	if demoSquareX > maxX {
		demoSquareX = maxX
		demoVelocity.X = 0
	}

	// Límite superior
	if demoSquareY < squareSize/2 {
		demoSquareY = squareSize / 2
		demoVelocity.Y = 0
	}

	// =========================================================================
	// DIBUJAR CUADRADO
	// =========================================================================

	// Color según input
	squareColor := ColorHeroPrimary
	if g.controller.IsAttackHeld() {
		squareColor = color.RGBA{255, 0, 0, 255}
	} else if g.controller.IsDashHeld() {
		squareColor = color.RGBA{255, 255, 0, 255}
	}

	// Dibujar
	square := ebiten.NewImage(int(squareSize), int(squareSize))
	square.Fill(squareColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(demoSquareX-squareSize/2, demoSquareY-squareSize/2)
	screen.DrawImage(square, op)

	// Texto
	ebitenutil.DebugPrintAt(
		screen,
		"▲ Muéveme con WASD o Stick!",
		int(demoSquareX)-80,
		int(demoSquareY)+40,
	)

	// =========================================================================
	// ESTADO Y DEBUG
	// =========================================================================

	// Hitbox final para verificaciones
	finalRect := utils.NewRectangle(
		demoSquareX-squareSize/2,
		demoSquareY-squareSize/2,
		squareSize,
		squareSize,
	)

	onGround := g.arena.IsOnGround(finalRect)
	touchingWall, wallSide := g.arena.IsTouchingWall(finalRect)

	statusMsg := ""
	if onGround {
		statusMsg = "En el suelo"
	}
	if touchingWall {
		if wallSide == -1 {
			statusMsg += " | Pared IZQ"
		} else {
			statusMsg += " | Pared DER"
		}
	}

	if statusMsg != "" {
		ebitenutil.DebugPrintAt(screen, statusMsg, int(demoSquareX)-80, int(demoSquareY)+55)
	}

	// DEBUG
	if g.config.ShowDebugInfo {
		debugMsg := fmt.Sprintf("Pos: (%.0f, %.0f) | Suelo: %v | Vel: (%.1f, %.1f)",
			demoSquareX, demoSquareY, onGround, demoVelocity.X, demoVelocity.Y)
		ebitenutil.DebugPrintAt(screen, debugMsg, 10, 370)
	}
}

func (g *Game) drawInputIndicators(screen *ebiten.Image) {
	startX := 400.0
	startY := 550.0
	spacing := 120.0

	g.drawButton(screen, startX, startY, "JUMP", g.controller.IsJumpHeld())
	g.drawButton(screen, startX+spacing, startY, "ATTACK", g.controller.IsAttackHeld())
	g.drawButton(screen, startX+spacing*2, startY, "DASH", g.controller.IsDashHeld())
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
// ANIMACIÓN ORIGINAL (ya no se usa, pero la dejamos por si acaso)
// ============================================================================
/*
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
*/
