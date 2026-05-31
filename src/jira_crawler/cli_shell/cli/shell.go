package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"jira_crawler/cli_shell/app"
	"os"
	"strings"
)

type Shell struct {
	service *app.Service
}

func NewShell(service *app.Service) *Shell {
	return &Shell{
		service: service,
	}
}

func (s *Shell) Run(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Issue Crawler CLI")
	fmt.Println("Type 'help' to see available commands.")
	fmt.Println()

	for {
		fmt.Print("crawler> ")

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if line == "exit" || line == "quit" {
			return nil
		}

		if err := s.handleCommand(ctx, line); err != nil {
			fmt.Println("error:", err)
		}
	}

	return scanner.Err()
}

func (s *Shell) handleCommand(ctx context.Context, line string) error {
	args := strings.Fields(line)

	if len(args) == 0 {
		return nil
	}

	switch args[0] {
	case "help":
		printHelp()
		return nil

	case "crawler":
		return s.handleCrawlerCommand(args[1:])

	case "session":
		return s.handleSessionCommand(ctx, args[1:])

	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func (s *Shell) handleCrawlerCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: crawler list | crawler show <id>")
	}

	switch args[0] {
	case "list":
		definitions := s.service.ListCrawlerDefinitions()

		if len(definitions) == 0 {
			fmt.Println("no crawler definitions registered")
			return nil
		}

		for _, d := range definitions {
			fmt.Printf("- %s | %s | %s\n", d.ID, d.Kind, d.Name)
		}

		return nil

	case "show":
		if len(args) < 2 {
			return fmt.Errorf("usage: crawler show <id>")
		}

		definition, err := s.service.ShowCrawlerDefinition(args[1])
		if err != nil {
			return err
		}

		printJSON(definition)
		return nil

	default:
		return fmt.Errorf("unknown crawler command %q", args[0])
	}
}

func (s *Shell) handleSessionCommand(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: session add|update|remove|list|show|run|rerun|export")
	}

	switch args[0] {
	case "add":
		values := parseFlags(args[1:])

		name := values["name"]
		crawlerID := values["crawler"]
		url := values["url"]

		if err := s.service.AddSession(name, crawlerID, url); err != nil {
			return err
		}

		fmt.Println("session added:", name)
		return nil

	case "update":
		if len(args) < 2 {
			return fmt.Errorf("usage: session update <name> --name <new-name> --url <new-url>")
		}

		sessionName := args[1]
		values := parseFlags(args[2:])

		if err := s.service.UpdateSession(
			sessionName,
			values["name"],
			values["url"],
		); err != nil {
			return err
		}

		fmt.Println("session updated:", sessionName)
		return nil

	case "remove":
		if len(args) < 2 {
			return fmt.Errorf("usage: session remove <name>")
		}

		if err := s.service.RemoveSession(args[1]); err != nil {
			return err
		}

		fmt.Println("session removed:", args[1])
		return nil

	case "list":
		sessions := s.service.ListSessions()

		if len(sessions) == 0 {
			fmt.Println("no sessions registered")
			return nil
		}

		for _, session := range sessions {
			fmt.Printf(
				"- %s | crawler=%s | status=%s | url=%s\n",
				session.Name,
				session.CrawlerDefinitionID,
				session.Status,
				session.URL,
			)
		}

		return nil

	case "show":
		if len(args) < 2 {
			return fmt.Errorf("usage: session show <name>")
		}

		session, err := s.service.ShowSession(args[1])
		if err != nil {
			return err
		}

		printJSON(session)
		return nil

	case "run":
		if len(args) < 2 {
			return fmt.Errorf("usage: session run <name>")
		}

		if err := s.service.RunSession(ctx, args[1]); err != nil {
			return err
		}

		fmt.Println("session run completed:", args[1])
		return nil

	case "rerun":
		if len(args) < 2 {
			return fmt.Errorf("usage: session rerun <name>")
		}

		if err := s.service.RerunSession(ctx, args[1]); err != nil {
			return err
		}

		fmt.Println("session rerun completed:", args[1])
		return nil

	case "export":
		if len(args) < 2 {
			return fmt.Errorf("usage: session export <name> --format csv|xlsx --out <path>")
		}

		sessionName := args[1]
		values := parseFlags(args[2:])

		format := values["format"]
		outputPath := values["out"]

		if format == "" {
			return fmt.Errorf("export format is required")
		}

		if outputPath == "" {
			return fmt.Errorf("output path is required")
		}

		if err := s.service.ExportSession(
			sessionName,
			format,
			outputPath,
		); err != nil {
			return err
		}

		fmt.Println("export completed:", outputPath)
		return nil

	default:
		return fmt.Errorf("unknown session command %q", args[0])
	}
}

func parseFlags(args []string) map[string]string {
	values := make(map[string]string)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if !strings.HasPrefix(arg, "--") {
			continue
		}

		key := strings.TrimPrefix(arg, "--")

		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
			values[key] = args[i+1]
			i++
		} else {
			values[key] = "true"
		}
	}

	return values
}

func printJSON(value any) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println(value)
		return
	}

	fmt.Println(string(data))
}

func printHelp() {
	fmt.Println(
		`
Crawler definition commands:
  crawler list
      List predefined crawler definitions.
  crawler show <id>
      Show one predefined crawler definition.

Session commands:
  session add --name <name> --crawler <crawler-id> --url <url>
      Add a new crawl session.
  session update <name> --name <new-name> --url <new-url>
      Update a session name and/or URL.
  session remove <name>
      Remove a session.
  session list
      List all sessions.
  session show <name>
      Show session details and latest result.
  session run <name>
      Run a session for the first time or refresh its result.
  session rerun <name>
      Rerun the session and overwrite the latest result.
  session export <name> --format csv --out <output-path>
      Export latest successful session result to CSV.
  session export <name> --format xlsx --out <output-path>
      Export latest successful session result to XLSX.

  exit
      Exit the CLI.

Example:

  crawler list
  session add --name CAMEL-10597 --crawler jira-issue --url https://issues.apache.org/jira/browse/CAMEL-10597
  session run CAMEL-10597
  session export CAMEL-10597 --format csv --out output/camel_10597.csv
  session export CAMEL-10597 --format xlsx --out output/camel_10597.xlsx
`,
	)
}
