package registry_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	registry "github.com/nhatthm/plugin-registry"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationFsRegistry_Enable(t *testing.T) {
	t.Parallel()

	cfg := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: "https://example.org"
        version: v1.2.0
        enabled: false
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
        tags:
            - tag1
`, runtime.GOOS, runtime.GOARCH)

	registryDir := t.TempDir()
	configFile := filepath.Join(registryDir, "config.yaml")

	fs := afero.NewOsFs()
	err := afero.WriteFile(fs, configFile, []byte(cfg), 0755)
	require.NoError(t, err)

	// Enable plugin.
	r, err := registry.NewRegistry(registryDir)
	require.NoError(t, err)

	err = r.Enable("my-plugin")
	require.NoError(t, err)

	// Verify result.
	expected := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: https://example.org
        version: v1.2.0
        description: ""
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
