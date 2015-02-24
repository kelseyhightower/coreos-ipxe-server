package profile

import (
	"bytes"
	"fmt"
	"text/template"
)

// Profile represents and CoreOS iPXE profile.
type Profile struct {
	CloudConfig     string
	Console         []string
	CoreOSAutologin string
	Name            string
	ReleaseChannel  string
	Root            string
	RootFstype      string
	SSHKey          string
}

var outputTemplate = template.Must(template.New("profile").Parse(`
name             : {{.Name}}
cloud-config     : {{.CloudConfig}}
console          : {{.Console}}
coreos-autologin : {{.CoreOSAutologin}}
release-channel  : {{.ReleaseChannel}}
root             : {{.Root}}
root-fstype      : {{.RootFstype}}
ssh-key          : {{.SSHKey}}
`))

func (p *Profile) String() string {
	var buf bytes.Buffer
	err := outputTemplate.Execute(&buf, p)
	if err != nil {
		fmt.Println("Error generating profile output:", err)
	}
	return buf.String()
}
