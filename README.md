# CoreOS iPXE Server

The CoreOS iPXE Server attempts to automate as much of the [Booting CoreOS via iPXE](https://coreos.com/docs/running-coreos/bare-metal/booting-with-ipxe/) process as possible, mainly generating iPXE boot scripts and serving CoreOS PXE boot images.

## Table of Contents

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

Dynamically generate a CoreOS iPXE boot script.

**Parameters**

Name | Type | Description 
-----|------|------------
state | boolean | If `true`, generate iPXE boot script without `state=tmpfs:` kernel parameter. Default: `false`
version | string | The CoreOS PXE image version to boot. Default: `latest`

```
GET http://coreos.ipxe.example.com
```

**Response**

```
HTTP/1.1 200 OK
```

```
set coreos-version latest
set base-url http://coreos.ipxe.example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: state=tmpfs: sshkey="ssh-rsa AAAAB3Nza..."
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
```

#### Set the CoreOS version to 268.1.0 and use a state partition

```
GET http://coreos.ipxe.example.com?version=268.1.0&state=true
```

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
export COREOS_IPXE_SERVER_BASE_DIR="/opt/coreos-ipxe-server"
export COREOS_IPXE_SERVER_BASE_URL="10.0.1.10:8080"
export COREOS_IPXE_SERVER_LISTEN_PORT="8080"
```

### Data Directory

Create the data directory which will hold the CoreOS images and the SSH public key:

```
mkdir -p $COREOS_IPXE_SERVER_BASE_DIR/coreos/amd64-generic
```

#### Download CoreOS images

```
mkdir $COREOS_IPXE_SERVER_BASE_DIR/coreos/amd64-generic/268.1.0
cd $COREOS_IPXE_SERVER_BASE_DIR/coreos/amd64-generic/268.1.0
wget http://storage.core-os.net/coreos/amd64-generic/268.1.0/coreos_production_pxe.vmlinuz
wget http://storage.core-os.net/coreos/amd64-generic/268.1.0/coreos_production_pxe_image.cpio.gz
```

#### Create a symlink to the default version

By default CoreOS iPXE boot scripts will be generated with the CoreOS version set to `latest`. Add a symlink to ensure this works.

```
ln -s $COREOS_IPXE_SERVER_BASE_DIR/coreos/amd64-generic/268.1.0 $COREOS_IPXE_SERVER_BASE_DIR/coreos/amd64-generic/latest
```

#### Add a SSH public key

```
cp ~/.ssh/id_rsa.pub $COREOS_IPXE_SERVER_BASE_DIR/coreos/coreos.pub
```

#### Example

```
/opt/coreos-ipxe-server
└── coreos
    ├── amd64-generic
    │   ├── 268.1.0
    │   │   ├── coreos_production_pxe.vmlinuz
    │   │   └── coreos_production_pxe_image.cpio.gz
    │   └── latest -> /opt/coreos-ipxe-server/coreos/amd64-generic/268.1.0
    └── coreos.pub
```

### SSH Public Key

The SSH public key must be in place before starting the CoreOS iPXE server. The SSH public will be used when generating the CoreOS iPXE boot scripts.

> The SSH public must exist under `$COREOS_IPXE_SERVER_BASE_DIR` as `coreos.pub`.
