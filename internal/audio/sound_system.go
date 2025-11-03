package audio

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

// SoundSystem maneja la reproducción de sonidos (estructura para futuro)
type SoundSystem struct {
	volume     float64
	isMuted    bool
	soundQueue []SoundType
}

// NewSoundSystem crea un nuevo sistema de audio
func NewSoundSystem() *SoundSystem {
	return &SoundSystem{
		volume:     1.0,
		isMuted:    false,
		soundQueue: make([]SoundType, 0, 10),
	}
}

// PlaySound añade un sonido a la cola (placeholder)
func (ss *SoundSystem) PlaySound(soundType SoundType) {
	if ss.isMuted {
		return
	}

	// TODO: En Módulo 10 implementaremos audio real con Ebitengine Audio
	ss.soundQueue = append(ss.soundQueue, soundType)
}

// SetVolume ajusta el volumen
func (ss *SoundSystem) SetVolume(volume float64) {
	if volume < 0 {
		volume = 0
	}
	if volume > 1 {
		volume = 1
	}
	ss.volume = volume
}

// ToggleMute alterna el mute
func (ss *SoundSystem) ToggleMute() {
	ss.isMuted = !ss.isMuted
}

// ClearQueue limpia la cola de sonidos
func (ss *SoundSystem) ClearQueue() {
	ss.soundQueue = ss.soundQueue[:0]
}
