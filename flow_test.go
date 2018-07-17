package work_test

import (
	"context"
	"errors"
	"testing"

	"github.com/freman/work"
)

func TestWorkflow(t *testing.T) {
	failOn := 3
	seq := 0

	s := &work.Flow{}

	for i := 0; i < 10; i++ {
		func(i int) {
			s.Add(func(next work.Task) work.Task {
				return work.TaskFunc(func(ctx context.Context) error {
					if seq != i {
						t.Errorf("Expected sequence %d but got sequence %d", i, seq)
					}

					if failOn == i {
						return errors.New("failing as ordered")
					}

					seq++

					return next.Execute(ctx)
				})
			})
		}(i)
	}

	run, err := s.Execute(context.Background())
	if err == nil {
		t.Error("Expected an error, didn't get one")
	}

	failOn = 0

	err = run.Resume()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if seq != 10 {
		t.Errorf("Expected 10 got %d", seq)
	}
}
