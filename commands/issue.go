// Package commands issue command
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	JiraCommand.AddCommand(issueCommand)
}

var issueCommand = &cobra.Command{
	Use:   "issue",
	Short: "Display issue details",
	Long:  `JiraCLI Issue Command`,

	Args: func(cc *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing issue ID argument")
		}
		return nil
	},

	Run: func(cc *cobra.Command, args []string) {
		readConfig()

		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Missing argument\n")
			os.Exit(1)
		}

		arg := args[0]

		verbose("Fetching issue %q\n", arg)

		j := newJira()
		issue, err := j.FetchIssue(arg)
		if err != nil {
			fmt.Printf("Failed to fetch issue: %v\n", err)
			os.Exit(1)
		}

		printIssue(issue, true)
	},
}
