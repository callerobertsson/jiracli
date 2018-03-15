// Package commands epic command for showing epic belonging to a project/team
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	JiraCommand.AddCommand(epicsCommand)
}

var epicsCommand = &cobra.Command{
	Use:   "epics",
	Short: "List epics in a project/team",
	Long: `JiraCLI Epics Command
  Display a list of Epics belonging to a Project, i.e. team.`,

	Args: func(cc *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing project/team argument")
		}
		return nil
	},

	Run: func(cc *cobra.Command, args []string) {
		readConfig()

		arg := args[0]
		verbose("Search Epics in %v", arg)

		j := newJira()
		jql := "project = " + arg + " AND issueType = epic AND status != Done"

		issueList, err := j.SearchIssues(jql, viper.GetInt("app.skip"), viper.GetInt("app.limit"))
		if err != nil {
			fmt.Printf("Error searching for Epics in %v: %v\n", arg, err)
			os.Exit(1)
		}

		for _, issue := range issueList.Issues {
			printIssueRow(issue, "")
		}

		printIssueListSummary(issueList)

	},
}
