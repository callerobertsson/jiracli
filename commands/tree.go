// Package commands tree command
package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {

	JiraCommand.AddCommand(treeCommand)
}

var treeCommand = &cobra.Command{
	Use:   "tree",
	Short: "List an issue and all the sub-tasks",
	Long:  `JiraCLI Tree Command`,

	Args: func(cc *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing issue ID argument")
		}
		return nil
	},

	Run: func(cc *cobra.Command, args []string) {
		readConfig()

		id := args[0]
		verbose("Tree on issue %q\n", id)

		// Fetch Epic or Task (or Sub-task?)
		j := newJira()

		issue, err := j.FetchIssue(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching issue %q: %v\n", id, err)
			os.Exit(0)
		}
		printIssueRow(issue, "")

		// Default JQL for Epic
		jql := fmt.Sprintf(`\"Epic Link\"=%v`, id)
		if issue.Fields.IssueType.Name != "Epic" {
			jql = fmt.Sprintf("parent=%v", id)
		}

		verbose("Searching JQL: %q\n", jql)

		issueList, err := j.SearchIssues(jql, viper.GetInt("app.skip"), viper.GetInt("app.limit"))
		if err != nil {
			fmt.Printf("Error searching issues: %v\n", err)
			os.Exit(1)
		}

		numTasks := 0
		numSubTasks := 0
		remainingEstimateSeconds := 0
		remainingEstimateMisses := 0

		for _, issue := range issueList.Issues {
			numTasks++
			printIssueRow(&issue, "  ")

			if issue.Fields.TimeTracking != nil {
				remainingEstimateSeconds += issue.Fields.TimeTracking.RemainingEstimateSeconds
			} else {
				remainingEstimateMisses++
			}

			jql = "parent = " + issue.Key
			subIssueList, err := j.SearchIssues(jql, viper.GetInt("app.skip"), viper.GetInt("app.limit"))
			if err != nil {
				fmt.Printf("Error searching sub issues: %v\n", err)
				os.Exit(1)
			}
			for _, subissue := range subIssueList.Issues {
				numSubTasks++
				printIssueRow(&subissue, "    ")
				if subissue.Fields.TimeTracking != nil {
					remainingEstimateSeconds += subissue.Fields.TimeTracking.RemainingEstimateSeconds
				} else {
					remainingEstimateMisses++
				}
			}
		}

		fmt.Printf("%d tasks and %d subtasks\n", numTasks, numSubTasks)

		sumRemaining := time.Duration(remainingEstimateSeconds) * time.Second
		fmt.Printf("Sum of remaining estimate: %v\n", toHumanDuration(sumRemaining))
		if remainingEstimateMisses > 0 {
			fmt.Printf("Warning: %d tasks without any remaining estimation!\n", remainingEstimateMisses)
		}
	},
}
