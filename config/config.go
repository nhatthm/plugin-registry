package config

import "github.com/nhatthm/plugin-registry/plugin"

// Configuration represents configuration of the registry.
type Configuration struct {
	Plugins map[string]plugin.Plugin `yaml:"plugins"`
}
