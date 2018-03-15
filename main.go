// Package main implements JiraCLI - a simple command line interface to Jira
// By Calle Robertsson <calle@robcon.se> 2018
package main

import (
	"fmt"
	"os"

	"github.com/callerobertsson/jiracli/commands"
)

func main() {
	// Just execute the cobra root command
	if err := commands.JiraCommand.Execute(); err != nil {
		fmt.Printf("Jira command failed: %v\n", err)
		os.Exit(1)
	}
}
