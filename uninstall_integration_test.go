package registry_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nhatthm/aferocopy"
	registry "github.com/nhatthm/plugin-registry"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationFsRegistry_Uninstall(t *testing.T) {
	t.Parallel()

	cfg := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: "https://example.org"
        version: v1.2.0
        enabled: true
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
`, runtime.GOOS, runtime.GOARCH)

	registryDir := t.TempDir()
	configFile := filepath.Join(registryDir, "config.yaml")
	pluginDir := filepath.Join(registryDir, "my-plugin")

	fs := afero.NewOsFs()
	err := afero.WriteFile(fs, configFile, []byte(cfg), 0755)
	require.NoError(t, err)

	err = fs.Mkdir(pluginDir, 0755)
	require.NoError(t, err)

	err = aferocopy.Copy(pluginDir, "resources/fixtures")
	require.NoError(t, err)

	// Verify the condition before taking action.
	fi, err := fs.Stat(pluginDir)
	require.NoError(t, err)
	assert.True(t, fi.IsDir())

	// Uninstall plugin.
	r, err := registry.NewRegistry(registryDir)
	require.NoError(t, err)

	err = r.Uninstall("my-plugin")
	require.NoError(t, err)

	// Verify result.
	exists, err := afero.Exists(fs, pluginDir)
	require.NoError(t, err)
	assert.False(t, exists)

	expected := "plugins: {}\n"
	actual, err := afero.ReadFile(fs, configFile)
	require.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}
