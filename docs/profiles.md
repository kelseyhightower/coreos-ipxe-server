# Profiles

iPXE profiles are used to define CoreOS iPXE boot parameters. Profiles are identified by id and are stored under the $COREOS_IPXE_SERVER_DATA_DIR/profiles directory.

## File Format

The iPXE profile file uses the JSON file format.

A iPXE profile file should contain an associative array which has zero or more of the following keys:

* cloud_config
* console
* coreos_autologin
* rootfstype
* sshkey
* version

Example:

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
