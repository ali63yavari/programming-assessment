package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

const url string = "https://issues.apache.org/jira/browse/CAMEL-10597"

func CleanString(text string) string {
	s := strings.Map(
		func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			if unicode.IsControl(r) {
				return -1
			}

			return r
		},
		text,
	)

	return s
}
func ExampleScrape() {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("container:", doc.Find("#issue_actions_container").Length())

	fmt.Println("comments any:", doc.Find(".activity-comment").Length())

	fmt.Println(
		"comments descendant:",
		doc.Find("#issue_actions_container .activity-comment").Length(),
	)

	fmt.Println(
		"comments direct:",
		doc.Find("#issue_actions_container > .activity-comment").Length(),
	)

}

type JiraIssue struct {
	Comments []Comment `sq:"selector=#issue_actions_container > .activity-comment"`
}

type Comment struct {
	Author  string `sq:"selector=.user-hover"`
	Created string `sq:"selector=time; mode=attr; attr=datetime"`
	Body    string `sq:"selector=.action-body"`
}

func main() {
	ExampleScrape()
}
