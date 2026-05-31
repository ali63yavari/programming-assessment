package jira

import (
	"encoding/csv"
	"fmt"
	"jira_crawler/utils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func ExportJiraIssueToCSVTemplate(
	templatePath string,
	outputPath string,
	issue JiraIssue,
) error {
	in, err := os.Open(templatePath)
	if err != nil {
		return fmt.Errorf("open csv template: %w", err)
	}
	defer in.Close()

	reader := csv.NewReader(in)
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv template: %w", err)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create csv output: %w", err)
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	defer writer.Flush()

	placeholders := jiraIssuePlaceholders(issue)

	for _, row := range rows {
		if isCSVCommentTemplateRow(row) {
			if len(issue.Comments) == 0 {
				emptyCommentRow := replaceCSVRowPlaceholders(
					row, map[string]string{
						"{{CommentNo}}":           "",
						"{{CommentAuthor}}":       "",
						"{{CommentCreated}}":      "",
						"{{CommentCreatedEpoch}}": "",
						"{{CommentBody}}":         "",
					},
				)

				if err := writer.Write(emptyCommentRow); err != nil {
					return fmt.Errorf("write empty comment row: %w", err)
				}

				continue
			}

			for i, comment := range issue.Comments {
				pt, _ := utils.ParseStringToTimeAndEpoch(comment.Created)
				commentPlaceholders := map[string]string{
					"{{CommentNo}}":           strconv.Itoa(i + 1),
					"{{CommentAuthor}}":       cleanCell(comment.Author),
					"{{CommentCreated}}":      cleanCell(comment.Created),
					"{{CommentCreatedEpoch}}": epochCell(pt.Epoch),
					"{{CommentBody}}":         cleanCell(comment.Body),
				}

				commentRow := replaceCSVRowPlaceholders(row, commentPlaceholders)

				if err := writer.Write(commentRow); err != nil {
					return fmt.Errorf("write comment row %d: %w", i+1, err)
				}
			}

			continue
		}

		replaced := replaceCSVRowPlaceholders(row, placeholders)

		if err := writer.Write(replaced); err != nil {
			return fmt.Errorf("write csv row: %w", err)
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush csv writer: %w", err)
	}

	return nil
}

func ExportJiraIssueToXLSXTemplate(
	templatePath string,
	outputPath string,
	issue JiraIssue,
) error {
	const dashboardSheet = "Jira Issue Dashboard"

	file, err := excelize.OpenFile(templatePath)
	if err != nil {
		return fmt.Errorf("open xlsx template: %w", err)
	}
	defer file.Close()

	placeholders := jiraIssuePlaceholders(issue)

	if err := replaceSheetPlaceholders(
		file,
		dashboardSheet,
		placeholders,
	); err != nil {
		return err
	}

	if err := fillXLSXComments(file, dashboardSheet, issue.Comments); err != nil {
		return err
	}

	if err := file.SaveAs(outputPath); err != nil {
		return fmt.Errorf("save xlsx output: %w", err)
	}

	return nil
}

func fillXLSXComments(file *excelize.File, sheet string, comments []Comment) error {
	const firstCommentRow = 22
	const templateCommentRows = 5

	if len(comments) > templateCommentRows {
		extraRows := len(comments) - templateCommentRows

		for i := 0; i < extraRows; i++ {
			insertAt := firstCommentRow + templateCommentRows + i

			if err := file.InsertRows(sheet, insertAt, 1); err != nil {
				return fmt.Errorf(
					"insert extra comment row at %d: %w",
					insertAt,
					err,
				)
			}

			sourceRow := firstCommentRow + templateCommentRows - 1 + i
			if sourceRow < firstCommentRow {
				sourceRow = firstCommentRow
			}

			if err := copyRowStyle(
				file,
				sheet,
				sourceRow,
				insertAt,
				[]string{"B", "C", "D", "E", "F"},
			); err != nil {
				return fmt.Errorf("copy comment row style: %w", err)
			}
		}
	}

	for i, comment := range comments {
		row := firstCommentRow + i
		pt, _ := utils.ParseStringToTimeAndEpoch(comment.Created)
		values := map[string]any{
			fmt.Sprintf("F%d", row): i + 1,
			fmt.Sprintf("B%d", row): cleanCell(comment.Author),
			fmt.Sprintf("C%d", row): cleanCell(comment.Created),
			fmt.Sprintf("D%d", row): epochCell(pt.Epoch),
			fmt.Sprintf("E%d", row): cleanCell(comment.Body),
		}

		for cell, value := range values {
			if err := file.SetCellValue(sheet, cell, value); err != nil {
				return fmt.Errorf("set %s: %w", cell, err)
			}
		}
	}

	for i := len(comments); i < templateCommentRows; i++ {
		row := firstCommentRow + i

		for _, cell := range []string{
			fmt.Sprintf("F%d", row),
			fmt.Sprintf("B%d", row),
			fmt.Sprintf("C%d", row),
			fmt.Sprintf("D%d", row),
			fmt.Sprintf("E%d", row),
		} {
			if err := file.SetCellValue(sheet, cell, ""); err != nil {
				return fmt.Errorf("clear unused comment cell %s: %w", cell, err)
			}
		}
	}

	return nil
}

