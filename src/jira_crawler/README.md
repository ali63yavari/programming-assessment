# Issue Crawler

> A Go-based programming assessment project that crawls rendered Jira issue pages, extracts structured issue data, manages crawl sessions, and exports the latest successful result to CSV or XLSX.

## Overview

Issue Crawler is a small but layered extraction tool focused on Apache Jira issue pages. The current implementation provides one predefined crawler, `jira-issue`, which renders a Jira issue page, extracts issue metadata, description, and comments, and exports the result into spreadsheet-friendly CSV and styled XLSX outputs.

The broader architectural direction is that crawler definitions represent extraction models: they describe how a class of pages should be rendered, parsed, mapped, and exported. A future version could load those definitions from YAML so users can define custom page structures and extraction rules without changing Go code. In this version, only one built-in Jira issue crawler definition is registered when the CLI starts.

The main domain concepts are:

- **Crawler Definition**: the extraction model, engine, and schema for a class of pages. The implemented definition is `jira-issue`.
- **Crawl Session**: a named target URL associated with a crawler definition.
- **Crawl Result**: the latest extracted result stored on a session after a successful run.
- **Export**: a CSV or XLSX file generated from the latest successful session result.

## Assessment Context

This project was created as part of a programming assessment. The implementation emphasizes correctness, maintainability, extensibility, and clear engineering boundaries.

The solution avoids placing one-off scraping logic directly in the CLI. Instead, the code separates rendering, extraction, session management, crawler orchestration, and export generation into distinct packages. This makes the current Jira crawler easier to review and also leaves a practical path toward additional crawler definitions later.

## Key Features

- Predefined `jira-issue` crawler registered at application startup.
- Interactive CLI shell for managing crawler definitions and crawl sessions.
- In-memory registry for crawler engines and sessions.
- Session lifecycle operations: add, update, remove, list, show, run, rerun, and export.
- Rod-based browser rendering for dynamic Jira pages.
- goquery-based parsing of the final rendered HTML.
- Reflection and `sq` struct-tag based extraction.
- Extraction support for scalar fields, attributes, HTML content, nested structs, and slices of structs.
- Jira issue extraction for key, summary, type, status, priority, resolution, assignee, reporter, dates, description, and comments.
- Comment extraction from Jira activity comment blocks.
- CSV export using a built-in CSV dashboard template.
- XLSX export using a built-in styled workbook template.
- Date parsing and Unix epoch conversion for issue and comment timestamps.
- Runnable example under `examples/jira_issue`.

## Architecture

The project is organized into focused packages:

- `cmd/issuecrawler`: application entry point and adapter between the CLI shell engine interface and the Jira crawler implementation.
- `cli_shell/cli`: interactive command parser and shell.
- `cli_shell/app`: application service coordinating definitions, sessions, crawls, reruns, and exports.
- `cli_shell/store`: in-memory storage for crawler engines and crawl sessions.
- `cli_shell/crawler`: shell-facing crawler engine interfaces and Jira engine wrapper.
- `cli_shell/model`: shell/session-facing models.
- `crawlerengine`: generic crawler and exporter abstractions.
- `crawlerengine/jira`: Jira issue model, crawler wrapper, and CSV/XLSX exporters.
- `structquery`: Rod renderer, goquery parsing, render options, tag parsing, and reflection mapper.
- `utils`: date parsing and epoch conversion helpers.
- `templates`: built-in CSV and XLSX export templates.
- `examples`: runnable demonstration code.

## Why Rod Rendering Is Used

Apache Jira pages can render sections dynamically after the initial HTML response. The activity and comments area is especially important for this project, and static HTML parsing alone may not see the final page state.

Rod is used as a rendering layer. It launches a browser, navigates to the issue URL, waits according to the model's render options, and returns the final page HTML. After rendering, the project still uses goquery and the reflection-based mapper for extraction.

This is an intentional separation: Rod does not know which Jira fields to extract. It only provides the rendered DOM. Field extraction remains selector-driven and model-driven.

