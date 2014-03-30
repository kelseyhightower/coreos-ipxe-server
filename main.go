// Copyright 2014 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var (
	dataDir         string
	baseUrl         string
	defaultSSHKeyId = "coreos"
)

const ipxeBootScript = `#!ipxe
set coreos-version {{.Version}}
set base-url http://{{.BaseUrl}}/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: {{if not .State}}state=tmpfs: {{end}}sshkey="{{.SSHKey}}"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

// SetDataDir sets the data directory.
func SetDataDir(dir string) {
	dataDir = dir
}

// SetBaseUrl sets the base url.
func SetBaseUrl(url string) {
	baseUrl = url
}

func ipxeBootScriptServer(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	version := v.Get("version")
	if version == "" {
		version = "latest"
	}
	state := v.Get("state")
	if state != "true" {
		state = ""
	}
	sshKeyId := v.Get("sshkey")
	if sshKeyId == "" {
		sshKeyId = defaultSSHKeyId
	}
	sshKeyPath := filepath.Join(dataDir, fmt.Sprintf("sshkeys/%s.pub", sshKeyId))
	sshKey, err := sshKeyFromFile(sshKeyPath)
	if err != nil {
		log.Printf("Error reading ssh publickey from %s: %s", sshKeyPath, err)
		http.Error(w, err.Error(), 500)
		return
	}

	t, err := template.New("ipxebootscript").Parse(ipxeBootScript)
	if err != nil {
		log.Print("Error generating iPXE boot script: " + err.Error())
		http.Error(w, "Error generating the iPXE boot script", 500)
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
		log.Print("Error generating iPXE boot script: " + err.Error())
		http.Error(w, "Error generating the iPXE boot script", 500)
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
	listenPort := os.Getenv("COREOS_IPXE_SERVER_LISTEN_PORT")
	if listenPort == "" {
		log.Fatal("COREOS_IPXE_SERVER_LISTEN_PORT must be set and non-empty")
	}
	hostPort := net.JoinHostPort(listenHost, listenPort)

	http.HandleFunc("/", ipxeBootScriptServer)
	// Serve kernel and pxe boot images
	staticFilePath := filepath.Join(dataDir, "coreos")
	http.Handle("/coreos/", http.StripPrefix("/coreos/", http.FileServer(http.Dir(staticFilePath))))
	log.Fatal(http.ListenAndServe(hostPort, nil))
}
