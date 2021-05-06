# Plugin Registry for Golang

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/plugin-registry)](https://github.com/nhatthm/plugin-registry/releases/latest)
[![Build Status](https://github.com/nhatthm/plugin-registry/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/plugin-registry/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/plugin-registry/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/plugin-registry)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/plugin-registry)](https://goreportcard.com/report/github.com/nhatthm/plugin-registry)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/plugin-registry)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

Install and manage plugins for Golang application.

## Prerequisites

- `Go >= 1.15`

## Install

```bash
go get github.com/nhatthm/plugin-registry
```

## Usage

`plugin-registry` helps to install plugins to a container at your choice. There are 4 tasks to manage them:

- Install
- Uninstall
- Enable
- Disable

`plugin-registry` is backed by [spf13/afero](https://github.com/spf13/afero) so feel free to use it with your favorite
backend file system by using `WithFs(fs afero.Fs)` option. For example

```go
package mypackage

import (
	registry "github.com/nhatthm/plugin-registry"
	_ "github.com/nhatthm/plugin-registry-github" // Add github installer.
	"github.com/spf13/afero"
)

func createRegistry() (registry.Registry, error) {
	return registry.NewRegistry("~/plugins", registry.WithFs(afero.NewMemMapFs()))
}
```

By default, `plugin-registry` will record the installed plugin in `config.yaml` file in the given container, you can
change it to a new place of your choice by using `WithConfigFile(path string)` option, for example:

```go
package mypackage

import (
	registry "github.com/nhatthm/plugin-registry"
	_ "github.com/nhatthm/plugin-registry-github" // Add github installer.
)

func createRegistry() (registry.Registry, error) {
	return registry.NewRegistry("/usr/local/bin/plugins", registry.WithConfigFile("~/plugins/config.yaml"))
}
```

If you want to manage the plugins differently, you can write your own `Configurator` and use `WithConfigurator()` option
to set it, for example:

```go
package mypackage

import (
	registry "github.com/nhatthm/plugin-registry"
	_ "github.com/nhatthm/plugin-registry-github" // Add github installer.
	"github.com/nhatthm/plugin-registry/config"
)

var _ config.Configurator = (*MyConfigurator)(nil)

type MyConfigurator struct{}

func createConfigurator() *MyConfigurator {
	var c MyConfigurator

	// init c.

	return &c
}

func createRegistry() (registry.Registry, error) {
	return registry.NewRegistry("/usr/local/bin/plugins", registry.WithConfigurator(createConfigurator()))
}
```

## Installer

There is no installer provided by this library, you need to install and import it in your project.

Known 3rd party installers:

- https://github.com/nhatthm/plugin-registry-fs: Support binary, folder, `.tar.gz`, `.gz`, `zip` plugin.
- https://github.com/nhatthm/plugin-registry-github: Support plugin from github,

## Examples

```go
package mypackage

import (
	"context"

	registry "github.com/nhatthm/plugin-registry"
	_ "github.com/nhatthm/plugin-registry-github" // Add github installer.
)

var defaultRegistry = mustCreateRegistry()

func mustCreateRegistry() registry.Registry {
	r, err := createRegistry()
	if err != nil {
		panic(err)
	}

	return r
}

func createRegistry() (registry.Registry, error) {
	return registry.NewRegistry("~/plugins")
}

func installPlugin(source string) error {
	return defaultRegistry.Install(context.Background(), source)
}

```

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
