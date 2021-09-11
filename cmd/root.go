package cmd

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed embed/help
var helpTxt string

func Root(args []string) error {
	if len(args) < 1 || args[0] == "-h" || args[0] == "--help" {
		printHelp()
		os.Exit(0)
	}

	cmds := []Runner{
		NewEditCmd(),
		NewAttachCmd(),
		NewSendCmd(),
		NewGetCmd(),
		NewCatCmd(),
		NewListCmd(),
		NewInitCmd(),
		NewEnvCmd(),
		NewKillCmd(),
	}

	subcommand := os.Args[1]

	cmdCtx, err := NewCmdContext()
	if err != nil {
		return err
	}

	for _, cmd := range cmds {
		if cmd.Name() == subcommand || containsString(cmd.Alias(), subcommand) {
			if err := cmd.Init(os.Args[2:], *cmdCtx); err != nil {
				return err
			}
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Print(helpTxt)
}
