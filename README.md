# driver-lima

Kutti driver for Lima

[![Go Report Card](https://goreportcard.com/badge/github.com/kuttiproject/driver-lima)](https://goreportcard.com/report/github.com/kuttiproject/driver-lima)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/kuttiproject/driver-lima)](https://pkg.go.dev/github.com/kuttiproject/driver-lima)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/kuttiproject/driver-lima?include_prereleases)

## Images

This driver depends on qcow2 images published via the [kuttiproject/driver-lima-images](https://github.com/kuttiproject/driver-lima-images) repository. The details of the driver-to-VM interface are documented there.

The releases of that repository are the default source for this driver. The list of available/deprecated images and the images themselves are published there. The releases of that repository follow the major and minor versions of this repository, but sometimes may lag by one version. The `ImagesVersion` constant specifies the version of the images repository that is used by a particular version of this driver.

## Apple Silicon Mac Only

This driver only works on macOS running on Apple silicon.
