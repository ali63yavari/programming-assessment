package crawlerengine

import (
	"context"
	"jira_crawler/structquery"
	"time"
)

type IssueCrawler struct {
	crawler  structquery.Crawler
	synced   bool
	lastSync time.Time
}

func (ic *IssueCrawler) Crawl(ctx context.Context, url string, model any) (
	any,
	error,
) {
	err := ic.crawler.Crawl(ctx, url, model)
	if err != nil {
		ic.synced = false
		return nil, err
	}

	ic.synced = true
	ic.lastSync = time.Now()

	return model, nil
}

func (ic *IssueCrawler) IsSynced() bool {
	return ic.synced
}

func (ic *IssueCrawler) LastSync() time.Time {
	return ic.lastSync
}

func NewIssueCrawler() (
	*IssueCrawler,
	error,
) {
	cr, err := structquery.NewCrawler()
	if err != nil {
		return nil, err
	}

	return &IssueCrawler{
		crawler:  cr,
		synced:   false,
		lastSync: time.Time{},
	}, nil
}
