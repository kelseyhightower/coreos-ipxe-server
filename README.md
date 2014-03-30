# CoreOS iPXE Server

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

### Example:

```
export COREOS_IPXE_SERVER_BASE_DIR="/var/lib/cis/"
export COREOS_IPXE_SERVER_BASE_URL="10.0.1.10:8080"
export COREOS_IPXE_SERVER_LISTEN_PORT="8080"
```

## Data Directory

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

## SSH Public Key

The SSH public must exist under `$COREOS_IPXE_SERVER_BASE_DIR` as `coreos.pub`.
