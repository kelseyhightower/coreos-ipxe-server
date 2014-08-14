FROM google/golang

RUN mkdir -p /gopath/src/github.com/kelseyhightower

WORKDIR /gopath/src/github.com/kelseyhightower
RUN git clone https://github.com/kelseyhightower/coreos-ipxe-server.git

WORKDIR /gopath/src/github.com/kelseyhightower/coreos-ipxe-server
RUN go install

RUN mkdir -p /opt/coreos-ipxe-server
ADD configs /opt/coreos-ipxe-server/configs
ADD images /opt/coreos-ipxe-server/images
ADD profiles /opt/coreos-ipxe-server/profiles
ADD sshkeys /opt/coreos-ipxe-server/sshkeys

ENV COREOS_IPXE_SERVER_DATA_DIR /opt/coreos-ipxe-server
# URL has to be substituted with correct value. Can be overwritten during docker run
ENV COREOS_IPXE_SERVER_BASE_URL coreos.ipxe.example.com:4777
ENV COREOS_IPXE_SERVER_LISTEN_ADDR 0.0.0.0:4777

CMD /gopath/bin/coreos-ipxe-server
