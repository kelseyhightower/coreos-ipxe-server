package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
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

var ipxeBootScriptDefaultOut = `#!ipxe
set coreos-version latest
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: state=tmpfs: sshkey="ssh-rsa AAAAB3Ncoreos"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var ipxeBootScriptCustomSSHKeyOut = `#!ipxe
set coreos-version latest
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: state=tmpfs: sshkey="ssh-rsa AAAAB3Ncustom"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var ipxeBootScriptStateTrueOut = `#!ipxe
set coreos-version latest
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: sshkey="ssh-rsa AAAAB3Ncoreos"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var ipxeBootScriptVersionSetOut = `#!ipxe
set coreos-version 268.1.0
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: state=tmpfs: sshkey="ssh-rsa AAAAB3Ncoreos"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var ipxeBootScriptVersionSetAndStateTrueOut = `#!ipxe
set coreos-version 268.1.0
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: sshkey="ssh-rsa AAAAB3Ncoreos"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var iPxeBootScriptTests = []struct {
	body string
	code int
	url  string
}{
	{ipxeBootScriptDefaultOut, 200, "http://example.com"},
	{ipxeBootScriptCustomSSHKeyOut, 200, "http://example.com?sshkey=custom"},
	{ipxeBootScriptStateTrueOut, 200, "http://example.com?state=true"},
	{ipxeBootScriptVersionSetOut, 200, "http://example.com?version=268.1.0"},
	{ipxeBootScriptVersionSetAndStateTrueOut, 200, "http://example.com?state=true&version=268.1.0"},
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

	SetDataDir(testDataDir)
	SetBaseUrl("example.com")
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
