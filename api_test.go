// Copyright 2014 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/kelseyhightower/coreos-ipxe-server/config"
	"github.com/kelseyhightower/coreos-ipxe-server/kernel"
)

func createTestData(profiles map[string]*kernel.Options, sshKeys map[string]string) (string, error) {
	d, err := ioutil.TempDir("", "coreos-ipxe-server")
	if err != nil {
		return "", err
	}
	sshKeyDir := filepath.Join(d, "sshkeys")
	err = os.Mkdir(sshKeyDir, 0755)
	if err != nil {
		return "", err
	}
	for k, v := range sshKeys {
		sshKeyPath := filepath.Join(sshKeyDir, fmt.Sprintf("%s.pub", k))
		err := ioutil.WriteFile(sshKeyPath, []byte(v), 0644)
		if err != nil {
			return "", err
		}
	}

	profileDir := filepath.Join(d, "profiles")
	err = os.Mkdir(profileDir, 0755)
	if err != nil {
		return "", err
	}
	for k, v := range profiles {
		profilePath := filepath.Join(profileDir, fmt.Sprintf("%s.json", k))
		data, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		err = ioutil.WriteFile(profilePath, data, 0644)
		if err != nil {
			return "", err
		}
	}
	return d, nil
}

var profileAOut = `#!ipxe
set coreos-release-channel alpha
set base-url http://example.com/images/${coreos-release-channel}
kernel ${base-url}/coreos_production_pxe.vmlinuz
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var profileBOut = `#!ipxe
set coreos-release-channel stable
set base-url http://example.com/images/${coreos-release-channel}
kernel ${base-url}/coreos_production_pxe.vmlinuz rootfstype=btrfs console=tty0 console=ttyS0 cloud-config-url=http://example.com/configs/b.yml coreos.autologin=ttyS0 sshkey="ssh-rsa AAAAB3Ncoreos" root=/dev/sda1
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var iPxeBootScriptTests = []struct {
	name    string
	body    string
	code    int
	baseUrl string
	url     string
}{
	{"a", profileAOut, 200, "", "http://example.com?profile=a"},
	{"b", profileBOut, 200, "example.com", "http://example.com?profile=b"},
	{"c", "", 500, "example.com", "http://example.com?profile=c"},
	{"d", "", 500, "example.com", "http://example.com?profile=d"},
}

func TestIPxeBootScriptServer(t *testing.T) {
	sshkeys := map[string]string{
		"coreos": "ssh-rsa AAAAB3Ncoreos",
	}

	profiles := map[string]*kernel.Options{
		"a": &kernel.Options{
			CloudConfig:     "",
			Console:         []string{},
			CoreOSAutologin: "",
			Root:            "",
			RootFstype:      "",
			SSHKey:          "",
			ReleaseChannel:  "alpha",
		},
		"b": &kernel.Options{
			CloudConfig:     "b",
			Console:         []string{"tty0", "ttyS0"},
			CoreOSAutologin: "ttyS0",
			Root:            "/dev/sda1",
			RootFstype:      "btrfs",
			SSHKey:          "coreos",
			ReleaseChannel:  "stable",
		},
		"c": &kernel.Options{
			CloudConfig:     "c",
			Console:         []string{"tty0", "ttyS0"},
			CoreOSAutologin: "ttyS0",
			Root:            "/dev/sda1",
			RootFstype:      "btrfs",
			SSHKey:          "imabadkey",
			ReleaseChannel:  "stable",
		},
	}

	testDataDir, err := createTestData(profiles, sshkeys)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(testDataDir)

	config.DataDir = testDataDir
	for _, v := range iPxeBootScriptTests {
		config.BaseUrl = v.baseUrl
		req, err := http.NewRequest("GET", v.url, nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		ipxeBootScriptServer(w, req)
		if w.Code == 200 && (v.name == "a" || v.name == "b") {
			if w.Body.String() != v.body {
				t.Errorf("expected %s\ngot %s\n", v.body, w.Body.String())
			}
		} else if (v.name == "c" || v.name == "d") && w.Code != 500 {
			t.Errorf("expected %d\ngot %d\n", v.code, w.Code)
		}
	}
}

var SSHKeyServerTests = []struct {
	name    string
	body    string
	code    int
	baseUrl string
	url     string
}{
	{"a", `[{"key": "ssh-rsa AAAAB3Ncoreos"}]`, 200, "", "http://example.com/keys?name=coreos"},
	{"b", `[{"key": "ssh-rsa AAAAB3Nfoo"}]`, 200, "example.com", "http://example.com/keys?name=foo"},
	{"c", "", 500, "example.com", "http://example.com/keys?name=badkey"},
}

func TestSSHKeyServer(t *testing.T) {
	sshkeys := map[string]string{
		"coreos": "ssh-rsa AAAAB3Ncoreos",
		"foo":    "ssh-rsa AAAAB3Nfoo",
	}

	testDataDir, err := createTestData(nil, sshkeys)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(testDataDir)

	config.DataDir = testDataDir
	for _, v := range SSHKeyServerTests {
		config.BaseUrl = v.baseUrl
		req, err := http.NewRequest("GET", v.url, nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		sshKeyServer(w, req)
		if w.Code == 200 && (v.name == "a" || v.name == "b") {
			if w.Body.String() != v.body {
				t.Errorf("expected %s\ngot %s\n", v.body, w.Body.String())
			}
		} else if (v.name == "c") && w.Code != 500 {
			t.Errorf("expected %d\ngot %d\n", v.code, w.Code)
		}
	}
}
