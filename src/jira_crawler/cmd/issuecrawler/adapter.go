package main

import (
	"context"
	"jira_crawler/cli_shell/model"
	"jira_crawler/crawlerengine"
	"jira_crawler/crawlerengine/jira"
	"jira_crawler/utils"
)

type CrawlerEngineAdapter struct {
	engine *jira.IssueCrawler
}

func NewCrawlerEngineAdapter() (*CrawlerEngineAdapter, error) {
	engine, err := jira.NewJiraIssueCrawler()
	if err != nil {
		return nil, err
	}
	return &CrawlerEngineAdapter{
		engine: engine,
	}, nil
}

func (a *CrawlerEngineAdapter) Crawl(
	ctx context.Context,
	url string,
) (*model.JiraIssue, error) {
	ji, err := a.engine.Crawl(ctx, url)
	if err != nil {
		return nil, err
	}

	pc, _ := utils.ParseStringToTimeAndEpoch(ji.Created)
	pu, _ := utils.ParseStringToTimeAndEpoch(ji.Updated)
	pr, _ := utils.ParseStringToTimeAndEpoch(ji.Resolved)
	var comments []model.JiraComment

	for _, c := range ji.Comments {
		comments = append(comments, fromEngineComment(c))
	}

	return &model.JiraIssue{
		Key:           ji.Key,
		Summary:       ji.Summary,
		Type:          ji.Type,
		Status:        ji.Status,
		Priority:      ji.Priority,
		Resolution:    ji.Resolution,
		Assignee:      ji.Assignee,
		Reporter:      ji.Reporter,
		Created:       ji.Created,
		CreatedEpoch:  pc.Epoch,
		Updated:       ji.Updated,
		UpdatedEpoch:  pu.Epoch,
		Resolved:      ji.Resolved,
		ResolvedEpoch: pr.Epoch,
		Description:   ji.Description,
		Comments:      comments,
	}, nil
}

func (a *CrawlerEngineAdapter) ExportJiraIssueToCSVTemplate(
	outputPath string,
	output any,
) error {
	return a.engine.Export(crawlerengine.ExportCSV, outputPath, output)
}

func (a *CrawlerEngineAdapter) ExportJiraIssueToXLSXTemplate(
	outputPath string,
	output any,
) error {
	return a.engine.Export(crawlerengine.ExportExcel, outputPath, output)
}

func fromEngineComment(comment jira.Comment) model.JiraComment {
	pc, _ := utils.ParseStringToTimeAndEpoch(comment.Created)
	return model.JiraComment{
		ID:           comment.ID,
		Author:       comment.Author,
		Created:      comment.Created,
		CreatedEpoch: pc.Epoch,
		Body:         comment.Body,
	}
}
