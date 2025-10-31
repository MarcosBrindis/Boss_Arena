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
	baseFloorY := float64(config.FloorY) - 50 // Piso en Y=600
	wallHeight := 450.0                       // Altura de las paredes

	// ========================================================================
	// PARED IZQUIERDA SÓLIDA (sin escalones por ahora)
	// ========================================================================
	a.wallsLeft = []*WallSegment{
		NewWallSegment(
			0,                     // X: Pegada al borde izquierdo
			baseFloorY-wallHeight, // Y: Desde arriba
			wallThickness,         // Ancho de la pared
			wallHeight,            // Altura completa
			a.wallColor,
			a.wallBorderColor,
		),
	}

	// ========================================================================
	// PARED DERECHA SÓLIDA (sin escalones por ahora)
	// ========================================================================
	a.wallsRight = []*WallSegment{
		NewWallSegment(
			float64(a.width)-wallThickness, // X: Pegada al borde derecho
			baseFloorY-wallHeight,          // Y: Desde arriba
			wallThickness,                  // Ancho de la pared
			wallHeight,                     // Altura completa
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

func (a *Arena) IsTouchingWall(rect utils.Rectangle) (bool, int) {
	testRectLeft := utils.NewRectangle(rect.X-2, rect.Y, rect.Width, rect.Height)
	testRectRight := utils.NewRectangle(rect.X+2, rect.Y, rect.Width, rect.Height)

	for _, wall := range a.wallsLeft {
		if wall.Intersects(testRectLeft) {
			return true, -1
		}
	}

	for _, wall := range a.wallsRight {
		if wall.Intersects(testRectRight) {
			return true, 1
		}
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
