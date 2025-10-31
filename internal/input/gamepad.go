package input

// GamepadButton representa los botones del gamepad PS5 DualSense
type GamepadButton int

const (
	// Face Buttons (PlayStation layout)
	ButtonCross    GamepadButton = 0 // □ / ⬜
	ButtonCircle   GamepadButton = 1 // X / ✕
	ButtonSquare   GamepadButton = 2 // O / ⚪
	ButtonTriangle GamepadButton = 3 // △ / 🔺

	// Shoulder Buttons
	ButtonL1 GamepadButton = 4
	ButtonR1 GamepadButton = 5
	ButtonL2 GamepadButton = 6
	ButtonR2 GamepadButton = 7

	// Special Buttons
	ButtonShare   GamepadButton = 8  // Share / Create
	ButtonOptions GamepadButton = 9  // Options
	ButtonL3      GamepadButton = 10 // Left Stick Click
	ButtonR3      GamepadButton = 11 // Right Stick Click

	// D-Pad
	ButtonDPadUp    GamepadButton = 12
	ButtonDPadDown  GamepadButton = 13
	ButtonDPadLeft  GamepadButton = 14
	ButtonDPadRight GamepadButton = 15

	// PlayStation Button
	ButtonPS GamepadButton = 16 // PlayStation button
)

// GamepadAxis representa los ejes analógicos
type GamepadAxis int

const (
	AxisLeftStickX  GamepadAxis = 0
	AxisLeftStickY  GamepadAxis = 1
	AxisRightStickX GamepadAxis = 2
	AxisRightStickY GamepadAxis = 3
	AxisL2Trigger   GamepadAxis = 4 // Analog trigger
	AxisR2Trigger   GamepadAxis = 5 // Analog trigger
)

// GamepadMapping contiene el mapeo de botones para el juego
type GamepadMapping struct {
	Jump    GamepadButton
	Attack  GamepadButton
	Dash    GamepadButton
	Special GamepadButton
}

// PS5Mapping retorna el mapeo para PS5 DualSense
func PS5Mapping() *GamepadMapping {
	return &GamepadMapping{
		Jump:    ButtonCross,    // X para saltar
		Attack:  ButtonSquare,   // Square para atacar
		Dash:    ButtonCircle,   // Circle para dash
		Special: ButtonTriangle, // Triangle para especial
	}
}

// XboxMapping retorna el mapeo para Xbox controller
func XboxMapping() *GamepadMapping {
	return &GamepadMapping{
		Jump:    ButtonCross,    // A para saltar
		Attack:  ButtonSquare,   // X para atacar
		Dash:    ButtonCircle,   // B para dash
		Special: ButtonTriangle, // Y para especial
	}
}
