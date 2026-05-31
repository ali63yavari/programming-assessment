package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const url = "https://issues.apache.org/jira/browse/CAMEL-10597"

func main() {
	bodyBytes := step1DownloadHTML()
	step2CheckRawHTML(bodyBytes)

	doc := step3BuildGoqueryDocument(bodyBytes)
	step4CheckBasicSelectors(doc)
	step5CheckActivitySelectors(doc)
	step6CheckScripts(doc)

	activityHTML, ok := step7ExtractBigPipeActivityHTML(doc)
	if !ok {
		fmt.Println("\nSTOP: BigPipe activity HTML was not extracted.")
		return
	}

	activityDoc := step8ParseActivityHTML(activityHTML)
	step9CheckActivityDocSelectors(activityDoc)
	step10ExtractComments(activityDoc)
}

func step1DownloadHTML() []byte {
	fmt.Println("\n========== STEP 1: DOWNLOAD HTML ==========")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set(
		"Accept",
		"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	)
	req.Header.Set("Accept-Language", "en-GB,en;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("HTTP status:", resp.Status)
	fmt.Println("Downloaded bytes:", len(bodyBytes))

	return bodyBytes
}

func step2CheckRawHTML(bodyBytes []byte) {
	fmt.Println("\n========== STEP 2: CHECK RAW HTML ==========")

	body := string(bodyBytes)

	checks := []string{
		"<title>",
		"CAMEL-10597",
		"issue_actions_container",
		"activity-comment",
		"comment-15748543",
		`WRM._unparsedData["activity-panel-pipe-id"]`,
	}

	for _, check := range checks {
		fmt.Printf("contains %-45q => %v\n", check, strings.Contains(body, check))
	}

	if idx := strings.Index(body, "activity-comment"); idx >= 0 {
		start := max(0, idx-150)
		end := min(len(body), idx+250)

		fmt.Println("\nRaw HTML snippet around first activity-comment:")
		fmt.Println(body[start:end])
	}
}

func step3BuildGoqueryDocument(bodyBytes []byte) *goquery.Document {
	fmt.Println("\n========== STEP 3: BUILD GOQUERY DOCUMENT ==========")

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Document created: true")
	fmt.Println("Title:", strings.TrimSpace(doc.Find("title").Text()))
	fmt.Println("Script count:", doc.Find("script").Length())

	return doc
}

func step4CheckBasicSelectors(doc *goquery.Document) {
	fmt.Println("\n========== STEP 4: CHECK BASIC ISSUE SELECTORS ==========")

	checks := []string{
		"#key-val",
		"#summary-val",
		"#type-val",
		"#status-val",
		"#priority-val",
		"#resolution-val",
		"#description-val",
		"#assignee-val",
		"#reporter-val",
		"#created-val time",
	}

	for _, selector := range checks {
		sel := doc.Find(selector)
		fmt.Printf("%-35s => %d", selector, sel.Length())

		if sel.Length() > 0 {
			fmt.Printf(" | text=%q", cleanText(sel.First().Text()))
		}

		fmt.Println()
	}
}

func step5CheckActivitySelectors(doc *goquery.Document) {
	fmt.Println("\n========== STEP 5: CHECK ACTIVITY SELECTORS IN MAIN DOCUMENT ==========")

	checks := []string{
		"#activitymodule",
		"#issue_actions_container",
		".activity-comment",
		"#issue_actions_container .activity-comment",
		"#issue_actions_container > .activity-comment",
		".issue-data-block",
		".actionContainer",
		".action-body",
		".user-hover",
		"time.livestamp",
	}

	for _, selector := range checks {
		fmt.Printf("%-50s => %d\n", selector, doc.Find(selector).Length())
	}
}

func step6CheckScripts(doc *goquery.Document) {
	fmt.Println("\n========== STEP 6: CHECK WHETHER COMMENTS ARE INSIDE SCRIPT ==========")

	found := false

	doc.Find("script").Each(
		func(i int, s *goquery.Selection) {
			text := s.Text()

			if strings.Contains(text, "activity-comment") {
				found = true
				fmt.Println("activity-comment found inside script index:", i)
				fmt.Println("script length:", len(text))

				idx := strings.Index(text, "activity-comment")
				start := max(0, idx-150)
				end := min(len(text), idx+250)

				fmt.Println("\nScript snippet around first activity-comment:")
				fmt.Println(text[start:end])
			}
		},
	)

	fmt.Println("activity-comment found inside any script:", found)
}

func step7ExtractBigPipeActivityHTML(doc *goquery.Document) (string, bool) {
	fmt.Println("\n========== STEP 7: EXTRACT BIGPIPE ACTIVITY HTML ==========")

	var scriptText string

	doc.Find("script").EachWithBreak(
		func(i int, s *goquery.Selection) bool {
			text := s.Text()

			if strings.Contains(
				text,
				`WRM._unparsedData["activity-panel-pipe-id"]`,
			) {
				scriptText = text
				fmt.Println("Found activity-panel-pipe-id script index:", i)
				fmt.Println("Script length:", len(scriptText))
				return false
			}

			return true
		},
	)

	if scriptText == "" {
		fmt.Println("activity-panel-pipe-id script not found")
		return "", false
	}

	key := `WRM._unparsedData["activity-panel-pipe-id"]="`

	start := strings.Index(scriptText, key)
	if start == -1 {
		fmt.Println("Assignment start not found")
		return "", false
	}

	start += len(key)

	end := strings.Index(scriptText[start:], `";`)
	if end == -1 {
		fmt.Println("Assignment end not found")
		return "", false
	}

	escapedHTML := scriptText[start : start+end]

	fmt.Println("Escaped activity HTML length:", len(escapedHTML))
	fmt.Println(
		"Escaped contains activity-comment:",
		strings.Contains(escapedHTML, "activity-comment"),
	)
	fmt.Println(
		"Escaped contains issue_actions_container:",
		strings.Contains(escapedHTML, "issue_actions_container"),
	)
	fmt.Println("Escaped contains \\/:", strings.Contains(escapedHTML, `\/`))
	fmt.Println("Escaped contains \\\\/:", strings.Contains(escapedHTML, `\\/`))
	fmt.Println("Escaped contains \\':", strings.Contains(escapedHTML, `\'`))
	fmt.Println("Escaped contains \\u003c:", strings.Contains(escapedHTML, `\u003c`))
	fmt.Println(
		"Escaped contains \\\\u003c:",
		strings.Contains(escapedHTML, `\\u003c`),
	)
	fmt.Println("Escaped contains \\n:", strings.Contains(escapedHTML, `\n`))
	fmt.Println("Escaped contains \\\\n:", strings.Contains(escapedHTML, `\\n`))

	current := escapedHTML

	for pass := 1; pass <= 4; pass++ {
		fmt.Printf("\nTrying custom JavaScript unescape pass %d\n", pass)

		decoded, err := jsUnescape(current)
		if err != nil {
			fmt.Println("decode failed:", err)
			return "", false
		}

		fmt.Println("Decoded length:", len(decoded))
		fmt.Println(
			"Decoded contains activity-comment:",
			strings.Contains(decoded, "activity-comment"),
		)
		fmt.Println(
			"Decoded contains issue_actions_container:",
			strings.Contains(decoded, "issue_actions_container"),
		)
		fmt.Println(
			"Decoded contains literal <div:",
			strings.Contains(decoded, "<div"),
		)
		fmt.Println(
			"Decoded contains escaped \\u003c:",
			strings.Contains(decoded, `\u003c`),
		)
		fmt.Println(
			"Decoded contains escaped \\\\u003c:",
			strings.Contains(decoded, `\\u003c`),
		)
		fmt.Println(
			"Decoded contains escaped \\\" :",
			strings.Contains(decoded, `\"`),
		)
		fmt.Println(
			"Decoded contains escaped \\\\\" :",
			strings.Contains(decoded, `\\"`),
		)

		if idx := strings.Index(decoded, "activity-comment"); idx >= 0 {
			snippetStart := max(0, idx-150)
			snippetEnd := min(len(decoded), idx+250)

			fmt.Println("\nDecoded snippet around first activity-comment:")
			fmt.Println(decoded[snippetStart:snippetEnd])
		}

		// Test if this pass is already parseable as real HTML.
		testDoc, err := goquery.NewDocumentFromReader(strings.NewReader(decoded))
		if err != nil {
			fmt.Println("goquery parse failed in this pass:", err)
		} else {
			fmt.Println(
				"Test selector .activity-comment:",
				testDoc.Find(".activity-comment").Length(),
			)
			fmt.Println(
				"Test selector #issue_actions_container:",
				testDoc.Find("#issue_actions_container").Length(),
			)

			if testDoc.Find(".activity-comment").Length() > 0 {
				fmt.Println("SUCCESS: activity comments are now real DOM nodes")
				return decoded, true
			}
		}

		current = decoded
	}

	fmt.Println("Could not decode activity HTML into selectable DOM after 4 passes")
	return "", false
}

func jsUnescape(input string) (string, error) {
	var b strings.Builder

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if ch != '\\' {
			b.WriteByte(ch)
			continue
		}

		if i+1 >= len(input) {
			b.WriteByte('\\')
			continue
		}

		i++
		next := input[i]

		switch next {
		case '\\':
			b.WriteByte('\\')

		case '"':
			b.WriteByte('"')

		case '\'':
			b.WriteByte('\'')

		case '/':
			b.WriteByte('/')

		case 'n':
			b.WriteByte('\n')

		case 'r':
			b.WriteByte('\r')

		case 't':
			b.WriteByte('\t')

		case 'b':
			b.WriteByte('\b')

		case 'f':
			b.WriteByte('\f')

		case 'u':
			if i+4 >= len(input) {
				return "", fmt.Errorf("invalid unicode escape at index %d", i-1)
			}

			hexPart := input[i+1 : i+5]

			r, err := strconv.ParseInt(hexPart, 16, 32)
			if err != nil {
				return "", fmt.Errorf(
					"invalid unicode escape \\u%s at index %d: %w",
					hexPart,
					i-1,
					err,
				)
			}

			b.WriteRune(rune(r))
			i += 4

		case 'x':
			if i+2 >= len(input) {
				return "", fmt.Errorf("invalid hex escape at index %d", i-1)
			}

			hexPart := input[i+1 : i+3]

			r, err := strconv.ParseInt(hexPart, 16, 8)
			if err != nil {
				return "", fmt.Errorf(
					"invalid hex escape \\x%s at index %d: %w",
					hexPart,
					i-1,
					err,
				)
			}

			b.WriteByte(byte(r))
			i += 2

		default:
			// JavaScript allows some non-standard escape behavior.
			// Instead of failing, keep the escaped character.
			b.WriteByte(next)
		}
	}

	return b.String(), nil
}

func errorIndexDebug(value string) {

	fmt.Println("\nDebugging normalized escaped string...")
	for i := 0; i < len(value); i++ {
		if value[i] != '\\' {
			continue
		}
		if i+1 >= len(value) {
			fmt.Println("Backslash at end of string at index:", i)
			return
		}
		next := value[i+1]
		valid := strings.ContainsRune(`abfnrtv\"`, rune(next)) ||
			next == 'u' ||
			next == 'U' ||
			next == 'x' ||
			(next >= '0' && next <= '7')
		if !valid {
			start := max(0, i-80)
			end := min(len(value), i+120)
			fmt.Println("Possibly invalid Go escape at index:", i)
			fmt.Println("Escape sequence:", value[i:i+2])
			fmt.Println("Context:")
			fmt.Println(value[start:end])
			return
		}
	}
	fmt.Println("No obvious invalid escape found")

}
func step8ParseActivityHTML(activityHTML string) *goquery.Document {
	fmt.Println("\n========== STEP 8: PARSE ACTIVITY HTML AS NEW DOCUMENT ==========")

	activityDoc, err := goquery.NewDocumentFromReader(strings.NewReader(activityHTML))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Activity document created: true")
	fmt.Println("Activity document text length:", len(activityDoc.Text()))

	return activityDoc
}

func step9CheckActivityDocSelectors(activityDoc *goquery.Document) {
	fmt.Println("\n========== STEP 9: CHECK SELECTORS IN ACTIVITY DOCUMENT ==========")

	checks := []string{
		"#issue_actions_container",
		".activity-comment",
		"#issue_actions_container .activity-comment",
		"#issue_actions_container > .activity-comment",
		".issue-data-block",
		".actionContainer",
		".twixi-wrap.verbose",
		".twixi-wrap.concise",
		".twixi-wrap.verbose .action-details > a.user-hover",
		".twixi-wrap.verbose time.livestamp",
		".twixi-wrap.verbose > .action-body",
	}

	for _, selector := range checks {
		fmt.Printf("%-65s => %d\n", selector, activityDoc.Find(selector).Length())
	}
}

func step10ExtractComments(activityDoc *goquery.Document) {
	fmt.Println("\n========== STEP 10: EXTRACT COMMENTS ==========")

	comments := activityDoc.Find(".activity-comment")

	fmt.Println("Total comments:", comments.Length())

	comments.Each(
		func(i int, s *goquery.Selection) {
			id, _ := s.Attr("id")

			author := cleanText(
				s.Find(".twixi-wrap.verbose .action-details > a.user-hover").First().Text(),
			)

			created, _ := s.Find(".twixi-wrap.verbose time.livestamp").First().Attr("datetime")

			body := cleanText(
				s.Find(".twixi-wrap.verbose > .action-body").First().Text(),
			)

			fmt.Println("\n----- COMMENT", i+1, "-----")
			fmt.Println("ID:", id)
			fmt.Println("Author:", author)
			fmt.Println("Created:", created)
			fmt.Println("Body:", body)
		},
	)
}

func cleanText(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Join(strings.Fields(value), " ")
	return value
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
