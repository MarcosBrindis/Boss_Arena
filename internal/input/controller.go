package input

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputMethod representa el método de input activo
type InputMethod int

const (
	InputKeyboard InputMethod = iota
	InputGamepad
)

// Controller maneja todos los inputs del juego (teclado + gamepad)
type Controller struct {
	// Estado actual
	inputMethod InputMethod
	gamepadID   ebiten.GamepadID

	// Configuración
	deadzone float64

	// Estados de teclas/botones (para detectar "JustPressed")
	// Movimiento
	leftPressedLastFrame  bool
	rightPressedLastFrame bool
	upPressedLastFrame    bool
	downPressedLastFrame  bool

	// Acciones
	jumpPressedLastFrame    bool
	attackPressedLastFrame  bool
	dashPressedLastFrame    bool
	specialPressedLastFrame bool

	// Buffers (implementados en input_buffer.go)
	jumpBuffer  *InputBuffer
	coyoteTimer *CoyoteTimer

	// Valores actuales de ejes analógicos
	horizontalAxis float64
	verticalAxis   float64
}

// NewController crea un nuevo controlador de input
func NewController(deadzone float64, jumpBufferFrames, coyoteFrames int) *Controller {
	return &Controller{
		inputMethod: InputKeyboard,
		gamepadID:   -1,
		deadzone:    deadzone,
		jumpBuffer:  NewInputBuffer(jumpBufferFrames),
		coyoteTimer: NewCoyoteTimer(coyoteFrames),
	}
}

// Update actualiza el estado del controlador (llamar cada frame)
func (c *Controller) Update() {
	// Detectar gamepad conectado
	c.detectGamepad()

	// Actualizar ejes analógicos
	c.updateAxes()

	// Actualizar buffers
	c.jumpBuffer.Update()
	c.coyoteTimer.Update()

	// Actualizar estados de teclas para la próxima frame
	c.updateButtonStates()
}

// detectGamepad detecta si hay un gamepad conectado
func (c *Controller) detectGamepad() {
	ids := ebiten.AppendGamepadIDs(nil)
	if len(ids) > 0 {
		c.gamepadID = ids[0]
		c.inputMethod = InputGamepad
	} else {
		c.gamepadID = -1
		c.inputMethod = InputKeyboard
	}
}

// updateAxes actualiza los valores de los ejes analógicos
func (c *Controller) updateAxes() {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		// Stick izquierdo - Eje horizontal (Axis 0)
		axisX := ebiten.GamepadAxisValue(c.gamepadID, 0)
		if abs(axisX) < c.deadzone {
			c.horizontalAxis = 0
		} else {
			c.horizontalAxis = axisX
		}

		// Stick izquierdo - Eje vertical (Axis 1)
		axisY := ebiten.GamepadAxisValue(c.gamepadID, 1)
		if abs(axisY) < c.deadzone {
			c.verticalAxis = 0
		} else {
			c.verticalAxis = axisY
		}
	} else {
		// Teclado - simular ejes digitales
		c.horizontalAxis = 0
		c.verticalAxis = 0

		if c.IsLeftHeld() {
			c.horizontalAxis = -1
		}
		if c.IsRightHeld() {
			c.horizontalAxis = 1
		}
		if c.IsUpHeld() {
			c.verticalAxis = -1
		}
		if c.IsDownHeld() {
			c.verticalAxis = 1
		}
	}
}

// updateButtonStates guarda el estado actual para la próxima frame
func (c *Controller) updateButtonStates() {
	c.leftPressedLastFrame = c.IsLeftHeld()
	c.rightPressedLastFrame = c.IsRightHeld()
	c.upPressedLastFrame = c.IsUpHeld()
	c.downPressedLastFrame = c.IsDownHeld()
	c.jumpPressedLastFrame = c.IsJumpHeld()
	c.attackPressedLastFrame = c.IsAttackHeld()
	c.dashPressedLastFrame = c.IsDashHeld()
	c.specialPressedLastFrame = c.IsSpecialHeld()
}

// ============================================================================
// MÉTODOS DE CONSULTA - MOVIMIENTO
// ============================================================================

// GetHorizontalAxis retorna el eje horizontal [-1, 1]
func (c *Controller) GetHorizontalAxis() float64 {
	return c.horizontalAxis
}

// GetVerticalAxis retorna el eje vertical [-1, 1]
func (c *Controller) GetVerticalAxis() float64 {
	return c.verticalAxis
}

// IsLeftHeld retorna true si izquierda está presionada
func (c *Controller) IsLeftHeld() bool {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		// D-Pad Left
		return ebiten.IsGamepadButtonPressed(c.gamepadID, ebiten.GamepadButton14)
	}
	return ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
}

// IsRightHeld retorna true si derecha está presionada
func (c *Controller) IsRightHeld() bool {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		// D-Pad Right
		return ebiten.IsGamepadButtonPressed(c.gamepadID, ebiten.GamepadButton15)
	}
	return ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight)
}

// IsUpHeld retorna true si arriba está presionada
func (c *Controller) IsUpHeld() bool {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		// D-Pad Up
		return ebiten.IsGamepadButtonPressed(c.gamepadID, ebiten.GamepadButton12)
	}
	return ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp)
}

// IsDownHeld retorna true si abajo está presionada
func (c *Controller) IsDownHeld() bool {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		// D-Pad Down
		return ebiten.IsGamepadButtonPressed(c.gamepadID, ebiten.GamepadButton13)
	}
	return ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown)
}

// ============================================================================
// MÉTODOS DE CONSULTA - ACCIONES
// ============================================================================

