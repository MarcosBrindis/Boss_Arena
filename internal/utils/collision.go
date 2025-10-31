package utils

// Rectangle representa un rectángulo AABB (Axis-Aligned Bounding Box)
type Rectangle struct {
	X      float64 // Posición X (esquina superior izquierda)
	Y      float64 // Posición Y (esquina superior izquierda)
	Width  float64
	Height float64
}

// NewRectangle crea un nuevo rectángulo
func NewRectangle(x, y, width, height float64) Rectangle {
	return Rectangle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Center retorna el centro del rectángulo
func (r Rectangle) Center() Vector2 {
	return Vector2{
		X: r.X + r.Width/2,
		Y: r.Y + r.Height/2,
	}
}

// Left retorna la coordenada X izquierda
func (r Rectangle) Left() float64 {
	return r.X
}

// Right retorna la coordenada X derecha
func (r Rectangle) Right() float64 {
	return r.X + r.Width
}

// Top retorna la coordenada Y superior
func (r Rectangle) Top() float64 {
	return r.Y
}

// Bottom retorna la coordenada Y inferior
func (r Rectangle) Bottom() float64 {
	return r.Y + r.Height
}

// Intersects verifica si dos rectángulos se intersectan (AABB collision)
func (r Rectangle) Intersects(other Rectangle) bool {
	return r.Left() < other.Right() &&
		r.Right() > other.Left() &&
		r.Top() < other.Bottom() &&
		r.Bottom() > other.Top()
}

// Contains verifica si un punto está dentro del rectángulo
func (r Rectangle) Contains(point Vector2) bool {
	return point.X >= r.Left() &&
		point.X <= r.Right() &&
		point.Y >= r.Top() &&
		point.Y <= r.Bottom()
}

// Overlaps retorna el área de solapamiento con otro rectángulo
func (r Rectangle) Overlaps(other Rectangle) Rectangle {
	if !r.Intersects(other) {
		return Rectangle{}
	}

	x := Max(r.Left(), other.Left())
	y := Max(r.Top(), other.Top())
	width := Min(r.Right(), other.Right()) - x
	height := Min(r.Bottom(), other.Bottom()) - y

	return Rectangle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// GetPenetration retorna el vector de penetración para resolver colisión
func (r Rectangle) GetPenetration(other Rectangle) Vector2 {
	if !r.Intersects(other) {
		return Zero()
	}

	// Calcular profundidad de penetración en cada eje
	overlapX := Min(r.Right()-other.Left(), other.Right()-r.Left())
	overlapY := Min(r.Bottom()-other.Top(), other.Bottom()-r.Top())

	// Retornar el eje con menor penetración
	if overlapX < overlapY {
		// Resolver en X
		if r.Center().X < other.Center().X {
			return Vector2{X: -overlapX, Y: 0}
		}
		return Vector2{X: overlapX, Y: 0}
	}

	// Resolver en Y
	if r.Center().Y < other.Center().Y {
		return Vector2{X: 0, Y: -overlapY}
	}
	return Vector2{X: 0, Y: overlapY}
}

// Move mueve el rectángulo
func (r Rectangle) Move(offset Vector2) Rectangle {
	return Rectangle{
		X:      r.X + offset.X,
		Y:      r.Y + offset.Y,
		Width:  r.Width,
		Height: r.Height,
	}
}

// Scale escala el rectángulo desde el centro
func (r Rectangle) Scale(factor float64) Rectangle {
	center := r.Center()
	newWidth := r.Width * factor
	newHeight := r.Height * factor

	return Rectangle{
		X:      center.X - newWidth/2,
		Y:      center.Y - newHeight/2,
		Width:  newWidth,
		Height: newHeight,
	}
}
