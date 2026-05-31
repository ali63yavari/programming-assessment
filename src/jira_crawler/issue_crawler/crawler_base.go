package issue_crawler

import (
	"context"
	"fmt"
	"jira_crawler/structquery"
	"reflect"
	"time"
)

type IssueCrawler struct {
	url      string
	crawler  structquery.Crawler
	model    any
	synced   bool
	lastSync time.Time
	ctx      context.Context
}

func (ic *IssueCrawler) Crawl() error {
	err := ic.crawler.Crawl(ic.ctx, ic.model)
	if err != nil {
		ic.synced = false
		return err
	}

	ic.synced = true
	ic.lastSync = time.Now()
	return nil
}

func (ic *IssueCrawler) GetModel() any {
	return ic.model
}

func (ic *IssueCrawler) IsSynced() bool {
	return ic.synced
}

func (ic *IssueCrawler) LastSync() time.Time {
	return ic.lastSync
}

func NewIssueCrawler(ctx context.Context, url string, model any) (
	*IssueCrawler,
	error,
) {
	if model == nil {
		return nil, fmt.Errorf("model should be a valid pointer of struct")
	}
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("model should be a valid pointer of struct")
	}

	if v.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("model should be a valid pointer of struct")
	}

	cr, err := structquery.NewCrawler(url)
	if err != nil {
		return nil, err
	}

	return &IssueCrawler{
		url:      url,
		crawler:  cr,
		ctx:      ctx,
		model:    model,
		synced:   false,
		lastSync: time.Time{},
	}, nil
}
