package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
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
	cloudConfigUrl := v.Get("cloudconfig")
	if cloudConfigUrl == "" {
		cloudConfigUrl = defaultCloudConfigUrl
	}
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
