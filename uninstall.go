package registry

import (
	"path/filepath"
)

// Uninstall uninstalls a plugin.
func (r *FsRegistry) Uninstall(name string) error {
	if err := r.fs.RemoveAll(filepath.Join(r.path, name)); err != nil {
		return err
	}

	return r.config.RemovePlugin(name)
}
