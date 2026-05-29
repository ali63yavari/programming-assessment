package main

import "errors"

type WordSearchResult struct {
	RowIndex    int
	ColumnIndex int
	Letter      string
	Order       int
}

type WordSearch interface {
	Search(string) (bool, []WordSearchResult)
}

type wordSearch struct {
	matrix   [][]byte
	rowCount int
	colCount int
}

func NewWordSearch(matrix [][]byte, rowCount int, colCount int) (WordSearch, error) {
	if rowCount <= 0 || colCount <= 0 {
		return nil, errors.New("rowCount and colCount must be greater than zero")
	}
	if len(matrix) != rowCount {
		return nil, errors.New("matrix length must be equal to row count")
	}
	if len(matrix[0]) != colCount {
		return nil, errors.New("matrix length must be equal to col count")
	}

	return &wordSearch{
		matrix:   matrix,
		rowCount: rowCount,
		colCount: colCount,
	}, nil
}

func makeNewVisitedMatrix(rowCount int, colCount int) [][]bool {
	visited := make([][]bool, rowCount)
	for i := range visited {
		visited[i] = make([]bool, colCount)
	}
	return visited
}

func (ws *wordSearch) Search(word string) (bool, []WordSearchResult) {
	results := make(map[int]WordSearchResult)

	for i, row := range ws.matrix {
		for j, _ := range row {
			visited := makeNewVisitedMatrix(ws.rowCount, ws.colCount)
			found := ws.Dfs(i, j, 0, []byte(word), &visited, &results)
			if found {
				valuesResult := make([]WordSearchResult, 0)
				for _, v := range results {
					valuesResult = append(valuesResult, v)
				}
				return found, valuesResult
			}
		}
	}

	return false, make([]WordSearchResult, 0)
}

func (ws *wordSearch) Dfs(
	row, col, index int, word []byte, visited *[][]bool,
	results *map[int]WordSearchResult,
) bool {
	if index == len(word) {
		return true
	}

	if row < 0 || row >= ws.rowCount || col < 0 || col >= ws.
		colCount || (*visited)[row][col] || ws.matrix[row][col] != word[index] {
		return false
	}

	(*visited)[row][col] = true
	(*results)[index] = WordSearchResult{
		RowIndex:    row,
		ColumnIndex: col,
		Letter:      string(word[index]),
	}

	found := ws.Dfs(row, col+1, index+1, word, visited, results) ||
		ws.Dfs(row, col-1, index+1, word, visited, results) ||
		ws.Dfs(row+1, col, index+1, word, visited, results) ||
		ws.Dfs(row-1, col, index+1, word, visited, results)

	if !found {
		(*visited)[row][col] = false
		delete(*results, index)
	}

	return found
}
