package bootstrap

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/freman/work"
	"golang.org/x/sys/unix"
)

type Bootstrap struct {
	work.Flow
}

// MkdirAll adds an os.MkDirAll step to your strap
func (b *Bootstrap) MkdirAll(path string, perm os.FileMode) *Bootstrap {
	b.Add(func(next work.Task) work.Task {
		return work.LabelFunc(fmt.Sprintf("create path %s", path), func(ctx context.Context) error {
			if err := os.MkdirAll(path, perm); err != nil {
				return err
			}
			return next.Execute(ctx)
		})
	})

	return b
}

// IsWritable adds an check to see if a path is writable
func (b *Bootstrap) IsWritable(path string) *Bootstrap {
	b.Add(func(next work.Task) work.Task {
		return work.LabelFunc(fmt.Sprintf("path is writable %s", path), func(ctx context.Context) error {
			if _, err := os.Stat(path); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
				path = filepath.Dir(path)
				if _, err := os.Stat(path); err != nil {
					return err
				}
			}

			if unix.Access(path, unix.W_OK) != nil {
				return os.ErrPermission
			}

			return next.Execute(ctx)
		})
	})

	return b
}
