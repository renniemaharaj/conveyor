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

- `minWorkersAllowed`: Minimum workers to keep running.
- `maxWorkersAllowed`: Maximum workers allowed.
- `safeQueueLength`: Queue length threshold to add workers.
- `ticker`: Interval at which the manager checks the queue.

### Worker

A worker executes jobs pulled from its assigned `CONVEYOR_BELT`. Each worker runs in its own go routine until:

- Stopped by the manager.
- Program terminates.

### Job

Represents a unit of work:

```go
Job {
    Context context.Context `json:"context"`
	Param   any             `json:"params"` // Paramter for consumption

	// Consume function for the worker
	Consume func(params any) error `json:"consume"`
	// On success callback function
	OnSuccess func(w Worker, j *Job) `json:"onSuccess"`
	// On error callback function
	OnError func(w Worker, j *Job) `json:"onError"`
}
```

### Conveyor Belt

Manager specific job queues

```go
// The conveyor belt struct containing a single channel of jobs
type ConveyorBelt struct {
	C chan Job
}

// Creates and returns a new ConveyorBelt with initialized channel
func NewConveyorBelt() *ConveyorBelt {
	return &ConveyorBelt{C: make(chan Job, 100)}
}

// Pushes a job to the conveyor belt
func (b *ConveyorBelt) Push(j *Job) {
	b.C <- *j
}

// Takes a job from the conveyor belt
func (b *ConveyorBelt) Take() *Job {
	j := <-b.C
	return &j
}
```

Workers continuously pull jobs from their assigned channels.

## Usage

### Creating and Starting a Manager

```go
m := conveyor.CreateManager() // Creates a default manager
// or configure manually:

m := conveyor.CreateManager().SetMinWorkers(1).SetMaxWorkers(100).
		SetSafeQueueLength(10).SetTimePerTicker(time.Second / 4)

	m.B = NewConveyorBelt()
	m.quit = make(chan struct{})

m.Start()
```

### Adding Jobs

```go
// Define a jobParam struct
// This can be anything, (any)
type JobParam struct {
	A string
}

jobParam := &JobParam{A: "Hello World"}

// Without success, error functions
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

// With success, error functions (optional)
    for range 100 {
		manager.B.Push(&Job{
			Context: context.Background(),
			Consume: func(j any) error {
				time.Sleep(time.Second)
				return fmt.Errorf("en error") // Switch this
                // return nil
			},
			OnSuccess: func(w Worker, j *Job) {
				fmt.Println("Test job completed")
			},
			OnError: func(w Worker, j *Job) {
				fmt.Println("Test job failed")
			},
		})
	}
```

### Stopping the Manager

```go
m.Stop() // Stops the ticker and gracefully shuts down workers
```

## Example Execution

1. Start the manager.
2. Push jobs to manager's `CONVEYOR_BELT`.
```go
manager.B.Push(&conveyor.Job{})
```
3. Manager automatically adjusts the worker count (up & down) based on job backlog.
4. Stop the manager when done.

## When to Use

- Processing tasks where load fluctuates.
- Systems requiring efficient resource use.
- Scenarios where idle workers should be reduced automatically.

---

This design allows a balance between responsiveness and resource efficiency by only running the number of workers needed to handle the current workload.
