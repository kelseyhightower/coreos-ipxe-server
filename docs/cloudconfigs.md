### Cloud Configs

Cloud configs can be used to automate the configuration of your CoreOS install. The `cloud-config-url` is configured via the cloud-config-url boot parameter, which is part of the CoreOS iPXE boot script. Cloud configs are identified by id and are stored under the `$COREOS_IPXE_SERVER_DATA_DIR/configs` directory.

Example:

```
$COREOS_IPXE_SERVER_DATA_DIR/configs/development.yml
```
