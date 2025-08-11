package conveyor

import (
	"log"
	"time"
)

func NewManager() *Manager {
	m := &Manager{}
	return m
}

// Create a new manager with default configuration
func CreateManager() *Manager {
	m := NewManager().SetMinWorkers(1).SetMaxWorkers(100).
		SetStepUpAt(10).SetStepDownAt(10).SetTimePerTicker(time.Second / 4)

	log.Printf(`
	Creating a manager. Will allow %d min, and %d max workers.
	Will scale up at %d, and scale down at %d jobs`+"\n",
		m.minWorkers, m.maxWorkers, m.stepUpAt, m.stepDownAt)

	return m
}

// Set the manager's min workers allowed
func (m *Manager) SetMinWorkers(mw int) *Manager {
	m.minWorkers = mw
	return m
}

// Set the manager's max workers allowed
func (m *Manager) SetMaxWorkers(mw int) *Manager {
	m.maxWorkers = mw
	return m
}

// Set the manager's stepUpAt for threshold of jobs
func (m *Manager) SetStepUpAt(s int) *Manager {
	m.stepUpAt = s
	return m
}

// Set the manager's stepDownAt for threshold of jobs
func (m *Manager) SetStepDownAt(s int) *Manager {
	m.stepDownAt = s
	return m
}

// Set the manager's time per ticker (dynamically changing)
func (m *Manager) SetTimePerTicker(t time.Duration) *Manager {
	if m.ticker != nil {
		m.ticker.Stop()
	}
	m.ticker = time.NewTicker(t)
	return m
}
