package machine

import (
	"bytes"
	"fmt"
	"net"
	"text/template"
)

type Machine struct {
	Name         string
	HardwareAddr net.HardwareAddr
	Profile      string
}

type Data map[string]string

func New(name string, mac string) (*Machine, error) {
	macAddr, err := net.ParseMAC(mac)
	if err != nil {
		return nil, err
	}
	m := &Machine{
		Name:         name,
		HardwareAddr: macAddr,
	}
	return m, nil
}

func (m *Machine) SetProfile(name string) error {
	m.Profile = name
	return nil
}

var outputTemplate = template.Must(template.New("machine").Parse(`
name        : {{.Name}}
mac-address : {{.HardwareAddr}}
profile     : {{.Profile}}
`))

func (m *Machine) String() string {
	var buf bytes.Buffer
	err := outputTemplate.Execute(&buf, m)
	if err != nil {
		fmt.Println("Error generating profile output:", err)
	}
	return buf.String()
}
