package effects

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
)

// ScreenShake maneja el efecto de sacudida de pantalla (THREAD-SAFE)
type ScreenShake struct {
	intensity float64
	duration  int
	elapsed   int
	offset    utils.Vector2
	decay     float64
	isActive  bool
	rng       *rand.Rand
	mu        sync.Mutex
}

// NewScreenShake crea un nuevo screen shake
func NewScreenShake() *ScreenShake {
	return &ScreenShake{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Start inicia una sacudida de pantalla (THREAD-SAFE)
func (ss *ScreenShake) Start(intensity float64, duration int) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.intensity = intensity
	ss.duration = duration
	ss.elapsed = 0
	ss.decay = 0.95
	ss.isActive = true
}

// Update actualiza el screen shake (THREAD-SAFE)
func (ss *ScreenShake) Update() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if !ss.isActive {
		ss.offset = utils.Zero()
		return
	}

	ss.elapsed++

	// Calcular intensidad actual con decay
	currentIntensity := ss.intensity * math.Pow(ss.decay, float64(ss.elapsed))

	// Generar offset aleatorio
	angle := ss.rng.Float64() * 2 * math.Pi
	ss.offset = utils.Vector2{
		X: math.Cos(angle) * currentIntensity,
		Y: math.Sin(angle) * currentIntensity,
	}

	// Terminar si cumplió la duración
	if ss.elapsed >= ss.duration {
		ss.isActive = false
		ss.offset = utils.Zero()
	}
}

// GetOffset retorna el offset actual de la cámara (THREAD-SAFE)
func (ss *ScreenShake) GetOffset() utils.Vector2 {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	return ss.offset
}

// IsActive retorna si el shake está activo (THREAD-SAFE)
func (ss *ScreenShake) IsActive() bool {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	return ss.isActive
}

// Stop detiene el shake inmediatamente (THREAD-SAFE)
func (ss *ScreenShake) Stop() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.isActive = false
	ss.offset = utils.Zero()
}
