package projectiles

import (
	"context"
	"sync"
)

// WorkerPool procesa proyectiles usando múltiples workers
type WorkerPool struct {
	// Número de workers
	numWorkers int

	// Channel de trabajo
	workQueue chan *Projectile

	// Control de concurrencia
	wg sync.WaitGroup

	// Contexto para cancelación
	ctx    context.Context
	cancel context.CancelFunc
}

// NewWorkerPool crea un nuevo worker pool
func NewWorkerPool(numWorkers int, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		numWorkers: numWorkers,
		workQueue:  make(chan *Projectile, queueSize),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start inicia los workers
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker es una goroutine que procesa proyectiles
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	_ = id

	for {
		select {
		case projectile := <-wp.workQueue:
			// Procesar proyectil
			if projectile != nil && projectile.IsActive {
				projectile.Update()
			}

		case <-wp.ctx.Done():
			// Cancelación recibida
			return
		}
	}
}

// Submit envía un proyectil para procesamiento
func (wp *WorkerPool) Submit(projectile *Projectile) {
	select {
	case wp.workQueue <- projectile:
		// Enviado a la cola
	default:
		// Cola llena, procesar en el hilo actual (fallback)
		projectile.Update()
	}
}

// Stop detiene el worker pool
func (wp *WorkerPool) Stop() {
	wp.cancel()
	wp.wg.Wait()
	close(wp.workQueue)
}
