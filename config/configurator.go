package config

import (
	"github.com/nhatthm/plugin-registry/plugin"
)

// Configurator is a configuration manager.
type Configurator interface {
	Config() (Configuration, error)
	SetPlugin(plugin plugin.Plugin) error
	RemovePlugin(name string) error
	EnablePlugin(name string) error
	DisablePlugin(name string) error
}
