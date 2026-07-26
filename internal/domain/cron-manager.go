package domain

import (
	"context"
	"log"

	"github.com/go-co-op/gocron/v2"
)

// JobManager manages all scheduled jobs.
type JobManager struct {
	scheduler gocron.Scheduler
	jobs      []Job
	context   context.Context
}

// NewJobManager creates a new JobManager instance.
func NewJobManager(ctx context.Context) (*JobManager, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &JobManager{
		scheduler: s,
		jobs:      []Job{},
		context:   ctx,
	}, nil
}

// Scheduler returns the underlying gocron.Scheduler instance.
func (jm *JobManager) Scheduler() gocron.Scheduler {
	return jm.scheduler
}

// RegisterJob adds a job to the job manager.
func (jm *JobManager) RegisterJob(job Job) {
	jm.jobs = append(jm.jobs, job)
}

// StartScheduler starts the cron scheduler to execute jobs at their scheduled times.
func (jm *JobManager) StartScheduler() {
	for _, job := range jm.jobs {
		duration := job.Duration()

		jm.scheduler.NewJob(
			gocron.DurationJob(duration),
			gocron.NewTask(
				func() {
					if err := job.Run(jm.context); err != nil {
						log.Printf("Error in job %s: %v", job.Name(), err)
					} else {
						log.Printf("Job %s executed successfully", job.Name())
					}
				}),
			gocron.WithName(job.Name()),
		)
	}
	jm.scheduler.Start()
}
