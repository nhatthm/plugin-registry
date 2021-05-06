package registry_test

import (
	"context"
	"errors"
	"testing"

	registry "github.com/nhatthm/plugin-registry"
	"github.com/nhatthm/plugin-registry/installer"
	configuratorMock "github.com/nhatthm/plugin-registry/mock/configurator"
	installerMock "github.com/nhatthm/plugin-registry/mock/installer"
	"github.com/nhatthm/plugin-registry/plugin"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFsRegistry_Install(t *testing.T) {
	t.Parallel()

	registerInstaller := func(caseName string, construct func(t *testing.T) installer.Installer) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			installer.Register(caseName, func(_ context.Context, source string) bool {
				return source == caseName
			}, func(afero.Fs) installer.Installer {
				return construct(t)
			})
		}
	}

	registerFailInstaller := registerInstaller("INSTALL_FAIL", func(t *testing.T) installer.Installer {
		t.Helper()

		return installerMock.Mock(func(i *installerMock.Installer) {
			i.On("Install", context.Background(), "/tmp", "INSTALL_FAIL").
				Return(nil, errors.New("install error"))
		})(t)
	})

	registerSuccessInstaller := registerInstaller("INSTALL_SUCCESS", func(t *testing.T) installer.Installer {
		t.Helper()

		return installerMock.Mock(func(i *installerMock.Installer) {
			i.On("Install", context.Background(), "/tmp", "INSTALL_SUCCESS").
				Return(&plugin.Plugin{Name: "my-plugin"}, nil)
		})(t)
	})

	testCases := []struct {
		scenario          string
		registerInstaller func(t *testing.T)
		mockConfig        configuratorMock.Mocker
		source            string
		expectedError     string
	}{
		{
			scenario:      "could not find installer",
			source:        "UNKNOWN",
			expectedError: "no supported installer",
		},
		{
			scenario:          "could not install",
			registerInstaller: registerFailInstaller,
			source:            "INSTALL_FAIL",
			expectedError:     "install error",
		},
		{
			scenario:          "could not update configuration",
			registerInstaller: registerSuccessInstaller,
			mockConfig: configuratorMock.Mock(func(c *configuratorMock.Configurator) {
				c.On("SetPlugin", plugin.Plugin{Name: "my-plugin"}).
					Return(errors.New("config error"))
			}),
			source:        "INSTALL_SUCCESS",
			expectedError: "config error",
		},
		{
			scenario:          "success",
			registerInstaller: registerSuccessInstaller,
			mockConfig: configuratorMock.Mock(func(c *configuratorMock.Configurator) {
				c.On("SetPlugin", plugin.Plugin{Name: "my-plugin"}).
					Return(nil)
			}),
			source: "INSTALL_SUCCESS",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			if tc.registerInstaller != nil {
				tc.registerInstaller(t)
			}

			if tc.mockConfig == nil {
				tc.mockConfig = configuratorMock.NoMock
			}

			r, err := registry.NewRegistry("/tmp", registry.WithConfigurator(tc.mockConfig(t)))
			require.NoError(t, err)

			err = r.Install(context.Background(), tc.source)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
