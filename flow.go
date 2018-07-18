package work

import (
	"context"
)

type Task interface {
	Execute(ctx context.Context) error
}

type TaskFunc func(ctx context.Context) error

// Execute calls f(ctx).
func (f TaskFunc) Execute(ctx context.Context) error {
	return f(ctx)
}

type Job func(Task) Task

// Flow encapsulates the list of jobs to be executed
type Flow struct {
	jobs []Job
}

// Add a job
func (f *Flow) Add(job Job) *Flow {
	f.jobs = append(f.jobs, job)
	return f
}

// Execute the jobs
func (f *Flow) Execute(ctx context.Context) (*Run, error) {
	run := &Run{}
	var last Task = TaskFunc(run.last)
	for i := len(f.jobs) - 1; i >= 0; i-- {
		last = func(i int, last Task) Task {
			return run.queue(f.jobs[i], last)
		}(i, last)
	}

	return run, last.Execute(ctx)
}
