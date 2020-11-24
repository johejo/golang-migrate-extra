# golang-migrate-extra

[![ci](https://github.com/johejo/golang-migrate-extra/workflows/ci/badge.svg?branch=main)](https://github.com/johejo/golang-migrate-extra/actions?query=workflow%3Aci)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/johejo/golang-migrate-extra)](https://pkg.go.dev/github.com/johejo/golang-migrate-extra)
[![codecov](https://codecov.io/gh/johejo/golang-migrate-extra/branch/main/graph/badge.svg)](https://codecov.io/gh/johejo/golang-migrate-extra)
[![Go Report Card](https://goreportcard.com/badge/github.com/johejo/golang-migrate-extra)](https://goreportcard.com/report/github.com/johejo/golang-migrate-extra)

## Description

Extra source.Driver for [golang-migrate/migrate](https://github.com/golang-migrate/migrate).

This module is based on [golang-migrate/migrate#472](https://github.com/golang-migrate/migrate/pull/472).

- `source/iofs` source.Driver for [io/fs#FS](https://tip.golang.org/pkg/io/fs/).
- `source/file` source.Driver for local file system using `source/iofs`.

## Requirements

- Go 1.16 or higher.

## License

Inherit the license of `golang-migrate/migrate`.

https://github.com/golang-migrate/migrate/blob/master/LICENSE
