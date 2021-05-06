package registry

// Disable disables a plugin by name.
func (r *FsRegistry) Disable(name string) error {
	return r.config.DisablePlugin(name)
}
