package crawler

import (
	"context"
	"jira_crawler/cli_shell/model"
)

type Engine interface {
	Definition() model.CrawlerDefinition

	Crawl(ctx context.Context, url string) (any, error)

	ExportCSV(outputPath string, output any) error

	ExportXLSX(outputPath string, output any) error
}
