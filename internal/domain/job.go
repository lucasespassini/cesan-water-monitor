package domain

import (
	"context"
	"time"
)

// Job interface defines the structure for any job that will be scheduled.
type Job interface {
	Name() string
	Duration() time.Duration
	Run(ctx context.Context) error
}