// IsJumpHeld retorna true si el botón de salto está presionado
func (c *Controller) IsJumpHeld() bool {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		// PS5: X/Cross (Button 1)
		return ebiten.IsGamepadButtonPressed(c.gamepadID, ebiten.GamepadButton1)
	}
	return ebiten.IsKeyPressed(ebiten.KeySpace) ||
		ebiten.IsKeyPressed(ebiten.KeyW) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowUp)
}

// IsJumpPressed retorna true solo en el frame que se presiona (con buffer)
func (c *Controller) IsJumpPressed() bool {
	var pressed bool

	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		pressed = inpututil.IsGamepadButtonJustPressed(c.gamepadID, ebiten.GamepadButton1)
	} else {
		pressed = inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			inpututil.IsKeyJustPressed(ebiten.KeyW) ||
			inpututil.IsKeyJustPressed(ebiten.KeyArrowUp)
	}

	// Si se presionó, activar el buffer
	if pressed {
		c.jumpBuffer.Activate()
	}

	return pressed
}

// ConsumeJumpBuffer retorna true si hay un salto en el buffer y lo consume
func (c *Controller) ConsumeJumpBuffer() bool {
	return c.jumpBuffer.Consume()
}

// IsAttackHeld retorna true si el botón de ataque está presionado
func (c *Controller) IsAttackHeld() bool {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		// PS5: Square (Button 0)
		return ebiten.IsGamepadButtonPressed(c.gamepadID, ebiten.GamepadButton0)
	}
	return ebiten.IsKeyPressed(ebiten.KeyZ) || ebiten.IsKeyPressed(ebiten.KeyJ)
}

// IsAttackPressed retorna true solo en el frame que se presiona
func (c *Controller) IsAttackPressed() bool {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		return inpututil.IsGamepadButtonJustPressed(c.gamepadID, ebiten.GamepadButton0)
	}
	return inpututil.IsKeyJustPressed(ebiten.KeyZ) || inpututil.IsKeyJustPressed(ebiten.KeyJ)
}

// IsDashHeld retorna true si el botón de dash está presionado
func (c *Controller) IsDashHeld() bool {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		// PS5: Circle (Button 2) O R2 (Button 7)
		return ebiten.IsGamepadButtonPressed(c.gamepadID, ebiten.GamepadButton2) ||
			ebiten.IsGamepadButtonPressed(c.gamepadID, ebiten.GamepadButton7) // R2
	}
	// Teclado: X/K O SHIFT (solo hay uno en Ebiten)
	return ebiten.IsKeyPressed(ebiten.KeyX) ||
		ebiten.IsKeyPressed(ebiten.KeyK) ||
		ebiten.IsKeyPressed(ebiten.KeyShift)
}

// IsDashPressed retorna true solo en el frame que se presiona
func (c *Controller) IsDashPressed() bool {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		// Circle O R2
		return inpututil.IsGamepadButtonJustPressed(c.gamepadID, ebiten.GamepadButton2) ||
			inpututil.IsGamepadButtonJustPressed(c.gamepadID, ebiten.GamepadButton7) // R2
	}
	// X/K O SHIFT
	return inpututil.IsKeyJustPressed(ebiten.KeyX) ||
		inpututil.IsKeyJustPressed(ebiten.KeyK) ||
		inpututil.IsKeyJustPressed(ebiten.KeyShift)
}

// IsSpecialHeld retorna true si el botón especial está presionado
func (c *Controller) IsSpecialHeld() bool {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		// PS5: Triangle (Button 3)
		return ebiten.IsGamepadButtonPressed(c.gamepadID, ebiten.GamepadButton3)
	}
	return ebiten.IsKeyPressed(ebiten.KeyC) || ebiten.IsKeyPressed(ebiten.KeyL)
}

// IsSpecialPressed retorna true solo en el frame que se presiona
func (c *Controller) IsSpecialPressed() bool {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		return inpututil.IsGamepadButtonJustPressed(c.gamepadID, ebiten.GamepadButton3)
	}
	return inpututil.IsKeyJustPressed(ebiten.KeyC) || inpututil.IsKeyJustPressed(ebiten.KeyL)
}

// ============================================================================
// COYOTE TIME
// ============================================================================

// StartCoyoteTime inicia el contador de coyote time
func (c *Controller) StartCoyoteTime() {
	c.coyoteTimer.Start()
}

// HasCoyoteTime retorna true si aún hay tiempo de coyote disponible
func (c *Controller) HasCoyoteTime() bool {
	return c.coyoteTimer.IsActive()
}

// ============================================================================
// VIBRACIÓN (GAMEPAD)
// ============================================================================

// Vibrate hace vibrar el gamepad (solo si está conectado)
func (c *Controller) Vibrate(durationMS int, strength float64) {
	if c.inputMethod == InputGamepad && c.gamepadID >= 0 {
		// Crear opciones de vibración
		options := &ebiten.VibrateGamepadOptions{
			Duration:        time.Duration(durationMS) * time.Millisecond,
			StrongMagnitude: strength,
			WeakMagnitude:   strength,
		}
		ebiten.VibrateGamepad(c.gamepadID, options)
	}
}

// ============================================================================
// INFORMACIÓN
// ============================================================================

// GetInputMethod retorna el método de input actual
func (c *Controller) GetInputMethod() InputMethod {
	return c.inputMethod
}

// IsGamepadConnected retorna true si hay un gamepad conectado
func (c *Controller) IsGamepadConnected() bool {
	return c.inputMethod == InputGamepad && c.gamepadID >= 0
}

// GetGamepadName retorna el nombre del gamepad conectado
func (c *Controller) GetGamepadName() string {
	if c.gamepadID >= 0 {
		return ebiten.GamepadName(c.gamepadID)
	}
	return "No gamepad"
}

// ============================================================================
// UTILIDADES
// ============================================================================

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
