package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var (
	baseUrl string
	sshKey  string
)

const ipxeBootScript = `#!ipxe
set coreos-version {{.Version}}
set base-url http://{{.BaseUrl}}/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: {{if not .State}}state=tmpfs: {{end}}sshkey="{{.SSHKey}}"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

func ipxeBootScriptServer(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	version := v.Get("version")
	if version == "" {
		version = "latest"
	}
	state := v.Get("state")

	t, err := template.New("ipxebootscript").Parse(ipxeBootScript)
	if err != nil {
		http.Error(w, "cannot generate ipxe boot script", 500)
		return
	}

	data := map[string]string{
		"BaseUrl": baseUrl,
		"SSHKey":  sshKey,
		"State":   state,
		"Version": version,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "cannot generate ipxe boot script", 500)
		return
	}
	return
}

func sshKeyFromFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(b)), nil
}

func main() {
	var err error
	
	// Set the base directory where the coreos directory containing
	// the ssh public key, kernal and boot images.
	baseDir := os.Getenv("COREOS_IPXE_SERVER_DATA_DIR")
	if baseDir == "" {
		log.Fatal("COREOS_IPXE_SERVER_DATA_DIR must be set and non-empty")
	}

	// Set the base URL used by the iPXE boot script.
	baseUrl = os.Getenv("COREOS_IPXE_SERVER_BASE_URL")
	if baseUrl == "" {
		log.Fatal("COREOS_IPXE_SERVER_BASE_URL must be set and non-empty")
	}

	// Set the host:port to listen for HTTP requests.
	listenHost := os.Getenv("COREOS_IPXE_SERVER_LISTEN_HOST")
	listenPort := os.Getenv("COREOS_IPXE_SERVER_LISTEN_PORT")
	if listenPort == "" {
		log.Fatal("COREOS_IPXE_SERVER_LISTEN_PORT must be set and non-empty")
	}
	hostPort := net.JoinHostPort(listenHost, listenPort)

	// The ssh public must exist as $baseDir/coreos/coreos.pub
	sshKeyPath := filepath.Join(baseDir, "coreos/coreos.pub")
	sshKey, err = sshKeyFromFile(sshKeyPath)
	if err != nil {
		log.Fatal("error reading ssh public key from " + sshKeyPath)
	}

	http.HandleFunc("/", ipxeBootScriptServer)

	// Serve kernel and pxe boot images
	staticFilePath := filepath.Join(baseDir, "coreos")
	http.Handle("/coreos/", http.StripPrefix("/coreos/", http.FileServer(http.Dir(staticFilePath))))
	log.Fatal(http.ListenAndServe(hostPort, nil))
}
