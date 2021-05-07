package registry

import (
	"context"

	fsCtx "github.com/nhatthm/plugin-registry/context"
	"github.com/nhatthm/plugin-registry/installer"
	"github.com/nhatthm/plugin-registry/plugin"
)

func (r *FsRegistry) makeInstaller(ctx context.Context, src string) (installer.Installer, error) {
	i, err := installer.Find(fsCtx.WithFs(ctx, r.fs), src)
	if err != nil {
		return nil, err
	}

	return r.installer(i), nil
}

func (r *FsRegistry) installer(i installer.Installer) installer.CallbackInstaller {
	return func(ctx context.Context, dest, src string) (*plugin.Plugin, error) {
		p, err := i.Install(ctx, dest, src)
		if err != nil {
			return nil, err
		}

		oldPlugin, err := r.GetPlugin(p.Name)
		if err != nil {
			return nil, err
		}

		// Do not accidentally enable the disabled plugin.
		if oldPlugin != nil {
			p.Enabled = oldPlugin.Enabled
		}

		return p, r.config.SetPlugin(*p)
	}
}

// Install installs plugin from a url.
func (r *FsRegistry) Install(ctx context.Context, src string) error {
	i, err := r.makeInstaller(ctx, src)
	if err != nil {
		return err
	}

	_, err = i.Install(ctx, r.path, src)

	return err
}
