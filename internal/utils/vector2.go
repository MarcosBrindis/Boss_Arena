package utils

import "math"

// Vector2 representa un vector 2D
type Vector2 struct {
	X float64
	Y float64
}

// NewVector2 crea un nuevo vector
func NewVector2(x, y float64) Vector2 {
	return Vector2{X: x, Y: y}
}

// Add suma dos vectores
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

// Sub resta dos vectores
func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

// Mul multiplica el vector por un escalar
func (v Vector2) Mul(scalar float64) Vector2 {
	return Vector2{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

// Div divide el vector por un escalar
func (v Vector2) Div(scalar float64) Vector2 {
	if scalar == 0 {
		return v
	}
	return Vector2{
		X: v.X / scalar,
		Y: v.Y / scalar,
	}
}

// Length retorna la longitud del vector
func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize retorna el vector normalizado
func (v Vector2) Normalize() Vector2 {
	length := v.Length()
	if length == 0 {
		return Vector2{X: 0, Y: 0}
	}
	return v.Div(length)
}

// Distance retorna la distancia entre dos vectores
func (v Vector2) Distance(other Vector2) float64 {
	return v.Sub(other).Length()
}

// Dot producto punto
func (v Vector2) Dot(other Vector2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Lerp interpolación lineal entre dos vectores
func (v Vector2) Lerp(other Vector2, t float64) Vector2 {
	return Vector2{
		X: v.X + (other.X-v.X)*t,
		Y: v.Y + (other.Y-v.Y)*t,
	}
}

// Rotate rota el vector por un ángulo (en radianes)
func (v Vector2) Rotate(angle float64) Vector2 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return Vector2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

// Clamp limita el vector a un rango
func (v Vector2) Clamp(min, max Vector2) Vector2 {
	return Vector2{
		X: Clamp(v.X, min.X, max.X),
		Y: Clamp(v.Y, min.Y, max.Y),
	}
}

// Zero retorna un vector cero
func Zero() Vector2 {
	return Vector2{X: 0, Y: 0}
}

// One retorna un vector (1, 1)
func One() Vector2 {
	return Vector2{X: 1, Y: 1}
}

// Up retorna un vector hacia arriba (0, -1)
func Up() Vector2 {
	return Vector2{X: 0, Y: -1}
}

// Down retorna un vector hacia abajo (0, 1)
func Down() Vector2 {
	return Vector2{X: 0, Y: 1}
}

// Left retorna un vector hacia la izquierda (-1, 0)
func Left() Vector2 {
	return Vector2{X: -1, Y: 0}
}

// Right retorna un vector hacia la derecha (1, 0)
func Right() Vector2 {
	return Vector2{X: 1, Y: 0}
}

// ============================================================================
// UTILIDADES MATEMÁTICAS
// ============================================================================

// Clamp limita un valor entre min y max
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Lerp interpolación lineal
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// Sign retorna el signo de un número (-1, 0, 1)
func Sign(value float64) float64 {
	if value > 0 {
		return 1
	}
	if value < 0 {
		return -1
	}
	return 0
}

// Abs valor absoluto
func Abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

// Min retorna el mínimo entre dos valores
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Max retorna el máximo entre dos valores
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
