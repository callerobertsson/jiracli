// Package commands search command
package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	JiraCommand.AddCommand(searchCommand)
}

var searchCommand = &cobra.Command{
	Use:   "search",
	Short: "Do a JQL query",
	Long:  `JiraCLI Search Command`,

	Args: func(cc *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing JQL query")
		}
		return nil
	},

	Run: func(cc *cobra.Command, args []string) {
		readConfig()

		verbose("Searching JQL %q\n", args)

		j := newJira()

		issueList, err := j.SearchIssues(strings.Join(args, " "), viper.GetInt("app.skip"), viper.GetInt("app.limit"))
		if err != nil {
			fmt.Printf("Error searching issues: %v\n", err)
			os.Exit(1)
		}

		for _, issue := range issueList.Issues {
			printIssueRow(&issue, "")
		}

		printIssueListSummary(issueList)

	},
}
