package configurator

import (
	"testing"

	"github.com/nhatthm/plugin-registry/config"
	"github.com/nhatthm/plugin-registry/plugin"
	"github.com/stretchr/testify/assert"
	tMock "github.com/stretchr/testify/mock"
)

// Mocker is Configurator mocker.
type Mocker func(tb testing.TB) *Configurator

// NoMock is no tMock Configurator.
var NoMock = Mock()

var _ config.Configurator = (*Configurator)(nil)

// Configurator is a config.Configurator.
type Configurator struct {
	tMock.Mock
}

// Config satisfies config.Configurator interface.
func (c *Configurator) Config() (config.Configuration, error) {
	ret := c.Called()

	return ret.Get(0).(config.Configuration), ret.Error(1)
}

// SetPlugin satisfies config.Configurator interface.
func (c *Configurator) SetPlugin(plugin plugin.Plugin) error {
	return c.Called(plugin).Error(0)
}

// RemovePlugin satisfies config.Configurator interface.
func (c *Configurator) RemovePlugin(name string) error {
	return c.Called(name).Error(0)
}

// EnablePlugin satisfies config.Configurator interface.
func (c *Configurator) EnablePlugin(name string) error {
	return c.Called(name).Error(0)
}

// DisablePlugin satisfies config.Configurator interface.
func (c *Configurator) DisablePlugin(name string) error {
	return c.Called(name).Error(0)
}

// New mocks config.Configurator interface.
func New(mocks ...func(c *Configurator)) *Configurator {
	c := &Configurator{}

	for _, m := range mocks {
		m(c)
	}

	return c
}

// Mock creates Configurator tMock with cleanup to ensure all the expectations are met.
func Mock(mocks ...func(c *Configurator)) Mocker {
	return func(tb testing.TB) *Configurator {
		tb.Helper()

		c := New(mocks...)

		tb.Cleanup(func() {
			assert.True(tb, c.Mock.AssertExpectations(tb))
		})

		return c
	}
}
