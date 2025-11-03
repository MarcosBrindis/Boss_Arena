package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MarcosBrindis/boss-arena-go/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Crear el juego
	game := core.NewGame()

	// Setup para limpiar recursos al cerrar
	setupCleanup(game)

	// Configurar ventana
	ebiten.SetWindowSize(core.ScreenWidth, core.ScreenHeight)
	ebiten.SetWindowTitle("Titan's Arena - Boss Rush Demo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Ejecutar el juego
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// setupCleanup configura la limpieza de recursos al cerrar
func setupCleanup(game *core.Game) {
	// Capturar señales de cierre
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		// Limpiar recursos
		game.Cleanup()
		os.Exit(0)
	}()
}
