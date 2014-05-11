# CoreOS iPXE Server

[![Build Status](https://drone.io/github.com/kelseyhightower/coreos-ipxe-server/status.png)](https://drone.io/github.com/kelseyhightower/coreos-ipxe-server/latest)

The CoreOS iPXE Server attempts to automate as much of the [Booting CoreOS via iPXE](https://coreos.com/docs/running-coreos/bare-metal/booting-with-ipxe/) process as possible, mainly generating iPXE boot scripts and serving CoreOS PXE boot images.

## Table of Contents

- [Installation](#installation)
- [Getting Started](docs/getting_started.md)
- [API](docs/api.md)

## Installation

### Binary Release

```
curl -L https://github.com/kelseyhightower/coreos-ipxe-server/releases/download/v0.3.0/coreos-ipxe-server-0.3.0-darwin-amd64 -o coreos-ipxe-server
chmod +x coreos-ipxe-server
```

### Source

#### Clone

```
mkdir -p ${GOPATH}/src/github.com/kelseyhightower
cd ${GOPATH}/src/github.com/kelseyhightower
git clone git@github.com:kelseyhightower/coreos-ipxe-server.git
```

#### Build

```
cd ${GOPATH}/src/github.com/kelseyhightower/coreos-ipxe-server
go build .
```
