package model

type JiraIssue struct {
	Key           string
	Summary       string
	Type          string
	Status        string
	Priority      string
	Resolution    string
	Assignee      string
	Reporter      string
	Created       string
	CreatedEpoch  int64
	Updated       string
	UpdatedEpoch  int64
	Resolved      string
	ResolvedEpoch int64
	Description   string
	Comments      []JiraComment
}

type JiraComment struct {
	ID           string `json:"id"`
	Author       string `json:"author"`
	Created      string `json:"created"`
	CreatedEpoch int64  `json:"created_epoch"`
	Body         string `json:"body"`
}
