# CoreOS iPXE Server

The CoreOS iPXE Server attempts to automate as much of the [Booting CoreOS via iPXE](https://coreos.com/docs/running-coreos/bare-metal/booting-with-ipxe/) as possible. Currently the following features are supported:

 - Creating dynamic iPXE boot scripts for running CoreOS 
 - Serve CoreOS pxe boot images

- [**API**](#api)
  - [iPXE Boot Script](#ipxe-boot-script)
    - [Set the CoreOS version](#set-the-coreos-version)
    - [Use a state partition](#use-a-state-partition)
<p></p>
- [**Configuration**](#configuration)
  - [Environment Variables](#environment-variables)
  - [Data Directory](#data-directory)
  - [SSH Public Key](#ssh-public-key)

## API

### iPXE Boot Script

```
GET /
```

**Response**

```
HTTP/1.1 200 OK
```

```
set coreos-version latest
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: state=tmpfs: sshkey="ssh-rsa AAAAB3Nza..."
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
```

#### Set the CoreOS version

```
GET http://example.com?version=268.1.0
```

**Response**

```
HTTP/1.1 200 OK
```

```
set coreos-version 268.1.0
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: state=tmpfs: sshkey="ssh-rsa AAAAB3Nza..."
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
```

> Notice the change in the `set coreos-version` line.

#### Use a state partition

```
GET http://example.com?state=1
```

**Response**

```
HTTP/1.1 200 OK
```

```
set coreos-version latest
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: sshkey="ssh-rsa AAAAB3Nza..."
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
```

> Notice `state=tmpfs:` missing from the kernel boot parameters

## Configuration

### Environment Variables

#### Required:

```
COREOS_IPXE_SERVER_BASE_DIR
COREOS_IPXE_SERVER_BASE_URL
COREOS_IPXE_SERVER_LISTEN_PORT
```

#### Optional:

```
COREOS_IPXE_SERVER_LISTEN_HOST
```

#### Example:

```
export COREOS_IPXE_SERVER_BASE_DIR="/var/lib/cis/"
export COREOS_IPXE_SERVER_BASE_URL="10.0.1.10:8080"
export COREOS_IPXE_SERVER_LISTEN_PORT="8080"
```

### Data Directory

```
tree data/
data/
└── coreos
    ├── amd64-generic
    │   └── 268.1.0
    │       ├── coreos_production_pxe.vmlinuz
    │       └── coreos_production_pxe_image.cpio.gz
    ├── amd64-usr
    │   └── 268.1.0
    │       ├── coreos_production_pxe.vmlinuz
    │       └── coreos_production_pxe_image.cpio.gz
    └── coreos.pub
```

### SSH Public Key

The SSH public must exist under `$COREOS_IPXE_SERVER_BASE_DIR` as `coreos.pub`.
