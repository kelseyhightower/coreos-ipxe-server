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
  - [iPXE Profiles](docs/profiles.md)

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

## Configuration

### Environment Variables

#### Optional:

```
COREOS_IPXE_SERVER_DATA_DIR
COREOS_IPXE_SERVER_BASE_URL
COREOS_IPXE_SERVER_LISTEN_ADDR
```

#### Example:

```
export COREOS_IPXE_SERVER_DATA_DIR="/opt/coreos-ipxe-server"
export COREOS_IPXE_SERVER_BASE_URL="coreos.ipxe.example.com:4777"
export COREOS_IPXE_SERVER_LISTEN_ADDR="0.0.0.0:4777"
```

### Data Directory

Create the data directory which will hold the CoreOS images, SSH public keys, and cloud configs:

```
mkdir -p $COREOS_IPXE_SERVER_DATA_DIR/{configs,images,profiles,sshkeys}
```

#### Download CoreOS images

```
mkdir -p $COREOS_IPXE_SERVER_DATA_DIR/images/amd64-usr/310.1.0
cd $COREOS_IPXE_SERVER_DATA_DIR/images/amd64-usr/310.1.0
wget http://storage.core-os.net/coreos/amd64-usr/310.1.0/coreos_production_pxe_image.cpio.gz
wget http://storage.core-os.net/coreos/amd64-usr/310.1.0/coreos_production_pxe.vmlinuz
```

#### Add a SSH public key

```
$COREOS_IPXE_SERVER_DATA_DIR/sshkeys/coreos.pub
```

#### Add a cloud config file

```
$COREOS_IPXE_SERVER_DATA_DIR/configs/development.yml
```

#### Add a iPXE profile

```
$COREOS_IPXE_SERVER_DATA_DIR/profiles/development.json
```

#### Example

```
/opt/coreos-ipxe-server/
├── configs
│   └── development.yml
├── images
│   └── amd64-usr
│       └── 310.1.0
│           ├── coreos_production_pxe.vmlinuz
│           └── coreos_production_pxe_image.cpio.gz
├── profiles
│   └── development.json
└── sshkeys
    └── coreos.pub
```

### SSH Public Keys

SSH keys are configured via the sshkey boot parameter, which is part of the CoreOS iPXE boot script. SSH keys are identified by id and are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/sshkeys` directory. 

Example:

```
$COREOS_IPXE_SERVER_DATA_DIR/sshkeys/coreos.pub
```

### Cloud Configs

Cloud configs can be used to automate the configuration of your CoreOS install. The `cloud-config-url` is configured via the cloud-config-url boot parameter, which is part of the CoreOS iPXE boot script. Cloud configs are identified by id and are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/configs` directory.

Example:

```
$COREOS_IPXE_SERVER_DATA_DIR/configs/development.yml
```
