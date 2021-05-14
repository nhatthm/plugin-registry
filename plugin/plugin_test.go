package plugin

import (
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestPlugins_FilterByTag(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		plugins := Plugins{
			"plugin1": {Tags: Tags{"tag1", "tag2"}},
			"plugin2": {Tags: Tags{"tag2", "tag3"}},
			"plugin3": {Tags: Tags{"tag3", "tag4"}},
		}

		expected := Plugins{
			"plugin1": {Tags: Tags{"tag1", "tag2"}},
			"plugin2": {Tags: Tags{"tag2", "tag3"}},
		}

		assert.Equal(t, expected, plugins.FilterByTag("tag2"))
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		plugins := Plugins{
			"plugin1": {Tags: Tags{"tag1", "tag2"}},
			"plugin2": {Tags: Tags{"tag2", "tag3"}},
			"plugin3": {Tags: Tags{"tag3", "tag4"}},
		}

		expected := Plugins{}

		assert.Equal(t, expected, plugins.FilterByTag("tag5"))
	})
}

func TestPlugin_RuntimeArtifact(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		plugin   Plugin
		expected Artifact
	}{
		{
			scenario: "has os and arch",
			plugin: Plugin{Artifacts: map[ArtifactIdentifier]Artifact{
				RuntimeArtifactIdentifier(): {
					File: defaultFile,
				},
			}},
			expected: Artifact{File: "${name}-${version}-${os}-${arch}.tar.gz"},
		},
		{
			scenario: "has only os",
			plugin: Plugin{Artifacts: map[ArtifactIdentifier]Artifact{
				RuntimeArtifactIdentifierWithoutArch(): {
					File: defaultFile,
				},
			}},
			expected: Artifact{File: "${name}-${version}-${os}-${arch}.tar.gz"},
		},
		{
			scenario: "has nothing",
			plugin:   Plugin{Artifacts: map[ArtifactIdentifier]Artifact{}},
			expected: Artifact{File: "${name}-${version}-${os}-${arch}.tar.gz"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.plugin.RuntimeArtifact())
		})
	}
}

func TestPlugin_ResolveArtifact(t *testing.T) {
	t.Parallel()

	a := Artifact{File: "${name}-${version}-${os}-${arch}.tar.gz"}
	p := Plugin{
		Name:    "my-plugin",
		Version: "1.0.3",
	}

	expected := Artifact{
		File: fmt.Sprintf("my-plugin-1.0.3-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH),
	}

	assert.Equal(t, expected, p.ResolveArtifact(a))
}

func TestPlugin_MarshalYAML(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		plugin   Plugin
		expected string
	}{
		{
			scenario: "empty",
			expected: `name: ""
url: ""
version: ""
description: ""
enabled: false
hidden: false
artifacts: {}
tags: []
`,
		},
		{
			scenario: "artifact without arch",
			plugin: Plugin{
				Name:        "my-plugin",
				URL:         "https://example.org",
				Version:     "v1.2.0",
				Description: "my plugin",
				Enabled:     true,
				Hidden:      true,
				Artifacts: Artifacts{
					NewArtifactIdentifier("darwin", ""): {
						File: "my-plugin",
					},
				},
			},
			expected: `name: my-plugin
url: https://example.org
version: v1.2.0
description: my plugin
enabled: true
hidden: true
artifacts:
    darwin:
        file: my-plugin
tags: []
`,
		},
		{
			scenario: "artifact with arch",
			plugin: Plugin{
				Name:        "my-plugin",
				URL:         "https://example.org",
				Version:     "v1.2.0",
				Description: "my plugin",
				Enabled:     true,
				Hidden:      true,
				Artifacts: Artifacts{
					NewArtifactIdentifier("darwin", "amd64"): {
						File: "my-plugin",
					},
				},
			},
			expected: `name: my-plugin
url: https://example.org
version: v1.2.0
description: my plugin
enabled: true
hidden: true
artifacts:
    darwin/amd64:
        file: my-plugin
tags: []
`,
		},
		{
			scenario: "multiple artifacts",
			plugin: Plugin{
				Name:        "my-plugin",
				URL:         "https://example.org",
				Version:     "v1.2.0",
				Description: "my plugin",
				Enabled:     true,
				Hidden:      true,
				Artifacts: Artifacts{
					NewArtifactIdentifier("linux", "amd64"): {
						File: "my-plugin",
					},
					NewArtifactIdentifier("windows", "amd64"): {
						File: "my-plugin",
					},
					NewArtifactIdentifier("darwin", "amd64"): {
						File: "my-plugin",
					},
				},
			},
			expected: `name: my-plugin
url: https://example.org
version: v1.2.0
description: my plugin
enabled: true
hidden: true
artifacts:
    darwin/amd64:
        file: my-plugin
    linux/amd64:
        file: my-plugin
    windows/amd64:
        file: my-plugin
tags: []
`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result, err := yaml.Marshal(tc.plugin)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, string(result))
		})
	}
}

