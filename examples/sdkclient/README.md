# Simple Helm SDK `HelmClient` consumer

This is a very simple example application demonstrating usage of the storage-agnostic `HelmClient` interface.

## Downloading a chart from an HTTPS repository

```sh
go run main.go https://stefanprodan.github.io/podinfo podinfo
```

## Downloading a chart from an OCI registry

```sh
go run main.go oci://ghcr.io/stefanprodan/charts podinfo
```
