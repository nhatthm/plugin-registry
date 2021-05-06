package registry_test

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nhatthm/aferocopy"
	registry "github.com/nhatthm/plugin-registry"
	"github.com/nhatthm/plugin-registry/installer"
	"github.com/nhatthm/plugin-registry/plugin"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationFsRegistry_Install_NoConfigFile(t *testing.T) {
	t.Parallel()

	registryDir := t.TempDir()
	configFile := filepath.Join(registryDir, "config.yaml")
	pluginDir := filepath.Join(registryDir, "my-plugin")

	fs := afero.NewOsFs()

	err := fs.Mkdir(pluginDir, 0755)
	require.NoError(t, err)

	// Register installer.
	installer.Register(t.Name(), func(_ context.Context, source string) bool {
		return source == t.Name()
	}, func(afero.Fs) installer.Installer {
		return installer.CallbackInstaller(func(_ context.Context, dest, src string) (*plugin.Plugin, error) {
			p := &plugin.Plugin{
				Name:        "my-plugin",
				URL:         "https://example.org",
				Version:     "v1.2.0",
				Description: "my plugin",
				Enabled:     true,
				Hidden:      true,
				Artifacts: plugin.Artifacts{
					plugin.RuntimeArtifactIdentifier(): {
						File: "my-plugin",
					},
				},
				Tags: plugin.Tags{"tag1"},
			}

			return p, aferocopy.Copy(pluginDir, "resources/fixtures")
		})
	})

	// Verify the condition before taking action.
	fi, err := fs.Stat(pluginDir)
	require.NoError(t, err)
	assert.True(t, fi.IsDir())

	// Install plugin.
	r, err := registry.NewRegistry(registryDir)
	require.NoError(t, err)

	err = r.Install(context.Background(), t.Name())
	require.NoError(t, err)

	// Verify result
	isDir, err := afero.IsDir(fs, pluginDir)
	require.NoError(t, err)
	assert.True(t, isDir)

	expected := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: https://example.org
        version: v1.2.0
        description: my plugin
        enabled: true
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
        tags:
            - tag1
`, runtime.GOOS, runtime.GOARCH)
	actual, err := afero.ReadFile(fs, configFile)
	require.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}

func TestIntegrationFsRegistry_Install_HasConfigFile(t *testing.T) {
	t.Parallel()

	registryDir := t.TempDir()
	configFile := filepath.Join(registryDir, "config.yaml")
	pluginDir := filepath.Join(registryDir, "my-plugin")

	fs := afero.NewOsFs()
	err := afero.WriteFile(fs, configFile, []byte(`plugins: {}`), 0755)
	require.NoError(t, err)

	err = fs.Mkdir(pluginDir, 0755)
	require.NoError(t, err)

	// Register installer.
	installer.Register(t.Name(), func(_ context.Context, source string) bool {
		return source == t.Name()
	}, func(afero.Fs) installer.Installer {
		return installer.CallbackInstaller(func(_ context.Context, dest, src string) (*plugin.Plugin, error) {
			p := &plugin.Plugin{
				Name:        "my-plugin",
				URL:         "https://example.org",
				Version:     "v1.2.0",
				Description: "my plugin",
				Enabled:     true,
				Hidden:      true,
				Artifacts: plugin.Artifacts{
					plugin.RuntimeArtifactIdentifier(): {
						File: "my-plugin",
					},
				},
				Tags: plugin.Tags{"tag1"},
			}

			return p, aferocopy.Copy(pluginDir, "resources/fixtures")
		})
	})

	// Verify the condition before taking action.
	fi, err := fs.Stat(pluginDir)
	require.NoError(t, err)
	assert.True(t, fi.IsDir())

	// Install plugin.
	r, err := registry.NewRegistry(registryDir)
	require.NoError(t, err)

	err = r.Install(context.Background(), t.Name())
	require.NoError(t, err)

	// Verify result.
	isDir, err := afero.IsDir(fs, pluginDir)
	require.NoError(t, err)
	assert.True(t, isDir)

	expected := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: https://example.org
        version: v1.2.0
        description: my plugin
        enabled: true
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
        tags:
            - tag1
`, runtime.GOOS, runtime.GOARCH)
	actual, err := afero.ReadFile(fs, configFile)
	require.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}
