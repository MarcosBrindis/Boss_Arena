package core

import "time"

// Config contiene la configuración global del juego
type Config struct {
	// Performance
	EnableVSync     bool
	TargetFPS       int
	ShowDebugInfo   bool
	EnableProfiling bool

	// Gameplay
	DifficultyLevel int // 1 = Fácil, 2 = Normal, 3 = Difícil
	PlayerStartHP   int
	BossStartHP     int
	MeteorSpawnRate time.Duration

	// Concurrencia
	NumPhysicsWorkers    int
	NumProjectileWorkers int
	NumMeteorWorkers     int
	ChannelBufferSize    int

	// Input
	EnableGamepad    bool
	GamepadDeadzone  float64
	JumpBufferFrames int
	CoyoteTimeFrames int
}

// DefaultConfig retorna la configuración por defecto
func DefaultConfig() *Config {
	return &Config{
		// Performance
		EnableVSync:     true,
		TargetFPS:       60,
		ShowDebugInfo:   true, // true en desarrollo, false en producción
		EnableProfiling: false,

		// Gameplay
		DifficultyLevel: 2, // Normal
		PlayerStartHP:   100,
		BossStartHP:     500,
		MeteorSpawnRate: 5 * time.Second,

		// Concurrencia
		NumPhysicsWorkers:    4,
		NumProjectileWorkers: 4,
		NumMeteorWorkers:     8,
		ChannelBufferSize:    100,

		// Input
		EnableGamepad:    true,
		GamepadDeadzone:  0.2, // 20% deadzone para sticks analógicos
		JumpBufferFrames: 5,   // Buffer de 5 frames para salto
		CoyoteTimeFrames: 6,   // 6 frames de coyote time
	}
}

// DevConfig retorna configuración para desarrollo (más debug info)
func DevConfig() *Config {
	cfg := DefaultConfig()
	cfg.ShowDebugInfo = true
	cfg.EnableProfiling = true
	cfg.BossStartHP = 200                 // Boss más débil para testear rápido
	cfg.MeteorSpawnRate = 3 * time.Second // Meteoros más frecuentes
	return cfg
}

// ProductionConfig retorna configuración para producción
func ProductionConfig() *Config {
	cfg := DefaultConfig()
	cfg.ShowDebugInfo = false
	cfg.EnableProfiling = false
	return cfg
}
