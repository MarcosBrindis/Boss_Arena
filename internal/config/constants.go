package config

import "image/color"

// ============================================================================
// CONFIGURACIÓN DE PANTALLA
// ============================================================================

const (
	// Resolución base (HD)
	ScreenWidth  = 1280
	ScreenHeight = 720

	// Área jugable (con márgenes para UI)
	PlayableWidth  = 1100
	PlayableHeight = 650

	// Performance
	TargetTPS = 60 // Ticks por segundo (60 para 60 FPS)
)

// ============================================================================
// FÍSICA GLOBAL
// ============================================================================

const (
	Gravity           = 0.6  // Gravedad estándar
	MaxFallSpeed      = 10.0 // Velocidad máxima de caída
	GroundFriction    = 0.85 // Fricción en suelo
	AirFriction       = 0.98 // Fricción en aire
	WallSlideFriction = 0.95 // Fricción al deslizarse en pared
)

// ============================================================================
// ARENA - DIMENSIONES
// ============================================================================

const (
	// Piso
	FloorY      = 650 // Posición Y del piso
	FloorHeight = 70  // Altura del piso

	// Paredes escalonadas (izquierda)
	WallLeftX = 50

	// Paredes escalonadas (derecha)
	WallRightX = 1230

	// Escalones (altura desde el piso hacia arriba)
	StepHeight = 150
)

// ============================================================================
// COLORES - PALETA PRINCIPAL
// ============================================================================

var (
	// Arena
	ColorBackground  = color.RGBA{26, 32, 44, 255}   // #1a202c - Negro azulado
	ColorFloor       = color.RGBA{74, 85, 104, 255}  // #4a5568 - Gris medio
	ColorFloorGrid   = color.RGBA{45, 55, 72, 255}   // #2d3748 - Gris oscuro
	ColorWalls       = color.RGBA{90, 103, 216, 255} // #5a67d8 - Azul metálico
	ColorWallsBorder = color.RGBA{76, 81, 191, 255}  // #4c51bf - Azul oscuro

	// Héroe
	ColorHeroPrimary   = color.RGBA{0, 217, 255, 255}   // #00d9ff - Cyan brillante
	ColorHeroSecondary = color.RGBA{255, 255, 255, 255} // #ffffff - Blanco
	ColorHeroDash      = color.RGBA{0, 128, 255, 128}   // #0080ff - Azul eléctrico (semi-transparente)

	// Boss (cambia por fase)
	ColorBossPhase1 = color.RGBA{255, 69, 0, 255}  // #ff4500 - Rojo fuego
	ColorBossPhase2 = color.RGBA{255, 99, 71, 255} // #ff6347 - Naranja
	ColorBossPhase3 = color.RGBA{255, 165, 0, 255} // #ffa500 - Amarillo

	// UI
	ColorHPBarFull  = color.RGBA{72, 187, 120, 255}  // #48bb78 - Verde
	ColorHPBarMid   = color.RGBA{246, 173, 85, 255}  // #f6ad55 - Naranja
	ColorHPBarLow   = color.RGBA{252, 129, 129, 255} // #fc8181 - Rojo
	ColorStaminaBar = color.RGBA{66, 153, 225, 255}  // #4299e1 - Azul

	// Efectos
	ColorMeteor     = color.RGBA{139, 0, 0, 255}   // #8B0000 - Rojo sangre
	ColorExplosion  = color.RGBA{255, 140, 0, 255} // #ff8c00 - Naranja fuego
	ColorProjectile = color.RGBA{255, 0, 0, 255}   // #ff0000 - Rojo puro

	// Debug
	ColorDebugText = color.RGBA{0, 255, 0, 255} // Verde brillante
)

// ============================================================================
// ESTADOS DEL JUEGO
// ============================================================================

type GameState int

const (
	StateMainMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
	StateVictory
)

// ============================================================================
// VERSIÓN
// ============================================================================

const (
	GameVersion = "v0.1.0-alpha"
	GameTitle   = "Titan's Arena"
)
