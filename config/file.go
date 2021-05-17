package config

import (
	"context"
	"os"
	"sync"

	"github.com/bool64/ctxd"
	"github.com/nhatthm/plugin-registry/plugin"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

var _ Configurator = (*FileConfigurator)(nil)

// FileConfigurator is a file configurator.
type FileConfigurator struct {
	fs afero.Fs

	configFile string

	mu sync.Mutex
}

func (c *FileConfigurator) lock() {
	c.mu.Lock()
}

func (c *FileConfigurator) unlock() {
	c.mu.Unlock()
}

// WithFs sets file system for the configurator.
func (c *FileConfigurator) WithFs(fs afero.Fs) *FileConfigurator {
	c.lock()
	defer c.unlock()

	c.fs = fs

	return c
}

// Config returns the current configuration.
func (c *FileConfigurator) Config() (Configuration, error) {
	c.lock()
	defer c.unlock()

	return c.loadLocked()
}

// SetPlugin sets a plugin.
func (c *FileConfigurator) SetPlugin(p plugin.Plugin) error {
	c.lock()
	defer c.unlock()

	cfg, err := c.loadLocked()
	if err != nil {
		return err
	}

	if cfg.Plugins == nil {
		cfg.Plugins = make(map[string]plugin.Plugin)
	}

	cfg.Plugins[p.Name] = p

	return c.writeLocked(cfg)
}

// RemovePlugin removes a plugin.
func (c *FileConfigurator) RemovePlugin(name string) error {
	c.lock()
	defer c.unlock()

	cfg, err := c.loadLocked()
	if err != nil {
		return err
	}

	if !cfg.Plugins.Has(name) {
		return plugin.ErrPluginNotExist
	}

	delete(cfg.Plugins, name)

	return c.writeLocked(cfg)
}

// EnablePlugin disable a plugin by name.
func (c *FileConfigurator) EnablePlugin(name string) error {
	c.lock()
	defer c.unlock()

	cfg, err := c.loadLocked()
	if err != nil {
		return err
	}

	p, ok := cfg.Plugins[name]
	if !ok {
		return plugin.ErrPluginNotExist
	}

	p.Enabled = true
	cfg.Plugins[name] = p

	return c.writeLocked(cfg)
}

// DisablePlugin disable a plugin by name.
func (c *FileConfigurator) DisablePlugin(name string) error {
	c.lock()
	defer c.unlock()

	cfg, err := c.loadLocked()
	if err != nil {
		return err
	}

	p, ok := cfg.Plugins[name]
	if !ok {
		return plugin.ErrPluginNotExist
	}

	p.Enabled = false
	cfg.Plugins[name] = p

	return c.writeLocked(cfg)
}

func (c *FileConfigurator) writeLocked(cfg Configuration) error {
	f, err := c.fs.OpenFile(c.configFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck

	enc := yaml.NewEncoder(f)

	return enc.Encode(cfg)
}

// loadLocked loads the current configuration.
func (c *FileConfigurator) loadLocked() (Configuration, error) {
	if exists, err := afero.Exists(c.fs, c.configFile); !exists {
		return Configuration{}, err
	}

	f, err := c.fs.Open(c.configFile)
	if err != nil {
		return Configuration{}, err
	}
	defer f.Close() // nolint: errcheck

	var cfg Configuration

	dec := yaml.NewDecoder(f)

	if err := dec.Decode(&cfg); err != nil {
		return Configuration{}, ctxd.WrapError(context.Background(), err, "could not load configuration",
			"path", c.configFile,
		)
	}

	return cfg, nil
}

// NewFileConfigurator initiates a new FileConfigurator.
func NewFileConfigurator(configFile string) *FileConfigurator {
	return &FileConfigurator{
		fs:         afero.NewOsFs(),
		configFile: configFile,
	}
}
