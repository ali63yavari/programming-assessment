package main

import "errors"

type WordSearchResult struct {
	RowIndex    int
	ColumnIndex int
	Letter      string
}

type WordSearch interface {
	Search(string) (bool, []WordSearchResult)
}

type wordSearch struct {
	matrix   [][]byte
	rowCount int
	colCount int
}

func NewWordSearch(matrix [][]byte) (WordSearch, error) {
	rc := len(matrix)
	if rc <= 0 {
		return nil, errors.New("rowCount and colCount must be greater than zero")
	}
	cc := len(matrix[0])
	if cc <= 0 {
		return nil, errors.New("rowCount and colCount must be greater than zero")
	}

	return &wordSearch{
		matrix:   matrix,
		rowCount: rc,
		colCount: cc,
	}, nil
}

func makeNewVisitedMatrix(rowCount int, colCount int) [][]bool {
	visited := make([][]bool, rowCount)
	for i := range visited {
		visited[i] = make([]bool, colCount)
	}
	return visited
}

func (ws *wordSearch) checkLetterOccurance(word string) bool {
	bmap := make(map[byte]int)
	wmap := make(map[byte]int)

	for i := 0; i < len(word); i++ {
		wmap[word[i]]++
	}

	for i := 0; i < ws.rowCount; i++ {
		for j := 0; j < ws.colCount; j++ {
			bmap[ws.matrix[i][j]]++
		}
	}

	for k, v := range wmap {
		if bmap[k] < v {
			return false
		}
	}

	return true
}

func (ws *wordSearch) Search(word string) (bool, []WordSearchResult) {
	if len(word) == 0 {
		return false, []WordSearchResult{}
	}
	if len(word) > ws.rowCount*ws.colCount {
		return false, []WordSearchResult{}
	}
	if !ws.checkLetterOccurance(word) {
		return false, []WordSearchResult{}
	}

	results := make(map[int]WordSearchResult)

	for i, row := range ws.matrix {
		for j, _ := range row {
			visited := makeNewVisitedMatrix(ws.rowCount, ws.colCount)
			found := ws.Dfs(i, j, 0, []byte(word), &visited, &results)
			if found {
				valuesResult := make([]WordSearchResult, len(results))
				for i, v := range results {
					valuesResult[i] = v
				}
				return found, valuesResult
			}
		}
	}

	return false, []WordSearchResult{}
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
