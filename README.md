# Conveyor System - Manager, Job, and Worker

This project implements a **dynamic worker pool manager** called `Manager` that processes queued jobs from a shared channel (`CONVEYOR_BELT`). It is designed for scenarios where workload varies over time, and the system should automatically scale the number of active workers up or down based on queue size.

## Features

- **Automatic Scaling**: Dynamically increases or decreases workers based on job queue length.
- **Configurable Thresholds**: Define when to scale up or down.
- **Safe Shutdown**: Gracefully stops workers after finishing their current job.
- **Ticker-based Monitoring**: Regularly checks the queue and adjusts worker count.

## Components

### Manager

Responsible for:

- Starting the initial worker pool.
- Monitoring queue size.
- Scaling workers up or down.
- Stopping all workers safely.

Key parameters:

- `minWorkers`: Minimum workers to keep running.
- `maxWorkers`: Maximum workers allowed.
- `stepUpAt`: Queue length threshold to add workers.
- `stepDownAt`: Queue length threshold to remove workers.
- `ticker`: Interval at which the manager checks the queue.

### Worker

A worker executes jobs pulled from `CONVEYOR_BELT`. Each worker runs until:

- Stopped by the manager.
- Program terminates.

### Job

Represents a unit of work:

```go
Job {
    Context context.Context // For cancellation, timeouts, etc.
    Param   any             // Job-specific parameters
    Consume func(param any) error // Job execution logic
}
```

### Conveyor Belt

A global job queue:

```go
var CONVEYOR_BELT chan Job
```

Workers continuously pull jobs from here.

## Usage

### Creating and Starting a Manager

```go
m := conveyor.CreateManager() // Creates a default manager
// or configure manually:
// m := NewManager().
//     SetMinWorkers(1).
//     SetMaxWorkers(100).
//     SetStepUpAt(10).
//     SetStepDownAt(5).
//     SetTimePerTicker(time.Second / 4)

m.Start()
```

### Adding Jobs

```go
type JobParam struct {
    A string
}

jobParam := &JobParam{A: "Hello World"}

for i := 0; i < 100; i++ {
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
```

### Stopping the Manager

```go
m.Stop() // Stops the ticker and gracefully shuts down workers
```

## Example Execution

1. Start the manager.
2. Push jobs to `CONVEYOR_BELT`.
3. Manager automatically adjusts the worker count based on job backlog.
4. Stop the manager when done.

## When to Use

- Processing tasks where load fluctuates.
- Systems requiring efficient resource use.
- Scenarios where idle workers should be reduced automatically.

---

This design allows a balance between responsiveness and resource efficiency by only running the number of workers needed to handle the current workload.
