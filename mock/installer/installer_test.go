package installer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nhatthm/plugin-registry/plugin"
)

func TestInstaller_Install(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mock           Mocker
		expectedResult *plugin.Plugin
		expectedError  string
	}{
		{
			scenario: "plugin is nil",
			mock: Mock(func(i *Installer) {
				i.On("Install", context.Background(), "/tmp", "mock").
					Return(nil, errors.New("install error"))
			}),
			expectedError: `install error`,
		},
		{
			scenario: "plugin is not nil",
			mock: Mock(func(i *Installer) {
				i.On("Install", context.Background(), "/tmp", "mock").
					Return(&plugin.Plugin{Name: "my-plugin"}, nil)
			}),
			expectedResult: &plugin.Plugin{Name: "my-plugin"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p, err := tc.mock(t).Install(context.Background(), "/tmp", "mock")

			assert.Equal(t, tc.expectedResult, p)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
