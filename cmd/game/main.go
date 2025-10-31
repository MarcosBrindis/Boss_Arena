package main

import (
	"log"

	"github.com/MarcosBrindis/boss-arena-go/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Crear instancia del juego
	game := core.NewGame()

	// Configurar ventana de Ebiten
	ebiten.SetWindowSize(core.ScreenWidth, core.ScreenHeight)
	ebiten.SetWindowTitle("Titan's Arena - Boss Rush Demo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// CRÍTICO: Configurar para 60 FPS constantes
	ebiten.SetTPS(core.TargetTPS)    // Ticks por segundo (lógica)
	ebiten.SetVsyncEnabled(true)     // Sincronización vertical
	ebiten.SetMaxTPS(core.TargetTPS) // Limitar TPS

	// Permitir fullscreen con F11
	ebiten.SetFullscreen(false)

	// Iniciar el juego
	log.Println("Iniciando Titan's Arena...")
	log.Printf("Target: %d TPS / %d FPS\n", core.TargetTPS, core.TargetTPS)
	log.Println("Controles: WASD/Arrows = Mover | Space = Saltar | Z = Atacar | X = Dash")
	log.Println("Gamepad: Conecta tu PS5 DualSense para mejor experiencia")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Ejecutar el game loop de Ebiten
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Error al ejecutar el juego:", err)
	}
}
