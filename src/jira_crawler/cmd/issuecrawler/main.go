package main

import (
	"context"
	"jira_crawler/cli_shell/app"
	"jira_crawler/cli_shell/cli"
	"jira_crawler/cli_shell/crawler"
	"jira_crawler/cli_shell/store"
	"log"
)

func main() {
	memoryStore := store.NewMemoryStore()
	service := app.NewService(memoryStore)

	adapter, err := NewCrawlerEngineAdapter()
	if err != nil {
		panic(err)
	}

	jiraEngine := crawler.NewJiraIssueEngine(
		adapter,
		adapter,
		adapter,
	)

	if err := service.RegisterCrawlerEngine(jiraEngine); err != nil {
		log.Fatal(err)
	}

	shell := cli.NewShell(service)

	if err := shell.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
