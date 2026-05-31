package jira

import (
	"jira_crawler/structquery"
	"time"
)

type JiraIssue struct {
	Key         string    `sq:"selector=#key-val"`
	Summary     string    `sq:"selector=#summary-val"`
	Type        string    `sq:"selector=#type-val"`
	Status      string    `sq:"selector=#status-val"`
	Priority    string    `sq:"selector=#priority-val"`
	Resolution  string    `sq:"selector=#resolution-val"`
	Assignee    string    `sq:"selector=#assignee-val"`
	Reporter    string    `sq:"selector=#reporter-val"`
	Created     string    `sq:"selector=#created-val time; mode=attr; attr=datetime"`
	Updated     string    `sq:"selector=#updated-val time; mode=attr; attr=datetime"`
	Resolved    string    `sq:"selector=#resolutiondate-val time; mode=attr; attr=datetime"`
	Description string    `sq:"selector=#description-val;mode=html"`
	Comments    []Comment `sq:"selector=.activity-comment"`
	Url         string
}

func (JiraIssue) RenderOptions() structquery.RenderOptions {
	b := true
	return structquery.RenderOptions{
		Timeout:      45 * time.Second,
		WaitMode:     structquery.WaitSelector,
		WaitSelector: "#issue_actions_container",
		Headless:     &b,
	}
}

type Comment struct {
	ID      string `sq:"selector=.; mode=attr; attr=id"`
	Author  string `sq:"selector=.twixi-wrap.verbose .action-details > a.user-hover"`
	Created string `sq:"selector=.twixi-wrap.verbose time.livestamp; mode=attr; attr=datetime"`
	Body    string `sq:"selector=.twixi-wrap.verbose > .action-body"`
}
