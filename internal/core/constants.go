// internal/core/constants.go
package core

// Este archivo ahora solo re-exporta las constantes de config
// para mantener compatibilidad

import "github.com/MarcosBrindis/boss-arena-go/internal/config"

// Re-exportar constantes de pantalla
const (
	ScreenWidth    = config.ScreenWidth
	ScreenHeight   = config.ScreenHeight
	PlayableWidth  = config.PlayableWidth
	PlayableHeight = config.PlayableHeight
	TargetTPS      = config.TargetTPS
)

// Re-exportar constantes de física
const (
	Gravity           = config.Gravity
	MaxFallSpeed      = config.MaxFallSpeed
	GroundFriction    = config.GroundFriction
	AirFriction       = config.AirFriction
	WallSlideFriction = config.WallSlideFriction
)

// Re-exportar constantes de arena
const (
	FloorY      = config.FloorY
	FloorHeight = config.FloorHeight
	WallLeftX   = config.WallLeftX
	WallRightX  = config.WallRightX
	StepHeight  = config.StepHeight
)

// Re-exportar colores
var (
	ColorBackground  = config.ColorBackground
	ColorFloor       = config.ColorFloor
	ColorFloorGrid   = config.ColorFloorGrid
	ColorWalls       = config.ColorWalls
	ColorWallsBorder = config.ColorWallsBorder

	ColorHeroPrimary   = config.ColorHeroPrimary
	ColorHeroSecondary = config.ColorHeroSecondary
	ColorHeroDash      = config.ColorHeroDash

	ColorBossPhase1 = config.ColorBossPhase1
	ColorBossPhase2 = config.ColorBossPhase2
	ColorBossPhase3 = config.ColorBossPhase3

	ColorHPBarFull  = config.ColorHPBarFull
	ColorHPBarMid   = config.ColorHPBarMid
	ColorHPBarLow   = config.ColorHPBarLow
	ColorStaminaBar = config.ColorStaminaBar

	ColorMeteor     = config.ColorMeteor
	ColorExplosion  = config.ColorExplosion
	ColorProjectile = config.ColorProjectile

	ColorDebugText = config.ColorDebugText
)

// Re-exportar tipos
type GameState = config.GameState

const (
	StateMainMenu = config.StateMainMenu
	StatePlaying  = config.StatePlaying
	StatePaused   = config.StatePaused
	StateGameOver = config.StateGameOver
	StateVictory  = config.StateVictory
)

// Re-exportar versión
const (
	GameVersion = config.GameVersion
	GameTitle   = config.GameTitle
)
