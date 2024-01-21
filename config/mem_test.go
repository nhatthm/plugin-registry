package config_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nhatthm/plugin-registry/config"
	"github.com/nhatthm/plugin-registry/mock/configurator"
	"github.com/nhatthm/plugin-registry/plugin"
)

func TestNewMemConfigurator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockUpstream  configurator.Mocker
		expectedError string
	}{
		{
			scenario: "could not init",
			mockUpstream: configurator.Mock(func(c *configurator.Configurator) {
				c.On("Config").
					Return(config.Configuration{}, errors.New("upstream error"))
			}),
			expectedError: "upstream error",
		},
		{
			scenario: "success",
			mockUpstream: configurator.Mock(func(c *configurator.Configurator) {
				c.On("Config").
					Return(config.Configuration{}, nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			c, err := config.NewMemConfigurator(tc.mockUpstream(t))

			if tc.expectedError == "" {
				assert.NotNil(t, c)
				require.NoError(t, err)
			} else {
				assert.Nil(t, c)
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestMemConfigurator_Config(t *testing.T) {
	t.Parallel()

	upstream := configurator.Mock(func(c *configurator.Configurator) {
		c.On("Config").
			Return(
				config.Configuration{
					Plugins: map[string]plugin.Plugin{
						"my-plugin": {
							Name:    "my-plugin",
							Enabled: true,
							Hidden:  true,
							Artifacts: plugin.Artifacts{
								plugin.RuntimeArtifactIdentifier(): {
									File: "my-plugin",
								},
							},
						},
					},
				},
				nil,
			)
	})(t)

	c, err := config.NewMemConfigurator(upstream)
	require.NoError(t, err)

	expected := config.Configuration{
		Plugins: map[string]plugin.Plugin{
			"my-plugin": {
				Name:    "my-plugin",
				Enabled: true,
				Hidden:  true,
				Artifacts: plugin.Artifacts{
					plugin.RuntimeArtifactIdentifier(): {
						File: "my-plugin",
					},
				},
			},
		},
	}

	actual, err := c.Config()
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestMemConfigurator_SetPlugin(t *testing.T) {
	t.Parallel()

	cfg := config.Configuration{}
	mockConfig := func(c *configurator.Configurator) {
		c.On("Config").Return(cfg, nil)
	}

	p := plugin.Plugin{
		Name:    "my-plugin",
		Enabled: true,
		Hidden:  true,
		Artifacts: plugin.Artifacts{
			plugin.RuntimeArtifactIdentifier(): {
				File: "my-plugin",
			},
		},
	}

	testCases := []struct {
		scenario       string
		mockUpstream   configurator.Mocker
		expectedConfig config.Configuration
		expectedError  string
	}{
		{
			scenario: "failure",
			mockUpstream: configurator.Mock(mockConfig, func(c *configurator.Configurator) {
				c.On("SetPlugin", p).
					Return(errors.New("set error"))
			}),
			expectedError: "set error",
		},
		{
			scenario: "success",
			mockUpstream: configurator.Mock(mockConfig, func(c *configurator.Configurator) {
				c.On("SetPlugin", p).
					Return(nil)
			}),
			expectedConfig: config.Configuration{
				Plugins: map[string]plugin.Plugin{
					"my-plugin": {
						Name:    "my-plugin",
						Enabled: true,
						Hidden:  true,
						Artifacts: plugin.Artifacts{
							plugin.RuntimeArtifactIdentifier(): {
								File: "my-plugin",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			c, err := config.NewMemConfigurator(tc.mockUpstream(t))
			require.NoError(t, err)

			err = c.SetPlugin(p)
			actualCfg, cfgErr := c.Config()
			require.NoError(t, cfgErr)

			assert.Equal(t, tc.expectedConfig, actualCfg)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestMemConfigurator_RemovePlugin(t *testing.T) {
	t.Parallel()

	mockConfig := func(c *configurator.Configurator) {
		c.On("Config").Return(config.Configuration{
			Plugins: map[string]plugin.Plugin{
				"my-plugin": {
					Name:    "my-plugin",
					Enabled: true,
					Hidden:  true,
					Artifacts: plugin.Artifacts{
						plugin.RuntimeArtifactIdentifier(): {
							File: "my-plugin",
						},
					},
				},
			},
		}, nil)
	}

	testCases := []struct {
		scenario       string
		mockUpstream   configurator.Mocker
		expectedConfig config.Configuration
		expectedError  string
	}{
		{
			scenario: "failure",
			mockUpstream: configurator.Mock(mockConfig, func(c *configurator.Configurator) {
				c.On("RemovePlugin", "my-plugin").
					Return(errors.New("remove error"))
			}),
			expectedConfig: config.Configuration{
				Plugins: map[string]plugin.Plugin{
					"my-plugin": {
						Name:    "my-plugin",
						Enabled: true,
						Hidden:  true,
						Artifacts: plugin.Artifacts{
							plugin.RuntimeArtifactIdentifier(): {
								File: "my-plugin",
							},
						},
					},
				},
			},
			expectedError: "remove error",
		},
		{
			scenario: "success",
			mockUpstream: configurator.Mock(mockConfig, func(c *configurator.Configurator) {
				c.On("RemovePlugin", "my-plugin").
					Return(nil)
			}),
			expectedConfig: config.Configuration{
				Plugins: map[string]plugin.Plugin{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			c, err := config.NewMemConfigurator(tc.mockUpstream(t))
			require.NoError(t, err)

			err = c.RemovePlugin("my-plugin")
			actualCfg, cfgErr := c.Config()
			require.NoError(t, cfgErr)

			assert.Equal(t, tc.expectedConfig, actualCfg)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestMemConfigurator_EnablePlugin(t *testing.T) {
	t.Parallel()

	mockConfig := func(c *configurator.Configurator) {
		c.On("Config").Return(config.Configuration{
			Plugins: map[string]plugin.Plugin{
				"my-plugin": {
					Name:    "my-plugin",
					Enabled: false,
					Hidden:  true,
					Artifacts: plugin.Artifacts{
						plugin.RuntimeArtifactIdentifier(): {
							File: "my-plugin",
						},
					},
				},
			},
		}, nil)
	}

	testCases := []struct {
		scenario       string
		mockUpstream   configurator.Mocker
		pluginName     string
		expectedConfig config.Configuration
		expectedError  string
	}{
		{
			scenario: "failure",
			mockUpstream: configurator.Mock(mockConfig, func(c *configurator.Configurator) {
				c.On("EnablePlugin", "my-plugin").
					Return(errors.New("enable error"))
			}),
			expectedConfig: config.Configuration{
				Plugins: map[string]plugin.Plugin{
					"my-plugin": {
						Name:    "my-plugin",
						Enabled: false,
						Hidden:  true,
						Artifacts: plugin.Artifacts{
							plugin.RuntimeArtifactIdentifier(): {
								File: "my-plugin",
							},
						},
					},
				},
			},
			expectedError: "enable error",
		},
		{
			scenario: "not found",
			mockUpstream: configurator.Mock(mockConfig, func(c *configurator.Configurator) {
				c.On("EnablePlugin", "unknown").
					Return(nil)
			}),
			pluginName: "unknown",
			expectedConfig: config.Configuration{
				Plugins: map[string]plugin.Plugin{
					"my-plugin": {
						Name:    "my-plugin",
						Enabled: false,
						Hidden:  true,
						Artifacts: plugin.Artifacts{
							plugin.RuntimeArtifactIdentifier(): {
								File: "my-plugin",
							},
						},
					},
				},
			},
			expectedError: "plugin does not exist",
		},
		{
			scenario: "success",
			mockUpstream: configurator.Mock(mockConfig, func(c *configurator.Configurator) {
				c.On("EnablePlugin", "my-plugin").
					Return(nil)
			}),
			expectedConfig: config.Configuration{
				Plugins: map[string]plugin.Plugin{
					"my-plugin": {
						Name:    "my-plugin",
						Enabled: true,
						Hidden:  true,
						Artifacts: plugin.Artifacts{
							plugin.RuntimeArtifactIdentifier(): {
								File: "my-plugin",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			if tc.pluginName == "" {
				tc.pluginName = "my-plugin"
			}

			c, err := config.NewMemConfigurator(tc.mockUpstream(t))
			require.NoError(t, err)

			err = c.EnablePlugin(tc.pluginName)
			actualCfg, cfgErr := c.Config()
			require.NoError(t, cfgErr)

			assert.Equal(t, tc.expectedConfig, actualCfg)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestMemConfigurator_DisablePlugin(t *testing.T) {
	t.Parallel()

	mockConfig := func(c *configurator.Configurator) {
		c.On("Config").Return(config.Configuration{
			Plugins: map[string]plugin.Plugin{
				"my-plugin": {
					Name:    "my-plugin",
					Enabled: true,
					Hidden:  true,
					Artifacts: plugin.Artifacts{
						plugin.RuntimeArtifactIdentifier(): {
							File: "my-plugin",
						},
					},
				},
			},
		}, nil)
	}

	testCases := []struct {
		scenario       string
		mockUpstream   configurator.Mocker
		pluginName     string
		expectedConfig config.Configuration
		expectedError  string
	}{
		{
			scenario: "failure",
			mockUpstream: configurator.Mock(mockConfig, func(c *configurator.Configurator) {
				c.On("DisablePlugin", "my-plugin").
					Return(errors.New("disable error"))
			}),
			expectedConfig: config.Configuration{
				Plugins: map[string]plugin.Plugin{
					"my-plugin": {
						Name:    "my-plugin",
						Enabled: true,
						Hidden:  true,
						Artifacts: plugin.Artifacts{
							plugin.RuntimeArtifactIdentifier(): {
								File: "my-plugin",
							},
						},
					},
				},
			},
			expectedError: "disable error",
		},
		{
			scenario: "not found",
			mockUpstream: configurator.Mock(mockConfig, func(c *configurator.Configurator) {
				c.On("DisablePlugin", "unknown").
					Return(nil)
			}),
			pluginName: "unknown",
			expectedConfig: config.Configuration{
				Plugins: map[string]plugin.Plugin{
					"my-plugin": {
						Name:    "my-plugin",
						Enabled: true,
						Hidden:  true,
						Artifacts: plugin.Artifacts{
							plugin.RuntimeArtifactIdentifier(): {
								File: "my-plugin",
							},
						},
					},
				},
			},
			expectedError: "plugin does not exist",
		},
		{
			scenario: "success",
			mockUpstream: configurator.Mock(mockConfig, func(c *configurator.Configurator) {
				c.On("DisablePlugin", "my-plugin").
					Return(nil)
			}),
			expectedConfig: config.Configuration{
				Plugins: map[string]plugin.Plugin{
					"my-plugin": {
						Name:    "my-plugin",
						Enabled: false,
						Hidden:  true,
						Artifacts: plugin.Artifacts{
							plugin.RuntimeArtifactIdentifier(): {
								File: "my-plugin",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			if tc.pluginName == "" {
				tc.pluginName = "my-plugin"
			}

			c, err := config.NewMemConfigurator(tc.mockUpstream(t))
			require.NoError(t, err)

			err = c.DisablePlugin(tc.pluginName)
			actualCfg, cfgErr := c.Config()
			require.NoError(t, cfgErr)

			assert.Equal(t, tc.expectedConfig, actualCfg)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
