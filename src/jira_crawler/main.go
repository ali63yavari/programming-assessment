package main

import (
	"fmt"
	"jira_crawler/structquery"
	"log"
	"net/http"
	"reflect"
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

	typeVal := doc.Find("#type-val").Text()
	println(CleanString(typeVal))
	//// Find the review items
	//doc.Find("#type-val").Each(
	//	func(i int, s *goquery.Selection) {
	//		// For each item found, get the title
	//		title := s.Find("a").Text()
	//		fmt.Printf("Review %d: %s\n", i, title)
	//	},
	//)
}

type nestedTst struct {
	name string `gg:"nested name, nested struct"`
}
type tst struct {
	name   string    `gg:"string,int,others"`
	age    int       `gg:"format"`
	grades []float32 `gg:"float32"`
	nested nestedTst `gg:"nested"`
}

func testReflect(out any) error {
	v := reflect.ValueOf(out)

	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("out should be a pointer to a struct")
	}

	v = v.Elem()

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("output should be type of struct")
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := t.Field(i)

		fmt.Printf(
			"field type: %v and field value: %v\n",
			fieldType,
			fieldVal.String(),
		)
	}

	return nil
}

func main() {
	ExampleScrape()
	err := testReflect(&tst{})
	println(err)

	s1 := "selector =   rfgdfg  "
	s2 := " required "

	c1 := strings.Split(s1, "=")
	c2 := strings.Split(s2, "=")

	_ = c1
	_ = c2

	v, err := structquery.ExtractInlineArray[string](
		"enum=Bug|Task|Improvement",
		"enum", "|",
	)

	_ = v
}
