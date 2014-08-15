# Example of running coreos-ipxe-server on docker

### The script will
- Prepare directory structure
- download images for versions specified in VERSIONS
- create default config files
- build ipxe docker image
- run ipxe docker container

### Following has to be done before running the script
- update following env vars so that they suit your needs
 - ```
COREOS_IPXE_SERVER_BASE_URL ip/dns on which your docker will listen
COREOS_IPXE_SERVER_LISTEN_ADDR="0.0.0.0:4777"
VERSIONS=("367.1.0" "379.3.0")
```
- include your ssh-key
- run script
 - cd docker && ./setup_docker.sh

If everything runs sucessfully you should be able to see you container running among others using ``docker ps``

You can try to access iPXE config on ``<COREOS_IPXE_SERVER_BASE_URL>:4777/?profile=development_<VERSION>``
In my case it will be ``http://10.102.11.42:4777/?profile=development_367.1.0``
