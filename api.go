package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kelseyhightower/coreos-ipxe-server/config"
	"github.com/kelseyhightower/coreos-ipxe-server/kernel"
)

const ipxeBootScript = `#!ipxe
set coreos-release-channel {{.ReleaseChannel}}
set base-url http://{{.BaseUrl}}/images/${coreos-release-channel}
kernel ${base-url}/coreos_production_pxe.vmlinuz{{.Options}}
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var replacer = strings.NewReplacer(":", "-", " ", "-")

func ipxeBootScriptServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("creating boot script for %s", r.RemoteAddr)
	log.Println(url.QueryUnescape(r.URL.String()))

	baseUrl := config.BaseUrl
	if baseUrl == "" {
		baseUrl = r.Host
	}
	v := r.URL.Query()

	options := kernel.New()

	// Process the profile parameter.
	profile := v.Get("profile")
	mac := replacer.Replace(v.Get("mac"))
	asset := replacer.Replace(v.Get("asset"))
	serial := replacer.Replace(v.Get("serial"))

	for _, s := range []string{profile, mac, asset, serial, "default"} {
		if s != "" {
			profilePath := filepath.Join(config.DataDir, fmt.Sprintf("profiles/%s.json", s))
			log.Println("loading profile:", profilePath)
			if _, err := os.Stat(profilePath); os.IsNotExist(err) {
				continue
			}
			err := kernalOptionsFromFile(profilePath, options)
			if err != nil {
				log.Printf("Error reading kernal options from %s: %s", profilePath, err)
				http.Error(w, err.Error(), 500)
				return
			}
			break
		}
	}

	if options.CloudConfig != "" {
		options.SetCloudConfigUrl(fmt.Sprintf("http://%s/configs/%s.yml", baseUrl, options.CloudConfig))
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

	// Process the iPXE boot script template.
	t, err := template.New("ipxebootscript").Parse(ipxeBootScript)
	if err != nil {
		log.Print("Error generating iPXE boot script: " + err.Error())
		http.Error(w, "Error generating the iPXE boot script", 500)
		return
	}
	data := map[string]string{
		"BaseUrl":        baseUrl,
		"Options":        options.String(),
		"ReleaseChannel": options.ReleaseChannel,
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Print("Error generating iPXE boot script: " + err.Error())
		http.Error(w, "Error generating the iPXE boot script", 500)
		return
	}
	return
}

func sshKeyServer(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	keyName := v.Get("name")
	if keyName != "" {
		log.Printf("retrieving ssh key %s.pub for %s", keyName, r.RemoteAddr)
		sshKeyPath := filepath.Join(config.DataDir, fmt.Sprintf("sshkeys/%s.pub", keyName))
		sshKey, err := sshKeyFromFile(sshKeyPath)
		if err != nil {
			log.Printf("Error reading ssh publickey from %s: %s", sshKeyPath, err)
			http.Error(w, err.Error(), 500)
			return
		}
		data := []byte(fmt.Sprintf("[{\"key\": \"%s\"}]", sshKey))
		w.Write(data)
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
