package registry

// Enable enabled a plugin by name.
func (r *FsRegistry) Enable(name string) error {
	return r.config.EnablePlugin(name)
}