func TestPlugin_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	os := runtime.GOOS
	arch := runtime.GOARCH

	testCases := []struct {
		scenario       string
		data           string
		expectedResult Plugin
		expectedError  string
	}{
		{
			scenario: "default value",
			data:     "name: my-plugin",
			expectedResult: Plugin{
				Name:    "my-plugin",
				Enabled: true,
				Artifacts: map[ArtifactIdentifier]Artifact{
					RuntimeArtifactIdentifier(): {
						File: "${name}-${version}-${os}-${arch}.tar.gz",
					},
				},
			},
		},
		{
			scenario: "override value",
			data:     "name: my-plugin\nenabled: false",
			expectedResult: Plugin{
				Name:    "my-plugin",
				Enabled: false,
				Artifacts: map[ArtifactIdentifier]Artifact{
					RuntimeArtifactIdentifier(): {
						File: "${name}-${version}-${os}-${arch}.tar.gz",
					},
				},
			},
		},
		{
			scenario: "override artifact",
			data: fmt.Sprintf(`
name: my-plugin
enabled: false
artifacts:
    %s/%s:
        file: "another-file"
`, os, arch),
			expectedResult: Plugin{
				Name:    "my-plugin",
				Enabled: false,
				Artifacts: map[ArtifactIdentifier]Artifact{
					RuntimeArtifactIdentifier(): {
						File: "another-file",
					},
				},
			},
		},
		{
			scenario: "runtime artifact without arch",
			data: fmt.Sprintf(`
name: my-plugin
enabled: false
artifacts:
    %s:
        file: "another-file"
`, os),
			expectedResult: Plugin{
				Name:    "my-plugin",
				Enabled: false,
				Artifacts: map[ArtifactIdentifier]Artifact{
					NewArtifactIdentifier(os, ""): {
						File: "another-file",
					},
				},
			},
		},
		{
			scenario: "artifact anchor",
			data: fmt.Sprintf(`
name: my-plugin
enabled: false
artifacts:
    %s/%s: &default
        file: "${name}-${version}-${os}-${arch}.tar.gz"
    os/arch: *default
`, os, arch),
			expectedResult: Plugin{
				Name:    "my-plugin",
				Enabled: false,
				Artifacts: map[ArtifactIdentifier]Artifact{
					RuntimeArtifactIdentifier(): {
						File: "${name}-${version}-${os}-${arch}.tar.gz",
					},
					NewArtifactIdentifier("os", "arch"): {
						File: "${name}-${version}-${os}-${arch}.tar.gz",
					},
				},
			},
		},
		{
			scenario: "artifact without arch",
			data: `
name: my-plugin
enabled: false
artifacts:
    os:
        file: "${name}-${version}-${os}.tar.gz"
`,
			expectedResult: Plugin{
				Name:    "my-plugin",
				Enabled: false,
				Artifacts: map[ArtifactIdentifier]Artifact{
					RuntimeArtifactIdentifier(): {
						File: "${name}-${version}-${os}-${arch}.tar.gz",
					},
					NewArtifactIdentifier("os", ""): {
						File: "${name}-${version}-${os}.tar.gz",
					},
				},
			},
		},
		{
			scenario: "full info",
			data: `
name: my-plugin
url: "github.com/john.doe/my-plugin"
version: "v1.0.1"
description: my plugin 
enabled: false
hidden: true
artifacts:
    os:
        file: "${name}-${version}-${os}.tar.gz"
tags:
    - tag1
`,
			expectedResult: Plugin{
				Name:        "my-plugin",
				URL:         "github.com/john.doe/my-plugin",
				Version:     "v1.0.1",
				Description: "my plugin",
				Enabled:     false,
				Hidden:      true,
				Artifacts: map[ArtifactIdentifier]Artifact{
					RuntimeArtifactIdentifier(): {
						File: "${name}-${version}-${os}-${arch}.tar.gz",
					},
					NewArtifactIdentifier("os", ""): {
						File: "${name}-${version}-${os}.tar.gz",
					},
				},
				Tags: Tags{"tag1"},
			},
		},
		{
			scenario:      "invalid plugin",
			data:          `42`,
			expectedError: "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!int `42` into plugin.rawPlugin",
		},
		{
			scenario: "invalid artifact identifier",
			data: `
name: my-plugin
url: "github.com/john.doe/my-plugin"
version: "v1.0.1"
enabled: false
hidden: true
artifacts:
    {}:
        file: "${name}-${version}-${os}.tar.gz"
`,
			expectedError: "yaml: unmarshal errors:\n  line 8: cannot unmarshal !!map into string",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			var result Plugin

			err := yaml.Unmarshal([]byte(tc.data), &result)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestTags_Contains(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		tags := Tags{"tag1", "tag2", "tag3"}

		assert.True(t, tags.Contains("tag2"))
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		tags := Tags{"tag1", "tag2", "tag3"}

		assert.False(t, tags.Contains("tag4"))
	})
}

func TestRuntimeArtifactIdentifier(t *testing.T) {
	t.Parallel()

	expected := ArtifactIdentifier{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	assert.Equal(t, expected, RuntimeArtifactIdentifier())
}

func TestLoad(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		makeFs         func() afero.Fs
		expectedResult *Plugin
		expectedError  string
	}{
		{
			scenario:      "file does not exist",
			makeFs:        afero.NewMemMapFs,
			expectedError: `could not read metadata: open /tmp/.plugin.registry.yaml: file does not exist`,
		},
		{
			scenario: "could not decode",
			makeFs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/tmp/.plugin.registry.yaml", []byte("[]"), 0755) // nolint: errcheck

				return fs
			},
			expectedError: "could not read metadata: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!seq into plugin.rawPlugin",
		},
		{
			scenario: "success",
			makeFs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/tmp/.plugin.registry.yaml", []byte("name: my-plugin"), 0755) // nolint: errcheck

				return fs
			},
			expectedResult: &Plugin{
				Name:    "my-plugin",
				Enabled: true,
				Artifacts: Artifacts{
					RuntimeArtifactIdentifier(): {
						File: defaultFile,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result, err := Load(tc.makeFs(), "/tmp")

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestLoadError(t *testing.T) {
	t.Parallel()

	err := loadError(errors.New("error"), "/tmp")
	expected := `could not read metadata: error`

	assert.EqualError(t, err, expected)
}
