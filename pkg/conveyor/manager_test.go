package conveyor

import (
	"context"
	"testing"
	"time"
)

func TestManagerScaling(t *testing.T) {
	minWorkers := 1

	manager := CreateManager()
	manager.Start()

	// scale up scenario
	for range 10 {
		CONVEYOR_BELT <- Job{
			Context: context.Background(),
			Consume: func(j any) error {
				time.Sleep(time.Second)
				return nil
			},
		}
	}

	time.Sleep(5 * time.Second) // Let workers scale up

	// check that workers increased
	if len(manager.workers) <= minWorkers {
		t.Errorf("Expected workers to scale up, but only %d running", len(manager.workers))
	}

	// scale down scenario
	time.Sleep(15 * time.Second)

	// should have reduced workers by now
	if len(manager.workers) != minWorkers {
		t.Errorf("Expected workers to scale back down to %d, but got %d", minWorkers, len(manager.workers))
	}

	manager.Stop()
}
