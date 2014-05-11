# SSH Public Keys

SSH keys are configured via the sshkey boot parameter, which is part of the CoreOS iPXE boot script. SSH keys are identified by id and are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/sshkeys` directory.

Example:

```
$COREOS_IPXE_SERVER_DATA_DIR/sshkeys/coreos.pub
```
