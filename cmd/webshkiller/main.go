package main

import (
	"fmt"
	"github.com/KosukeShimofuji/webshkiller"
	"github.com/mitchellh/cli"
	_ "log"
	"os"
)

// Implement Core command
type CoreCommand struct {
}

func (c *CoreCommand) Synopsis() string {
	return `list core command`
}

func (c *CoreCommand) Help() string {
	return "Usage: webshkiller core [command]"
}

func (c *CoreCommand) Run(args []string) int {
	fmt.Printf("core init - initialize core instance\n")
	return 0
}

// Implement Core init command
type CoreInitCommand struct {
}

func (c *CoreInitCommand) Synopsis() string {
	return "initialize core instance"
}

func (c *CoreInitCommand) Help() string {
	return "Usage: webshkiller core init"
}

func (c *CoreInitCommand) Run(args []string) int {
	fmt.Printf("initialie core instance\n")
	return 0
}

// entry point
func main() {
	c := cli.NewCLI("webshkiller", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"core": func() (cli.Command, error) {
			return &CoreCommand{}, nil
		},
		"core init": func() (cli.Command, error) {
			return &CoreInitCommand{}, nil
		},
	}
	exitCode, err := c.Run()
	if err != nil {
		fmt.Printf("Failed to execute: %s\n", err.Error())
	}
	os.Exit(exitCode)
}
