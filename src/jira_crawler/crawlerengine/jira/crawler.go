package jira

import (
	"context"
	"fmt"
	"jira_crawler/crawlerengine"
)

type IssueCrawler struct {
	crawlerengine.IssueCrawler
	crawlerengine.IssueCrawlerExporter
}

func NewJiraIssueCrawler() (
	*IssueCrawler, error,
) {
	ic, err := crawlerengine.NewIssueCrawler()
	if err != nil {
		return nil, err
	}

	jic := &IssueCrawler{
		IssueCrawler: *ic,
	}

	jic.IssueCrawlerExporter = jic

	return jic, nil
}

func (jic *IssueCrawler) Crawl(ctx context.Context, url string) (JiraIssue, error) {
	ji := JiraIssue{}

	_, err := jic.IssueCrawler.Crawl(ctx, url, &ji)
	if err != nil {
		return JiraIssue{}, nil
	}

	return ji, nil
}

func (jic *IssueCrawler) Export(
	exportType crawlerengine.ExportType,
	outputPath string,
	output any,
) error {

	ji, ok := output.(JiraIssue)
	if !ok {
		ji = JiraIssue{}
	}
	switch exportType {
	case crawlerengine.ExportCSV:
		return ExportJiraIssueToCSVTemplate(
			"templates/jira_template.csv", outputPath,
			ji,
		)
	case crawlerengine.ExportExcel:
		return ExportJiraIssueToXLSXTemplate(
			"templates/jira_template.xlsx",
			outputPath,
			ji,
		)
	default:
		return fmt.Errorf("export type %s not supported", exportType)
	}
}
