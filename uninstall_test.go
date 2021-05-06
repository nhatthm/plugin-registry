package registry

import (
	"errors"
	"testing"

	"github.com/nhatthm/aferomock"
	configuratorMock "github.com/nhatthm/plugin-registry/mock/configurator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFsRegistry_Uninstall(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockFs        aferomock.FsMocker
		mockConfig    configuratorMock.Mocker
		expectedError string
	}{
		{
			scenario: "remove error",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("RemoveAll", "/tmp/my-plugin").
					Return(errors.New("remove error"))
			}),
			expectedError: "remove error",
		},
		{
			scenario: "config error",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("RemoveAll", "/tmp/my-plugin").
					Return(nil)
			}),
			mockConfig: configuratorMock.Mock(func(c *configuratorMock.Configurator) {
				c.On("RemovePlugin", "my-plugin").
					Return(errors.New("config error"))
			}),
			expectedError: "config error",
		},
		{
			scenario: "success",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("RemoveAll", "/tmp/my-plugin").
					Return(nil)
			}),
			mockConfig: configuratorMock.Mock(func(c *configuratorMock.Configurator) {
				c.On("RemovePlugin", "my-plugin").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			if tc.mockConfig == nil {
				tc.mockConfig = configuratorMock.NoMock
			}

			r, err := NewRegistry("/tmp",
				WithFs(tc.mockFs(t)),
				WithConfigurator(tc.mockConfig(t)),
			)
			require.NoError(t, err)

			err = r.Uninstall("my-plugin")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
