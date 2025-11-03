package combat

import (
	"image/color"
	"sync"

	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
)

// EffectType representa el tipo de efecto visual
type EffectType int

const (
	EffectHitSpark EffectType = iota
	EffectSlash
	EffectImpact
	EffectBlood
	EffectSmoke
	EffectExplosion
)

// VisualEffect representa un efecto visual temporal
type VisualEffect struct {
	Type     EffectType
	Position utils.Vector2
	Velocity utils.Vector2
	Color    color.RGBA
	Size     float64
	Lifetime int
	Age      int
	IsActive bool
}

// EffectManager maneja todos los efectos visuales (THREAD-SAFE)
type EffectManager struct {
	effects    []*VisualEffect
	maxEffects int
	mu         sync.Mutex // ← NUEVO: Mutex
}

// NewEffectManager crea un nuevo manejador de efectos
func NewEffectManager(maxEffects int) *EffectManager {
	return &EffectManager{
		effects:    make([]*VisualEffect, 0, maxEffects),
		maxEffects: maxEffects,
	}
}

// SpawnEffect crea un nuevo efecto (THREAD-SAFE)
func (em *EffectManager) SpawnEffect(effectType EffectType, position utils.Vector2, color color.RGBA) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if len(em.effects) >= em.maxEffects {
		return
	}

	effect := &VisualEffect{
		Type:     effectType,
		Position: position,
		Color:    color,
		IsActive: true,
	}

	switch effectType {
	case EffectHitSpark:
		effect.Size = 8
		effect.Lifetime = 8
		effect.Velocity = utils.NewVector2(0, -2)

	case EffectSlash:
		effect.Size = 30
		effect.Lifetime = 12

	case EffectImpact:
		effect.Size = 15
		effect.Lifetime = 10

	case EffectExplosion:
		effect.Size = 40
		effect.Lifetime = 15
		effect.Velocity = utils.NewVector2(0, -1)
	}

	em.effects = append(em.effects, effect)
}

// Update actualiza todos los efectos (THREAD-SAFE)
func (em *EffectManager) Update() {
	em.mu.Lock()
	defer em.mu.Unlock()

	for i := len(em.effects) - 1; i >= 0; i-- {
		effect := em.effects[i]

		if !effect.IsActive {
			em.effects = append(em.effects[:i], em.effects[i+1:]...)
			continue
		}

		effect.Age++
		effect.Position = effect.Position.Add(effect.Velocity)

		if effect.Age >= effect.Lifetime {
			effect.IsActive = false
		}
	}
}

// GetActiveEffects retorna una COPIA de los efectos activos (THREAD-SAFE)
func (em *EffectManager) GetActiveEffects() []*VisualEffect {
	em.mu.Lock()
	defer em.mu.Unlock()

	effects := make([]*VisualEffect, len(em.effects))
	copy(effects, em.effects)
	return effects
}

// Clear limpia todos los efectos (THREAD-SAFE)
func (em *EffectManager) Clear() {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.effects = em.effects[:0]
}
