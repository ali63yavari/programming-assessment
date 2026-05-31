package main

import (
	"context"
	"fmt"
	"jira_crawler/crawlerengine"
	"jira_crawler/crawlerengine/jira"
	"log"
)

const issueURL = "https://issues.apache.org/jira/browse/CAMEL-10597"
const issueURL2 = "https://issues.apache.org/jira/browse/CAMEL-23239"

func main() {
	crawler, err := jira.NewJiraIssueCrawler()

	if err != nil {
		panic(err)
	}

	issue, err := crawler.Crawl(context.Background(), issueURL)
	if err != nil {
		log.Fatal(err)
	}

	err = crawler.Export(crawlerengine.ExportCSV, "output/test.csv", issue)
	err = crawler.Export(crawlerengine.ExportExcel, "output/test.xlsx", issue)

	fmt.Println("")

	fmt.Println("Key:", issue.Key)
	fmt.Println("Summary:", issue.Summary)
	fmt.Println("Type:", issue.Type)
	fmt.Println("Status:", issue.Status)
	fmt.Println("Priority:", issue.Priority)
	fmt.Println("Resolution:", issue.Resolution)
	fmt.Println("Assignee:", issue.Assignee)
	fmt.Println("Reporter:", issue.Reporter)
	fmt.Println("Created:", issue.Created)
	fmt.Println("Updated:", issue.Updated)
	fmt.Println("Resolved:", issue.Resolved)
	fmt.Println("Description:", issue.Description)
	fmt.Println("Comments:", len(issue.Comments))

	for i, c := range issue.Comments {
		fmt.Println()
		fmt.Println("Comment:", i+1)
		fmt.Println("ID:", c.ID)
		fmt.Println("Author:", c.Author)
		fmt.Println("Created:", c.Created)
		fmt.Println("Body:", c.Body)
	}
}
