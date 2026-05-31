package model

import "time"

type CrawlerDefinition struct {
	ID          string
	Name        string
	Kind        string
	Description string
	CreatedAt   time.Time
}

type CrawlStatus string

const (
	CrawlStatusPending   CrawlStatus = "pending"
	CrawlStatusRunning   CrawlStatus = "running"
	CrawlStatusSucceeded CrawlStatus = "succeeded"
	CrawlStatusFailed    CrawlStatus = "failed"
)

type CrawlSession struct {
	ID                  string
	Name                string
	CrawlerDefinitionID string
	URL                 string
	Status              CrawlStatus
	CreatedAt           time.Time
	UpdatedAt           time.Time
	LastRunAt           *time.Time
	LastFinishedAt      *time.Time
	LastError           string
	Result              any
}
