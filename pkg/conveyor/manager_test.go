package conveyor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"testing"

	"github.com/stretchr/testify/require"
)

type OddEven struct {
	m    sync.Mutex
	odd  uint8
	even uint8
}

func (oe *OddEven) Odd() {
	oe.m.Lock()
	oe.odd++
	oe.m.Unlock()
}

func (oe *OddEven) Even() {
	oe.m.Lock()
	oe.even++
	oe.m.Unlock()
}

func TestManager(t *testing.T) {
	manager := CreateManager().SetTimePerTicker(time.Second / 10).SetMinWorkers(0).
		SetDebugging(true).Start()

	t.Cleanup(manager.Stop)

	oddEven := &OddEven{}
	// scale up scenario
	for i := range 100 {
		manager.B.Push(&Job{
			Context: context.Background(),
			Param:   i%2 == 0,
			Consume: func(p any) error {
				time.Sleep(time.Second)
				if param := p.(bool); param {
					return nil // represents success (even)
				}

				return fmt.Errorf("failure") // represents failure (odd)
			},
			OnSuccess: func(w Worker, j *Job) { oddEven.Even() },
			OnError:   func(w Worker, j *Job) { oddEven.Odd() },
		})
	}

	time.Sleep(5 * time.Second) // Let workers scale up

	// check that workers increased
	require.Eventually(t, func() bool {
		return len(manager.workers) > manager.minWorkersAllowed
	}, 5*time.Second, 100*time.Millisecond, "workers should scale up")

	require.Eventually(t, func() bool {
		return oddEven.even == 50 && oddEven.odd == 50
	}, 5*time.Second, 100*time.Millisecond, "Expected 50 odds and 50 evens")

	// scale down scenario
	time.Sleep(3 * time.Second)

	// should have reduced workers by now
	require.Eventually(t, func() bool {
		return len(manager.workers) == manager.minWorkersAllowed
	}, 5*time.Second, 100*time.Millisecond, "workers should scale down")
}
