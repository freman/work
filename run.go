package work

import (
	"context"
	"errors"
)

var ErrNoMoreTasks = errors.New("no more tasks")

// Run encompasses a single resumeable execution and it's current state
type Run struct {
	finished bool
	ctx      context.Context
	execute  func(ctx context.Context) error
	task     Task
	err      error
}

func (r *Run) last(_ context.Context) error {
	r.finished = true
	r.ctx = nil
	r.task = nil
	r.execute = nil
	return nil
}

func (r *Run) queue(next Job, last Task) Task {
	return TaskFunc(func(ctx context.Context) error {
		r.task = next(last)
		r.ctx = ctx
		r.err = r.task.Execute(r.ctx)
		return r.err
	})
}

// Resume attempt to resume where the last job failed
func (r *Run) Resume() error {
	if r.task == nil {
		return ErrNoMoreTasks
	}
	r.err = r.task.Execute(r.ctx)
	return r.err
}

// Task is the last failed task
func (r *Run) Task() Task {
	return r.task
}

// Err of the failed job
func (r *Run) Err() error {
	return r.err
}

// Finished is true if there are no more jobs
func (r *Run) Finished() bool {
	return r.finished
}