## Reflection-Based Extraction Model

The generic extraction layer in `structquery` maps HTML into Go structs using the `sq` struct tag. The Jira model demonstrates the approach:

```go
type JiraIssue struct {
    Key         string    `sq:"selector=#key-val"`
    Summary     string    `sq:"selector=#summary-val"`
    Created     string    `sq:"selector=#created-val time; mode=attr; attr=datetime"`
    Description string    `sq:"selector=#description-val;mode=html"`
    Comments    []Comment `sq:"selector=.activity-comment"`
}

type Comment struct {
    ID      string `sq:"selector=.; mode=attr; attr=id"`
    Author  string `sq:"selector=.twixi-wrap.verbose .action-details > a.user-hover"`
    Created string `sq:"selector=.twixi-wrap.verbose time.livestamp; mode=attr; attr=datetime"`
    Body    string `sq:"selector=.twixi-wrap.verbose > .action-body"`
}
```

Selectors describe where data comes from. Modes describe how data is extracted. The implemented modes are `text`, `attr`, `html`, and `outer_html`; `exists` and `count` are defined and validated but are not currently handled in the extraction switch.

The mapper supports:

- Scalar field conversion for strings, integers, unsigned integers, floats, and booleans.
- Attribute extraction with `mode=attr; attr=<name>`.
- Inner HTML extraction with `mode=html`.
- Outer HTML extraction with `mode=outer_html`.
- Nested structs.
- Slices of structs, where each selected node becomes one struct instance.
- `selector=.` to use the current selected node itself.
- Optional validation flags such as `required`, `nonempty`, `enum`, and `match`.
- Text trimming modes: `all`, `space`, `control`, and `none`.

## Crawler Definitions and Sessions

### Crawler Definition

A crawler definition describes how a class of pages should be rendered, parsed, and exported. In this version, the `jira-issue` crawler definition is predefined and registered at startup.

The registered definition has:

- ID: `jira-issue`
- Name: `Jira Issue Crawler`
- Kind: `jira_issue`

### Crawl Session

A crawl session is a user-registered URL that uses a crawler definition. Multiple sessions can use the same crawler definition.

Example:

```text
CAMEL-10597 -> https://issues.apache.org/jira/browse/CAMEL-10597
```

Supported session operations:

- Add session
- Update session name or URL
- Remove session
- List sessions
- Show session
- Run session
- Rerun session
- Export session

## CLI Usage

Run from the project root:

```bash
go mod tidy
go run ./cmd/issuecrawler
```

The command starts the interactive shell:

```text
Issue Crawler CLI
Type 'help' to see available commands.

crawler>
```

Complete example:

```text
crawler list
crawler show jira-issue
session add --name CAMEL-10597 --crawler jira-issue --url https://issues.apache.org/jira/browse/CAMEL-10597
session list
session run CAMEL-10597
session show CAMEL-10597
session export CAMEL-10597 --format csv --out output/camel_10597.csv
session export CAMEL-10597 --format xlsx --out output/camel_10597.xlsx
exit
```

The CLI export command does not accept a template path. It uses the built-in templates from the `templates` directory.

## Exporting Results

### CSV

CSV export uses the built-in `templates/jira_template.csv` file. The exporter reads the template, replaces issue placeholders, and repeats the comment template row for each extracted comment.

CSV is useful for spreadsheet tools and review, but it cannot preserve visual formatting such as colors, borders, merged cells, or fonts.

### XLSX

XLSX export uses the built-in `templates/jira_template.xlsx` workbook. The exporter replaces issue placeholders in the `Jira Issue Dashboard` sheet and writes comments into the comment table starting at row 22. If there are more than five comments, it inserts additional rows and copies styles from the template rows.

The XLSX output is intended as the more readable presentation format. It includes issue details, dates and epoch values, description, and comments.

Comments are repeated data, so both exporters handle them separately from the single-value issue fields.

## Templates

