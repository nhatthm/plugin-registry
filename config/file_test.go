package config_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/nhatthm/aferomock"
	"github.com/nhatthm/plugin-registry/config"
	"github.com/nhatthm/plugin-registry/plugin"
	"github.com/spf13/afero"
	"github.com/spf13/afero/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeConfigFile(data string) afero.File {
	f := mem.NewFileHandle(mem.CreateFile("config.yaml"))

	_, _ = f.Write([]byte(data))   // nolint: errcheck
	_, _ = f.Seek(0, io.SeekStart) // nolint: errcheck

	return f
}

func writeConfigFile(fs afero.Fs, data string) error {
	return afero.WriteFile(fs, "config.yaml", []byte(data), os.FileMode(0o644))
}

func assertConfigFile(t *testing.T, fs afero.Fs, name string, expected string) bool { // nolint: unparam
	t.Helper()

	actual, err := afero.ReadFile(fs, name)
	require.NoError(t, err)

	return assert.Equal(t, expected, string(actual))
}

func TestFileConfigurator_Config(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockFs         aferomock.FsMocker
		expectedConfig config.Configuration
		expectedError  string
	}{
		{
			scenario: "could not stat file",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(nil, errors.New("stat error"))
			}),
			expectedError: "stat error",
		},
		{
			scenario: "file does not exist",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(nil, os.ErrNotExist)
			}),
		},
		{
			scenario: "could not open file",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(aferomock.NewFileInfo(), nil)

				fs.On("Open", "config.yaml").
					Return(nil, errors.New("open error"))
			}),
			expectedError: "open error",
		},
		{
			scenario: "could not decode",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(aferomock.NewFileInfo(), nil)

				fs.On("Open", "config.yaml").
					Return(mem.NewFileHandle(mem.CreateFile("config.yaml")), nil)
			}),
			expectedError: "could not load configuration: EOF",
		},
		{
			scenario: "success",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				cfg := fmt.Sprintf(`
plugins:
    my-plugin:
        name: my-plugin
        url: "https://example.org"
        version: v1.2.0
        description: my plugin
        enabled: true
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
        tags: []
`, runtime.GOOS, runtime.GOARCH)

				fs.On("Stat", "config.yaml").
					Return(aferomock.NewFileInfo(), nil)

				fs.On("Open", "config.yaml").
					Return(makeConfigFile(cfg), nil)
			}),
			expectedConfig: config.Configuration{
				Plugins: map[string]plugin.Plugin{
					"my-plugin": {
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
						Tags: plugin.Tags{},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			c := config.NewFileConfigurator("config.yaml").WithFs(tc.mockFs(t))
			cfg, err := c.Config()

			assert.Equal(t, tc.expectedConfig, cfg)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestFileConfigurator_SetPlugin_Error(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockFs        aferomock.FsMocker
		expectedError string
	}{
		{
			scenario: "could not load config",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(nil, errors.New("load error"))
			}),
			expectedError: "load error",
		},
		{
			scenario: "could not open file for write",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(nil, os.ErrNotExist)

				fs.On("OpenFile", "config.yaml", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0o644)).
					Return(nil, errors.New("open error"))
			}),
			expectedError: "open error",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			c := config.NewFileConfigurator("config.yaml").WithFs(tc.mockFs(t))
			err := c.SetPlugin(plugin.Plugin{Name: "my-plugin"})

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestFileConfigurator_SetPlugin_Success(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	err := writeConfigFile(fs, "plugins:")
	require.NoError(t, err)

	c := config.NewFileConfigurator("config.yaml").WithFs(fs)
	err = c.SetPlugin(plugin.Plugin{
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
		Tags: plugin.Tags{},
	})
	require.NoError(t, err)

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
        tags: []
`, runtime.GOOS, runtime.GOARCH)

	assertConfigFile(t, fs, "config.yaml", expected)
}

func TestFileConfigurator_RemovePlugin_Error(t *testing.T) {
	t.Parallel()

	cfg := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: https://example.org
        version: v1.2.0
        enabled: true
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
        tags: []
`, runtime.GOOS, runtime.GOARCH)

	testCases := []struct {
		scenario      string
		mockFs        aferomock.FsMocker
		expectedError string
	}{
		{
			scenario: "could not load config",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(nil, errors.New("load error"))
			}),
			expectedError: "load error",
		},
		{
			scenario: "plugin not found",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(nil, os.ErrNotExist)
			}),
			expectedError: "plugin does not exist",
		},
		{
			scenario: "could not open file for write",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(aferomock.NewFileInfo(), nil)

				fs.On("Open", "config.yaml").
					Return(makeConfigFile(cfg), nil)

				fs.On("OpenFile", "config.yaml", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0o644)).
					Return(nil, errors.New("open error"))
			}),
			expectedError: "open error",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			c := config.NewFileConfigurator("config.yaml").WithFs(tc.mockFs(t))
			err := c.RemovePlugin("my-plugin")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestFileConfigurator_RemovePlugin_Success(t *testing.T) {
	t.Parallel()

	cfg := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: https://example.org
        version: v1.2.0
        enabled: true
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
        tags: []
`, runtime.GOOS, runtime.GOARCH)

	fs := afero.NewMemMapFs()
	err := writeConfigFile(fs, cfg)
	require.NoError(t, err)

	c := config.NewFileConfigurator("config.yaml").WithFs(fs)
	err = c.RemovePlugin("my-plugin")
	require.NoError(t, err)

	expected := "plugins: {}\n"

	assertConfigFile(t, fs, "config.yaml", expected)
}

func TestFileConfigurator_EnablePlugin_Error(t *testing.T) {
	t.Parallel()

	cfg := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: https://example.org
        version: v1.2.0
        enabled: false
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
        tags: []
`, runtime.GOOS, runtime.GOARCH)

	testCases := []struct {
		scenario      string
		mockFs        aferomock.FsMocker
		pluginName    string
		expectedError string
	}{
		{
			scenario: "could not load config",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(nil, errors.New("load error"))
			}),
			expectedError: "load error",
		},
		{
			scenario: "plugin not found",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(aferomock.NewFileInfo(), nil)

				fs.On("Open", "config.yaml").
					Return(makeConfigFile(cfg), nil)
			}),
			pluginName:    "unknown",
			expectedError: "plugin does not exist",
		},
		{
			scenario: "could not open file for write",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(aferomock.NewFileInfo(), nil)

				fs.On("Open", "config.yaml").
					Return(makeConfigFile(cfg), nil)

				fs.On("OpenFile", "config.yaml", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0o644)).
					Return(nil, errors.New("open error"))
			}),
			expectedError: "open error",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			if tc.pluginName == "" {
				tc.pluginName = "my-plugin"
			}

			c := config.NewFileConfigurator("config.yaml").WithFs(tc.mockFs(t))
			err := c.EnablePlugin(tc.pluginName)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestFileConfigurator_EnablePlugin_Success(t *testing.T) {
	t.Parallel()

	cfg := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: https://example.org
        version: v1.2.0
        enabled: false
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
        tags: []
`, runtime.GOOS, runtime.GOARCH)

	fs := afero.NewMemMapFs()
	err := writeConfigFile(fs, cfg)
	require.NoError(t, err)

	c := config.NewFileConfigurator("config.yaml").WithFs(fs)
	err = c.EnablePlugin("my-plugin")
	require.NoError(t, err)

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
        tags: []
`, runtime.GOOS, runtime.GOARCH)

	assertConfigFile(t, fs, "config.yaml", expected)
}

