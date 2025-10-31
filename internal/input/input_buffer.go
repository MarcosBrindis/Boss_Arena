package input

// InputBuffer implementa un buffer de inputs para saltos más permisivos
// Ejemplo: Si presionas saltar 5 frames antes de tocar el suelo, aún funciona
type InputBuffer struct {
	maxFrames int // Cantidad de frames que dura el buffer
	frames    int // Frames restantes
	active    bool
}

// NewInputBuffer crea un nuevo buffer de input
func NewInputBuffer(maxFrames int) *InputBuffer {
	return &InputBuffer{
		maxFrames: maxFrames,
		frames:    0,
		active:    false,
	}
}

// Activate activa el buffer (cuando se presiona el botón)
func (b *InputBuffer) Activate() {
	b.frames = b.maxFrames
	b.active = true
}

// Update actualiza el buffer (llamar cada frame)
func (b *InputBuffer) Update() {
	if b.active && b.frames > 0 {
		b.frames--
		if b.frames <= 0 {
			b.active = false
		}
	}
}

// IsActive retorna true si el buffer está activo
func (b *InputBuffer) IsActive() bool {
	return b.active && b.frames > 0
}

// Consume consume el buffer (cuando se usa el input)
func (b *InputBuffer) Consume() bool {
	if b.IsActive() {
		b.active = false
		b.frames = 0
		return true
	}
	return false
}

// Reset resetea el buffer
func (b *InputBuffer) Reset() {
	b.active = false
	b.frames = 0
}

// GetFramesRemaining retorna los frames restantes del buffer
func (b *InputBuffer) GetFramesRemaining() int {
	if b.active {
		return b.frames
	}
	return 0
}

// ============================================================================
// COYOTE TIME
// ============================================================================

// CoyoteTimer implementa "coyote time" - tiempo extra para saltar después de caer
// Ejemplo: Tienes 6 frames extra para saltar después de dejar una plataforma
type CoyoteTimer struct {
	maxFrames int
	frames    int
	active    bool
}

// NewCoyoteTimer crea un nuevo timer de coyote
func NewCoyoteTimer(maxFrames int) *CoyoteTimer {
	return &CoyoteTimer{
		maxFrames: maxFrames,
		frames:    0,
		active:    false,
	}
}

// Start inicia el coyote time
func (c *CoyoteTimer) Start() {
	c.frames = c.maxFrames
	c.active = true
}

// Update actualiza el timer (llamar cada frame)
func (c *CoyoteTimer) Update() {
	if c.active && c.frames > 0 {
		c.frames--
		if c.frames <= 0 {
			c.active = false
		}
	}
}

// IsActive retorna true si el coyote time está activo
func (c *CoyoteTimer) IsActive() bool {
	return c.active && c.frames > 0
}

// Stop detiene el coyote time
func (c *CoyoteTimer) Stop() {
	c.active = false
	c.frames = 0
}

// GetFramesRemaining retorna los frames restantes
func (c *CoyoteTimer) GetFramesRemaining() int {
	if c.active {
		return c.frames
	}
	return 0
}
