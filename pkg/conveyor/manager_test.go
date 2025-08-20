package conveyor

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestManagerScaling(t *testing.T) {
	minWorkers := 1

	manager := CreateManager().SetTimePerTicker(time.Second / 10)
	manager.Start()

	// scale up scenario
	for range 100 {
		manager.B.Push(&Job{
			Context: context.Background(),
			Consume: func(j any) error {
				time.Sleep(time.Second)
				return fmt.Errorf("en error")
			},
			OnSuccess: func(w Worker, j *Job) {
				fmt.Println("Test job completed")
			},
			OnError: func(w Worker, j *Job) {
				fmt.Println("Test job failed")
			},
		})
	}

	time.Sleep(5 * time.Second) // Let workers scale up

	// check that workers increased
	if len(manager.workers) <= minWorkers {
		t.Errorf("Expected workers to scale up, but only %d running", len(manager.workers))
	}

	// scale down scenario
	time.Sleep(3 * time.Second)

	// should have reduced workers by now
	if len(manager.workers) != minWorkers {
		t.Errorf("Expected workers to scale back down to %d, but got %d", minWorkers, len(manager.workers))
	}

	manager.Stop()
}
