// Package jira implements a Jira client API
package jira

import (
	"encoding/json"
	"fmt"
)

// Issue response data
type Issue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Fields Fields `json:"fields"`
}

// Fields response data
type Fields struct {
	IssueType    *IssueType    `json:"issuetype"`
	Status       *Status       `json:"status"`
	Summary      string        `json:"summary"`
	Assignee     *Assignee     `json:"assignee"`
	Parent       *Parent       `json:"parent"`
	TimeTracking *TimeTracking `json:"timetracking"`
}

// IssueType response data
type IssueType struct {
	Name string `json:"name"`
}

// Status response data
type Status struct {
	Name string `json:"name"`
}

// Assignee response data
type Assignee struct {
	Name string `json:"name"`
}

// Parent response data
type Parent struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

// TimeTracking response data
type TimeTracking struct {
	RemainingEstimateSeconds int `json:"remainingEstimateSeconds"`
}

// IssueList response data
type IssueList struct {
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

// FetchIssue gets an Issue from the Jira API
func (j *Jira) FetchIssue(key string) (*Issue, error) {
	data, err := j.apiGET("issue/" + key)
	if err != nil {
		return nil, err
	}

	issue := Issue{}
	err = json.Unmarshal(data, &issue)
	if err != nil {
		return nil, err
	}

	return &issue, nil
}

// SearchIssues finds Jira issues using a JQL query
func (j *Jira) SearchIssues(jql string, skip, limit int) (*IssueList, error) {
	body := fmt.Sprintf(`{
		"jql": "%v",
		"startAt": %v,
		"maxResults": %v,
		"fields": [
			"summary",
			"status",
			"assignee",
			"timetracking"
		],
		"fieldsByKeys": false
	}`, jql, skip, limit)

	j.dump("Search Body: %q", body)
	data, err := j.apiPOST("search", body)
	if err != nil {
		return nil, err
	}

	j.dump("Response: %v", string(data))

	issueList := IssueList{}
	err = json.Unmarshal(data, &issueList)
	if err != nil {
		return nil, err
	}

	return &issueList, nil
}
