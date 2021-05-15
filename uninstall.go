package registry

import (
	"path/filepath"
)

// Uninstall uninstalls a plugin.
func (r *FsRegistry) Uninstall(name string) error {
	if err := r.config.RemovePlugin(name); err != nil {
		return err
	}

	return r.fs.RemoveAll(filepath.Join(r.path, name))
}