func replaceSheetPlaceholders(
	file *excelize.File,
	sheet string,
	placeholders map[string]string,
) error {
	rows, err := file.GetRows(sheet)
	if err != nil {
		return fmt.Errorf("read sheet rows %q: %w", sheet, err)
	}

	for rowIndex, row := range rows {
		for colIndex, value := range row {
			if !strings.Contains(value, "{{") {
				continue
			}

			replaced := replacePlaceholders(value, placeholders)

			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			if err != nil {
				return fmt.Errorf(
					"cell coordinate row=%d col=%d: %w",
					rowIndex+1,
					colIndex+1,
					err,
				)
			}

			if err := file.SetCellValue(sheet, cell, replaced); err != nil {
				return fmt.Errorf("set placeholder cell %s: %w", cell, err)
			}
		}
	}

	return nil
}

func copyRowStyle(
	file *excelize.File,
	sheet string,
	sourceRow int,
	targetRow int,
	columns []string,
) error {
	for _, col := range columns {
		sourceCell := fmt.Sprintf("%s%d", col, sourceRow)
		targetCell := fmt.Sprintf("%s%d", col, targetRow)

		styleID, err := file.GetCellStyle(sheet, sourceCell)
		if err != nil {
			return fmt.Errorf("get style from %s: %w", sourceCell, err)
		}

		if err := file.SetCellStyle(
			sheet,
			targetCell,
			targetCell,
			styleID,
		); err != nil {
			return fmt.Errorf("set style to %s: %w", targetCell, err)
		}
	}

	height, err := file.GetRowHeight(sheet, sourceRow)
	if err == nil && height > 0 {
		_ = file.SetRowHeight(sheet, targetRow, height)
	}

	return nil
}

func jiraIssuePlaceholders(issue JiraIssue) map[string]string {
	ptCreate, _ := utils.ParseStringToTimeAndEpoch(issue.Created)
	ptUpdated, _ := utils.ParseStringToTimeAndEpoch(issue.Updated)
	ptResolved, _ := utils.ParseStringToTimeAndEpoch(issue.Resolved)

	return map[string]string{
		"{{url}}":           cleanCell(issue.Url),
		"{{GeneratedOn}}":   time.Now().Format("2006-01-02 15:04:05"),
		"{{IssueKey}}":      cleanCell(issue.Key),
		"{{Summary}}":       cleanCell(issue.Summary),
		"{{Type}}":          cleanCell(issue.Type),
		"{{Status}}":        cleanCell(issue.Status),
		"{{Priority}}":      cleanCell(issue.Priority),
		"{{Resolution}}":    cleanCell(issue.Resolution),
		"{{Assignee}}":      cleanCell(issue.Assignee),
		"{{Reporter}}":      cleanCell(issue.Reporter),
		"{{Created}}":       cleanCell(issue.Created),
		"{{CreatedEpoch}}":  epochCell(ptCreate.Epoch),
		"{{Updated}}":       cleanCell(issue.Updated),
		"{{UpdatedEpoch}}":  epochCell(ptUpdated.Epoch),
		"{{Resolved}}":      cleanCell(issue.Resolved),
		"{{ResolvedEpoch}}": epochCell(ptResolved.Epoch),
		"{{Description}}":   cleanCell(issue.Description),
	}
}

func isCSVCommentTemplateRow(row []string) bool {
	for _, cell := range row {
		if strings.Contains(cell, "{{CommentNo}}") ||
			strings.Contains(cell, "{{CommentAuthor}}") ||
			strings.Contains(cell, "{{CommentCreated}}") ||
			strings.Contains(cell, "{{CommentCreatedEpoch}}") ||
			strings.Contains(cell, "{{CommentBody}}") {
			return true
		}
	}

	return false
}

func replaceCSVRowPlaceholders(
	row []string,
	placeholders map[string]string,
) []string {
	replaced := make([]string, len(row))

	for i, cell := range row {
		replaced[i] = replacePlaceholders(cell, placeholders)
	}

	return replaced
}

func replacePlaceholders(value string, placeholders map[string]string) string {
	result := value

	for placeholder, replacement := range placeholders {
		result = strings.ReplaceAll(result, placeholder, replacement)
	}

	return result
}

func cleanCell(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Join(strings.Fields(value), " ")
	return value
}

func epochCell(value int64) string {
	if value == 0 {
		return ""
	}

	return strconv.FormatInt(value, 10)
}
