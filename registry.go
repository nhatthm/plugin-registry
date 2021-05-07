package registry

import (
	"context"
	"path/filepath"

	"github.com/nhatthm/plugin-registry/config"
	"github.com/nhatthm/plugin-registry/plugin"
	"github.com/spf13/afero"
)

// Option configures Registry.
type Option func(r *FsRegistry)

// Registry is a plugin registry.
type Registry interface {
	Enable(name string) error
	Disable(name string) error
	Install(ctx context.Context, src string) error
	Uninstall(name string) error
}

// FsRegistry is a file system plugin registry.
type FsRegistry struct {
	fs     afero.Fs
	config config.Configurator

	path       string
	configFile string
}

// Config returns the configuration of the registry.
func (r *FsRegistry) Config() (config.Configuration, error) {
	return r.config.Config()
}

// GetPlugin gets plugin by name.
func (r *FsRegistry) GetPlugin(name string) (*plugin.Plugin, error) {
	cfg, err := r.Config()
	if err != nil {
		return nil, err
	}

	p, ok := cfg.Plugins[name]
	if !ok {
		return nil, nil
	}

	return &p, nil
}

// NewRegistry initiates a new plugin registry.
func NewRegistry(path string, options ...Option) (*FsRegistry, error) {
	r := &FsRegistry{
		fs:         afero.NewOsFs(),
		path:       filepath.Clean(path),
		configFile: filepath.Join(path, "config.yaml"),
	}

	for _, o := range options {
		o(r)
	}

	if r.config == nil {
		c, err := config.NewMemConfigurator(
			config.NewFileConfigurator(r.configFile).
				WithFs(r.fs),
		)
		if err != nil {
			return nil, err
		}

		r.config = c
	}

	return r, nil
}

// WithFs sets filesystem.
func WithFs(fs afero.Fs) Option {
	return func(r *FsRegistry) {
		r.fs = fs
	}
}

// WithConfigurator sets configurator.
func WithConfigurator(c config.Configurator) Option {
	return func(r *FsRegistry) {
		r.config = c
	}
}

// WithConfigFile sets config file location.
func WithConfigFile(configFile string) Option {
	return func(r *FsRegistry) {
		r.configFile = filepath.Clean(configFile)
	}
}
