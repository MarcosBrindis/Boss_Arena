package effects

import "sync"

// HitStop maneja el efecto de "freeze frame" al golpear (THREAD-SAFE)
type HitStop struct {
	duration int
	elapsed  int
	isActive bool
	mu       sync.Mutex
}

// NewHitStop crea un nuevo hit stop
func NewHitStop() *HitStop {
	return &HitStop{}
}

// Start inicia un freeze frame (THREAD-SAFE)
func (hs *HitStop) Start(duration int) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	hs.duration = duration
	hs.elapsed = 0
	hs.isActive = true
}

// Update actualiza el hit stop (THREAD-SAFE)
func (hs *HitStop) Update() {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if !hs.isActive {
		return
	}

	hs.elapsed++

	if hs.elapsed >= hs.duration {
		hs.isActive = false
	}
}

// IsActive retorna si el freeze frame está activo (THREAD-SAFE)
func (hs *HitStop) IsActive() bool {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	return hs.isActive
}

// ShouldFreeze retorna true si el juego debe pausarse (THREAD-SAFE)
func (hs *HitStop) ShouldFreeze() bool {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	return hs.isActive
}
