// internal/audio/sound_system.go
package audio

import "sync"

// SoundType representa un tipo de sonido
type SoundType int

const (
	SoundHit SoundType = iota
	SoundSlash
	SoundExplosion
	SoundJump
	SoundDash
	SoundBossRoar
	SoundPlayerHurt
	SoundVictory
	SoundGameOver
)

// SoundSystem maneja la reproducción de sonidos (THREAD-SAFE)
type SoundSystem struct {
	volume     float64
	isMuted    bool
	soundQueue []SoundType
	mu         sync.Mutex // ← NUEVO: Mutex para proteger soundQueue
}

// NewSoundSystem crea un nuevo sistema de audio
func NewSoundSystem() *SoundSystem {
	return &SoundSystem{
		volume:     1.0,
		isMuted:    false,
		soundQueue: make([]SoundType, 0, 10),
	}
}

// PlaySound añade un sonido a la cola (THREAD-SAFE)
func (ss *SoundSystem) PlaySound(soundType SoundType) {
	ss.mu.Lock()         // ← LOCK
	defer ss.mu.Unlock() // ← UNLOCK

	if ss.isMuted {
		return
	}

	// TODO: En Módulo 10 implementaremos audio real con Ebitengine Audio
	ss.soundQueue = append(ss.soundQueue, soundType)
}

// SetVolume ajusta el volumen (THREAD-SAFE)
func (ss *SoundSystem) SetVolume(volume float64) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if volume < 0 {
		volume = 0
	}
	if volume > 1 {
		volume = 1
	}
	ss.volume = volume
}

// ToggleMute alterna el mute (THREAD-SAFE)
func (ss *SoundSystem) ToggleMute() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.isMuted = !ss.isMuted
}

// ClearQueue limpia la cola de sonidos (THREAD-SAFE)
func (ss *SoundSystem) ClearQueue() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.soundQueue = ss.soundQueue[:0]
}

// GetQueueSize retorna el tamaño de la cola (THREAD-SAFE)
func (ss *SoundSystem) GetQueueSize() int {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	return len(ss.soundQueue)
}
