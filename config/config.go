package config

import "github.com/nhatthm/plugin-registry/plugin"

// Configuration represents configuration of the registry.
type Configuration struct {
	Plugins plugin.Plugins `yaml:"plugins"`
}
