# Profiles

iPXE profiles are used to define CoreOS kernel options used during the PXE boot process. Profiles are identified by id and are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/profiles` directory.

## File Format

The iPXE profile file uses the JSON file format.

A iPXE profile file should contain an associative array which has zero or more of the following keys:

* cloud_config
* console
* coreos_autologin
* rootfstype
* sshkey
* version

See [Configuring pxelinux](https://coreos.com/docs/running-coreos/bare-metal/booting-with-pxe/#configuring-pxelinux) for more details.

### Example Profile

```
$COREOS_IPXE_SERVER_DATA_DIR/profiles/development.json
```

```
{
  "cloud_config": "development",
  "console": ["tty0", "tty1"],
  "coreos_autologin": "tty1",
  "rootfstype": "btrfs",
  "sshkey": "coreos",
  "version": "310.1.0"
}
```
