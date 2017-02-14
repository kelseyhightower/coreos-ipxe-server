package config

import (
	"os"
)

var (
	BaseUrl        string
	DataDir        string
	ListenAddr     string
	DefaultProfile string
)

var defaultDataDir = "/opt/coreos-ipxe-server"
var defaultListenAddr = "0.0.0.0:4777"

func init() {
	BaseUrl = os.Getenv("COREOS_IPXE_SERVER_BASE_URL")
	// Set the data directory where the coreos directory containing
	// the ssh public key, kernal and boot images.
	DataDir = os.Getenv("COREOS_IPXE_SERVER_DATA_DIR")
	if DataDir == "" {
		DataDir = defaultDataDir
	}
	ListenAddr = os.Getenv("COREOS_IPXE_SERVER_LISTEN_ADDR")
	if ListenAddr == "" {
		ListenAddr = defaultListenAddr
	}
	DefaultProfile = os.Getenv("COREOS_IPXE_SERVER_DEFAULT_PROFILE")
}
