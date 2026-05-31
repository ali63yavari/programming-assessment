package main

import (
	"context"
	"fmt"
	"jira_crawler/issue_crawler"
	"jira_crawler/issue_crawler/jira"
	"log"
)

const issueURL = "https://issues.apache.org/jira/browse/CAMEL-10597"
const issueURL2 = "https://issues.apache.org/jira/browse/CAMEL-23239"

func main() {
	crawler, err := jira.NewJiraIssueCrawler(
		context.Background(),
		issueURL2,
	)

	if err != nil {
		panic(err)
	}

	if err := crawler.Crawl(); err != nil {
		log.Fatal(err)
	}

	var issue = crawler.GetJiraIssue()

	err = crawler.Export(issue_crawler.ExportCSV, "output/test.csv")
	err = crawler.Export(issue_crawler.ExportExcel, "output/test.xlsx")

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
