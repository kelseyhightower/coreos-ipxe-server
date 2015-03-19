package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kelseyhightower/coreos-ipxe-server/store/backend/boltdb"
)

var (
	store            *boltdb.Store
	defaultStorePath = "/tmp/db"
)

// commands represent the list of available sub-commands.
var commands = map[string]func([]string) error{
	"machine": machineCmdHandler,
	"profile": profileCmdHandler,
}

var usageString = `Usage: %s <machine | profile> [flags]

commands:
  machine    Manage machines.
  profile    Manage iPXE profiles.
`

func main() {
	var err error

	fs := flag.NewFlagSet("main", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, usageString, os.Args[0])
		os.Exit(0)
	}
	fs.Parse(os.Args[1:])

	if len(os.Args) < 2 {
		fs.Usage()
	}

	subcommand := os.Args[1]
	args := os.Args[2:]

	command, ok := commands[subcommand]
	if !ok {
		fmt.Printf("unknown command: %s\n", subcommand)
		fs.Usage()
	}

	store, err = boltdb.New(defaultStorePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := command(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
