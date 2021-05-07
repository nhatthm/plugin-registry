package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/nhatthm/plugin-registry/config"
	"github.com/nhatthm/plugin-registry/plugin"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockRegistry   Mocker
		expectedConfig config.Configuration
		expectedError  string
	}{
		{
			scenario: "error",
			mockRegistry: MockRegistry(func(r *Registry) {
				r.On("Config").
					Return(config.Configuration{}, errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "no error",
			mockRegistry: MockRegistry(func(r *Registry) {
				r.On("Config").
					Return(config.Configuration{
						Plugins: plugin.Plugins{},
					}, nil)
			}),
			expectedConfig: config.Configuration{
				Plugins: plugin.Plugins{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			cfg, err := tc.mockRegistry(t).Config()

			assert.Equal(t, tc.expectedConfig, cfg)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestEnable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockRegistry  Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mockRegistry: MockRegistry(func(r *Registry) {
				r.On("Enable", "my-plugin").
					Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "no error",
			mockRegistry: MockRegistry(func(r *Registry) {
				r.On("Enable", "my-plugin").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mockRegistry(t).Enable("my-plugin")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestDisable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockRegistry  Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mockRegistry: MockRegistry(func(r *Registry) {
				r.On("Disable", "my-plugin").
					Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "no error",
			mockRegistry: MockRegistry(func(r *Registry) {
				r.On("Disable", "my-plugin").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mockRegistry(t).Disable("my-plugin")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestInstall(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockRegistry  Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mockRegistry: MockRegistry(func(r *Registry) {
				r.On("Install", context.Background(), "my-plugin").
					Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "no error",
			mockRegistry: MockRegistry(func(r *Registry) {
				r.On("Install", context.Background(), "my-plugin").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mockRegistry(t).Install(context.Background(), "my-plugin")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestUninstall(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockRegistry  Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mockRegistry: MockRegistry(func(r *Registry) {
				r.On("Uninstall", "my-plugin").
					Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "no error",
			mockRegistry: MockRegistry(func(r *Registry) {
				r.On("Uninstall", "my-plugin").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mockRegistry(t).Uninstall("my-plugin")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
