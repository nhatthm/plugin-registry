package registry

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	configuratorMock "github.com/nhatthm/plugin-registry/mock/configurator"
)

func TestRegistry_Enable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockConfig    configuratorMock.Mocker
		expectedError string
	}{
		{
			scenario: "failure",
			mockConfig: configuratorMock.Mock(func(c *configuratorMock.Configurator) {
				c.On("EnablePlugin", "my-plugin").
					Return(errors.New("disable error"))
			}),
			expectedError: "disable error",
		},
		{
			scenario: "success",
			mockConfig: configuratorMock.Mock(func(c *configuratorMock.Configurator) {
				c.On("EnablePlugin", "my-plugin").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			r, err := NewRegistry("/tmp", WithConfigurator(tc.mockConfig(t)))
			require.NoError(t, err)

			err = r.Enable("my-plugin")

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
