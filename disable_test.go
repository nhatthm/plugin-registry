package registry

import (
	"errors"
	"testing"

	configuratorMock "github.com/nhatthm/plugin-registry/mock/configurator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_Disable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockConfig    configuratorMock.Mocker
		expectedError string
	}{
		{
			scenario: "failure",
			mockConfig: configuratorMock.Mock(func(c *configuratorMock.Configurator) {
				c.On("DisablePlugin", "my-plugin").
					Return(errors.New("disable error"))
			}),
			expectedError: "disable error",
		},
		{
			scenario: "success",
			mockConfig: configuratorMock.Mock(func(c *configuratorMock.Configurator) {
				c.On("DisablePlugin", "my-plugin").
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

			err = r.Disable("my-plugin")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
