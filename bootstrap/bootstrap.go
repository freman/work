package bootstrap

import (
	"context"
	"fmt"
	"os"

	"github.com/freman/work"
	"golang.org/x/sys/unix"
)

type Bootstrap struct {
	*work.Flow
}

// MkdirAll adds an os.MkDirAll step to your strap
func (b *Bootstrap) MkdirAll(path string, perm os.FileMode) {
	b.Add(func(next work.Task) work.Task {
		return work.LabelFunc(fmt.Sprintf("create path %s", path), func(ctx context.Context) error {
			err := os.MkdirAll(path, perm)
			if err != nil {
				return err
			}
			return next.Execute(ctx)
		})
	})
}

// IsWritable adds an check to see if a path is writable
func (b *Bootstrap) IsWritable(path string) {
	b.Add(func(next work.Task) work.Task {
		return work.LabelFunc(fmt.Sprintf("path is writable %s", path), func(ctx context.Context) error {
			_, err := os.Stat(path)
			if err != nil {
				return err
			}

			if unix.Access(path, unix.W_OK) != nil {
				return os.ErrPermission
			}

			return next.Execute(ctx)
		})
	})
}