func TestFileConfigurator_DisablePlugin_Error(t *testing.T) {
	t.Parallel()

	cfg := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: https://example.org
        version: v1.2.0
        enabled: true
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
        tags: []
`, runtime.GOOS, runtime.GOARCH)

	testCases := []struct {
		scenario      string
		mockFs        aferomock.FsMocker
		pluginName    string
		expectedError string
	}{
		{
			scenario: "could not load config",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(nil, errors.New("load error"))
			}),
			expectedError: "load error",
		},
		{
			scenario: "plugin not found",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(aferomock.NewFileInfo(), nil)

				fs.On("Open", "config.yaml").
					Return(makeConfigFile(cfg), nil)
			}),
			pluginName:    "unknown",
			expectedError: "plugin does not exist",
		},
		{
			scenario: "could not open file for write",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "config.yaml").
					Return(aferomock.NewFileInfo(), nil)

				fs.On("Open", "config.yaml").
					Return(makeConfigFile(cfg), nil)

				fs.On("OpenFile", "config.yaml", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0o644)).
					Return(nil, errors.New("open error"))
			}),
			expectedError: "open error",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			if tc.pluginName == "" {
				tc.pluginName = "my-plugin"
			}

			c := config.NewFileConfigurator("config.yaml").WithFs(tc.mockFs(t))
			err := c.DisablePlugin(tc.pluginName)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestFileConfigurator_DisablePlugin_Success(t *testing.T) {
	t.Parallel()

	cfg := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: https://example.org
        version: v1.2.0
        enabled: true
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
        tags: []
`, runtime.GOOS, runtime.GOARCH)

	fs := afero.NewMemMapFs()
	err := writeConfigFile(fs, cfg)
	require.NoError(t, err)

	c := config.NewFileConfigurator("config.yaml").WithFs(fs)
	err = c.DisablePlugin("my-plugin")
	require.NoError(t, err)

	expected := fmt.Sprintf(`plugins:
    my-plugin:
        name: my-plugin
        url: https://example.org
        version: v1.2.0
        description: ""
        enabled: false
        hidden: true
        artifacts:
            %s/%s:
                file: my-plugin
        tags: []
`, runtime.GOOS, runtime.GOARCH)

	assertConfigFile(t, fs, "config.yaml", expected)
}
