# Getting Started

- [Installation](#installation)
- [Configuration](#configuration)
  - [Environment Variables](#environment-variables)
- [Data Directory](#data-directory) 
  - [Cloud Config Files](#cloud-config-files)
  - [SSH Keys](#sshkeys)
  - [iPXE Profiles](#ipxe-profiles)

## Installation

```
curl -L https://github.com/kelseyhightower/coreos-ipxe-server/releases/download/v0.3.0/coreos-ipxe-server-0.3.0-darwin-amd64 -o coreos-ipxe-server
chmod +x coreos-ipxe-server
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

#### SSH Keys

```
$COREOS_IPXE_SERVER_DATA_DIR/sshkeys/coreos.pub
```

#### Cloud Config Files

```
$COREOS_IPXE_SERVER_DATA_DIR/configs/development.yml
```

#### iPXE Profiles

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
