package work

import "context"

type Labelled struct {
	label string
	task  Task
}

func (l *Labelled) Execute(ctx context.Context) error {
	return l.task.Execute(ctx)
}

func LabelFunc(label string, task TaskFunc) Task {
	return &Labelled{label, task}
}
