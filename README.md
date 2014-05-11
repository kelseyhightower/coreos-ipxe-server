# CoreOS iPXE Server

[![Build Status](https://drone.io/github.com/kelseyhightower/coreos-ipxe-server/status.png)](https://drone.io/github.com/kelseyhightower/coreos-ipxe-server/latest)

The CoreOS iPXE Server attempts to automate as much of the [Booting CoreOS via iPXE](https://coreos.com/docs/running-coreos/bare-metal/booting-with-ipxe/) process as possible, mainly generating iPXE boot scripts and serving CoreOS PXE boot images.

## Table of Contents

- [**API**](#api)
  - [iPXE Boot Script](#ipxe-boot-script)
<p></p>
- [**Getting Started**](docs/getting_started.md)

## API

### iPXE Boot Script

Dynamically generate a CoreOS iPXE boot script.

```
GET http://coreos.ipxe.example.com
```

**Parameters**

Name | Type | Description 
-----|------|------------
profile | string | The CoreOS iPXE profile to use.


#### Generate Boot Script with a specific profile

```
GET http://coreos.ipxe.example.com?profile=development
```

**Response**

```
HTTP/1.1 200 OK
```

```
set coreos-version 310.1.0
set base-url http://coreos.ipxe.example.com/images/amd64-usr/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz console=tty0 rootfstype=btrfs cloud-config-url=http://coreos.ipxe.example.com/configs/development.yml
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
```
