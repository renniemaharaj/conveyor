package main

import (
	"context"
	"fmt"
	"time"

	"github.com/renniemaharaj/conveyor/pkg/conveyor"
)

type JobParam struct {
	A string
}

func main() {
	manager := conveyor.CreateManager() // Use the default manager or build a custom

	// manager := BlankManager().SetMinWorkers(1).SetMaxWorkers(100).
	// 	SetSafeQueueLength(10).SetTimePerTicker(time.Second / 4)

	// manager.B = NewConveyorBelt()
	// manager.quit = make(chan struct{})

	manager.Start()

	// Unopinionated job param
	jobParam := &JobParam{
		A: "Hello World",
	}

	// Adding workers scenario
	for range 100 {
		manager.B.Push(&conveyor.Job{
			Context: context.Background(),
			Param:   jobParam,
			Consume: func(param any) error {
				time.Sleep(time.Second)
				jParam := param.(*JobParam)
				fmt.Println(jParam.A)
				return nil
			},
		})
	}

	select {}
}
