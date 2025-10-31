package input

import "github.com/hajimehoshi/ebiten/v2"

// KeyboardLayout define un layout de teclado
type KeyboardLayout struct {
	// Movimiento
	Left  []ebiten.Key
	Right []ebiten.Key
	Up    []ebiten.Key
	Down  []ebiten.Key

	// Acciones
	Jump    []ebiten.Key
	Attack  []ebiten.Key
	Dash    []ebiten.Key
	Special []ebiten.Key

	// Sistema
	Pause      []ebiten.Key
	Fullscreen []ebiten.Key
	Debug      []ebiten.Key
}

// DefaultKeyboardLayout retorna el layout por defecto
func DefaultKeyboardLayout() *KeyboardLayout {
	return &KeyboardLayout{
		// Movimiento (WASD + Arrows)
		Left:  []ebiten.Key{ebiten.KeyA, ebiten.KeyArrowLeft},
		Right: []ebiten.Key{ebiten.KeyD, ebiten.KeyArrowRight},
		Up:    []ebiten.Key{ebiten.KeyW, ebiten.KeyArrowUp},
		Down:  []ebiten.Key{ebiten.KeyS, ebiten.KeyArrowDown},

		// Acciones
		Jump:    []ebiten.Key{ebiten.KeySpace, ebiten.KeyW, ebiten.KeyArrowUp},
		Attack:  []ebiten.Key{ebiten.KeyZ, ebiten.KeyJ},
		Dash:    []ebiten.Key{ebiten.KeyX, ebiten.KeyK},
		Special: []ebiten.Key{ebiten.KeyC, ebiten.KeyL},

		// Sistema
		Pause:      []ebiten.Key{ebiten.KeyEscape},
		Fullscreen: []ebiten.Key{ebiten.KeyF11},
		Debug:      []ebiten.Key{ebiten.KeyF3},
	}
}

// IsAnyKeyPressed verifica si alguna tecla del slice está presionada
func IsAnyKeyPressed(keys []ebiten.Key) bool {
	for _, key := range keys {
		if ebiten.IsKeyPressed(key) {
			return true
		}
	}
	return false
}
