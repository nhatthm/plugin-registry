package registry

import (
	"context"
	"testing"

	registry "github.com/nhatthm/plugin-registry"
	"github.com/stretchr/testify/assert"
	tMock "github.com/stretchr/testify/mock"
)

// Mocker is Registry mocker.
type Mocker func(tb testing.TB) *Registry

// NoMock is no tMock Registry.
var NoMock = MockRegistry()

var _ registry.Registry = (*Registry)(nil)

// Registry is a registry.Registry.
type Registry struct {
	tMock.Mock
}

// Enable satisfies registry.Registry.
func (r *Registry) Enable(name string) error {
	return r.Called(name).Error(0)
}

// Disable satisfies registry.Registry.
func (r *Registry) Disable(name string) error {
	return r.Called(name).Error(0)
}

// Install satisfies registry.Registry.
func (r *Registry) Install(ctx context.Context, source string) error {
	return r.Called(ctx, source).Error(0)
}

// Uninstall satisfies registry.Registry.
func (r *Registry) Uninstall(name string) error {
	return r.Called(name).Error(0)
}

// New mocks registry.Registry interface.
func New(mocks ...func(r *Registry)) *Registry {
	r := &Registry{}

	for _, m := range mocks {
		m(r)
	}

	return r
}

// MockRegistry creates Registry tMock with cleanup to ensure all the expectations are met.
func MockRegistry(mocks ...func(r *Registry)) Mocker {
	return func(tb testing.TB) *Registry {
		tb.Helper()

		r := New(mocks...)

		tb.Cleanup(func() {
			assert.True(tb, r.Mock.AssertExpectations(tb))
		})

		return r
	}
}
