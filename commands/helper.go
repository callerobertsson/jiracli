// Package commands implements the Cobra root command for JiraCLI
package commands

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/callerobertsson/jiracli/jira"
	"github.com/spf13/viper"
)

func printIssue(issue *jira.Issue, dump bool) {
	if dump {
		bs, err := json.MarshalIndent(issue, "", "  ")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(bs))
		return
	}

	format := "> %-14s: %s\n"

	fmt.Printf("%s\n", issue.ID)
	//fmt.Printf(format, "self", issue.Self)
	fmt.Printf(format, "key", issue.Key)
	//fmt.Printf(format, "expand", issue.Expand)
	//fields := issue.Fields
	//fmt.Printf(format, "summary", fields.Summary)
	//fmt.Printf(format, "reporter", fields.Reporter.Name)
	//fmt.Printf(format, "assignee", assigneeOrDefault(fields, "Unassigned"))
	////fmt.Printf("type: %#v\n", fields.IssueType)
	//fmt.Printf(format, "type", fields.IssueType.Name)
	//fmt.Printf(format, "description", fields.IssueType.Description)
	//fmt.Printf(format, "is subtask?", strconv.FormatBool(fields.IssueType.Subtask))
	// name, descr, subtask

}

func printIssueRow(issue *jira.Issue, indent string) {

	fields := issue.Fields

	assigneeName := "unassigned"
	if fields.Assignee != nil {
		assigneeName = issue.Fields.Assignee.Name
	}
	statusName := "unknown"
	if fields.Status != nil {
		statusName = issue.Fields.Status.Name
	}
	remainingEstimateSeconds := -1
	//if fields.TimeTracking != nil && fields.TimeTracking.RemainingEstimateSeconds != nil {
	if fields.TimeTracking != nil {
		remainingEstimateSeconds = fields.TimeTracking.RemainingEstimateSeconds
	}

	maxSummary := 50
	remainingTime := "undef"
	if remainingEstimateSeconds != -1 {
		remainingTime = toHumanDuration(time.Duration(remainingEstimateSeconds) * time.Second)
	}

	fmt.Printf(indent+"%-7v %-"+strconv.Itoa(maxSummary)+"q [%v, %v] remaining: %v\n",
		issue.Key,
		limitString(issue.Fields.Summary, maxSummary),
		statusName,
		assigneeName,
		remainingTime,
	)

	if viper.GetBool("app.verbose") {
		fmt.Printf("       Summary: %v\n", issue.Fields.Summary)
	}

}

func limitString(s string, max int) string {
	if max > len(s) {
		max = len(s)
	}
	return s[0:max]
}

func printIssueListSummary(issueList *jira.IssueList) {
	fmt.Printf("%d of %d from %d limit %d\n", len(issueList.Issues), issueList.Total, issueList.StartAt, issueList.MaxResults)
}

func toHumanDuration(t time.Duration) (res string) {

	s := t.Seconds()

	// Weeks, 1 week = 5 days
	w := math.Trunc(s / (5 * 6 * 60 * 60))
	s -= w * (5 * 6 * 60 * 60)
	if w >= 1.0 {
		res += fmt.Sprintf("%vw ", w)
	}

	// Days, 1 day = 6 hours
	d := math.Trunc(s / (6 * 60 * 60))
	s -= d * (6 * 60 * 60)
	if d >= 1.0 {
		res += fmt.Sprintf("%vd ", d)
	}

	// Hours, show 0h if week and day are < 1
	h := math.Trunc(s / (60 * 60))
	s -= h * (60 * 60)
	if h >= 1.0 || res == "" {
		res += fmt.Sprintf("%vh ", h)
	}

	// Show hours sum only if week or days are shown
	if w >= 1.0 || d >= 1.0 {
		return res + fmt.Sprintf("(%vh)", t.Hours())
	}

	return res
}

func newJira() *jira.Jira {

	// Ask for user and pass, if not in viper
	resolveUserAndPass()

	// Create a new Jira API client
	url := viper.GetString("api.host") + viper.GetString("api.path")
	return jira.New(url, jira.Auth{
		User: viper.GetString("api.user"),
		Pass: viper.GetString("api.pass"),
	}, viper.GetBool("dbg.dump-response"))
}
