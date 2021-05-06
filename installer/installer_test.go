package installer_test

import (
	"context"
	"testing"

	fsCtx "github.com/nhatthm/plugin-registry/context"
	"github.com/nhatthm/plugin-registry/installer"
	installerMock "github.com/nhatthm/plugin-registry/mock/installer"
	"github.com/nhatthm/plugin-registry/plugin"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallbackInstaller_Install(t *testing.T) {
	t.Parallel()

	i := installer.CallbackInstaller(func(context.Context, string, string) (*plugin.Plugin, error) {
		return &plugin.Plugin{Name: "my-plugin"}, nil
	})

	actual, err := i.Install(context.Background(), "", "")

	expected := &plugin.Plugin{Name: "my-plugin"}

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestNew(t *testing.T) {
	t.Parallel()

	expected := installerMock.NoMock(t)

	installer.Register("TestNew", func(_ context.Context, src string) bool {
		return src == "TestNew"
	}, func(afero.Fs) installer.Installer {
		return expected
	})

	fs := afero.NewMemMapFs()
	ctx := fsCtx.WithFs(context.Background(), fs)
	actual, err := installer.New(ctx, "TestNew")
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestNew_NotFound(t *testing.T) {
	t.Parallel()

	ctx := fsCtx.WithFs(context.Background(), afero.NewMemMapFs())
	actual, err := installer.New(ctx, "unknown")

	expected := `unknown installer`

	assert.Nil(t, actual)
	assert.EqualError(t, err, expected)
}

func TestFind(t *testing.T) {
	t.Parallel()

	expected := installerMock.NoMock(t)

	installer.Register("TestFind", func(_ context.Context, src string) bool {
		return src == "TestFind"
	}, func(afero.Fs) installer.Installer {
		return expected
	})

	fs := afero.NewMemMapFs()
	ctx := fsCtx.WithFs(context.Background(), fs)
	actual, err := installer.Find(ctx, "TestFind")
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestFind_Error(t *testing.T) {
	t.Parallel()

	installer.Register("TestFind_Error", func(context.Context, string) bool {
		return false
	}, func(afero.Fs) installer.Installer {
		return installerMock.NoMock(t)
	})

	ctx := fsCtx.WithFs(context.Background(), afero.NewMemMapFs())
	actual, err := installer.Find(ctx, "unknown")

	expected := `no supported installer`

	assert.Nil(t, actual)
	assert.EqualError(t, err, expected)
}
