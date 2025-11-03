package combat

import (
	"math/rand"
	"time"
)

// DamageType representa el tipo de daño
type DamageType int

const (
	DamagePhysical DamageType = iota
	DamageMagic
	DamageTrue // Daño verdadero (ignora defensa)
)

// DamageCalculator calcula el daño con modificadores
type DamageCalculator struct {
	rng *rand.Rand
}

// NewDamageCalculator crea un nuevo calculador de daño
func NewDamageCalculator() *DamageCalculator {
	return &DamageCalculator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CalculateDamage calcula el daño final con todos los modificadores
func (dc *DamageCalculator) CalculateDamage(
	baseDamage int,
	damageType DamageType,
	isCritical bool,
	comboMultiplier float64,
) int {
	damage := float64(baseDamage)

	// Aplicar multiplicador de combo
	damage *= comboMultiplier

	// Aplicar crítico
	if isCritical {
		damage *= 1.5
	}

	// Variación aleatoria ±10%
	variance := 0.9 + dc.rng.Float64()*0.2
	damage *= variance

	finalDamage := int(damage)
	if finalDamage < 1 {
		finalDamage = 1
	}

	return finalDamage
}

// RollCritical determina si un ataque es crítico
func (dc *DamageCalculator) RollCritical(critChance float64) bool {
	return dc.rng.Float64() < critChance
}

// CalculateKnockback calcula el knockback según el daño
func (dc *DamageCalculator) CalculateKnockback(damage int, baseKnockback float64) float64 {
	// Más daño = más knockback
	multiplier := 1.0 + float64(damage)/100.0
	return baseKnockback * multiplier
}
