package context

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFs(t *testing.T) {
	t.Parallel()

	t.Run("not in context", func(t *testing.T) {
		t.Parallel()

		fs := Fs(context.Background())
		expected := afero.NewOsFs()

		assert.Equal(t, expected, fs)
	})

	t.Run("in context", func(t *testing.T) {
		t.Parallel()

		ctx := WithFs(context.Background(), afero.NewMemMapFs())
		fs := Fs(ctx)
		expected := afero.NewMemMapFs()

		assert.Equal(t, expected, fs)
	})
}
