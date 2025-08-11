package conveyor

import (
	"log"
	"sync"
	"time"
)

// The manager shape
type Manager struct {
	mu         sync.Mutex
	workers    []*Worker
	ticker     *time.Ticker
	quit       chan struct{}
	minWorkers int
	maxWorkers int
	stepUpAt   int // if queue length > thresholdUp, scale up
	stepDownAt int // if queue length < thresholdDown, scale down
}

// Manager start function
func (m *Manager) Start() {
	// Initialize min workers
	for range m.minWorkers {
		m.scaleWorkersUp()
	}

	// Routine
	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.routineCheck()
			case <-m.quit:
				m.stopAll()
				return
			}
		}
	}()
}

// Manager
func (m *Manager) Stop() {
	close(m.quit)
	m.ticker.Stop()
}

// Manager's routine check
func (m *Manager) routineCheck() {
	m.mu.Lock()
	defer m.mu.Unlock()

	queueLen := len(CONVEYOR_BELT)

	if queueLen > m.stepUpAt && len(m.workers) < m.maxWorkers {
		log.Printf("At %d current workers, scaling up. Jobs queued: %d\n", len(m.workers), queueLen)
		m.scaleWorkersUp()
	} else if queueLen < m.stepDownAt && len(m.workers) > m.minWorkers {
		log.Printf("At %d current workers, scaling down. Jobs queued: %d\n", len(m.workers), queueLen)
		m.scaleWorkersDown()
	}
}

// scaleWorkersUp internal function, creates and starts a new worker
func (m *Manager) scaleWorkersUp() {
	w := &Worker{}
	m.workers = append(m.workers, w)
	go w.Start()
}

// scaleWorkersDown internal function, stops the last worker safely and removes
func (m *Manager) scaleWorkersDown() {
	if len(m.workers) == 0 {
		return
	}

	last := m.workers[len(m.workers)-1]
	last.Stop() // Worker will complete its current job before
	m.workers = m.workers[:len(m.workers)-1]
}

// stopAll function safely stops all works
func (m *Manager) stopAll() {
	for _, w := range m.workers {
		w.Stop()
	}
	m.workers = nil
}
