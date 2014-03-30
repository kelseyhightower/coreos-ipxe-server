# CoreOS iPXE Server

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
