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
cloudconfig | string | The cloud config id to use. The key must exist as `$COREOS_IPXE_SERVER_DATA_DIR/configs/$cloudconfig.yml`
sshkey | string | The ssh key id to use. The key must exist as `$COREOS_IPXE_SERVER_DATA_DIR/sshkeys/$sshkey.pub`
version | string | The CoreOS PXE image version to boot. Default: `latest`


**Response**

```
HTTP/1.1 200 OK
```

```
set coreos-version latest
set base-url http://coreos.ipxe.example.com/coreos/amd64-usr/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz sshkey="ssh-rsa AAAAB3Nza..."
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
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

Create the data directory which will hold the CoreOS images, SSH public keys, and cloud configs:

```
mkdir $COREOS_IPXE_SERVER_DATA_DIR/images
mkdir $COREOS_IPXE_SERVER_DATA_DIR/sshkeys
mkdir $COREOS_IPXE_SERVER_DATA_DIR/configs
```

#### Download CoreOS images

```
mkdir $COREOS_IPXE_SERVER_DATA_DIR/images/amd64-usr/298.0.0
cd $COREOS_IPXE_SERVER_DATA_DIR/images/amd64-usr/298.0.0
wget http://storage.core-os.net/coreos/amd64-usr/298.0.0/coreos_production_pxe_image.cpio.gz
wget http://storage.core-os.net/coreos/amd64-usr/298.0.0/coreos_production_pxe.vmlinuz
```

#### Create a symlink to the default version

By default CoreOS iPXE boot scripts will be generated with the CoreOS version set to `latest`. Add a symlink to ensure this works.

```
ln -s $COREOS_IPXE_SERVER_DATA_DIR/images/amd64-usr/298.0.0 $COREOS_IPXE_SERVER_DATA_DIR/images/amd64-usr/latest
```

#### Add a SSH public key

```
cp ~/.ssh/id_rsa.pub $COREOS_IPXE_SERVER_DATA_DIR/sshkeys/coreos.pub
```

#### Add a cloud config file

cp cloud-config.yml $COREOS_IPXE_SERVER_DATA_DIR/configs/cloud-config.yml

#### Example

```
/opt/coreos-ipxe-server/
├── configs
│   └── cloud-config.yml
├── images
│   └── amd64-usr
│       ├── coreos_production_pxe.vmlinuz
│       ├── coreos_production_pxe_image.cpio.gz
│       └── latest -> /opt/coreos-ipxe-server/images/amd64-usr/298.0.0
└── sshkeys
    └── coreos.pub
```

### SSH Public Keys

SSH public keys are required to log into your CoreOS system. SSH keys are configured via the sshkey boot parameter, which is part of the CoreOS iPXE boot script. SSH keys are identified by id and are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/sshkeys` directory. 

Example:

```
$COREOS_IPXE_SERVER_DATA_DIR/sshkeys/coreos.pub
```

### Cloud Configs

Cloud configs can be used to automate the configuration of your CoreOS install. The `cloud-config-url` is configured via the cloud-config-url boot parameter, which is part of the CoreOS iPXE boot script. Cloud configs are identified by id and are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/cloud-configs` directory.

Example:

```
$COREOS_IPXE_SERVER_DATA_DIR/configs/cloud-config.yml
```
