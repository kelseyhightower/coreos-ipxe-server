package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"os"

	"github.com/kelseyhightower/coreos-ipxe-server/machine"
)

var machineUsageString = `Usage: %s machines <action>

actions:
  create   Create a new machine.
  get      Get a machine.
  list     List all machines.
`

func machineCmdHandler(args []string) error {
	fs := flag.NewFlagSet("machine", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, machineUsageString, os.Args[0])
	}
	fs.Parse(args)
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "missing action\n")
		fs.Usage()
		os.Exit(1)
	}
	action := args[0]
	switch action {
	case "create":
		return createMachine(args[1:])
	case "get":
		return getMachine(args[1:])
	case "list":
		return listMachines(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown action\n")
		fs.Usage()
		os.Exit(1)
	}
	return nil
}

func createMachine(arguments []string) error {
	var (
		name    string
		mac     string
		profile string
	)
	fs := flag.NewFlagSet("createMachine", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s machine create -mac <mac> [flags]\n", os.Args[0])
		fs.PrintDefaults()
	}
	fs.StringVar(&name, "name", "", "machine name")
	fs.StringVar(&mac, "mac", "", "machine mac address")
	fs.StringVar(&profile, "profile", "default", "machine profile")

	err := fs.Parse(arguments)
	if err == flag.ErrHelp {
		fs.Usage()
		os.Exit(0)
	}
	if err != nil {
		return err
	}

	if mac == "" {
		fmt.Fprintf(os.Stderr, "A valid (non-empty) mac address is required. Use the -mac flag.\n")
		os.Exit(1)
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	m, err := machine.New(name, mac)
	if err != nil {
		return err
	}
	m.SetProfile(profile)

	err = enc.Encode(m)
	if err != nil {
		return err
	}

	err = store.Save("machines", name, buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func listMachines(args []string) error {
	ms, err := store.List("machines")
	if err != nil {
		return err
	}
	fmt.Println("machines:\n")
	for k := range ms {
		fmt.Println(k)
	}
	return nil
}

func getMachine(args []string) error {
	var (
		name string
	)
	fs := flag.NewFlagSet("findMachine", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s machine find -name <name>\n", os.Args[0])
		fs.PrintDefaults()
	}
	fs.StringVar(&name, "name", "", "machine name")

	err := fs.Parse(args)
	if err == flag.ErrHelp {
		fs.Usage()
		os.Exit(0)
	}
	if err != nil {
		return err
	}

	if name == "" {
		fmt.Fprintf(os.Stderr, "A valid (non-empty) name is required. Use the -name flag.\n")
		os.Exit(1)
	}

	rawMachine, err := store.Get("machines", name)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(bytes.NewReader(rawMachine))
	var node machine.Machine
	if err := dec.Decode(&node); err != nil {
		return err
	}

	fmt.Println(&node)

	return nil
}
