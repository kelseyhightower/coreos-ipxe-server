package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kelseyhightower/coreos-ipxe-server/config"
	"github.com/kelseyhightower/coreos-ipxe-server/kernel"
)

const ipxeBootScript = `#!ipxe
set coreos-version {{.Version}}
set base-url http://{{.BaseUrl}}/images/amd64-usr/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz {{.Options}}
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

func ipxeBootScriptServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("creating boot script for %s", r.RemoteAddr)

	baseUrl := config.BaseUrl
	if baseUrl == "" {
		baseUrl = r.Host
	}
	v := r.URL.Query()

	options := kernel.New()

	// Process the profile parameter.
	profile := v.Get("profile")
	if profile != "" {
		profilePath := filepath.Join(config.DataDir, fmt.Sprintf("profiles/%s.json", profile))
		err := kernalOptionsFromFile(profilePath, options)
		if err != nil {
			log.Printf("Error reading kernal options from %s: %s", profilePath, err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	// Process the root parameter.
	root := v.Get("root")
	if root != "" {
		options.Root = root
	}
	// Process the fstype parameter.
	fstype := v.Get("fstype")
	if fstype != "" {
		options.RootFstype = fstype
	}

	// Process the console parameter.
	console := v.Get("console")
	if console != "" {
		options.Console = strings.Split(console, ",")
	}

	// Process the cloudconfig parameter.
	cloudConfigId := v.Get("cloudconfig")
	if cloudConfigId != "" {
		options.CloudConfig = cloudConfigId
	}
	if options.CloudConfig != "" {
		options.SetCloudConfigUrl(fmt.Sprintf("http://%s/configs/%s.yml",
			baseUrl, options.CloudConfig))
	}

	// Process the sshkey paramter.
	sshKeyId := v.Get("sshkey")
	if sshKeyId != "" {
		options.SSHKey = sshKeyId
	}
	if options.SSHKey != "" {
		sshKeyPath := filepath.Join(config.DataDir, fmt.Sprintf("sshkeys/%s.pub", options.SSHKey))
		sshKey, err := sshKeyFromFile(sshKeyPath)
		if err != nil {
			log.Printf("Error reading ssh publickey from %s: %s", sshKeyPath, err)
			http.Error(w, err.Error(), 500)
			return
		}
		options.SSHKey = sshKey
	}

	// Process the version parameter.
	version := v.Get("version")
	if version != "" {
		options.Version = version
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
		"Version": options.Version,
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

func kernalOptionsFromFile(filename string, options *kernel.Options) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, options)
	if err != nil {
		return err
	}
	return nil
}
