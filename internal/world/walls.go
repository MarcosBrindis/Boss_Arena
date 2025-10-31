package world

import (
	"image/color"

	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// WallSegment representa un segmento de pared
type WallSegment struct {
	Rect        utils.Rectangle
	Color       color.RGBA
	BorderColor color.RGBA
}

// NewWallSegment crea un nuevo segmento de pared
func NewWallSegment(x, y, width, height float64, wallColor, borderColor color.RGBA) *WallSegment {
	return &WallSegment{
		Rect:        utils.NewRectangle(x, y, width, height),
		Color:       wallColor,
		BorderColor: borderColor,
	}
}

// Draw dibuja el segmento de pared
func (w *WallSegment) Draw(screen *ebiten.Image) {
	// Dibujar relleno
	vector.DrawFilledRect(
		screen,
		float32(w.Rect.X),
		float32(w.Rect.Y),
		float32(w.Rect.Width),
		float32(w.Rect.Height),
		w.Color,
		false,
	)

	// Dibujar borde (4 líneas)
	borderWidth := float32(2)

	// Borde superior
	vector.StrokeLine(
		screen,
		float32(w.Rect.Left()),
		float32(w.Rect.Top()),
		float32(w.Rect.Right()),
		float32(w.Rect.Top()),
		borderWidth,
		w.BorderColor,
		false,
	)

	// Borde derecho
	vector.StrokeLine(
		screen,
		float32(w.Rect.Right()),
		float32(w.Rect.Top()),
		float32(w.Rect.Right()),
		float32(w.Rect.Bottom()),
		borderWidth,
		w.BorderColor,
		false,
	)

	// Borde inferior
	vector.StrokeLine(
		screen,
		float32(w.Rect.Right()),
		float32(w.Rect.Bottom()),
		float32(w.Rect.Left()),
		float32(w.Rect.Bottom()),
		borderWidth,
		w.BorderColor,
		false,
	)

	// Borde izquierdo
	vector.StrokeLine(
		screen,
		float32(w.Rect.Left()),
		float32(w.Rect.Bottom()),
		float32(w.Rect.Left()),
		float32(w.Rect.Top()),
		borderWidth,
		w.BorderColor,
		false,
	)
}

// Intersects verifica colisión con otro rectángulo
func (w *WallSegment) Intersects(rect utils.Rectangle) bool {
	return w.Rect.Intersects(rect)
}

// GetPenetration retorna el vector de penetración para resolver colisión
func (w *WallSegment) GetPenetration(rect utils.Rectangle) utils.Vector2 {
	return w.Rect.GetPenetration(rect)
}
