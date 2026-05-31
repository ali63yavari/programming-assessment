package jira

import (
	"context"
	"fmt"
	"jira_crawler/issue_crawler"
)

type JiraIssueCrawler struct {
	issue_crawler.IssueCrawler
	issue_crawler.IssueCrawlerExporter
}

func NewJiraIssueCrawler(ctx context.Context, url string) (
	*JiraIssueCrawler, error,
) {

	ji := JiraIssue{}
	ic, err := issue_crawler.NewIssueCrawler(ctx, url, &ji)
	if err != nil {
		return nil, err
	}

	jic := &JiraIssueCrawler{
		IssueCrawler: *ic,
	}

	jic.IssueCrawlerExporter = jic

	return jic, nil
}

func (jic *JiraIssueCrawler) GetJiraIssue() JiraIssue {
	m := jic.GetModel()
	issue, ok := m.(*JiraIssue)
	if ok {
		return *issue
	}
	return JiraIssue{}
}

func (jic *JiraIssueCrawler) Export(
	exportType issue_crawler.ExportType,
	outputPath string,
) error {
	switch exportType {
	case issue_crawler.ExportCSV:
		return ExportJiraIssueToCSVTemplate(
			"templates/jira_template.csv", outputPath,
			jic.GetJiraIssue(),
		)
	case issue_crawler.ExportExcel:
		return ExportJiraIssueToXLSXTemplate(
			"templates/jira_template.xlsx",
			outputPath,
			jic.GetJiraIssue(),
		)
	default:
		return fmt.Errorf("export type %s not supported", exportType)
	}
}
