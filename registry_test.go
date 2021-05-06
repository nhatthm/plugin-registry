package registry

import (
	"errors"
	"testing"

	"github.com/nhatthm/aferomock"
	"github.com/nhatthm/plugin-registry/config"
	configuratorMock "github.com/nhatthm/plugin-registry/mock/configurator"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistry(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		makeFs        func() afero.Fs
		expectedError string
	}{
		{
			scenario: "failure",
			makeFs: func() afero.Fs {
				return aferomock.MockFs(func(fs *aferomock.Fs) {
					fs.On("Stat", "/tmp/config.yaml").
						Return(nil, errors.New("read error"))
				})(t)
			},
			expectedError: "read error",
		},
		{
			scenario: "success",
			makeFs:   afero.NewMemMapFs,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			r, err := NewRegistry("/tmp", WithFs(tc.makeFs()))

			if tc.expectedError == "" {
				assert.NotNil(t, r)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, r)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestWithFs(t *testing.T) {
	t.Parallel()

	expected := afero.NewMemMapFs()
	r, err := NewRegistry("/tmp", WithFs(expected))
	require.NoError(t, err)

	assert.Equal(t, expected, r.fs)
}

func TestWithConfigurator(t *testing.T) {
	t.Parallel()

	expected := configuratorMock.NoMock(t)
	r, err := NewRegistry("/tmp", WithConfigurator(expected))
	require.NoError(t, err)

	assert.Equal(t, expected, r.config)
}

func TestWithConfigFile(t *testing.T) {
	t.Parallel()

	expected, err := config.NewMemConfigurator(
		config.NewFileConfigurator("/tmp/plugins/config.yaml"),
	)
	require.NoError(t, err)

	r, err := NewRegistry("/tmp", WithConfigFile("/tmp/plugins/config.yaml"))
	require.NoError(t, err)

	assert.Equal(t, expected, r.config)
}
