package combat

import (
	"sync"
	"time"

	"github.com/MarcosBrindis/boss-arena-go/internal/utils"
)

// EventType representa el tipo de evento de combate
type EventType int

const (
	EventDamageDealt EventType = iota
	EventDamageTaken
	EventAttackLanded
	EventAttackMissed
	EventKill
	EventComboIncreased
	EventCriticalHit
	EventBlock
	EventParry
	EventDodge
)

// CombatEvent representa un evento de combate
type CombatEvent struct {
	Type       EventType
	Timestamp  time.Time
	Damage     int
	Position   utils.Vector2
	Attacker   string // "player" o "boss"
	Target     string // "player" o "boss"
	IsCritical bool
	ComboCount int
	Metadata   map[string]interface{} // Datos adicionales
}

// EventSystem maneja los eventos de combate usando channels
type EventSystem struct {
	// Channels para eventos
	eventChannel chan CombatEvent
	doneChannel  chan bool

	// Listeners (funciones que reaccionan a eventos)
	listeners map[EventType][]func(CombatEvent)

	// Estadísticas
	stats   *CombatStats
	statsMu sync.Mutex

	// Control
	isRunning bool
}

// CombatStats guarda estadísticas de combate
type CombatStats struct {
	PlayerDamageDealt   int
	PlayerDamageTaken   int
	BossDamageDealt     int
	BossDamageTaken     int
	PlayerAttacksLanded int
	PlayerAttacksMissed int
	BossAttacksLanded   int
	BossAttacksMissed   int
	HighestCombo        int
	TotalHits           int
	CriticalHits        int
	TotalEvents         int
}

// NewEventSystem crea un nuevo sistema de eventos
func NewEventSystem(bufferSize int) *EventSystem {
	return &EventSystem{
		eventChannel: make(chan CombatEvent, bufferSize),
		doneChannel:  make(chan bool),
		listeners:    make(map[EventType][]func(CombatEvent)),
		stats:        &CombatStats{},
		isRunning:    false,
	}
}

// Start inicia el sistema de eventos (goroutine)
func (es *EventSystem) Start() {
	if es.isRunning {
		return
	}
	es.isRunning = true
	go func() {
		for {
			select {
			case event := <-es.eventChannel:
				// Procesar evento
				es.processEvent(event)

			case <-es.doneChannel:
				// Terminar goroutine
				return
			}
		}
	}()
}

// Stop detiene el sistema de eventos
func (es *EventSystem) Stop() {
	if !es.isRunning {
		return
	}

	es.isRunning = false
	es.doneChannel <- true
	close(es.eventChannel)
}

// EmitEvent envía un evento al sistema (non-blocking)
func (es *EventSystem) EmitEvent(event CombatEvent) {
	event.Timestamp = time.Now()

	// Enviar al channel de forma no bloqueante
	select {
	case es.eventChannel <- event:
		// Evento enviado exitosamente
	default:
		// Channel lleno, ignorar evento (evita bloqueos)
	}
}

// processEvent procesa un evento (llamado por la goroutine)
func (es *EventSystem) processEvent(event CombatEvent) {
	// Actualizar estadísticas
	es.updateStats(event)

	// Notificar a los listeners
	if listeners, exists := es.listeners[event.Type]; exists {
		for _, listener := range listeners {
			listener(event)
		}
	}
}

// updateStats actualiza las estadísticas según el evento (THREAD-SAFE)
func (es *EventSystem) updateStats(event CombatEvent) {
	es.statsMu.Lock()         // ← LOCK
	defer es.statsMu.Unlock() // ← UNLOCK

	es.stats.TotalEvents++

	switch event.Type {
	case EventDamageDealt:
		if event.Attacker == "player" {
			es.stats.PlayerDamageDealt += event.Damage
		} else {
			es.stats.BossDamageDealt += event.Damage
		}

	case EventDamageTaken:
		if event.Target == "player" {
			es.stats.PlayerDamageTaken += event.Damage
		} else {
			es.stats.BossDamageTaken += event.Damage
		}

	case EventAttackLanded:
		es.stats.TotalHits++
		if event.Attacker == "player" {
			es.stats.PlayerAttacksLanded++
		} else {
			es.stats.BossAttacksLanded++
		}

	case EventAttackMissed:
		if event.Attacker == "player" {
			es.stats.PlayerAttacksMissed++
		} else {
			es.stats.BossAttacksMissed++
		}

	case EventComboIncreased:
		if event.ComboCount > es.stats.HighestCombo {
			es.stats.HighestCombo = event.ComboCount
		}

	case EventCriticalHit:
		es.stats.CriticalHits++
	}
}

// AddListener añade un listener para un tipo de evento
func (es *EventSystem) AddListener(eventType EventType, listener func(CombatEvent)) {
	es.listeners[eventType] = append(es.listeners[eventType], listener)
}

// GetStats retorna una COPIA de las estadísticas actuales (THREAD-SAFE)
func (es *EventSystem) GetStats() CombatStats {
	es.statsMu.Lock()
	defer es.statsMu.Unlock()

	// Retornar copia (no puntero)
	return *es.stats
}

// ResetStats resetea las estadísticas (THREAD-SAFE)
func (es *EventSystem) ResetStats() {
	es.statsMu.Lock()
	defer es.statsMu.Unlock()

	es.stats = &CombatStats{}
}
