# CoreOS iPXE Server

[![Build Status](https://drone.io/github.com/kelseyhightower/coreos-ipxe-server/status.png)](https://drone.io/github.com/kelseyhightower/coreos-ipxe-server/latest)

The CoreOS iPXE Server attempts to automate as much of the [Booting CoreOS via iPXE](https://coreos.com/docs/running-coreos/bare-metal/booting-with-ipxe/) process as possible, mainly generating iPXE boot scripts and serving CoreOS PXE boot images.

## Table of Contents

- [**API**](#api)
  - [iPXE Boot Script](#ipxe-boot-script)
<p></p>
- [**Configuration**](#configuration)
  - [Environment Variables](#environment-variables)
  - [Data Directory](#data-directory)
  - [SSH Public Keys](#ssh-public-keys)

## API

### iPXE Boot Script

Dynamically generate a CoreOS iPXE boot script.

```
GET http://coreos.ipxe.example.com
```

**Parameters**

Name | Type | Description 
-----|------|------------
sshkey | string | The ssh key id to use. The key must exist as `$COREOS_IPXE_SERVER_DATA_DIR/sshkeys/$sshkey.pub`. Default: `coreos`
state | boolean | If `true`, generate iPXE boot script without `state=tmpfs:` kernel parameter. Default: `false`
version | string | The CoreOS PXE image version to boot. Default: `latest`


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
COREOS_IPXE_SERVER_DATA_DIR
COREOS_IPXE_SERVER_BASE_URL
COREOS_IPXE_SERVER_LISTEN_PORT
```

#### Optional:

```
COREOS_IPXE_SERVER_LISTEN_HOST
```

#### Example:

```
export COREOS_IPXE_SERVER_DATA_DIR="/opt/coreos-ipxe-server"
export COREOS_IPXE_SERVER_BASE_URL="coreos.ipxe.example.com"
export COREOS_IPXE_SERVER_LISTEN_PORT="80"
```

### Data Directory

Create the data directory which will hold the CoreOS images and the SSH public key:

```
mkdir -p $COREOS_IPXE_SERVER_DATA_DIR/coreos/amd64-generic
```

#### Download CoreOS images

```
mkdir $COREOS_IPXE_SERVER_DATA_DIR/coreos/amd64-generic/268.1.0
cd $COREOS_IPXE_SERVER_DATA_DIR/coreos/amd64-generic/268.1.0
wget http://storage.core-os.net/coreos/amd64-generic/268.1.0/coreos_production_pxe.vmlinuz
wget http://storage.core-os.net/coreos/amd64-generic/268.1.0/coreos_production_pxe_image.cpio.gz
```

#### Create a symlink to the default version

By default CoreOS iPXE boot scripts will be generated with the CoreOS version set to `latest`. Add a symlink to ensure this works.

```
ln -s $COREOS_IPXE_SERVER_DATA_DIR/coreos/amd64-generic/268.1.0 $COREOS_IPXE_SERVER_DATA_DIR/coreos/amd64-generic/latest
```

#### Add a SSH public key

```
cp ~/.ssh/id_rsa.pub $COREOS_IPXE_SERVER_DATA_DIR/sshkeys/coreos.pub
```

#### Example

```
/opt/coreos-ipxe-server/
├── coreos
│   └── amd64-generic
│       ├── 268.1.0
│       │   ├── coreos_production_pxe.vmlinuz
│       │   └── coreos_production_pxe_image.cpio.gz
│       └── latest -> /opt/coreos-ipxe-server/coreos/amd64-generic/268.1.0
└── sshkeys
    └── coreos.pub
```

### SSH Public Keys

SSH public keys are required to log into your CoreOS system. SSH keys are configured via the sshkey boot parameter, which is part of the CoreOS iPXE boot script. SSH keys are identified by id and are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/sshkeys` directory. 

Example:

```
$COREOS_IPXE_SERVER_DATA_DIR/sshkeys/$sshkeyid.pub
```
