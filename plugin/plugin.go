package plugin

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bool64/ctxd"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const (
	// MetadataFile is the plugin metadata file.
	MetadataFile = ".plugin.registry.yaml"

	defaultFile = "${name}-${version}-${os}-${arch}.tar.gz"
)

// ErrPluginNotExist indicates that the plugin does not exist.
var ErrPluginNotExist = errors.New("plugin does not exist")

// Plugins is a map of plugins.
type Plugins map[string]Plugin

// FilterByTag filter the plugins by tags.
func (p Plugins) FilterByTag(tag string) Plugins {
	result := make(Plugins, len(p))

	for k, v := range p {
		if v.Tags.Contains(tag) {
			result[k] = v
		}
	}

	return result
}

// Plugin represents metadata of a plugin.
type Plugin struct {
	Name        string    `yaml:"name"`
	URL         string    `yaml:"url"`
	Version     string    `yaml:"version"`
	Description string    `yaml:"description"`
	Enabled     bool      `yaml:"enabled"`
	Hidden      bool      `yaml:"hidden"`
	Artifacts   Artifacts `yaml:"artifacts"`
	Tags        Tags      `yaml:"tags"`
}

// RuntimeArtifact returns the artifact of current arch.
func (p *Plugin) RuntimeArtifact() Artifact {
	if a, ok := p.Artifacts[RuntimeArtifactIdentifier()]; ok {
		return a
	}

	if a, ok := p.Artifacts[RuntimeArtifactIdentifierWithoutArch()]; ok {
		return a
	}

	return Artifact{File: defaultFile}
}

// ResolveArtifact replaces all placeholders in artifact definition by real values.
func (p *Plugin) ResolveArtifact(a Artifact) Artifact {
	r := strings.NewReplacer(
		"${name}", p.Name,
		"${name}", p.Name,
		"${version}", p.Version,
		"${os}", runtime.GOOS,
		"${arch}", runtime.GOARCH,
	)

	a.File = r.Replace(a.File)

	return a
}

// UnmarshalYAML satisfies yaml.Unmarshaler.
func (p *Plugin) UnmarshalYAML(value *yaml.Node) error {
	type rawPlugin Plugin

	raw := rawPlugin(defaultPluginConfig())

	if err := value.Decode(&raw); err != nil {
		return err
	}

	if !raw.Artifacts.Has(RuntimeArtifactIdentifier()) &&
		!raw.Artifacts.Has(RuntimeArtifactIdentifierWithoutArch()) {
		raw.Artifacts[RuntimeArtifactIdentifier()] = Artifact{File: defaultFile}
	}

	*p = Plugin(raw)

	return nil
}

// Tags is a list of string tag.
type Tags []string

// Contains checks whether a tag is in the list or not.
func (t Tags) Contains(tag string) bool {
	for _, v := range t {
		if v == tag {
			return true
		}
	}

	return false
}

// Artifacts is a map of Artifact, identified by os and arch.
type Artifacts map[ArtifactIdentifier]Artifact

// Has checks whether the artifact is in the list.
func (a *Artifacts) Has(id ArtifactIdentifier) bool {
	_, ok := (*a)[id]

	return ok
}

// ArtifactIdentifier represents information to identify an artifact.
type ArtifactIdentifier struct {
	OS   string
	Arch string
}

// String satisfies fmt.Stringer.
func (a ArtifactIdentifier) String() string {
	var s string

	if a.Arch == "" {
		s = a.OS
	} else {
		s = fmt.Sprintf("%s/%s", a.OS, a.Arch)
	}

	return s
}

// MarshalYAML satisfies yaml.Marshaler.
func (a ArtifactIdentifier) MarshalYAML() (interface{}, error) { // nolint: unparam
	return a.String(), nil
}

// UnmarshalYAML satisfies yaml.Unmarshaler.
func (a *ArtifactIdentifier) UnmarshalYAML(value *yaml.Node) error {
	var raw string

	if err := value.Decode(&raw); err != nil {
		return err
	}

	parts := strings.SplitN(raw, "/", 2)

	*a = ArtifactIdentifier{
		OS:   parts[0],
		Arch: "",
	}

	if len(parts) > 1 {
		a.Arch = parts[1]
	}

	return nil
}

// Artifact represents all information about an artifact of a plugin.
type Artifact struct {
	File string `yaml:"file"`
}

func defaultPluginConfig() Plugin {
	return Plugin{
		Enabled:   true,
		Artifacts: map[ArtifactIdentifier]Artifact{},
	}
}

// RuntimeArtifactIdentifier returns the system's identifier.
func RuntimeArtifactIdentifier() ArtifactIdentifier {
	return ArtifactIdentifier{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

// RuntimeArtifactIdentifierWithoutArch returns the system's identifier without arch.
func RuntimeArtifactIdentifierWithoutArch() ArtifactIdentifier {
	return ArtifactIdentifier{
		OS:   runtime.GOOS,
		Arch: "",
	}
}

// NewArtifactIdentifier creates a new ArtifactIdentifier.
func NewArtifactIdentifier(os, arch string) ArtifactIdentifier {
	return ArtifactIdentifier{
		OS:   os,
		Arch: arch,
	}
}

// Load loads plugin metadata.
func Load(fs afero.Fs, path string) (*Plugin, error) {
	r, err := fs.Open(filepath.Join(path, MetadataFile))
	if err != nil {
		return nil, loadError(err, path)
	}

	var p Plugin

	dec := yaml.NewDecoder(r)

	if err := dec.Decode(&p); err != nil {
		return nil, loadError(err, path)
	}

	return &p, nil
}

func loadError(err error, path string) error {
	return ctxd.WrapError(context.Background(), err, "could not read metadata", "path", path)
}
