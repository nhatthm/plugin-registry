package config

import (
	"sync"

	"github.com/nhatthm/plugin-registry/plugin"
)

var _ Configurator = (*MemConfigurator)(nil)

// MemConfigurator is a memory configurator.
type MemConfigurator struct {
	upstream Configurator

	config Configuration

	mu sync.Mutex
}

func (c *MemConfigurator) init() error {
	cfg, err := c.upstream.Config()
	if err != nil {
		return err
	}

	c.config = cfg

	return nil
}

// Config returns the current configuration.
func (c *MemConfigurator) Config() (Configuration, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config, nil
}

// SetPlugin sets a plugin.
func (c *MemConfigurator) SetPlugin(p plugin.Plugin) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.upstream.SetPlugin(p); err != nil {
		return err
	}

	if c.config.Plugins == nil {
		c.config.Plugins = make(map[string]plugin.Plugin)
	}

	c.config.Plugins[p.Name] = p

	return nil
}

// RemovePlugin removes a plugin.
func (c *MemConfigurator) RemovePlugin(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.upstream.RemovePlugin(name); err != nil {
		return err
	}

	delete(c.config.Plugins, name)

	return nil
}

// EnablePlugin disable a plugin by name.
func (c *MemConfigurator) EnablePlugin(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.upstream.EnablePlugin(name); err != nil {
		return err
	}

	p, ok := c.config.Plugins[name]
	if !ok {
		return plugin.ErrPluginNotExist
	}

	p.Enabled = true

	c.config.Plugins[name] = p

	return nil
}

// DisablePlugin disable a plugin by name.
func (c *MemConfigurator) DisablePlugin(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.upstream.DisablePlugin(name); err != nil {
		return err
	}

	p, ok := c.config.Plugins[name]
	if !ok {
		return plugin.ErrPluginNotExist
	}

	p.Enabled = false

	c.config.Plugins[name] = p

	return nil
}

// NewMemConfigurator initiates a new MemConfigurator.
func NewMemConfigurator(upstream Configurator) (*MemConfigurator, error) {
	c := &MemConfigurator{upstream: upstream}

	if err := c.init(); err != nil {
		return nil, err
	}

	return c, nil
}
