# API

## iPXE Boot Script

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

## Retrieve an SSH Public Key

Produce an SSH public key suitable for use by `coreos-cloudconfig` to insert into
a user's `authorized_keys`.

```
GET http://coreos.ipxe.example.com/keys
```

**Parameters**

Name | Type | Description
-----|------|------------
name | string | The name of the key (not including .pub)


#### Generate a key for a specified name

```
GET http://coreos.ipxe.example.com/keys?name=foo
```

**Response**

```
HTTP/1.1 200 OK
```

```
[{"key":"ssh-rsa AAAAB3N... foo@bar"}]
```
