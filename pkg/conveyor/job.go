package conveyor

import (
	"context"
)

// Job represents a job on the conveyor belt
type Job struct {
	Context context.Context        `json:"context"`
	Consume func(params any) error `json:"consume"`
	Param   any                    `json:"params"`
}