Templates live in:

```text
templates/jira_template.csv
templates/jira_template.xlsx
```

## Running Examples

The repository includes one runnable example:

```bash
go run ./examples/jira_issue
```

## Future Improvements
-[ ] YAML-based crawler definition import.
-[ ] Persistent storage such as SQLite or BoltDB.
-[ ] One-shot CLI commands after persistence is introduced.
-[ ] Crawl run history.
-[ ] Stronger crawler definition validation.
-[ ] Additional crawler definitions, such as GitHub Issues or GitLab Issues.
-[ ] Configurable render policies per crawler/session.
-[ ] Authentication and session-cookie support.
-[ ] Batch crawling.
-[ ] Parallel crawling.
-[ ] Configurable template paths.
-[ ] Richer export customization.
-[ ] Automated tests and CI.
-[ ] Propagation of lower-level crawl errors from the Jira crawler wrapper.

## Project Structure

```text
.
|-- cli_shell/
|   |-- app/
|   |   `-- service.go
|   |-- cli/
|   |   `-- shell.go
|   |-- crawler/
|   |   |-- engine.go
|   |   `-- jira_engine.go
|   |-- model/
|   |   |-- jira.go
|   |   `-- session.go
|   `-- store/
|       `-- memory.go
|-- cmd/
|   `-- issuecrawler/
|       |-- adapter.go
|       `-- main.go
|-- crawlerengine/
|   |-- crawler_base.go
|   |-- exporter_base.go
|   `-- jira/
|       |-- crawler.go
|       |-- exporter.go
|       `-- model.go
|-- examples/
|   `-- jira_issue/
|       `-- main.go
|-- output/
|   |-- test.csv
|   `-- test.xlsx
|-- structquery/
|   |-- crawler.go
|   |-- field_tag_config.go
|   |-- render_provider.go
|   |-- renderer_options.go
|   |-- rod_renderer.go
|   |-- tag_parser.go
|   `-- trim_text.go
|-- templates/
|   |-- jira_template.csv
|   `-- jira_template.xlsx
|-- utils/
|   `-- time_helper.go
|-- go.mod
|-- go.sum
`-- tt.xlsx
```

Important folders:

- `structquery` contains the generic rendering and extraction engine.
- `crawlerengine/jira` contains the Jira-specific struct tags and exporters.
- `cli_shell` contains the interactive shell, application service, store, and shell models.
- `templates` contains the built-in export templates.
- `examples/jira_issue` demonstrates direct use of the Jira crawler without the shell.

## Dependencies

Major dependencies from `go.mod`:

- `github.com/go-rod/rod`: headless browser rendering for dynamic pages.
- `github.com/PuerkitoBio/goquery`: CSS-selector based querying over rendered HTML.
- `github.com/xuri/excelize/v2`: XLSX template reading, cell updates, row insertion, styling, and saving.
- `golang.org/x/net`: HTML-related support used by the module dependency graph.

Important indirect dependencies include Cascadia for selector support through goquery and Rod support packages for browser launching and communication.

## How to Review the Project

Suggested review path:

1. Read the architecture overview in this README.
2. Review `crawlerengine/jira/model.go` to see the tag-based Jira extraction model.
3. Review `structquery/crawler.go` and `structquery/rod_renderer.go` to see rendering and reflection-based mapping.
4. Start the CLI with `go run ./cmd/issuecrawler`.
5. Register a Jira issue session with `session add`.
6. Run the crawler with `session run`.
7. Export CSV and XLSX files with `session export`.
8. Inspect the generated output files.
9. Review `crawlerengine/jira/exporter.go` to see template-based export handling.

## Final Notes
> The current implementation prioritizes clarity, separation of concerns, and extensibility. The Jira crawler is intentionally treated as the first predefined crawler definition rather than as ad hoc CLI logic. The architecture is prepared for future dynamic definitions, persistence, and additional export or crawler types, while the present version remains focused on a reviewable Jira issue crawling workflow.
