package world

import (
	"image/color"

	"github.com/MarcosBrindis/boss-arena-go/internal/config"
	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Arena representa el escenario completo de batalla
type Arena struct {
	width  int
	height int

	// Componentes
	background *Background
	floor      *WallSegment
	wallsLeft  []*WallSegment
	wallsRight []*WallSegment

	// Colores
	floorColor      color.RGBA
	floorGridColor  color.RGBA
	wallColor       color.RGBA
	wallBorderColor color.RGBA
}

// NewArena crea una nueva arena
func NewArena(width, height int) *Arena {
	arena := &Arena{
		width:  width,
		height: height,

		// Colores (usando config en lugar de core)
		floorColor:      config.ColorFloor,
		floorGridColor:  config.ColorFloorGrid,
		wallColor:       config.ColorWalls,
		wallBorderColor: config.ColorWallsBorder,
	}

	// Crear componentes
	arena.background = NewBackground(width, height)
	arena.createFloor()
	arena.createWalls()

	return arena
}

// createFloor crea el piso de la arena
func (a *Arena) createFloor() {
	// El piso debe ser GRUESO para que el cuadrado no caiga por gaps
	floorY := float64(config.FloorY) - 50 // 600
	floorHeight := 120.0                  // MÁS GRUESO

	a.floor = NewWallSegment(
		0,
		floorY,
		float64(a.width),
		floorHeight,
		a.floorColor,
		a.floorGridColor,
	)
}

// createWalls crea las paredes escalonadas (estilo Mega Man X)
func (a *Arena) createWalls() {
	wallThickness := 50.0
	baseFloorY := float64(config.FloorY) - 50 // 600
	wallHeight := baseFloorY

	// ========================================================================
	// PARED IZQUIERDA SÓLIDA
	// ========================================================================
	a.wallsLeft = []*WallSegment{
		NewWallSegment(
			0,                     // X = 0 (pegada al borde)
			baseFloorY-wallHeight, // Y = 150
			wallThickness,         // Width = 50
			wallHeight,            // Height = 450
			a.wallColor,
			a.wallBorderColor,
		),
	}

	// ========================================================================
	// PARED DERECHA SÓLIDA
	// ========================================================================
	a.wallsRight = []*WallSegment{
		NewWallSegment(
			float64(a.width)-wallThickness, // X = 1230 (1280 - 50)
			baseFloorY-wallHeight,          // Y = 150
			wallThickness,                  // Width = 50
			wallHeight,                     // Height = 450
			a.wallColor,
			a.wallBorderColor,
		),
	}
}

// Update actualiza la arena
func (a *Arena) Update() {
	a.background.Update()
}

// Draw dibuja la arena completa
func (a *Arena) Draw(screen *ebiten.Image) {
	a.background.Draw(screen)
	a.floor.Draw(screen)
	a.drawFloorGrid(screen)

	for _, wall := range a.wallsLeft {
		wall.Draw(screen)
	}

	for _, wall := range a.wallsRight {
		wall.Draw(screen)
	}
}

// drawFloorGrid dibuja el grid decorativo del piso
func (a *Arena) drawFloorGrid(screen *ebiten.Image) {
	gridSpacing := 40.0
	floorY := float64(config.FloorY)
	floorBottom := floorY + float64(config.FloorHeight)

	// Líneas verticales
	for x := 0.0; x < float64(a.width); x += gridSpacing {
		vector.StrokeLine(
			screen,
			float32(x),
			float32(floorY),
			float32(x),
			float32(floorBottom),
			1,
			a.floorGridColor,
			false,
		)
	}

	// Líneas horizontales
	for y := floorY; y < floorBottom; y += gridSpacing {
		vector.StrokeLine(
			screen,
			0,
			float32(y),
			float32(a.width),
			float32(y),
			1,
			a.floorGridColor,
			false,
		)
	}
}

// ============================================================================
// DETECCIÓN DE COLISIONES
// ============================================================================

func (a *Arena) CheckCollision(rect utils.Rectangle) (bool, utils.Vector2) {
	if a.floor.Intersects(rect) {
		return true, a.floor.GetPenetration(rect)
	}

	for _, wall := range a.wallsLeft {
		if wall.Intersects(rect) {
			return true, wall.GetPenetration(rect)
		}
	}

	for _, wall := range a.wallsRight {
		if wall.Intersects(rect) {
			return true, wall.GetPenetration(rect)
		}
	}

	return false, utils.Zero()
}

func (a *Arena) IsOnGround(rect utils.Rectangle) bool {
	testRect := utils.NewRectangle(
		rect.X,
		rect.Y,
		rect.Width,
		rect.Height+2,
	)

	if a.floor.Intersects(testRect) {
		return true
	}

	for _, wall := range a.wallsLeft {
		if wall.Intersects(testRect) {
			return true
		}
	}

	for _, wall := range a.wallsRight {
		if wall.Intersects(testRect) {
			return true
		}
	}

	return false
}

// IsTouchingWall verifica si un rectángulo está tocando una pared lateral
func (a *Arena) IsTouchingWall(rect utils.Rectangle) (bool, int) {
	// Configuración de detección
	margin := 8.0

	// Crear áreas de detección expandidas
	testRectLeft := utils.NewRectangle(
		rect.X-margin,
		rect.Y+5,
		rect.Width+margin,
		rect.Height-10,
	)

	testRectRight := utils.NewRectangle(
		rect.X,
		rect.Y+5,
		rect.Width+margin,
		rect.Height-10,
	)

	// =========================================================================
	// MÉTODO 1: Verificar contra las paredes en los arrays
	// =========================================================================

	// Verificar paredes izquierdas
	if len(a.wallsLeft) > 0 {
		for _, wall := range a.wallsLeft {
			if wall != nil && wall.Intersects(testRectLeft) {
				return true, -1
			}
		}
	}

	// Verificar paredes derechas
	if len(a.wallsRight) > 0 {
		for _, wall := range a.wallsRight {
			if wall != nil && wall.Intersects(testRectRight) {
				return true, 1
			}
		}
	}

	// =========================================================================
	// MÉTODO 2: Fallback - Verificar contra posiciones fijas (más robusto)
	// =========================================================================

	// Pared izquierda: X entre 0 y 60
	if rect.Left() <= 60 {
		return true, -1
	}

	// Pared derecha: X entre 1220 y 1280
	if rect.Right() >= float64(a.width)-60 {
		return true, 1
	}

	return false, 0
}

func (a *Arena) GetFloorY() float64 {
	return float64(config.FloorY)
}

func (a *Arena) GetBounds() utils.Rectangle {
	return utils.NewRectangle(
		110,
		0,
		float64(a.width-220),
		float64(config.FloorY),
	)
}
