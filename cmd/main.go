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
	m := conveyor.CreateManager() // Use the default manager or build a custom
	// m := NewManager().SetMinWorkers(1).SetMaxWorkers(100).
	// 	SetStepUpAt(10).SetStepDownAt(10).SetTimePerTicker(time.Second / 4)

	m.Start()

	jobParam := &JobParam{
		A: "Hello World",
	}

	// Adding workers scenario
	for range 100 {
		conveyor.CONVEYOR_BELT <- conveyor.Job{
			Context: context.Background(),
			Param:   jobParam,
			Consume: func(param any) error {
				time.Sleep(time.Second)
				jParam := param.(*JobParam)
				fmt.Println(jParam.A)
				return nil
			},
		}
	}

	select {}
}
