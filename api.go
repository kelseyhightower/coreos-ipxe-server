package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/kelseyhightower/coreos-ipxe-server/kernel"
)

const ipxeBootScript = `#!ipxe
set coreos-version {{.Version}}
set base-url http://{{.BaseUrl}}/coreos/amd64-usr/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz {{.Options}}"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

func ipxeBootScriptServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("creating boot script for %s", r.RemoteAddr)
	v := r.URL.Query()

	options := kernel.New()

	// Process the cloudconfig parameter.
	cloudConfigId := v.Get("cloudconfig")
	if cloudConfigId != "" {
		// cloudConfigPath := filepath.Join(dataDir, fmt.Sprintf("cloud-configs/%s.yml", cloudConfigId))
	}

	// Process the sshkey paramter.
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
	options.SSHKey = sshKey

	// Process the version parameter.
	version := v.Get("version")
	if version == "" {
		version = "latest"
	}

	// Process the iPXE boot script template.
	t, err := template.New("ipxebootscript").Parse(ipxeBootScript)
	if err != nil {
		log.Print("Error generating iPXE boot script: " + err.Error())
		http.Error(w, "Error generating the iPXE boot script", 500)
		return
	}
	data := map[string]string{
		"BaseUrl": baseUrl,
		"Options": options.String(),
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
