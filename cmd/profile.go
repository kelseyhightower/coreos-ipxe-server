package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"os"

	"github.com/kelseyhightower/coreos-ipxe-server/profile"
)

var profileUsageString = `Usage: %s profile <action>

actions:
  create   Create a new profile.
  get      Get a profile.
  list     List all profiles.
`

func profileCmdHandler(args []string) error {
	fs := flag.NewFlagSet("profile", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, profileUsageString, os.Args[0])
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
		return createProfile(args[1:])
	case "list":
		return listProfiles(args[1:])
	case "get":
		return getProfile(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown action\n")
		fs.Usage()
		os.Exit(1)
	}
	return nil
}

func createProfile(args []string) error {
	var (
		cloudConfig     string
		console         string
		coreOSAutologin string
		name            string
		releaseChannel  string
		root            string
		rootFstype      string
		sshKey          string
	)

	fs := flag.NewFlagSet("createProfile", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s profile create -name <name> [flags]\n", os.Args[0])
		fs.PrintDefaults()
	}

	fs.StringVar(&cloudConfig, "cloud-config", "default", "name of the cloud-config file")
	fs.StringVar(&console, "console", "", "the consoles to enable kernel output and login prompt on")
	fs.StringVar(&coreOSAutologin, "coreos-auto-login", "", "the consoles to enable shell without password")
	fs.StringVar(&name, "name", "", "profile name")
	fs.StringVar(&releaseChannel, "release-channel", "stable", "CoreOS release channel")
	fs.StringVar(&root, "root", "", "local filesystem name")
	fs.StringVar(&rootFstype, "root-fs-type", "", "file system type for the writable root filesystem")
	fs.StringVar(&sshKey, "ssh-key", "", "the ssh key name for the core user")

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

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	p := &profile.Profile{
		CloudConfig:     cloudConfig,
		Console:         []string{console},
		CoreOSAutologin: coreOSAutologin,
		Name:            name,
		ReleaseChannel:  releaseChannel,
		Root:            root,
		RootFstype:      rootFstype,
		SSHKey:          sshKey,
	}

	if err := enc.Encode(p); err != nil {
		return err
	}

	if err := store.Save("profiles", name, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func listProfiles(args []string) error {
	ps, err := store.List("profiles")
	if err != nil {
		return err
	}
	fmt.Println("profiles:\n")
	for k := range ps {
		fmt.Printf("  %s\n", k)
	}
	return nil
}

func getProfile(args []string) error {
	var (
		name string
	)

	fs := flag.NewFlagSet("getProfile", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s profile get -name <name> [flags]\n", os.Args[0])
		fs.PrintDefaults()
	}

	fs.StringVar(&name, "name", "", "profile name")

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

	rawProfile, err := store.Get("profiles", name)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(bytes.NewReader(rawProfile))
	var p profile.Profile
	if err := dec.Decode(&p); err != nil {
		return err
	}

	fmt.Println(&p)

	return nil
}
