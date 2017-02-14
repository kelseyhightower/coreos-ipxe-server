FROM alpine:latest

WORKDIR /

RUN mkdir -p /opt/coreos-ipxe-server/images/amd64-usr/

ADD coreos/amd64-usr /opt/coreos-ipxe-server/images/amd64-usr/

RUN ln -s /opt/coreos-ipxe-server/images/amd64-usr/1235.9.0/ /opt/coreos-ipxe-server/images/amd64-usr/stable

ADD coreos-ipxe-server /usr/local/bin/coreos-ipxe-server

CMD /usr/local/bin/coreos-ipxe-server
