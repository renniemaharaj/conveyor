package conveyor

import (
	"context"
)

// Job represents a job on the conveyor belt
type Job struct {
	Context context.Context `json:"context"`
	Param   any             `json:"params"` // Paramter for consumption

	// Consume function for the worker
	Consume func(params any) error `json:"consume"`
	// On success callback function
	OnSuccess func(w Worker, j *Job) `json:"onSuccess"`
	// On error callback function
	OnError func(w Worker, j *Job) `json:"onError"`
}
