package kernel

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	FstypeTmpfs = "tmpfs"
	FstypeBtrfs = "btrfs"
)

var (
	ErrInvalidFstype = errors.New("kernal: invalid fstype")
)

type Options struct {
	cloudConfigUrl  string
	console         []string
	CoreOSAutoLogin string
	Root            string
	rootFSType      string
	SSHKey          string
}

func New() *Options {
	o := &Options{
		console:    []string{"tty0"},
		rootFSType: FstypeTmpfs,
	}
	return o
}

func (o *Options) SetCloudConfigUrl(url string) {
	o.cloudConfigUrl = url
}

func (o *Options) SetRootFSType(fstype string) error {
	if fstype != FstypeTmpfs || fstype != FstypeBtrfs {
		return ErrInvalidFstype
	}
	o.rootFSType = fstype
	return nil
}

func (o *Options) SetConsole(console []string) {
	o.console = console
}

func (o *Options) String() string {
	var options bytes.Buffer
	options.WriteString(fmt.Sprintf("rootfstype=%s", o.rootFSType))
	for _, c := range o.console {
		options.WriteString(fmt.Sprintf(" console=%s", c))
	}
	if o.cloudConfigUrl != "" {
		options.WriteString(fmt.Sprintf(" cloud-config-url=%s", o.cloudConfigUrl))
	}
	if o.CoreOSAutoLogin != "" {
		options.WriteString(fmt.Sprintf(" coreos.autologin=%s", o.CoreOSAutoLogin))
	}
	if o.SSHKey != "" {
		options.WriteString(fmt.Sprintf(" sshkey=%s", o.SSHKey))
	}
	if o.Root != "" {
		options.WriteString(fmt.Sprintf(" root=%s", o.Root))
	}
	return options.String()
}
