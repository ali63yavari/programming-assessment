package crawler

import (
	"context"
	"fmt"
	"jira_crawler/cli_shell/model"
	"time"
)

type JiraIssueCrawler interface {
	Crawl(ctx context.Context, url string) (*model.JiraIssue, error)
}

type JiraCSVExporter interface {
	ExportJiraIssueToCSVTemplate(
		outputPath string,
		output any,
	) error
}

type JiraXLSXExporter interface {
	ExportJiraIssueToXLSXTemplate(
		outputPath string,
		output any,
	) error
}

type JiraIssueEngine struct {
	crawler      JiraIssueCrawler
	csvExporter  JiraCSVExporter
	xlsxExporter JiraXLSXExporter
}

func NewJiraIssueEngine(
	crawler JiraIssueCrawler,
	csvExporter JiraCSVExporter,
	xlsxExporter JiraXLSXExporter,
) *JiraIssueEngine {
	return &JiraIssueEngine{
		crawler:      crawler,
		csvExporter:  csvExporter,
		xlsxExporter: xlsxExporter,
	}
}

func (e *JiraIssueEngine) Definition() model.CrawlerDefinition {
	return model.CrawlerDefinition{
		ID:          "jira-issue",
		Name:        "Jira Issue Crawler",
		Kind:        "jira_issue",
		Description: "Predefined crawler for Apache Jira issue pages. Extracts issue details, dates, people, description, and comments.",
		CreatedAt:   time.Now(),
	}
}

func (e *JiraIssueEngine) Crawl(ctx context.Context, url string) (any, error) {
	if e.crawler == nil {
		return nil, fmt.Errorf("jira crawler is not configured")
	}

	result, err := e.crawler.Crawl(ctx, url)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (e *JiraIssueEngine) ExportCSV(
	outputPath string,
	output any,
) error {
	if e.csvExporter == nil {
		return fmt.Errorf("jira csv exporter is not configured")
	}

	return e.csvExporter.ExportJiraIssueToCSVTemplate(
		outputPath,
		output,
	)
}

func (e *JiraIssueEngine) ExportXLSX(
	outputPath string,
	output any,
) error {
	if e.xlsxExporter == nil {
		return fmt.Errorf("jira xlsx exporter is not configured")
	}

	return e.xlsxExporter.ExportJiraIssueToXLSXTemplate(
		outputPath,
		output,
	)
}
