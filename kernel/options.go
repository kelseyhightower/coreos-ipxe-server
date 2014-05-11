package kernel

import (
	"bytes"
	"fmt"
)

type Options struct {
	CloudConfig     string   `json:"cloud_config"`
	Console         []string `json:"console"`
	CoreOSAutologin string   `json:"coreos_autologin"`
	Root            string   `json:"root"`
	RootFstype      string   `json:"rootfstype"`
	SSHKey          string   `json:"sshkey"`
	Version         string   `json:"version"`
	cloudConfigUrl  string
}

func New() *Options {
	return &Options{}
}

func (o *Options) SetCloudConfigUrl(url string) {
	o.cloudConfigUrl = url
}

func (o *Options) String() string {
	var options bytes.Buffer
	if o.RootFstype != "" {
		options.WriteString(fmt.Sprintf(" rootfstype=%s", o.RootFstype))
	}
	for _, c := range o.Console {
		options.WriteString(fmt.Sprintf(" console=%s", c))
	}
	if o.cloudConfigUrl != "" {
		options.WriteString(fmt.Sprintf(" cloud-config-url=%s", o.cloudConfigUrl))
	}
	if o.CoreOSAutologin != "" {
		options.WriteString(fmt.Sprintf(" coreos.autologin=%s", o.CoreOSAutologin))
	}
	if o.SSHKey != "" {
		options.WriteString(fmt.Sprintf(" sshkey=\"%s\"", o.SSHKey))
	}
	if o.Root != "" {
		options.WriteString(fmt.Sprintf(" root=%s", o.Root))
	}
	return options.String()
}
