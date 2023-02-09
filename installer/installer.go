package installer

import (
	"context"
	"errors"
	"sync"

	"github.com/spf13/afero"

	fsCtx "github.com/nhatthm/plugin-registry/context"
	"github.com/nhatthm/plugin-registry/plugin"
)

var (
	// ErrUnknownInstaller indicates that the installer is not registered.
	ErrUnknownInstaller = errors.New("unknown installer")
	// ErrNoInstaller indicates that the plugin has no supported installer.
	ErrNoInstaller = errors.New("no supported installer")
)

var (
	installersMu sync.Mutex

	installers = map[string]metadata{}
)

type metadata struct {
	validate  Validity
	construct Constructor
}

// Validity checks whether a plugin is supported by the plugin or not.
type Validity func(ctx context.Context, source string) bool

// Constructor is to construct a new installer.
type Constructor func(fs afero.Fs) Installer

// Installer installs a plugin from its URL.
type Installer interface {
	Install(ctx context.Context, dest, src string) (*plugin.Plugin, error)
}

// CallbackInstaller is a callback installer.
type CallbackInstaller func(ctx context.Context, dest, src string) (*plugin.Plugin, error)

// Install installs the plugin.
func (f CallbackInstaller) Install(ctx context.Context, dest, src string) (*plugin.Plugin, error) {
	return f(ctx, dest, src)
}

// Register registers a plugin installer.
func Register(name string, validity Validity, constructor Constructor) {
	installersMu.Lock()
	defer installersMu.Unlock()

	installers[name] = metadata{
		validate:  validity,
		construct: constructor,
	}
}

// New creates a new installer.
func New(ctx context.Context, name string) (Installer, error) {
	installersMu.Lock()
	defer installersMu.Unlock()

	m, ok := installers[name]
	if !ok {
		return nil, ErrUnknownInstaller
	}

	return m.construct(fsCtx.Fs(ctx)), nil
}

// Find finds an installer for the given plugin url.
func Find(ctx context.Context, src string) (Installer, error) {
	installersMu.Lock()
	defer installersMu.Unlock()

	for _, m := range installers {
		if m.validate(ctx, src) {
			return m.construct(fsCtx.Fs(ctx)), nil
		}
	}

	return nil, ErrNoInstaller
}
