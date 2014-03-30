package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var ipxeBootScriptDefaultOut = `#!ipxe
set coreos-version latest
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: state=tmpfs: sshkey="ssh-rsa AAAAB3NzaC1yc2"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var ipxeBootScriptStateTrueOut = `#!ipxe
set coreos-version latest
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: sshkey="ssh-rsa AAAAB3NzaC1yc2"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var ipxeBootScriptVersionSetOut = `#!ipxe
set coreos-version 268.1.0
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: state=tmpfs: sshkey="ssh-rsa AAAAB3NzaC1yc2"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var ipxeBootScriptVersionSetAndStateTrueOut = `#!ipxe
set coreos-version 268.1.0
set base-url http://example.com/coreos/amd64-generic/${coreos-version}
kernel ${base-url}/coreos_production_pxe.vmlinuz root=squashfs: sshkey="ssh-rsa AAAAB3NzaC1yc2"
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
`

var iPxeBootScriptTests = []struct {
	body string
	code int
	url  string
}{
	{ipxeBootScriptDefaultOut, 200, "http://example.com"},
	{ipxeBootScriptStateTrueOut, 200, "http://example.com?state=true"},
	{ipxeBootScriptVersionSetOut, 200, "http://example.com?version=268.1.0"},
	{ipxeBootScriptVersionSetAndStateTrueOut, 200, "http://example.com?state=true&version=268.1.0"},
}

func TestIPxeBootScriptServer(t *testing.T) {
	SetSSHKey("ssh-rsa AAAAB3NzaC1yc2")
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
