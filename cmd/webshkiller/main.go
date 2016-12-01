package main

import (
	"flag"
	"fmt"
	_ "github.com/KosukeShimofuji/WebSHKiller"
	"github.com/KosukeShimofuji/WebSHKiller/logger"
	"github.com/mitchellh/cli"
	"log"
	"os"
)

// global variable settings
var DEBUG_FLAG = false

// Implement Core command
type ControlCommand struct {
}

func (c *ControlCommand) Synopsis() string {
	return `list core command`
}

func (c *ControlCommand) Help() string {
	return "Usage: webshkiller core [command]"
}

func (c *ControlCommand) Run(args []string) int {
	fmt.Printf("core init - initialize core instance\n")
	return 0
}

// Implement Core init command
type ControlInitCommand struct {
}

func (c *ControlInitCommand) Synopsis() string {
	return "initialize core instance"
}

func (c *ControlInitCommand) Help() string {
	return "Usage: webshkiller core init [-debug]"
}

func (c *ControlInitCommand) Run(args []string) int {
	// check credentials of openstack
	var debug bool

	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.BoolVar(&debug, "debug", false, "Run as DEBUG mode")

	if err := flags.Parse(args); err != nil {
		log.Fatal(err)
	}

	DEBUG_FLAG = debug

	logger.Debug("Check openstack's credentials", DEBUG_FLAG)

	return 0
}

// entry point
func main() {
	// settings log package for output line number
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// settings sub command
	c := cli.NewCLI("webshkiller", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"core": func() (cli.Command, error) {
			return &ControlCommand{}, nil
		},
		"core init": func() (cli.Command, error) {
			return &ControlInitCommand{}, nil
		},
	}
	exitCode, err := c.Run()
	if err != nil {
		fmt.Printf("Failed to execute: %s\n", err.Error())
	}
	os.Exit(exitCode)
}
