# Titan's Arena - Boss Rush Demo

Demo de un juego estilo **metroidvania/souls-like** en Go usando concurrencia y Ebiten.

## Requisitos del Proyecto

**Asignatura**: Programación Concurrente  
**Tecnología**: Go 1.20+ con EbitenEngine  
**Objetivo**: Demostrar patrones de concurrencia en un juego 2D interactivo

## Características

- ✅ **Concurrencia**: Múltiples goroutines, channels, worker pools
- ✅ **Sincronización**: Mutex, WaitGroup, Channels
- ✅ **Patrones**: Producer-Consumer, Worker Pool, Pipeline, Fan-out/Fan-in
- ✅ **Juego**: Boss fight con mecánicas tipo Hollow Knight/Dark Souls
- ✅ **Performance**: 60 FPS constantes, sin race conditions

## Instalación

### 1. Clonar el repositorio
```bash
git clone https://github.com/MarcosBrindis/boss-arena-go.git
cd boss-arena-go

## controles
GAMEPAD PS5:
├─ Button 0 → Square   (⬜) → ATTACK
├─ Button 1 → Cross    (✕) → JUMP
├─ Button 2 → Circle   (⚪) → DASH
├─ Button 3 → Triangle (🔺) → SPECIAL
└─ Button 7 → R2            → DASH (alternativo)