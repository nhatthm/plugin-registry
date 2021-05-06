package installer

import (
	"context"
	"testing"

	"github.com/nhatthm/plugin-registry/installer"
	"github.com/nhatthm/plugin-registry/plugin"
	"github.com/stretchr/testify/assert"
	tMock "github.com/stretchr/testify/mock"
)

// Mocker is Installer mocker.
type Mocker func(tb testing.TB) *Installer

// NoMock is no mock Installer.
var NoMock = Mock()

var _ installer.Installer = (*Installer)(nil)

// Installer is a installer.Installer.
type Installer struct {
	tMock.Mock
}

// Install satisfies installer.Installer interface.
func (i *Installer) Install(ctx context.Context, dest, pluginURL string) (*plugin.Plugin, error) {
	ret := i.Called(ctx, dest, pluginURL)

	p := ret.Get(0)
	err := ret.Error(1)

	if p == nil {
		return nil, err
	}

	return p.(*plugin.Plugin), err
}

// New mocks installer.Installer interface.
func New(mocks ...func(i *Installer)) *Installer {
	i := &Installer{}

	for _, m := range mocks {
		m(i)
	}

	return i
}

// Mock creates Installer mock with cleanup to ensure all the expectations are met.
func Mock(mocks ...func(i *Installer)) Mocker {
	return func(tb testing.TB) *Installer {
		tb.Helper()

		i := New(mocks...)

		tb.Cleanup(func() {
			assert.True(tb, i.Mock.AssertExpectations(tb))
		})

		return i
	}
}
