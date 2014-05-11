// Copyright 2014 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/kelseyhightower/coreos-ipxe-server/config"
)

type testSSHKey struct {
	id  string
	key string
}

func createTestData(sshKeys []testSSHKey) (string, error) {
	d, err := ioutil.TempDir("", "coreos-ipxe-server")
	if err != nil {
		return "", err
	}
	sshKeyDir := filepath.Join(d, "sshkeys")
	err = os.Mkdir(sshKeyDir, 0755)
	if err != nil {
		return "", err
	}
	for _, s := range sshKeys {
		sshKeyPath := filepath.Join(sshKeyDir, fmt.Sprintf("%s.pub", s.id))
		err := ioutil.WriteFile(sshKeyPath, []byte(s.key), 0644)
		if err != nil {
			return "", err
		}
	}
	return d, nil
}

var defaultParameterOut = `#!ipxe
set coreos-version latest
set base-url http://example.com/images/amd64-usr/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz rootfstype=tmpfs console=tty0
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var SSHKeyParameterOut = `#!ipxe
set coreos-version latest
set base-url http://example.com/images/amd64-usr/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz rootfstype=tmpfs console=tty0 sshkey="ssh-rsa AAAAB3Ncustom"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var versionParameterOut = `#!ipxe
set coreos-version 298.0.0
set base-url http://example.com/images/amd64-usr/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz rootfstype=tmpfs console=tty0
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var consoleParameterOut = `#!ipxe
set coreos-version latest
set base-url http://example.com/images/amd64-usr/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz rootfstype=tmpfs console=tty0 console=ttyS0
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var cloudConfigParameterOut = `#!ipxe
set coreos-version latest
set base-url http://example.com/images/amd64-usr/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz rootfstype=tmpfs console=tty0 cloud-config-url=http://example.com/configs/cloud-config.yml
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var versionSSHKeyParameterOut = `#!ipxe
set coreos-version 298.0.0
set base-url http://example.com/images/amd64-usr/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz rootfstype=tmpfs console=tty0 sshkey="ssh-rsa AAAAB3Ncustom"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var iPxeBootScriptTests = []struct {
	body string
	code int
	url  string
}{
	{defaultParameterOut, 200, "http://example.com"},
	{SSHKeyParameterOut, 200, "http://example.com?sshkey=custom"},
	{versionParameterOut, 200, "http://example.com?version=298.0.0"},
	{versionSSHKeyParameterOut, 200, "http://example.com?version=298.0.0&sshkey=custom"},
	{consoleParameterOut, 200, "http://example.com?console=tty0,ttyS0"},
	{cloudConfigParameterOut, 200, "http://example.com?cloudconfig=cloud-config"},
}

func TestIPxeBootScriptServer(t *testing.T) {
	coreosSSHKey := testSSHKey{
		id:  "coreos",
		key: "ssh-rsa AAAAB3Ncoreos",
	}
	customSSHKey := testSSHKey{
		id:  "custom",
		key: "ssh-rsa AAAAB3Ncustom",
	}

	testDataDir, err := createTestData([]testSSHKey{coreosSSHKey, customSSHKey})
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(testDataDir)

	config.DataDir = testDataDir
	config.BaseUrl = "example.com"
	for _, v := range iPxeBootScriptTests {
		req, err := http.NewRequest("GET", v.url, nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		ipxeBootScriptServer(w, req)
		if w.Body.String() != v.body {
			t.Errorf("expected %s\ngot %s\n", v.body, w.Body.String())
		}
	}
}
