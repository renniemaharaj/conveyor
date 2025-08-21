package conveyor

import (
	"log"
	"sync"
	"time"
)

// The manager shape
type Manager struct {
	mu        sync.Mutex
	workers   []*Worker
	ticker    *time.Ticker
	debugging bool          // Whether manager should log
	quit      chan struct{} // Manager's quit signal channel

	minWorkersAllowed int
	maxWorkersAllowed int
	safeQueueLength   int // if queue length > thresholdUp, scale up

	B *ConveyorBelt // The manager's conveyor belt
}

// Internal blank manager function
func blankManager() *Manager {
	m := &Manager{}
	return m
}

// Create a new manager with default configuration
func CreateManager() *Manager {
	m := blankManager().SetMinWorkers(1).SetMaxWorkers(100).
		SetSafeQueueLength(10).SetTimePerTicker(time.Second / 4)
	m.B = NewConveyorBelt()
	m.quit = make(chan struct{}) // initialize the quit channel

	log.Printf(`
	Creating a manager. Will allow %d min, and %d max workers.
	Will scale resources above %d jobs`+"\n",
		m.minWorkersAllowed, m.maxWorkersAllowed, m.safeQueueLength)

	return m
}

// Manager start function
func (m *Manager) Start() *Manager {
	// initialize min workers
	for range m.minWorkersAllowed {
		m.scaleWorkersUp()
	}

	// routine dynamically scales the manager's workers
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
	return m
}

// Manager stop function will close cleanup its channel and stop ticker
func (m *Manager) Stop() {
	close(m.quit)
	m.ticker.Stop()
}

// Manager's routine check
func (m *Manager) routineCheck() {
	m.mu.Lock()
	defer m.mu.Unlock()

	queueLen := len(m.B.C)
	if queueLen > m.safeQueueLength && len(m.workers) < m.maxWorkersAllowed {
		m.scaleWorkersUp()
		if m.debugging {
			log.Printf("Workers: %d, after scaling up. Jobs queued: %d\n", len(m.workers), queueLen)
		}
	} else if queueLen <= m.safeQueueLength && len(m.workers) > m.minWorkersAllowed {
		m.scaleWorkersDown()
		if m.debugging {
			log.Printf("Workers: %d, after scaling down. Jobs queued: %d\n", len(m.workers), queueLen)
		}
	}
}

// scaleWorkersUp internal function, creates and starts a new worker
func (m *Manager) scaleWorkersUp() {
	w := NewWorker(m.B, len(m.workers)) // Create a worker by assigning the manager's conveyor belt
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
