package main

import (
	"log"
	"net"
	"os"
)

var (
	dataDir  string
	baseUrl  string
	hostPort string
)

// SetDataDir sets the data directory.
func SetDataDir(dir string) {
	dataDir = dir
}

// SetBaseUrl sets the base url.
func SetBaseUrl(url string) {
	baseUrl = url
}

func config() {
	// Set the data directory where the coreos directory containing
	// the ssh public key, kernal and boot images.
	dataDir = os.Getenv("COREOS_IPXE_SERVER_DATA_DIR")
	if dataDir == "" {
		log.Fatal("COREOS_IPXE_SERVER_DATA_DIR must be set and non-empty")
	}

	// Set the base URL used by the iPXE boot script.
	baseUrl = os.Getenv("COREOS_IPXE_SERVER_BASE_URL")
	if baseUrl == "" {
		log.Fatal("COREOS_IPXE_SERVER_BASE_URL must be set and non-empty")
	}

	// Set the host:port to listen for HTTP requests.
	listenHost := os.Getenv("COREOS_IPXE_SERVER_LISTEN_HOST")
	if listenHost == "" {
		listenHost = "0.0.0.0"
	}
	listenPort := os.Getenv("COREOS_IPXE_SERVER_LISTEN_PORT")
	if listenPort == "" {
		log.Fatal("COREOS_IPXE_SERVER_LISTEN_PORT must be set and non-empty")
	}
	hostPort = net.JoinHostPort(listenHost, listenPort)
}
