package projectiles

import (
	"context"
	"sync"

	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

// ProjectileManager maneja todos los proyectiles del juego
type ProjectileManager struct {
	// Proyectiles activos
	projectiles []*Projectile
	mu          sync.Mutex

	// Object Pool
	pool *ProjectilePool

	// Worker Pool (opcional, para muchos proyectiles)
	workerPool    *WorkerPool
	useWorkerPool bool

	// Contexto
	ctx    context.Context
	cancel context.CancelFunc

	// Estadísticas
	activeCount int
}

// NewProjectileManager crea un nuevo manager de proyectiles
func NewProjectileManager(poolSize int, useWorkerPool bool) *ProjectileManager {
	ctx, cancel := context.WithCancel(context.Background())

	pm := &ProjectileManager{
		projectiles:   make([]*Projectile, 0, 100),
		pool:          NewProjectilePool(poolSize),
		ctx:           ctx,
		cancel:        cancel,
		useWorkerPool: useWorkerPool,
	}

	// Iniciar worker pool si está habilitado
	if useWorkerPool {
		pm.workerPool = NewWorkerPool(4, 100) // 4 workers
		pm.workerPool.Start()
	}

	return pm
}

// Spawn crea un nuevo proyectil
func (pm *ProjectileManager) Spawn(
	projectileType ProjectileType,
	position, direction utils.Vector2,
	owner string,
) *Projectile {
	// Obtener del pool
	projectile := pm.pool.Get(projectileType, position, direction, owner)

	// Añadir a la lista
	pm.mu.Lock()
	pm.projectiles = append(pm.projectiles, projectile)
	pm.activeCount++
	pm.mu.Unlock()

	return projectile
}

// Update actualiza todos los proyectiles
func (pm *ProjectileManager) Update() {
	pm.mu.Lock()
	// Usar workers solo si HAY SUFICIENTES PROYECTILES

	const WORKER_POOL_THRESHOLD = 30 // Umbral mínimo

	activeCount := 0
	for _, p := range pm.projectiles {
		if p.IsActive {
			activeCount++
		}
	}

	// Usar worker pool solo si hay MUCHOS proyectiles
	useWorkers := pm.useWorkerPool && activeCount >= WORKER_POOL_THRESHOLD

	if useWorkers {
		// MODO PARALELO: Muchos proyectiles (30+)

		var wg sync.WaitGroup

		for _, p := range pm.projectiles {
			if p.IsActive {
				wg.Add(1)
				p := p // Capturar variable
				go func() {
					defer wg.Done()
					p.Update()
				}()
			}
		}

		pm.mu.Unlock() // Desbloquear antes de esperar
		wg.Wait()      // Esperar a que todos terminen
		pm.mu.Lock()   // Volver a bloquear

	} else {
		// MODO SECUENCIAL: Pocos proyectiles (<30)

		for _, p := range pm.projectiles {
			if p.IsActive {
				p.Update()
			}
		}
	}

	// Limpiar proyectiles inactivos y devolverlos al pool
	activeProj := pm.projectiles[:0]
	for _, p := range pm.projectiles {
		if p.IsActive {
			activeProj = append(activeProj, p)
		} else {
			pm.pool.Put(p)
		}
	}
	pm.projectiles = activeProj
	pm.activeCount = len(pm.projectiles)

	pm.mu.Unlock()
}

// Draw dibuja todos los proyectiles
func (pm *ProjectileManager) Draw(screen *ebiten.Image) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, p := range pm.projectiles {
		if p.IsActive {
			p.Draw(screen)
		}
	}
}

// GetActiveProjectiles retorna una copia de los proyectiles activos
func (pm *ProjectileManager) GetActiveProjectiles() []*Projectile {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Retornar copia para evitar data races
	projectiles := make([]*Projectile, len(pm.projectiles))
	copy(projectiles, pm.projectiles)
	return projectiles
}

// GetProjectilesByOwner retorna proyectiles de un propietario específico
func (pm *ProjectileManager) GetProjectilesByOwner(owner string) []*Projectile {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	result := make([]*Projectile, 0)
	for _, p := range pm.projectiles {
		if p.IsActive && p.Owner == owner {
			result = append(result, p)
		}
	}
	return result
}

// Clear limpia todos los proyectiles
func (pm *ProjectileManager) Clear() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Devolver todos al pool
	for _, p := range pm.projectiles {
		pm.pool.Put(p)
	}

	pm.projectiles = pm.projectiles[:0]
	pm.activeCount = 0
}

// Cleanup limpia recursos
func (pm *ProjectileManager) Cleanup() {
	pm.cancel()

	if pm.workerPool != nil {
		pm.workerPool.Stop()
	}

	pm.Clear()
	pm.pool.Clear()
}

// GetStats retorna estadísticas
func (pm *ProjectileManager) GetStats() (active, created, reused int) {
	pm.mu.Lock()
	active = pm.activeCount
	pm.mu.Unlock()

	created, reused = pm.pool.GetStats()
	return
}
