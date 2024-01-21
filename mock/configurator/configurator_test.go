package configurator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nhatthm/plugin-registry/config"
	"github.com/nhatthm/plugin-registry/plugin"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockConfig     Mocker
		expectedConfig config.Configuration
		expectedError  string
	}{
		{
			scenario: "error",
			mockConfig: Mock(func(c *Configurator) {
				c.On("Config").
					Return(config.Configuration{}, errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "no error",
			mockConfig: Mock(func(c *Configurator) {
				c.On("Config").
					Return(config.Configuration{Plugins: map[string]plugin.Plugin{}}, nil)
			}),
			expectedConfig: config.Configuration{Plugins: map[string]plugin.Plugin{}},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			cfg, err := tc.mockConfig(t).Config()

			assert.Equal(t, tc.expectedConfig, cfg)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestSetPlugin(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockConfig    Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mockConfig: Mock(func(c *Configurator) {
				c.On("SetPlugin", plugin.Plugin{Name: "my-plugin"}).
					Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "no error",
			mockConfig: Mock(func(c *Configurator) {
				c.On("SetPlugin", plugin.Plugin{Name: "my-plugin"}).
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mockConfig(t).SetPlugin(plugin.Plugin{Name: "my-plugin"})

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestRemovePlugin(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockConfig    Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mockConfig: Mock(func(c *Configurator) {
				c.On("RemovePlugin", "my-plugin").
					Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "no error",
			mockConfig: Mock(func(c *Configurator) {
				c.On("RemovePlugin", "my-plugin").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mockConfig(t).RemovePlugin("my-plugin")

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestEnablePlugin(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockConfig    Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mockConfig: Mock(func(c *Configurator) {
				c.On("EnablePlugin", "my-plugin").
					Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "no error",
			mockConfig: Mock(func(c *Configurator) {
				c.On("EnablePlugin", "my-plugin").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mockConfig(t).EnablePlugin("my-plugin")

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestDisablePlugin(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockConfig    Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mockConfig: Mock(func(c *Configurator) {
				c.On("DisablePlugin", "my-plugin").
					Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "no error",
			mockConfig: Mock(func(c *Configurator) {
				c.On("DisablePlugin", "my-plugin").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mockConfig(t).DisablePlugin("my-plugin")

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
