package world

import (
	"image/color"
	"math"

	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Background representa el fondo con efecto parallax
type Background struct {
	width  int
	height int
	frame  uint64

	// Colores de las capas
	skyColor       color.RGBA
	mountainColor1 color.RGBA
	mountainColor2 color.RGBA
}

// NewBackground crea un nuevo fondo
func NewBackground(width, height int) *Background {
	return &Background{
		width:  width,
		height: height,
		frame:  0,

		// Colores del fondo
		skyColor:       color.RGBA{26, 32, 44, 255},  // Azul oscuro
		mountainColor1: color.RGBA{45, 55, 72, 100},  // Gris oscuro transparente
		mountainColor2: color.RGBA{74, 85, 104, 150}, // Gris medio transparente
	}
}

// Update actualiza el fondo (para animaciones)
func (b *Background) Update() {
	b.frame++
}

// Draw dibuja el fondo con efecto parallax
func (b *Background) Draw(screen *ebiten.Image) {
	// Capa 1: Cielo base (sin parallax)
	screen.Fill(b.skyColor)

	// Capa 2: Estrellas parpadeantes (parallax muy lento)
	b.drawStars(screen)

	// Capa 3: Montañas lejanas (parallax lento)
	b.drawDistantMountains(screen)

	// Capa 4: Montañas cercanas (parallax medio)
	b.drawNearMountains(screen)
}

// drawStars dibuja estrellas parpadeantes
func (b *Background) drawStars(screen *ebiten.Image) {
	starColor := color.RGBA{255, 255, 255, 200}

	// Usar el frame para crear variación
	for i := 0; i < 50; i++ {
		// Posición pseudo-aleatoria basada en el índice
		x := float32((i*137 + 50) % b.width)
		y := float32((i*219 + 30) % (b.height / 2))

		// Parpadeo sutil
		brightness := float32(math.Sin(float64(b.frame)*0.02+float64(i)*0.5)*0.3 + 0.7)
		starColor.A = uint8(brightness * 200)

		// Dibujar estrella pequeña
		vector.DrawFilledCircle(screen, x, y, 1.5, starColor, false)
	}
}

// drawDistantMountains dibuja montañas lejanas (parallax lento)
func (b *Background) drawDistantMountains(screen *ebiten.Image) {
	// Offset de parallax muy lento (simulación de cámara)
	offset := float32(math.Sin(float64(b.frame)*0.001) * 10)

	// Montañas triangulares simples
	for i := 0; i < 5; i++ {
		baseX := float32(i*300) + offset - 100
		baseY := float32(b.height / 2)
		peakX := baseX + 150
		peakY := baseY - 200
		rightX := baseX + 300

		// Dibujar triángulo
		b.drawTriangle(screen, baseX, baseY, peakX, peakY, rightX, baseY, b.mountainColor1)
	}
}

// drawNearMountains dibuja montañas cercanas (parallax medio)
func (b *Background) drawNearMountains(screen *ebiten.Image) {
	// Offset de parallax medio
	offset := float32(math.Sin(float64(b.frame)*0.003) * 20)

	for i := 0; i < 4; i++ {
		baseX := float32(i*350) + offset - 150
		baseY := float32(b.height/2 + 50)
		peakX := baseX + 180
		peakY := baseY - 150
		rightX := baseX + 360

		b.drawTriangle(screen, baseX, baseY, peakX, peakY, rightX, baseY, b.mountainColor2)
	}
}

// drawTriangle dibuja un triángulo simple
func (b *Background) drawTriangle(screen *ebiten.Image, x1, y1, x2, y2, x3, y3 float32, col color.RGBA) {
	// Dibujar líneas del triángulo
	vector.StrokeLine(screen, x1, y1, x2, y2, 2, col, false)
	vector.StrokeLine(screen, x2, y2, x3, y3, 2, col, false)
	vector.StrokeLine(screen, x3, y3, x1, y1, 2, col, false)

	// Rellenar (aproximado con líneas verticales)
	// Esto es una simplificación; para un relleno real se necesitaría vertex arrays
	steps := int(utils.Abs(float64(x3 - x1)))
	for i := 0; i < steps; i++ {
		t := float32(i) / float32(steps)

		// Interpolar bordes
		leftX := x1 + (x2-x1)*t
		leftY := y1 + (y2-y1)*t

		rightX := x1 + (x3-x1)*t
		rightY := y1 + (y3-y1)*t

		// Dibujar línea vertical de relleno
		vector.StrokeLine(screen, leftX, leftY, rightX, rightY, 1, col, false)
	}
}
