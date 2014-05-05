package kernel

import (
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	want := " rootfstype=tmpfs console=tty0"
	o := New()
	options := o.String()
	if options != want {
		t.Errorf("wanted %s, got %s", want, options)
	}
}

var optionstests = []struct {
	cloudConfigUrl  string
	console         []string
	coreOSAutoLogin string
	root            string
	rootFSType      string
	sshKey          string
	options         string
}{
	{
		"http://host/config.yml",
		[]string{"tty0", "ttyS0"},
		"ttyS0",
		"",
		FstypeTmpfs,
		"ssh-rsa AAAAB3Nza...",
		" rootfstype=tmpfs console=tty0 console=ttyS0 cloud-config-url=http://host/config.yml coreos.autologin=ttyS0 sshkey=ssh-rsa AAAAB3Nza...",
	},
	{
		"",
		nil,
		"",
		"",
		"",
		"ssh-rsa AAAAB3Nza...",
		" rootfstype=tmpfs console=tty0 sshkey=ssh-rsa AAAAB3Nza...",
	},
}

func TestOptions(t *testing.T) {
	for _, tt := range optionstests {
		o := New()
		o.SetCloudConfigUrl(tt.cloudConfigUrl)
		if tt.console != nil {
			o.SetConsole(tt.console)
		}
		o.SetRootFSType(tt.rootFSType)
		o.SSHKey = tt.sshKey
		o.CoreOSAutoLogin = tt.coreOSAutoLogin
		o.Root = tt.root
		got := o.String()
		if got != tt.options {
			t.Errorf("wanted %s, got %s", tt.options, got)
		}
	}
}
