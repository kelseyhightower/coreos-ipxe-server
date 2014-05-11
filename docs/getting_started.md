# Getting Started

- [Installation](#installation)
- [Configuration](#configuration)
  - [Environment Variables](#environment-variables)
- [Create the Data Directory](#create-the-data-directory) 
- [Download the CoreOS PXE Images](#download-the-coreos-pxe-images)
- [Add a Cloud Config File](#add-a-cloud-config-file)
- [Add a SSH Public Key](#add-an-ssh-public-key)
- [Add an iPXE Profile](#add-an-ipxe-profile)
- [Example Data Directory Layout](#example-data-directory-layout)

## Installation

```
curl -L https://github.com/kelseyhightower/coreos-ipxe-server/releases/download/v0.3.0/coreos-ipxe-server-0.3.0-darwin-amd64 -o coreos-ipxe-server
chmod +x coreos-ipxe-server
```

## Configuration

All configuration is handled via environment variables with sane defaults. See [Configuration](configuration.md) for more details.


## Create the Data Directory

The data directory is where the CoreOS images, SSH public keys, cloud configs and iPXE profiles are stored. The data directory defaults to `/opt/coreos-ipxe-server`; set it to a different directory via the `COREOS_IPXE_SERVER_DATA_DIR` environment variable:

```
export COREOS_IPXE_SERVER_DATA_DIR=/var/lib/coreos-ipxe-server
```

Next create the subdirectories that will hold the CoreOS iPXE assets: 

```
mkdir -p $COREOS_IPXE_SERVER_DATA_DIR/{configs,images,profiles,sshkeys}
```

## Download the CoreOS PXE Images

The CoreOS PXE images are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/images` directory.

```
mkdir -p $COREOS_IPXE_SERVER_DATA_DIR/images/amd64-usr/310.1.0
cd $COREOS_IPXE_SERVER_DATA_DIR/images/amd64-usr/310.1.0
wget http://storage.core-os.net/coreos/amd64-usr/310.1.0/coreos_production_pxe_image.cpio.gz
wget http://storage.core-os.net/coreos/amd64-usr/310.1.0/coreos_production_pxe.vmlinuz
```

## Add an SSH Public Key

SSH public keys are used to login to your CoreOS instance. SSH public keys are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/sshkeys` directory.

Edit `$COREOS_IPXE_SERVER_DATA_DIR/sshkeys/coreos.pub`

```
ssh-rsa AAAAB3Nza...
```

## Add a Cloud Config File

Cloud config files are used to automated the setup of your CoreOS instance. See [Customize with Cloud Config](https://coreos.com/docs/cluster-management/setup/cloudinit-cloud-config/) for more details. Cloud config files are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/configs` directory.

Edit `$COREOS_IPXE_SERVER_DATA_DIR/configs/development.yml`

```
#cloud-config

ssh_authorized _keys:
    - ssh-rsa AAAAB3Nza...
coreos:
  etcd:
    addr: $private_ipv4:4001
    peer-addr: $private_ipv4:7001
  units:
    - name: etcd.service
      command: start
    - name: fleet.service
      command: start
    - name: docker.socket
      command: start
  oem:
    id: coreos
    name: CoreOS Custom
    version-id: 310.1.0
    home-url: https://coreos.com
```

## Add an iPXE Profile

iPXE profiles are used to define CoreOS boot parameters. iPXE profiles are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/profiles` directory.

Edit `$COREOS_IPXE_SERVER_DATA_DIR/profiles/development.json` 

```
{
	"cloud_config": "development",
	"rootfstype": "btrfs",
	"version": "310.1.0"
}
```

## Example Data Directory Layout

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
