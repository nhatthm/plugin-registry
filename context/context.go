package context

import (
	"context"

	"github.com/spf13/afero"
)

type fsKey struct{}

// WithFs returns the context with file system.
func WithFs(ctx context.Context, fs afero.Fs) context.Context {
	return context.WithValue(ctx, fsKey{}, fs)
}

// Fs returns file system from context.
func Fs(ctx context.Context) afero.Fs {
	fs, ok := ctx.Value(fsKey{}).(afero.Fs)
	if !ok {
		return afero.NewOsFs()
	}

	return fs
}
